package solution

import "testing"
import "golang.org/x/text/unicode/norm"

var dictionary = []string{"PUSSGURKA", "DATORSPEL", "ABCDEFÅÄÖ"}

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

func TestSolution_NonNormalWordInDictionary_SolutionIsCorrect(t *testing.T) {
	solutions := NewSolutions(dictionary)

	nonNormalWord := Word(norm.NFD.String("ABCDEFÅÄÖ"))
	puzzle := Puzzle(string(nonNormalWord))

	if !solutions.Check(nonNormalWord, puzzle) {
		t.Fatalf("Word %s should be a correct solution for puzzle %s, even though it's not in normal form", nonNormalWord, puzzle)
	}
}

func TestSolution_CorrectWordWithNonNormalDictionary_SolutionIsCorrect(t *testing.T) {
	nonNormalDictionary := make([]string, len(dictionary))
	for i, w := range dictionary {
		nonNormalDictionary[i] = norm.NFD.String(w)
	}
	solutions := NewSolutions(nonNormalDictionary)

	normalWord := Word(norm.NFKC.String("ABCDEFÅÄÖ"))
	puzzle := Puzzle(string(normalWord))

	if !solutions.Check(normalWord, puzzle) {
		t.Fatalf("Word %s should match non-normal entry in dictionary for puzzle %s", normalWord, puzzle)
	}
}

func TestSolution_NonNormalPuzzleWithNormalWord_SolutionIsCorrect(t *testing.T) {
	solutions := NewSolutions(dictionary)

	normalWord := Word(norm.NFKC.String("ABCDEFÅÄÖ"))
	puzzle := Puzzle(norm.NFD.String(string(normalWord)))

	if !solutions.Check(normalWord, puzzle) {
		t.Fatalf("Word %s should be a correct solution for puzzle %s, even though the puzzle is not in normal form", normalWord, puzzle)
	}
}
