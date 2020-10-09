package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Node struct {
	pos  pixel.Vec
	accl pixel.Vec
}

// Create a node with no starting acceleration
func NewNode(x, y float64) *Node {
	pos := pixel.V(x, y)
	n := &Node{
		pos:  pos,
		accl: pixel.V(0, 0),
	}
	return n
}

// Draws a vector representing the acceleration within the node.
func (n *Node) Draw(win pixel.Target, imd *imdraw.IMDraw) {
	imd.Color = colornames.White
	imd.Push(n.pos, n.pos.Add(n.accl.Normal().Scaled(SCALE*50)))
	imd.Line(2)
	imd.Draw(win)
	imd.Reset()
	imd.Clear()
}

// Updates the acceleration using opensimplex noise
func (n *Node) Update() {
	nval := noise.Eval3(
		n.pos.X/SCALE*SIMILARITY,
		n.pos.Y/SCALE*SIMILARITY,
		depth,
	)

	angle := Scale(nval, MIN_ANGLE, MAX_ANGLE, 0, 1)
	n.accl = pixel.Unit(angle).Scaled(0.01)
}
