package game_test

import (
	"testing"

	"github.com/dandeliondeathray/nona/game"
	"golang.org/x/text/unicode/norm"
)

var dictionary = []string{"PUSSGURKA", "DATORSPEL", "ABCDEFÅÄÖ", "Å", "Ä", "Ö", "TETEATETE"}

func TestSolution_ProperWordAndPuzzleMatch_SolutionIsCorrect(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	properWord := game.Word("PUSSGURKA")
	puzzle := game.Puzzle("USSGURKAP")

	if !solutions.Check(properWord, puzzle) {
		t.Fatalf("Word %s should be a correct solution for puzzle %s", properWord, puzzle)
	}
}

func TestSolution_ImproperWordAndPuzzleMatch_SolutionIsIncorrect(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	improperWord := game.Word("NOTAWORD")
	puzzle := game.Puzzle("NOTAWORD")

	if solutions.Check(improperWord, puzzle) {
		t.Fatalf("Word %s is not in the dictionary, and should therefore NOT be a correct solution for puzzle %s", improperWord, puzzle)
	}
}

func TestSolution_ProperWordDoesNotMAtchPuzzle_SolutionIsIncorrect(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	properWord := game.Word("PUSSGURKA")
	puzzle := game.Puzzle("DATORSPEL")

	if solutions.Check(properWord, puzzle) {
		t.Fatalf("Word %s should NOT be a correct solution for puzzle %s", properWord, puzzle)
	}
}

func TestSolution_NonNormalWordInDictionary_SolutionIsCorrect(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	nonNormalWord := game.Word(norm.NFD.String("ABCDEFÅÄÖ"))
	puzzle := game.Puzzle(string(nonNormalWord))

	if !solutions.Check(nonNormalWord, puzzle) {
		t.Fatalf("Word %s should be a correct solution for puzzle %s, even though it's not in normal form", nonNormalWord, puzzle)
	}
}

func TestSolution_CorrectWordWithNonNormalDictionary_SolutionIsCorrect(t *testing.T) {
	nonNormalDictionary := make([]string, len(dictionary))
	for i, w := range dictionary {
		nonNormalDictionary[i] = norm.NFD.String(w)
	}
	solutions := game.NewSolutions(nonNormalDictionary)

	normalWord := game.Word(norm.NFKC.String("ABCDEFÅÄÖ"))
	puzzle := game.Puzzle(string(normalWord))

	if !solutions.Check(normalWord, puzzle) {
		t.Fatalf("Word %s should match non-normal entry in dictionary for puzzle %s", normalWord, puzzle)
	}
}

func TestSolution_NonNormalPuzzleWithNormalWord_SolutionIsCorrect(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	normalWord := game.Word(norm.NFKC.String("ABCDEFÅÄÖ"))
	puzzle := game.Puzzle(norm.NFD.String(string(normalWord)))

	if !solutions.Check(normalWord, puzzle) {
		t.Fatalf("Word %s should be a correct solution for puzzle %s, even though the puzzle is not in normal form", normalWord, puzzle)
	}
}

func TestSolution_LowercaseWord_SolutionIsCaseInsensitive(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	lowercaseWord := game.Word(norm.NFKC.String("pussgurka"))
	puzzle := game.Puzzle(norm.NFKC.String("PUSSGURKA"))

	if !solutions.Check(lowercaseWord, puzzle) {
		t.Fatalf("Word %s should be correct, even though it's lower case", lowercaseWord)
	}
}

func TestSolution_LowercasePuzzle_PuzzleIsCaseInsensitive(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	word := game.Word(norm.NFKC.String("PUSSGURKA"))
	lowercasePuzzle := game.Puzzle(norm.NFKC.String("pussgurka"))

	if !solutions.Check(word, lowercasePuzzle) {
		t.Fatalf("Word %s should be correct, even though it's lower case", word)
	}
}

func TestSpecialCharacters_SpaceInWord_SpaceIsRemovedWhenComparing(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	word := game.Word("PUSS GURKA")
	puzzle := game.Puzzle("PUSSGURKA")

	if !solutions.Check(word, puzzle) {
		t.Fatalf("Word %s should be correct, even though it has a space in it", word)
	}
}

func TestSpecialCharacters_Accents_AccentsAreRemoved(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	word := game.Word("TÊTEÀTÊTE")
	puzzle := game.Puzzle("TETEATETE")

	if !solutions.Check(word, puzzle) {
		t.Fatalf("Word %s should a correct solution for puzzle %s, even with the accents", word, puzzle)
	}
}

func TestSpecialCharacters_Dash_DashesAreRemoved(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	word := game.Word("DATOR-SPEL")
	puzzle := game.Puzzle("DATORSPEL")

	if !solutions.Check(word, puzzle) {
		t.Fatalf("Word %s should a correct solution for puzzle %s, even with the dash", word, puzzle)
	}
}

func TestSpecialCharacters_SwedishCharacterÅ_IsNotTheSameWithoutDiacritic(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	word := game.Word("Å")
	puzzle := game.Puzzle("A")

	if solutions.Check(word, puzzle) {
		t.Fatalf("Word %s should not be a correct solution for %s, since they're different in Swedish", word, puzzle)
	}
}

func TestSpecialCharacters_SwedishCharacterÄ_IsNotTheSameWithoutDiacritic(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	word := game.Word("Ä")
	puzzle := game.Puzzle("A")

	if solutions.Check(word, puzzle) {
		t.Fatalf("Word %s should not be a correct solution for %s, since they're different in Swedish", word, puzzle)
	}
}

func TestSpecialCharacters_SwedishCharacterÖ_IsNotTheSameWithoutDiacritic(t *testing.T) {
	solutions := game.NewSolutions(dictionary)

	word := game.Word("Ö")
	puzzle := game.Puzzle("O")

	if solutions.Check(word, puzzle) {
		t.Fatalf("Word %s should not be a correct solution for %s, since they're different in Swedish", word, puzzle)
	}
}
