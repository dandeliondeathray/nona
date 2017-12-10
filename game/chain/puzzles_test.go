package chain

import (
	"sort"
	"testing"

	"golang.org/x/text/unicode/norm"
)

var dictionary = []string{
	"ACCEPTANS",
	"BUKMUSKEL",
	"CHARKDISK",
	"DATORSPEL",
	"EFTERSMAK",
	"FAKTARUTA",
	"GARANTERA",
	"HYPERTEXT",
	"IMPERATIV",
	"JETSETARE",
	"KANELBRUN",
	"LEGOKLOSS",
	"MIKROCHIP",
	"NAVELLUDD",
	"OISOLERAD",
	"PARAMETER",
	"QATARISKA",
	"ROSRABATT",
	"SKILJEDOM"}

type RuneSlice []rune

func (r RuneSlice) Len() int {
	return len(r)
}

func (r RuneSlice) Less(i, j int) bool {
	return r[i] < r[j]

}
func (r RuneSlice) Swap(i, j int) {
	r[j], r[i] = r[i], r[j]
}

func sortLetters(word string) string {
	runes := RuneSlice(word)
	sort.Sort(runes)
	return string(runes)
}

func isInDictionary(word string) bool {
	for _, a := range dictionary {
		if a == word {
			return true
		}
	}
	return false
}

func findSolutionToPuzzle(puzzle string) (string, bool) {
	sortedPuzzle := sortLetters(puzzle)
	for _, a := range dictionary {
		if sortLetters(a) == sortedPuzzle {
			return a, true
		}
	}
	return "", false
}

//
// Getting the puzzle
//
func TestGetPuzzle_ForAGivenUserAndTeam_PuzzleIsNotEmpty(t *testing.T) {
	puzzles := NewPuzzles(dictionary, 0)

	myPuzzle := puzzles.Get(0)

	if myPuzzle == "" {
		t.Fatalf("The returned puzzle was empty.")
	}
}

func TestGetPuzzle_GetTwoPuzzles_NextPuzzleIsDifferent(t *testing.T) {
	puzzles := NewPuzzles(dictionary, 0)

	puzzle0 := puzzles.Get(0)
	puzzle1 := puzzles.Get(1)

	if puzzle0 == puzzle1 {
		t.Fatalf("Puzzles at index 0 and 1 are the same: %s", puzzle0)
	}
}

func TestGetPuzzle_GetIndex10Directly_NoNeedToGetPrecedingIndicesFirst(t *testing.T) {
	puzzles := NewPuzzles(dictionary, 0)

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
	puzzles1 := NewPuzzles(dictionary, 0)
	puzzle1 := puzzles1.Get(0)

	puzzles2 := NewPuzzles(dictionary, 1)
	puzzle2 := puzzles2.Get(0)

	if puzzle1 == puzzle2 {
		t.Fatalf("Puzzles for different seeds should be different, but were both: %s", puzzle1)
	}
}

func TestPseudoRandomness_GetIndex0Twice_SamePuzzleBothTimes(t *testing.T) {
	puzzles := NewPuzzles(dictionary, 0)
	puzzle1 := puzzles.Get(0)
	puzzle2 := puzzles.Get(0)

	if puzzle1 != puzzle2 {
		t.Fatalf("Puzzle at an index should remain constant during a round, but they were'%s' and '%s'",
			puzzle1, puzzle2)
	}
}

//
// Dictionary
//

func TestDictionary_GetPuzzle_PuzzleIsNotAWord(t *testing.T) {
	puzzles := NewPuzzles(dictionary, 0)

	puzzle := puzzles.Get(0)

	if isInDictionary(puzzle) {
		t.Fatalf("The puzzle %s was found in the dictionary. It should not have been.", puzzle)
	}
}

func TestDictionary_GetPuzzle_PuzzleIsAShuffledWord(t *testing.T) {
	puzzles := NewPuzzles(dictionary, 0)

	puzzle := puzzles.Get(0)

	if _, ok := findSolutionToPuzzle(puzzle); !ok {
		t.Fatalf("The puzzle %s should be a shuffled word in the dictionary.", puzzle)
	}
}

func TestDictionary_GetTwoPuzzles_DifferentSolutions(t *testing.T) {
	puzzles := NewPuzzles(dictionary, 0)

	puzzle0 := puzzles.Get(0)
	puzzle1 := puzzles.Get(1)

	solution0, _ := findSolutionToPuzzle(puzzle0)
	solution1, _ := findSolutionToPuzzle(puzzle1)

	if solution0 == solution1 {
		t.Fatalf("the puzzles at index 0 and 1 should (with high probability) have different "+
			"solutions, but were the same: %s", solution0)
	}
}

//
// Unicode normalization.
//

func TestUnicodeNormalization_DictionaryWordIsDecomposed_PuzzleIsNormalized(t *testing.T) {
	word := norm.NFKC.String("ÅÄÖ")
	thisDictionary := []string{norm.NFD.String(word)}
	puzzles := NewPuzzles(thisDictionary, 0)

	puzzle := puzzles.Get(0)

	sortedWord := []rune(sortLetters(word))
	sortedPuzzle := []rune(sortLetters(puzzle))
	if len(sortedWord) != len(sortedPuzzle) {
		t.Fatalf("Word and puzzle have different number of runes: word '%v', puzzle '%v'", sortedWord, sortedPuzzle)
	}
	for i := range sortedWord {
		if sortedWord[i] != sortedPuzzle[i] {
			t.Fatalf("Word and puzzle mismatch at index %d: Word %v, puzzle %v", i, sortedWord, sortedPuzzle)
		}
	}
}
