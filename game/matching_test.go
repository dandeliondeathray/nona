package game_test

import (
	"fmt"

	"github.com/dandeliondeathray/nona/game"
)

type differentPuzzlesMatcher struct {
	puzzle *game.Puzzle
}

func (m *differentPuzzlesMatcher) Matches(x interface{}) bool {
	p, ok := x.(game.Puzzle)
	if !ok {
		return false
	}
	if m.puzzle == nil {
		m.puzzle = &p
		return true
	}

	return *m.puzzle != p
}

func (m *differentPuzzlesMatcher) String() string {
	if m.puzzle == nil {
		return "<first puzzle not set>"
	}
	return fmt.Sprintf("puzzle different from %s", *m.puzzle)
}
