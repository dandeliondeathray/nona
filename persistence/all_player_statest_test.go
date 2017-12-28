package persistence_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dandeliondeathray/nona/game"
	"github.com/dandeliondeathray/nona/persistence"
)

func TestAllPlayerStatePersistence_OnePlayerAtIndex17_OnlyThatPlayerReturned(t *testing.T) {
	if testing.Short() {
		return
	}

	resolution := newAsyncAllPlayerStatesResolution()
	player := game.Player("U1")
	puzzleIndex := 17

	p := persistence.NewPersistence("konsulatet1", testingEndpoints)
	p.PlayerSolvedPuzzle(player, puzzleIndex)
	time.Sleep(time.Duration(1) * time.Second)

	// Act
	p.ResolveAllPlayerStates(resolution)

	playerStates, err := resolution.AwaitPlayerStates()
	if err != nil {
		t.Fatalf("%v", err)
	}

	checkPlayerState(t, playerStates, player, puzzleIndex)
}

func TestAllPlayerStatePersistence_TwoPlayers_BothPlayersReturned(t *testing.T) {
	if testing.Short() {
		return
	}

	resolution := newAsyncAllPlayerStatesResolution()
	player1 := game.Player("U1")
	puzzleIndex1 := 17
	player2 := game.Player("U2")
	puzzleIndex2 := 42

	p := persistence.NewPersistence("konsulatet2", testingEndpoints)
	p.PlayerSolvedPuzzle(player1, puzzleIndex1)
	p.PlayerSolvedPuzzle(player2, puzzleIndex2)
	time.Sleep(time.Duration(1) * time.Second)

	// Act
	p.ResolveAllPlayerStates(resolution)

	playerStates, err := resolution.AwaitPlayerStates()
	if err != nil {
		t.Fatalf("%v", err)
	}

	checkPlayerState(t, playerStates, player1, puzzleIndex1)
	checkPlayerState(t, playerStates, player2, puzzleIndex2)
}

//
// Test helpers
//

func checkPlayerState(t *testing.T, playerStates map[game.Player]game.PlayerState, player game.Player, index int) {
	actualState, ok := playerStates[player]
	if !ok {
		t.Fatalf("Player %s was not present in the resolved player states", player)
	}

	if actualState.PuzzleIndex != index {
		t.Fatalf("Player %s: Expected puzzle index %d, but got %d", player, index, actualState.PuzzleIndex)
	}
}

type asyncAllPlayerStatesResolution struct {
	chResult chan map[game.Player]game.PlayerState
}

func (p *asyncAllPlayerStatesResolution) AllPlayerStatesResolved(states map[game.Player]game.PlayerState) {
	p.chResult <- states
}

func newAsyncAllPlayerStatesResolution() *asyncAllPlayerStatesResolution {
	return &asyncAllPlayerStatesResolution{make(chan map[game.Player]game.PlayerState, 10)}
}

func (p *asyncAllPlayerStatesResolution) AwaitPlayerStates() (map[game.Player]game.PlayerState, error) {
	select {
	case states := <-p.chResult:
		return states, nil
	case <-time.After(time.Duration(1) * time.Second):
		return nil, fmt.Errorf("No full player states received within 1 second.")
	}
}
