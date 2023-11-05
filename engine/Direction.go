package engine

type Direction int

const (
	TopRight    Direction = 0
	Right                 = iota
	BottomRight           = iota
	BottomLeft            = iota
	Left                  = iota
	TopLeft               = iota
)

var Directions = [6]Direction{TopRight, Right, BottomRight, BottomLeft, Left, TopLeft}
