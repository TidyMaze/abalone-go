package engine

import (
	"abalone-go/helpers"
	"testing"
)

func TestPushSingle(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{0, 0, 0}, 1)

	gameCopy := game.Copy()
	_, err := gameCopy.Push(Coord3D{0, 0, 0}, Right)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := NewGame()
	expected.SetGrid(Coord3D{1, -1, 0}, 1)

	helpers.AssertEqual(showGrid(expected.grid), showGrid(gameCopy.grid))
}

func TestPushSingleBlocked(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{0, 0, 0}, 1)
	game.SetGrid(Coord3D{1, -1, 0}, 2)

	gameCopy := game.Copy()
	_, err := gameCopy.Push(Coord3D{0, 0, 0}, Right)

	helpers.AssertEqual("not enough marbles to push enemy (got 1, need 2)", err.Error())
}

func TestPushTwo(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{0, 0, 0}, 1)
	game.SetGrid(Coord3D{1, -1, 0}, 1)

	gameCopy := game.Copy()
	_, err := gameCopy.Push(Coord3D{0, 0, 0}, Right)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := NewGame()
	expected.SetGrid(Coord3D{1, -1, 0}, 1)
	expected.SetGrid(Coord3D{2, -2, 0}, 1)

	helpers.AssertEqual(showGrid(expected.grid), showGrid(gameCopy.grid))
}

func TestPushTwoBlocked(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{0, 0, 0}, 1)
	game.SetGrid(Coord3D{1, -1, 0}, 1)
	game.SetGrid(Coord3D{2, -2, 0}, 2)
	game.SetGrid(Coord3D{3, -3, 0}, 2)

	gameCopy := game.Copy()
	_, err := gameCopy.Push(Coord3D{0, 0, 0}, Right)

	helpers.AssertEqual("not enough marbles to push enemy (got 2, need 3)", err.Error())
}

func TestPushThree(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{0, 0, 0}, 1)
	game.SetGrid(Coord3D{1, -1, 0}, 1)
	game.SetGrid(Coord3D{2, -2, 0}, 1)

	gameCopy := game.Copy()
	_, err := gameCopy.Push(Coord3D{0, 0, 0}, Right)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := NewGame()
	expected.SetGrid(Coord3D{1, -1, 0}, 1)
	expected.SetGrid(Coord3D{2, -2, 0}, 1)
	expected.SetGrid(Coord3D{3, -3, 0}, 1)

	helpers.AssertEqual(showGrid(expected.grid), showGrid(gameCopy.grid))
}

func TestPushThreeBlocked(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{-2, 2, 0}, 1)
	game.SetGrid(Coord3D{-1, 1, 0}, 1)
	game.SetGrid(Coord3D{0, 0, 0}, 1)
	game.SetGrid(Coord3D{1, -1, 0}, 2)
	game.SetGrid(Coord3D{2, -2, 0}, 2)
	game.SetGrid(Coord3D{3, -3, 0}, 2)

	gameCopy := game.Copy()
	_, err := gameCopy.Push(Coord3D{-2, 2, 0}, Right)

	helpers.AssertEqual("too many enemy marbles to push (max 2, got 3)", err.Error())
}

func TestTwoPushOne(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{0, 0, 0}, 1)
	game.SetGrid(Coord3D{1, -1, 0}, 1)
	game.SetGrid(Coord3D{2, -2, 0}, 2)

	gameCopy := game.Copy()
	_, err := gameCopy.Push(Coord3D{0, 0, 0}, Right)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := NewGame()
	expected.SetGrid(Coord3D{1, -1, 0}, 1)
	expected.SetGrid(Coord3D{2, -2, 0}, 1)
	expected.SetGrid(Coord3D{3, -3, 0}, 2)

	helpers.AssertEqual(showGrid(expected.grid), showGrid(gameCopy.grid))
}

func TestThreePushOne(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{-1, 1, 0}, 1)
	game.SetGrid(Coord3D{0, 0, 0}, 1)
	game.SetGrid(Coord3D{1, -1, 0}, 1)
	game.SetGrid(Coord3D{2, -2, 0}, 2)

	gameCopy := game.Copy()
	_, err := gameCopy.Push(Coord3D{-1, 1, 0}, Right)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := NewGame()
	expected.SetGrid(Coord3D{0, 0, 0}, 1)
	expected.SetGrid(Coord3D{1, -1, 0}, 1)
	expected.SetGrid(Coord3D{2, -2, 0}, 1)
	expected.SetGrid(Coord3D{3, -3, 0}, 2)

	helpers.AssertEqual(showGrid(expected.grid), showGrid(gameCopy.grid))
}

func TestThreePushTwo(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{-2, 2, 0}, 1)
	game.SetGrid(Coord3D{-1, 1, 0}, 1)
	game.SetGrid(Coord3D{0, 0, 0}, 1)
	game.SetGrid(Coord3D{1, -1, 0}, 2)
	game.SetGrid(Coord3D{2, -2, 0}, 2)

	gameCopy := game.Copy()
	_, err := gameCopy.Push(Coord3D{-2, 2, 0}, Right)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := NewGame()
	expected.SetGrid(Coord3D{-1, 1, 0}, 1)
	expected.SetGrid(Coord3D{0, 0, 0}, 1)
	expected.SetGrid(Coord3D{1, -1, 0}, 1)
	expected.SetGrid(Coord3D{2, -2, 0}, 2)
	expected.SetGrid(Coord3D{3, -3, 0}, 2)

	helpers.AssertEqual(showGrid(expected.grid), showGrid(gameCopy.grid))
}

func TestSandwich(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{-2, 2, 0}, 1)
	game.SetGrid(Coord3D{-1, 1, 0}, 1)
	game.SetGrid(Coord3D{0, 0, 0}, 2)
	game.SetGrid(Coord3D{1, -1, 0}, 1)

	gameCopy := game.Copy()
	_, err := gameCopy.Push(Coord3D{-2, 2, 0}, Right)

	helpers.AssertEqual("my marbles are sandwiching enemy marbles", err.Error())
}

func TestPushSingleCaptured(t *testing.T) {
	game := NewGame()
	game.SetGrid(Coord3D{1, -1, 0}, 1)
	game.SetGrid(Coord3D{2, -2, 0}, 1)
	game.SetGrid(Coord3D{3, -3, 0}, 2)

	gameCopy := game.Copy()
	captured, err := gameCopy.Push(Coord3D{1, -1, 0}, Right)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	helpers.AssertEqual(true, captured)

	expected := NewGame()
	expected.SetGrid(Coord3D{2, -2, 0}, 1)
	expected.SetGrid(Coord3D{3, -3, 0}, 1)
}
