package main

import (
	"github.com/faiface/pixel"
	"math/rand"
)

type Particle struct {
	pos, vel pixel.Vec
	color    *pixel.Sprite
}

// Create a particle with a random starting position/color.
// Particles created by this function have no starting velocity
func NewParticle() *Particle {
	p := &Particle{
		pos: pixel.V(
			Scale(rand.Float64(), 0, WINDOW_X, 0, 1),
			Scale(rand.Float64(), 0, WINDOW_Y, 0, 1),
		),
		vel:   pixel.V(0, 0),
		color: sprites[rand.Intn(len(sprites))],
	}

	return p
}

// Move a particle based on it's current velocity.
func (p *Particle) Move() {
	p.pos = p.pos.Add(p.vel)

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

// Update the velocity of a particle with a given acceleration vector.
func (p *Particle) Update(accl pixel.Vec) {
	p.vel = p.vel.Add(accl)
	l := p.vel.Len()

	// Cap velocity to MAX_VELOCITY
	if l > MAX_VELOCITY {
		p.vel = pixel.Unit(p.vel.Angle())
		p.vel = p.vel.Scaled(MAX_VELOCITY)
	}
}

// Draw a particle to the screen
func (p *Particle) Draw(win pixel.Target) {
	p.color.Draw(win, pixel.IM.Moved(p.pos))
}
