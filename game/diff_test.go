package game_test

import (
	"testing"

	"github.com/dandeliondeathray/nona/game"
)

func TestDiff_WordAndPuzzleAreTheSame_NoDiff(t *testing.T) {
	tooMany, tooFew := game.Diff("ABCDEF", "ABCDEF")

	if tooMany != "" {
		t.Errorf("Word and puzzle are identical, but Diff reports too many %s", tooMany)
	}
	if tooFew != "" {
		t.Errorf("Word and puzzle are identical, but Diff reports too few %s", tooFew)
	}
}

func TestDiff_WordIsEmpty_TooFewSameAsPuzzle(t *testing.T) {
	puzzle := "ABCDEF"
	tooMany, tooFew := game.Diff("", puzzle)
	if tooMany != "" {
		t.Errorf("Diff reported too many characters %s, but it should be empty.", tooMany)
	}
	if tooFew != puzzle {
		t.Errorf("Diff reported tooFew '%s' but puzzle was %s and they should be the same", tooFew, puzzle)
	}
}

func TestDiff_PuzzleIsEmpty_TooManySameAsWord(t *testing.T) {
	word := "ABCDEF"
	tooMany, tooFew := game.Diff(word, "")

	if tooFew != "" {
		t.Errorf("Diff reported too few %s, but it should be empty", tooFew)
	}

	if tooMany != word {
		t.Errorf("Diff reported too many '%s', but it should be equal to the word %s", tooMany, word)
	}
}

func TestDiff_WordHasOneTooManyLetters_TooManyIsThatLetter(t *testing.T) {
	tooMany, tooFew := game.Diff("ABCDEFG", "ABCDEF")

	if tooMany != "G" {
		t.Errorf("Diff reported too many '%s', but should be G", tooMany)
	}
	if tooFew != "" {
		t.Errorf("Diff reported too few %s, but it should be empty", tooFew)
	}
}

func TestDiff_WordHasOneLetterRepeated_TooManyIsThatLetterRepeated(t *testing.T) {
	tooMany, tooFew := game.Diff("ABCDEFGGG", "ABCDEF")

	if tooMany != "GGG" {
		t.Errorf("Diff reported too many '%s', but it should be GGG", tooMany)
	}
	if tooFew != "" {
		t.Errorf("Diff reported too few %s, but it should be empty", tooFew)
	}
}

func TestDiff_WordHasDifferentLettersNotInOrder_ResultIsSorted(t *testing.T) {
	tooMany, _ := game.Diff("AZBCDKEFG", "ABCDEF")

	if tooMany != "GKZ" {
		t.Fatalf("Diff reported too many '%s', but result should be sorted like GKZ", tooMany)
	}
}
