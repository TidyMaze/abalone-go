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

func (d Direction) String() string {
	switch d {
	case TopRight:
		return "TopRight"
	case Right:
		return "Right"
	case BottomRight:
		return "BottomRight"
	case BottomLeft:
		return "BottomLeft"
	case Left:
		return "Left"
	case TopLeft:
		return "TopLeft"
	default:
		panic("Invalid direction")
	}
}
