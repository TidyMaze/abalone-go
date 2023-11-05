package engine

type Move interface{}

type PushLine struct {
	From      Coord3D
	Direction Direction
}
