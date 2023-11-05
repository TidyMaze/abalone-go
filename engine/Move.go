package engine

import "fmt"

type Move interface{}

type PushLine struct {
	From      Coord3D
	Direction Direction
}

func (p PushLine) String() string {
	return fmt.Sprintf("Action: PushLine, From: %v, Direction: %v", p.From, p.Direction)
}
