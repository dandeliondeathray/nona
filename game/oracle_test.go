package game_test

import (
	"fmt"
	"sort"

	"github.com/dandeliondeathray/nona/game"
	"golang.org/x/text/unicode/norm"
)

type gameOracle struct {
	solutions map[string][]string
}

func newGameOracle(dictionary []string) *gameOracle {
	solutions := make(map[string][]string)
	for _, nonNormalWord := range dictionary {
		word := normalize(nonNormalWord)
		sortedLetters := sortLetters(word)
		solutions[sortedLetters] = append(solutions[sortedLetters], word)
	}
	return &gameOracle{solutions: solutions}
}

func (g gameOracle) FindAllSolutionsFor(puzzle game.Puzzle) []string {
	sortedPuzzle := sortLetters(string(puzzle))
	solutions, ok := g.solutions[sortedPuzzle]
	if !ok {
		return []string{}
	}
	return solutions
}

func (g gameOracle) FindASolutionFor(puzzle game.Puzzle) string {
	solutions := g.FindAllSolutionsFor(puzzle)
	if len(solutions) == 0 {
		panic(fmt.Sprintf("Found no solutions for puzzle %s", puzzle))
	}
	return solutions[0]
}

// Letter transformations
func normalize(s string) string {
	return norm.NFKC.String(s)
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
