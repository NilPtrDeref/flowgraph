package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	vec "github.com/woodywood117/vector"
)

type Particle struct {
	pos, vel *vec.Vec
	color    *ebiten.Image
	history  []vec.Vec
}

// Create a particle with a random starting position/color.
// Particles created by this function have no starting velocity
func NewParticle() *Particle {
	p := &Particle{
		pos: vec.New(
			Scale(rand.Float64(), 0, WINDOW_X, 0, 1),
			Scale(rand.Float64(), 0, WINDOW_Y, 0, 1),
		),
		vel:   vec.Zero(),
		color: sprites[rand.Intn(len(sprites))],
	}

	return p
}

func (p *Particle) Update(grid [][]*Node) {
	width := int(WINDOW_X / SCALE)
	height := int(WINDOW_Y / SCALE)
	p.Move()

	x := math.Floor(p.pos.X / SCALE)
	y := math.Floor(p.pos.Y / SCALE)

	if int(x) < width && int(y) < height &&
		int(x) >= 0 && int(y) >= 0 {
		p.Accelerate(grid[int(x)][int(y)].accl)
	}
}

// Move a particle based on it's current velocity.
func (p *Particle) Move() {
	if DRAW_TRAIL {
		p.history = append(p.history, *p.pos)
		if len(p.history) > MAX_TRAIL_LEN {
			p.history = p.history[1:]
		}
	}

	p.pos.Add(p.vel)

	// Wrap around if necessary
	if p.pos.X > WINDOW_X {
		diff := p.pos.X - WINDOW_X
		p.pos.X = diff
	}
	if p.pos.X < 0 {
		p.pos.X = WINDOW_X + p.pos.X
	}
	if p.pos.Y > WINDOW_Y {
		diff := p.pos.Y - WINDOW_Y
		p.pos.Y = diff
	}
	if p.pos.Y < 0 {
		p.pos.Y = WINDOW_Y + p.pos.Y
	}
}

// Accelerate the velocity of a particle with a given acceleration vector.
func (p *Particle) Accelerate(accl *vec.Vec) {
	p.vel.Add(accl)
	l := p.vel.Magnitude()

	// Cap velocity to MAX_VELOCITY
	if l > MAX_VELOCITY {
		angle := p.vel.Angle()
		p.vel = vec.New(1, 0)
		p.vel.Rotate(angle)
		p.vel.Multiply(MAX_VELOCITY)
	}
}

// Draw a particle to the screen
func (p *Particle) Draw(screen *ebiten.Image) {
	// Draw the trail
	if DRAW_TRAIL {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(PARTICLE_WIDTH, PARTICLE_WIDTH)
		for trail := len(p.history) - 1; trail > -1; trail-- {
			op.ColorM.Scale(1, 1, 1, 0.85)
			op.GeoM.Translate(p.history[trail].X, p.history[trail].Y)
			screen.DrawImage(p.color, op)
			op.GeoM.Translate(-p.history[trail].X, -p.history[trail].Y)
		}
	}

	// Draw the particle itself
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(PARTICLE_WIDTH, PARTICLE_WIDTH)
	op.GeoM.Translate(p.pos.X, p.pos.Y)
	screen.DrawImage(p.color, op)
}
