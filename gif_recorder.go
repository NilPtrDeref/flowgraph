package main

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"sync"

	"github.com/hajimehoshi/ebiten"
)

type Recorder struct {
	Writer       io.Writer
	frames       int
	gif          *gif.GIF
	currentFrame int
	wg           sync.WaitGroup
}

func (r *Recorder) delay() int {
	delay := 100 / ebiten.MaxTPS()
	if delay < 2 {
		return 2
	}
	return delay
}

func (r *Recorder) Update(screen *ebiten.Image) error {
	if r.currentFrame == r.frames {
		return nil
	}

	s := image.NewNRGBA(screen.Bounds())
	draw.Draw(s, s.Bounds(), screen, screen.Bounds().Min, draw.Src)

	img := image.NewPaletted(s.Bounds(), palette.Plan9)
	f := r.currentFrame
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		draw.FloydSteinberg.Draw(img, img.Bounds(), s, s.Bounds().Min)
		r.gif.Image[f] = img
		r.gif.Delay[f] = r.delay()
	}()

	r.currentFrame++
	if r.currentFrame == r.frames {
		r.wg.Wait()
		if err := gif.EncodeAll(r.Writer, r.gif); err != nil {
			return err
		}
	}
	return nil
}

func NewRecorder(out io.Writer, frames int) *Recorder {
	return &Recorder{
		Writer: out,
		frames: frames,
		gif: &gif.GIF{
			Image:     make([]*image.Paletted, frames),
			Delay:     make([]int, frames),
			LoopCount: -1,
		},
	}
}
