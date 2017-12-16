package game

import (
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/unicode/norm"
)

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
	decomposed := norm.NFKD.String(s)
	åäöSet := runes.Predicate(isSwedishÅÄÖ)
	diacritics := runes.In(unicode.Diacritic)
	hyphenSet := runes.In(unicode.Hyphen)
	spaceSet := runes.In(unicode.White_Space)
	removeRuneSet := inAnySet{[]runes.Set{diacritics, hyphenSet, spaceSet}}
	t := runes.If(åäöSet, nil, runes.Remove(removeRuneSet))
	decomposedWithoutDiacritics := t.String(decomposed)
	return strings.ToUpper(norm.NFKC.String(decomposedWithoutDiacritics))
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

//
// Transforming characters
//

func isSwedishÅÄÖ(r rune) bool {
	åäö := norm.NFKD.String("ÅÄÖ")
	return strings.ContainsRune(åäö, r)
}

func isSeparator(r rune) bool {
	return strings.ContainsRune(" -", r)
}

type inAnySet struct {
	sets []runes.Set
}

func (s inAnySet) Contains(r rune) bool {
	for _, runeSet := range s.sets {
		if runeSet.Contains(r) {
			return true
		}
	}
	return false
}
