package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	vec "github.com/woodywood117/vector"
	"golang.org/x/image/colornames"
)

type Node struct {
	pos  *vec.Vec
	accl *vec.Vec
}

// Create a node with no starting acceleration
func NewNode(x, y float64) *Node {
	pos := vec.New(x, y)
	n := &Node{
		pos:  pos,
		accl: vec.Zero(),
	}
	return n
}

// Draws a vector representing the acceleration within the node.
func (n *Node) Draw(dst *ebiten.Image) {
	from := n.pos
	to := vec.Add(n.pos, vec.Multiply(n.accl, 100))
	ebitenutil.DrawLine(dst, from.X, from.Y, to.X, to.Y, colornames.White)
}

// Updates the acceleration using opensimplex noise
func (n *Node) Update() {
	nval := noise.Eval3(
		n.pos.X/SCALE*SIMILARITY,
		n.pos.Y/SCALE*SIMILARITY,
		depth,
	)

	angle := Scale(nval, MIN_ANGLE, MAX_ANGLE, 0, 1)
	n.accl = vec.Unit(angle)
	n.accl.Multiply(ACCELERATION_MAGNITUDE)
}
