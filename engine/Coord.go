package engine

import "abalone-go/helpers"

type Coord3D struct {
	X int
	Y int
	Z int
}

type Coord2D struct {
	X int
	Y int
}

func (c Coord2D) To3D() Coord3D {
	return Coord3D{c.X, -c.X - c.Y, c.Y}
}

func (c Coord3D) To2D() Coord2D {
	return Coord2D{c.X, c.Y}
}

func isValidCoord(c Coord3D) bool {
	return helpers.Between(c.X, -3, 3) &&
		helpers.Between(c.Y, -3, 3) &&
		helpers.Between(c.Z, -3, 3) &&
		c.X+c.Y+c.Z == 0
}
