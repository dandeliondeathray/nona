package solution

import (
	"sort"

	"golang.org/x/text/unicode/norm"
)

// Word represents a potential solution to a puzzle.
type Word string

// Puzzle represents a puzzle, which is a word with the letters shuffled.
type Puzzle string

// Solutions is a struct that maintains all functionality for checking solutions to puzzles.
type Solutions struct {
	dictionary map[Word]bool
}

// Check if a word is a correct solution to a given puzzle.
func (s *Solutions) Check(word Word, puzzle Puzzle) bool {
	normalWord := normalizeWord(word)
	normalPuzzle := normalizePuzzle(puzzle)
	sortedWord := sortLetters(string(normalWord))
	sortedPuzzle := sortLetters(string(normalPuzzle))
	matchesPuzzle := sortedWord == sortedPuzzle
	_, isInDictionary := s.dictionary[normalWord]
	return isInDictionary && matchesPuzzle
}

// NewSolutions creates a new Solutions struct from a dictionary.
func NewSolutions(dictionary []string) *Solutions {
	dictionarySet := make(map[Word]bool)
	for _, w := range dictionary {
		dictionarySet[normalizeWord(Word(w))] = true
	}
	return &Solutions{dictionarySet}
}

// Helpers

func normalize(s string) string {
	return norm.NFKC.String(s)
}

func normalizeWord(w Word) Word {
	return Word(normalize(string(w)))
}

func normalizePuzzle(p Puzzle) Puzzle {
	return Puzzle(normalize(string(p)))
}

type runeSlice []rune

func (r runeSlice) Len() int {
	return len(r)
}

func (r runeSlice) Less(i, j int) bool {
	return r[i] < r[j]

}
func (r runeSlice) Swap(i, j int) {
	r[j], r[i] = r[i], r[j]
}

func sortLetters(word string) string {
	runes := runeSlice(word)
	sort.Sort(runes)
	return string(runes)
}
