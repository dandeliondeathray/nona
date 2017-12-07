package solution

import "testing"

var dictionary = []string{"PUSSGURKA", "DATORSPEL"}

func TestSolution_ProperWordAndPuzzleMatch_SolutionIsCorrect(t *testing.T) {
	solutions := NewSolutions(dictionary)

	properWord := Word("PUSSGURKA")
	puzzle := Puzzle("USSGURKAP")

	if !solutions.Check(properWord, puzzle) {
		t.Fatalf("Word %s should be a correct solution for puzzle %s", properWord, puzzle)
	}
}

func TestSolution_ImproperWordAndPuzzleMatch_SolutionIsIncorrect(t *testing.T) {
	solutions := NewSolutions(dictionary)

	improperWord := Word("NOTAWORD")
	puzzle := Puzzle("NOTAWORD")

	if solutions.Check(improperWord, puzzle) {
		t.Fatalf("Word %s is not in the dictionary, and should therefore NOT be a correct solution for puzzle %s", improperWord, puzzle)
	}
}

func TestSolution_ProperWordDoesNotMAtchPuzzle_SolutionIsIncorrect(t *testing.T) {
	solutions := NewSolutions(dictionary)

	properWord := Word("PUSSGURKA")
	puzzle := Puzzle("DATORSPEL")

	if solutions.Check(properWord, puzzle) {
		t.Fatalf("Word %s should NOT be a correct solution for puzzle %s", properWord, puzzle)
	}
}
