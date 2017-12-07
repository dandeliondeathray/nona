package solution

import "sort"

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
	sortedWord := sortLetters(string(word))
	sortedPuzzle := sortLetters(string(puzzle))
	matchesPuzzle := sortedWord == sortedPuzzle
	_, isInDictionary := s.dictionary[word]
	return isInDictionary && matchesPuzzle
}

// NewSolutions creates a new Solutions struct from a dictionary.
func NewSolutions(dictionary []string) *Solutions {
	dictionarySet := make(map[Word]bool)
	for _, w := range dictionary {
		dictionarySet[Word(w)] = true
	}
	return &Solutions{dictionarySet}
}

// Helpers

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
