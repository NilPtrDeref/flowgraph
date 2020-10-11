package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	vec "github.com/ungerik/go3d/float64/vec2"
	"golang.org/x/image/colornames"
)

type Node struct {
	pos  *vec.T
	accl *vec.T
}

// Create a node with no starting acceleration
func NewNode(x, y float64) *Node {
	pos := &vec.T{x, y}
	n := &Node{
		pos:  pos,
		accl: &vec.T{0, 0},
	}
	return n
}

// Draws a vector representing the acceleration within the node.
func (n *Node) Draw(dst *ebiten.Image) {
	from := n.pos
	to := vec.Add(n.pos, n.accl.Scale(100))
	ebitenutil.DrawLine(dst, from[0], from[1], to[0], to[1], colornames.White)
}

// Updates the acceleration using opensimplex noise
func (n *Node) Update() {
	nval := noise.Eval3(
		n.pos[0]/SCALE*SIMILARITY,
		n.pos[1]/SCALE*SIMILARITY,
		depth,
	)

	angle := Scale(nval, MIN_ANGLE, MAX_ANGLE, 0, 1)
	n.accl = &vec.T{1, 0}
	n.accl = n.accl.Rotate(angle).Scale(ACCELERATION_MAGNITUDE)
}
