package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/ojrac/opensimplex-go"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"math/rand"
	"sync"
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

// How fast the depth changes.
// Bigger is more change.
const DELTA_Z = 0.0001

// How close the x/y value neighbors should be (for simplex noise)
// Smaller is more similar.
const SIMILARITY = 0.09

// Number of particles to flow.
const NUM_PARTICLES = 1000

// Number of threads to render the particles with.
// DO NOT MAKE THIS <= 0, IT WILL CRASH THE PROGRAM!
const PARTICLE_THREADS = 10

// Maximum angle for acceleration within a node.
const MAX_ANGLE = math.Pi * 1
const MIN_ANGLE = math.Pi * -1

const PARTICLE_WIDTH = 2

var noise = opensimplex.NewNormalized(time.Now().Unix())
var depth float64 = 0
var colors = []color.RGBA{
	{0xAF, 0xD2, 0xE9, 0xFF},
	{0x9D, 0x96, 0xB8, 0xFF},
	{0x9A, 0x71, 0x97, 0xFF},
	{0x88, 0x61, 0x76, 0xFF},
	{0x7C, 0x58, 0x69, 0xFF},
}
var picture *pixel.PictureData
var sprites []*pixel.Sprite

func init() {
	l := len(colors)
	picture = pixel.MakePictureData(pixel.R(0, 0, PARTICLE_WIDTH, PARTICLE_WIDTH*float64(l)))

	for i := range picture.Pix {
		picture.Pix[i] = colors[i/int(PARTICLE_WIDTH*PARTICLE_WIDTH)]
	}

	for c := range colors {
		sprites = append(sprites, pixel.NewSprite(
			picture,
			pixel.R(
				0, float64(c)*PARTICLE_WIDTH,
				PARTICLE_WIDTH, float64(c)*PARTICLE_WIDTH+PARTICLE_WIDTH,
			),
		))
	}
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
	for i := 0; i < NUM_PARTICLES; i++ {
		particles = append(particles, NewParticle())
	}

	// Create batch and drawer for fast drawing
	batches := [PARTICLE_THREADS]*pixel.Batch{}
	for i := 0; i < PARTICLE_THREADS; i++ {
		batches[i] = pixel.NewBatch(
			&pixel.TrianglesData{},
			picture,
		)
	}

	// Flow loop
	t := time.Now()
	for !win.Closed() {
		// Quit if escape pressed
		if win.JustPressed(pixelgl.KeyEscape) {
			return
		}

		// Restart graph
		if win.JustPressed(pixelgl.KeyR) {
			// Clear screen
			win.Clear(colornames.Black)

			noise = opensimplex.NewNormalized(time.Now().Unix())

			// Reset grid
			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {
					grid[x][y] = NewNode(
						float64(x)*SCALE+SCALE/2, float64(y)*SCALE+SCALE/2,
					)
				}
			}

			// Reset particles
			particles = []*Particle{}
			for i := 0; i < NUM_PARTICLES; i++ {
				particles = append(particles, NewParticle())
			}
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
				// grid[x][y].Draw(batches[0], imd[0])

				grid[x][y].Update()
			}
		}
		depth += DELTA_Z

		//Move and update the acceleration of the particles
		// Multithreaded to prevent slowdown if possible
		wg := &sync.WaitGroup{}
		for i := 0; i < PARTICLE_THREADS; i++ {
			wg.Add(1)
			UpdateAndDrawParticles(i, particles, wg, batches[i], grid)
		}
		wg.Wait()

		// Batch draw to the screen
		for _, batch := range batches {
			batch.Draw(win)
			batch.Clear()
		}

		win.Update()
		t = time.Now()
	}
}

func UpdateAndDrawParticles(i int, particles []*Particle, wg *sync.WaitGroup,
	batch *pixel.Batch, grid [][]*Node) {
	width := int(WINDOW_X / SCALE)
	height := int(WINDOW_Y / SCALE)
	start := i * (NUM_PARTICLES / PARTICLE_THREADS)
	stop := start + (NUM_PARTICLES / PARTICLE_THREADS)
	for j := start; j < stop; j++ {
		particles[j].Move()

		x := math.Floor(particles[j].pos.X / SCALE)
		y := math.Floor(particles[j].pos.Y / SCALE)

		if int(x) < width && int(y) < height &&
			int(x) >= 0 && int(y) >= 0 {
			particles[j].Update(grid[int(x)][int(y)].accl)
		}

		particles[j].Draw(batch)
	}

	wg.Done()
}

func main() {
	rand.Seed(time.Now().Unix())
	pixelgl.Run(run)
}
