package mock

import (
	"fmt"

	"github.com/dandeliondeathray/nona/game"
)

// FakePersistence is a stand-in for the persistence layer.
type FakePersistence struct {
	resolutions map[game.Player]game.PlayerStateResolution
	states      map[game.Player]game.PlayerState
	// This fake persistence only supports a single AllPlayerStatesResolution
	// at any given point. This does not necessarily reflect how the actual
	// persistence works though.
	allStatesResolution game.AllPlayerStatesResolution
	seed                int64
}

func (p *FakePersistence) ResolvePlayerState(player game.Player, resolution game.PlayerStateResolution) {
	p.resolutions[player] = resolution
	_, playerStateExists := p.states[player]
	if !playerStateExists {
		p.states[player] = game.NewPlayerState()
	}
}

func (p *FakePersistence) ResolveAllPlayerStates(resolution game.AllPlayerStatesResolution) {
	p.allStatesResolution = resolution
}

func (p *FakePersistence) PlayerSolvedPuzzle(player game.Player, newPuzzleIndex int) {
	state, ok := p.states[player]
	if !ok {
		panic(fmt.Sprintf("Player %s solved the puzzle, and the new state is %d, but no such state was found", player, newPuzzleIndex))
	}
	state.PuzzleIndex = newPuzzleIndex
	p.states[player] = state
}

func (p *FakePersistence) PlayerSkippedPuzzle(player game.Player, newPuzzleIndex int) {
	state, ok := p.states[player]
	if !ok {
		panic(fmt.Sprintf("Player %s solved the puzzle, and the new state is %d, but no such state was found", player, newPuzzleIndex))
	}
	state.PuzzleIndex = newPuzzleIndex
	state.Skipped++
	p.states[player] = state
}

func (p *FakePersistence) StoreNewRound(seed int64) {
	p.seed = seed
}

func (p *FakePersistence) FakePlayerStateResolved(player game.Player) {
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

func (p *FakePersistence) ResolveAllPlayers(states map[game.Player]game.PlayerState) {
	if p.allStatesResolution == nil {
		panic("There is no current request for all player states.")
	}
	p.allStatesResolution.AllPlayerStatesResolved(states)
}

func NewFakePersistence() *FakePersistence {
	return &FakePersistence{
		make(map[game.Player]game.PlayerStateResolution),
		make(map[game.Player]game.PlayerState),
		nil,
		0}
}
