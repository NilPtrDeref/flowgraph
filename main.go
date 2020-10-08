package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/ojrac/opensimplex-go"
	"image/color"
	"math"
	"math/rand"
	"time"
)

/* Global constants */
// Window width and height
const WINDOW_X, WINDOW_Y float64 = 600, 600

// Grid element width/height
const SCALE float64 = 5

// Minimum time per frame
const MIN_FRAMETIME time.Duration = 0

// Maximum velocity of the particles
const MAX_VELOCITY float64 = 0.5

// How fast the depth changes
const DELTA_Z = 0.00001

// How close the x/y value neighbors should be (for simplex noise)
// Closer is more similar.
const SIMILARITY = 0.01

// Maximum angle for acceleration within a node.
const MAX_ANGLE = math.Pi * 2

var noise = opensimplex.NewNormalized(time.Now().Unix())
var depth float64 = 0
var colors = []color.RGBA{
	{0xEE, 0xB4, 0xB3, 0xFF},
	{0xC1, 0x79, 0xB9, 0xFF},
	{0xA4, 0x2C, 0xD6, 0xFF},
	{0x50, 0x22, 0x74, 0xFF},
	{0x2F, 0x24, 0x2C, 0xFF},
}

// Scale a floating point number to fit between two values.
// minallow and maxallow are the new bounds.
// min and max are the current min and max values.
func Scale(val, minallow, maxallow, min, max float64) float64 {
	return (maxallow-minallow)*(val-min)/(max-min) + minallow
}

func run() {
	// Set up window
	cfg := pixelgl.WindowConfig{
		Title:  "Flowgraph",
		Bounds: pixel.R(0, 0, WINDOW_X, WINDOW_Y),
		VSync:  false,
		Icon: []pixel.Picture{
			pixel.MakePictureData(
				pixel.R(0, 0, 16, 16),
			),
		},
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

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
	for i := 0; i < 200; i++ {
		particles = append(particles, NewParticle())
	}

	// Create batch and drawer for fast drawing
	batch := pixel.NewBatch(
		&pixel.TrianglesData{},
		pixel.MakePictureData(pixel.R(0, 0, WINDOW_X, WINDOW_Y)),
	)
	imd := imdraw.New(nil)

	// Flow loop
	t := time.Now()
	for !win.Closed() {
		// Quit if escape pressed
		if win.JustPressed(pixelgl.KeyEscape) {
			return
		}

		// Enforce minimum frametime
		dt := time.Since(t)
		if dt < MIN_FRAMETIME {
			continue
		}

		// Uncomment to clear screen after every frame
		//win.Clear(colornames.Black)

		// Update nodes' acceleration
		for x := range grid {
			for y := range grid[x] {
				// Uncomment to draw vectors showing acceleration
				// grid[x][y].Draw(batch, imd)

				grid[x][y].Update()
			}
		}
		depth += DELTA_Z

		//Move, update the acceleration of and draw the particles
		for _, particle := range particles {
			particle.Move()

			x := math.Floor(particle.pos.X / SCALE)
			y := math.Floor(particle.pos.Y / SCALE)
			if int(x) < width && int(y) < height &&
				int(x) > 0 && int(y) > 0 {
				particle.Update(grid[int(x)][int(y)].accl)
			}

			particle.Draw(batch, imd)
		}

		// Batch draw to the screen
		batch.Draw(win)
		batch.Clear()

		win.Update()
		t = time.Now()
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	pixelgl.Run(run)
}
