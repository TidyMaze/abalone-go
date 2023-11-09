package engine

import (
	"abalone-go/helpers"
	"fmt"
)

type Coord3D struct {
	X int8
	Y int8
	Z int8
}

func (c Coord3D) String() string {
	return fmt.Sprintf("(x: %d, y: %d, z: %d)", c.X, c.Y, c.Z)
}

type Coord2D struct {
	X int8
	Y int8
}

func (c Coord2D) To3D() Coord3D {
	return Coord3D{c.X, -c.X - c.Y, c.Y}
}

func (c Coord3D) To2D() Coord2D {
	return Coord2D{c.X, c.Z}
}

func (c Coord3D) Add(direction Direction) Coord3D {
	switch direction {
	case TopRight:
		return Coord3D{c.X + 1, c.Y, c.Z - 1}
	case Right:
		return Coord3D{c.X + 1, c.Y - 1, c.Z}
	case BottomRight:
		return Coord3D{c.X, c.Y - 1, c.Z + 1}
	case BottomLeft:
		return Coord3D{c.X - 1, c.Y, c.Z + 1}
	case Left:
		return Coord3D{c.X - 1, c.Y + 1, c.Z}
	case TopLeft:
		return Coord3D{c.X, c.Y + 1, c.Z - 1}
	default:
		panic("Invalid direction")
	}

	return Coord3D{}
}

func IsValidCoord(c Coord3D) bool {
	return helpers.Between(c.X, -4, 4) &&
		helpers.Between(c.Y, -4, 4) &&
		helpers.Between(c.Z, -4, 4) &&
		c.X+c.Y+c.Z == 0
}
