package engine

import "fmt"

type Move struct {
	At Coord2D
}

func (m Move) String() string {
	return fmt.Sprintf("Move(%v)", m.At)
}
