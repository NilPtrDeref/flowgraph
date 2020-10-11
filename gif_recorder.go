package main

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten"
)

type Recorder struct {
	Writer       io.Writer
	frames       int
	gif          *gif.GIF
	currentFrame int
	wg           sync.WaitGroup
	pipe         chan order
}

type order struct {
	frame int
	img   *ebiten.Image
}

func NewRecorder(out io.Writer, frames int) *Recorder {
	r := &Recorder{
		Writer: out,
		frames: frames,
		gif: &gif.GIF{
			Image:     make([]*image.Paletted, frames),
			Delay:     make([]int, frames),
			LoopCount: -1,
		},
		pipe: make(chan order, frames),
	}

	go r.Record()

	return r
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

	clone, _ := ebiten.NewImageFromImage(screen, ebiten.FilterDefault)
	r.pipe <- order{r.currentFrame, clone}

	r.currentFrame++
	if r.currentFrame == r.frames {
		close(r.pipe)
	}
	return nil
}

func (r *Recorder) Record() {
	for ord := range r.pipe {
		r.wg.Add(1)

		go func(ord order) {
			img := image.NewPaletted(ord.img.Bounds(), palette.Plan9)
			draw.FloydSteinberg.Draw(img, img.Bounds(), ord.img, ord.img.Bounds().Min)
			r.gif.Image[ord.frame] = img
			r.gif.Delay[ord.frame] = 2
			r.wg.Done()
		}(ord)
	}

	r.wg.Wait()
	err := gif.EncodeAll(r.Writer, r.gif)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done outputting gif")
	os.Exit(0)
}
