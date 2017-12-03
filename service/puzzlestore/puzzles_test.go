package puzzlestore

import (
	"testing"
)

func TestGetPuzzle_ForAGivenUserAndTeam_PuzzleIsNotEmpty(t *testing.T) {
	puzzles := NewPuzzles(0)

	myPuzzle := puzzles.Get(0)

	if myPuzzle == "" {
		t.Fatalf("The returned puzzle was empty.")
	}
}

func TestGetPuzzle_GetTwoPuzzles_NextPuzzleIsDifferent(t *testing.T) {
	puzzles := NewPuzzles(0)

	puzzle0 := puzzles.Get(0)
	puzzle1 := puzzles.Get(1)

	if puzzle0 == puzzle1 {
		t.Fatalf("Puzzles at index 0 and 1 are the same: %s", puzzle0)
	}
}

func TestGetPuzzle_GetIndex10Directly_NoNeedToGetPrecedingIndicesFirst(t *testing.T) {
	puzzles := NewPuzzles(0)

	puzzle10 := puzzles.Get(10)
	puzzle10Again := puzzles.Get(10)

	if puzzle10 != puzzle10Again {
		t.Fatalf("Puzzle at index 10 should be the same, but were: %s and %s",
			puzzle10, puzzle10Again)
	}
}

//
// Test pseudo-randomness
//

func TestPseudoRandomness_Index0ForDifferentSeeds_PuzzlesAreDifferent(t *testing.T) {
	puzzles1 := NewPuzzles(0)
	puzzle1 := puzzles1.Get(0)

	puzzles2 := NewPuzzles(1)
	puzzle2 := puzzles2.Get(0)

	if puzzle1 == puzzle2 {
		t.Fatalf("Puzzles for different seeds should be different, but were both: %s", puzzle1)
	}
}

func TestPseudoRandomness_GetIndex0Twice_SamePuzzleBothTimes(t *testing.T) {
	puzzles := NewPuzzles(0)
	puzzle1 := puzzles.Get(0)
	puzzle2 := puzzles.Get(0)

	if puzzle1 != puzzle2 {
		t.Fatalf("Puzzle at an index should remain constant during a round, but they were'%s' and '%s'",
			puzzle1, puzzle2)
	}
}
