package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/ojrac/opensimplex-go"
)

/* Global constants */
// Window width and height
const WINDOW_X, WINDOW_Y float64 = 1200, 800

// Grid element width/height
const SCALE float64 = 5

// Maximum velocity of the particles
const MAX_VELOCITY float64 = 1.5
const ACCELERATION_MAGNITUDE = 0.1

// Whether the program should clear the frame each time
const CLEAR_EACH_FRAME = true

// How fast the depth changes.
// Bigger is more change.
const DELTA_Z = 0.001

// How close the x/y value neighbors should be (for simplex noise)
// Smaller is more similar.
const SIMILARITY = 0.1

// Number of particles to flow.
const NUM_PARTICLES = 20000

// Number of threads to render the particles with.
// DO NOT MAKE THIS <= 0, IT WILL CRASH THE PROGRAM!
const PARTICLE_THREADS = 10

// Maximum angle for acceleration within a node.
const MAX_ANGLE = math.Pi * 2
const MIN_ANGLE = math.Pi * -2

const PARTICLE_WIDTH = 2

const COLORSCHEME_SCALE = 50

// Number of frames to record.
// If it is greater than 0, the program will generate a gif at output.gif
// and then exit.
const RECORD_FRAMES = 0

var noise = opensimplex.NewNormalized(time.Now().Unix())
var depth float64 = 0
var colors = []color.RGBA{
	{0xAF, 0xD2, 0xE9, 0xFF},
	{0x9D, 0x96, 0xB8, 0xFF},
	{0x9A, 0x71, 0x97, 0xFF},
	{0x88, 0x61, 0x76, 0xFF},
	{0x7C, 0x58, 0x69, 0xFF},
}
var picture *ebiten.Image
var sprites []*ebiten.Image

func init() {
	l := len(colors)
	picture, _ = ebiten.NewImage(1, l, ebiten.FilterDefault)

	for i := 0; i < l; i++ {
		picture.Set(0, i, colors[i])
		sprites = append(sprites, picture.SubImage(
			image.Rect(0, i, 1, i+1),
		).(*ebiten.Image))
	}
}

// Scale a floating point number to fit between two values.
// minallow and maxallow are the new bounds.
// min and max are the current min and max values.
func Scale(val, minallow, maxallow, min, max float64) float64 {
	return (maxallow-minallow)*(val-min)/(max-min) + minallow
}

type Grid struct {
	width, height int
	grid          [][]*Node
	particles     []*Particle
	recorder      *Recorder
}

func NewGrid() *Grid {
	// Create the grid that contains all of the nodes
	width := int(WINDOW_X / SCALE)
	height := int(WINDOW_Y / SCALE)
	var grid [][]*Node
	for x := 0; x < width; x++ {
		grid = append(grid, make([]*Node, height))
		for y := 0; y < height; y++ {
			grid[x][y] = NewNode(
				float64(x)*SCALE+SCALE/2, float64(y)*SCALE+SCALE/2,
			)
		}
	}

	// Create particles that accelerate based on the vector of the node they are
	// currently in.
	var particles []*Particle
	for i := 0; i < NUM_PARTICLES; i++ {
		particles = append(particles, NewParticle())
	}

	return &Grid{
		width:     width,
		height:    height,
		grid:      grid,
		particles: particles,
	}
}

func (g *Grid) Restart() {
	noise = opensimplex.NewNormalized(time.Now().Unix())

	//Reset grid
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			g.grid[x][y] = NewNode(
				float64(x)*SCALE+SCALE/2,
				float64(y)*SCALE+SCALE/2,
			)
		}
	}

	//Reset particles
	g.particles = []*Particle{}
	for i := 0; i < NUM_PARTICLES; i++ {
		g.particles = append(g.particles, NewParticle())
	}
}

func (g *Grid) Update(screen *ebiten.Image) error {
	// Quit if escape pressed
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("Exiting game")
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.Restart()
	}

	// Accelerate nodes' acceleration
	for x := range g.grid {
		for y := range g.grid[x] {
			g.grid[x][y].Update()
		}
	}
	depth += DELTA_Z

	//Move and update the acceleration of the particles
	//Multithreaded to prevent slowdown if possible
	wg := &sync.WaitGroup{}
	for i := 0; i < PARTICLE_THREADS; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			start := i * (NUM_PARTICLES / PARTICLE_THREADS)
			stop := start + (NUM_PARTICLES / PARTICLE_THREADS)
			for j := start; j < stop; j++ {
				g.particles[j].Update(g.grid)
			}

			wg.Done()
		}(i, wg)
	}
	wg.Wait()

	return nil
}

func (g *Grid) Draw(screen *ebiten.Image) {
	//Uncomment to draw vectors showing acceleration
	//for x := range g.grid {
	//	for y := range g.grid[x] {
	//		g.grid[x][y].Draw(screen)
	//	}
	//}

	// Uncommment to draw colorscheme in top left
	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Scale(COLORSCHEME_SCALE,COLORSCHEME_SCALE)
	//op.GeoM.Translate(10, 10)
	//screen.DrawImage(picture, op)

	for _, particle := range g.particles {
		particle.Draw(screen)
	}

	// Uncomment to draw fps in top left
	//fps := fmt.Sprintf("Current FPS: %.1f", ebiten.CurrentFPS())
	//ebitenutil.DebugPrint(screen, fps)

	if g.recorder != nil {
		_ = g.recorder.Update(screen)
	}
}

func (g *Grid) Layout(_, _ int) (int, int) {
	return int(WINDOW_X), int(WINDOW_Y)
}

func main() {
	ebiten.SetWindowSize(int(WINDOW_X), int(WINDOW_Y))
	ebiten.SetWindowTitle("Flowgraph")
	ebiten.SetMaxTPS(60)
	ebiten.SetScreenClearedEveryFrame(CLEAR_EACH_FRAME)
	ebiten.SetRunnableOnUnfocused(true)

	grid := NewGrid()
	if RECORD_FRAMES > 0 {
		file, err := os.Create("output.gif")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		grid.recorder = NewRecorder(file, RECORD_FRAMES)
	}

	if err := ebiten.RunGame(grid); err != nil {
		log.Fatal(err)
	}
}
