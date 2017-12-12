package game_test

import (
	"fmt"

	"github.com/dandeliondeathray/nona/game"
)

// Check that puzzles are different.
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

// Check that puzzles are the same.
type samePuzzlesMatcher struct {
	puzzle *game.Puzzle
}

func (m *samePuzzlesMatcher) Matches(x interface{}) bool {
	p, ok := x.(game.Puzzle)
	if !ok {
		return false
	}
	if m.puzzle == nil {
		m.puzzle = &p
		return true
	}

	return *m.puzzle == p
}

func (m *samePuzzlesMatcher) String() string {
	if m.puzzle == nil {
		return "<first puzzle not set>"
	}
	return fmt.Sprintf("puzzle %s", *m.puzzle)
}

// Save the puzzle sent from Nona.
type puzzleSaver struct {
	puzzle *game.Puzzle
}

func (p *puzzleSaver) Matches(x interface{}) bool {
	puzzle, ok := x.(game.Puzzle)
	if !ok {
		return false
	}
	p.puzzle = &puzzle
	return true
}

func (p *puzzleSaver) String() string {
	s := "<no puzzle set>"
	if p.puzzle != nil {
		s = string(*p.puzzle)
	}
	return fmt.Sprintf("puzzleSaver{%s}", s)
}

// fakePersistence is a stand-in for the persistence layer.
type fakePersistence struct {
	resolutions map[game.Player]game.PlayerStateResolution
	states      map[game.Player]game.PlayerState
}

func (p *fakePersistence) ResolvePlayerState(player game.Player, resolution game.PlayerStateResolution) {
	p.resolutions[player] = resolution
	_, playerStateExists := p.states[player]
	if !playerStateExists {
		p.states[player] = game.NewPlayerState()
	}
}

func (p *fakePersistence) PlayerSolvedPuzzle(player game.Player, newPuzzleIndex int) {
	state, ok := p.states[player]
	if !ok {
		panic(fmt.Sprintf("Player %s solved the puzzle, and the new state is %d, but no such state was found", player, newPuzzleIndex))
	}
	state.PuzzleIndex = newPuzzleIndex
	p.states[player] = state
}

func (p *fakePersistence) playerStateResolved(player game.Player) {
	resolution, ok := p.resolutions[player]
	if !ok {
		panic(fmt.Sprintf("playerStateResolved: No call to ResolvePlayerState for player %s was made", player))
	}
	state, ok := p.states[player]
	if !ok {
		state = game.NewPlayerState()
	}
	resolution.PlayerStateResolved(state)
}

func newFakePersistence() *fakePersistence {
	return &fakePersistence{
		make(map[game.Player]game.PlayerStateResolution),
		make(map[game.Player]game.PlayerState)}
}
