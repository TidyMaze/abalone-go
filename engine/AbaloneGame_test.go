package engine

import (
	"abalone-go/helpers"
	"testing"
)

func TestPushSingle(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{0, 0, 0}, 1)

	gameCopy := game.Copy()
	gameCopy.Push(Coord3D{0, 0, 0}, Right)

	expected := NewGame()
	expected.SetGrid(Coord3D{1, -1, 0}, 1)

	helpers.AssertEqual(showGrid(expected.grid), showGrid(gameCopy.grid))
}

func TestPushSingleBlocked(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{0, 0, 0}, 1)
	game.SetGrid(Coord3D{1, -1, 0}, 2)

	gameCopy := game.Copy()
	err := gameCopy.Push(Coord3D{0, 0, 0}, Right)

	helpers.AssertEqual("destination is not empty", err.Error())
}
