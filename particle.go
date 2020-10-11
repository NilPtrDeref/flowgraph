package main

import (
	"github.com/hajimehoshi/ebiten"
	vec "github.com/ungerik/go3d/float64/vec2"
	"math"
	"math/rand"
)

type Particle struct {
	pos, vel *vec.T
	color    *ebiten.Image
}

// Create a particle with a random starting position/color.
// Particles created by this function have no starting velocity
func NewParticle() *Particle {
	p := &Particle{
		pos: &vec.T{
			Scale(rand.Float64(), 0, WINDOW_X, 0, 1),
			Scale(rand.Float64(), 0, WINDOW_Y, 0, 1),
		},
		vel:   &vec.T{0, 0},
		color: sprites[rand.Intn(len(sprites))],
	}

	return p
}

func (p *Particle) Update(grid [][]*Node) {
	width := int(WINDOW_X / SCALE)
	height := int(WINDOW_Y / SCALE)
	p.Move()

	x := math.Floor(p.pos[0] / SCALE)
	y := math.Floor(p.pos[1] / SCALE)

	if int(x) < width && int(y) < height &&
		int(x) >= 0 && int(y) >= 0 {
		p.Accelerate(grid[int(x)][int(y)].accl)
	}
}

// Move a particle based on it's current velocity.
func (p *Particle) Move() {
	p.pos = p.pos.Add(p.vel)

	// Wrap around if necessary
	if p.pos[0] > WINDOW_X {
		diff := p.pos[0] - WINDOW_X
		p.pos[0] = diff
	}
	if p.pos[0] < 0 {
		p.pos[0] = WINDOW_X + p.pos[0]
	}
	if p.pos[1] > WINDOW_Y {
		diff := p.pos[1] - WINDOW_Y
		p.pos[1] = diff
	}
	if p.pos[1] < 0 {
		p.pos[1] = WINDOW_Y + p.pos[1]
	}
}

// Accelerate the velocity of a particle with a given acceleration vector.
func (p *Particle) Accelerate(accl *vec.T) {
	p.vel = p.vel.Add(accl)
	l := p.vel.Length()

	// Cap velocity to MAX_VELOCITY
	if l > MAX_VELOCITY {
		angle := p.vel.Angle()
		p.vel = &vec.T{1, 0}
		p.vel = p.vel.Rotate(angle)
		p.vel = p.vel.Scale(MAX_VELOCITY)
	}
}

// Draw a particle to the screen
func (p *Particle) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(PARTICLE_WIDTH, PARTICLE_WIDTH)
	op.GeoM.Translate(p.pos[0], p.pos[1])
	screen.DrawImage(p.color, op)
}
