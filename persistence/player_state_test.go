package persistence_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dandeliondeathray/nona/game"
	"github.com/dandeliondeathray/nona/mock"
	"github.com/golang/mock/gomock"

	"github.com/dandeliondeathray/nona/persistence"
)

var testingEndpoints = []string{"localhost:2379", "localhost:22379", "localhost:32379"}

func TestPersistPlayerState_NewPlayer_PlayerStartsAtIndexZero(t *testing.T) {
	if testing.Short() {
		return
	}

	resolution := newAsyncPlayerStateResolution()
	player := game.Player("UNEWPLAYER")

	p := persistence.NewPersistence("konsulatet", testingEndpoints)
	p.ResolvePlayerState(player, resolution)

	err := resolution.AwaitPuzzleIndex(0)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPersistPlayerState_PlayerSolvedPuzzleNewIndex42_ResolvingPlayerStateTo42(t *testing.T) {
	if testing.Short() {
		return
	}

	resolution := newAsyncPlayerStateResolution()
	player := game.Player("U1")

	p := persistence.NewPersistence("konsulatet", testingEndpoints)
	p.PlayerSolvedPuzzle(player, 42)
	time.Sleep(time.Duration(1) * time.Second)

	// Act
	p.ResolvePlayerState(player, resolution)

	err := resolution.AwaitPuzzleIndex(42)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPersistPlayerState_PersistStateWithOneInstance_ResolveTheSameStateWithAnother(t *testing.T) {
	if testing.Short() {
		return
	}

	player := game.Player("U2")

	p1 := persistence.NewPersistence("konsulatet", testingEndpoints)
	p1.PlayerSolvedPuzzle(player, 43)

	time.Sleep(time.Duration(1) * time.Second)

	// Act
	resolution := newAsyncPlayerStateResolution()
	p2 := persistence.NewPersistence("konsulatet", testingEndpoints)
	p2.ResolvePlayerState(player, resolution)

	err := resolution.AwaitPuzzleIndex(43)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPersistPlayerState_NewRound_PlayerStateIsReset(t *testing.T) {
	if testing.Short() {
		return
	}

	player := game.Player("U3")

	p := persistence.NewPersistence("konsulatet", testingEndpoints)
	p.StoreNewRound(0)
	p.PlayerSolvedPuzzle(player, 43)

	time.Sleep(time.Duration(1) * time.Second)

	// Act
	resolution := newAsyncPlayerStateResolution()
	p.StoreNewRound(1)
	time.Sleep(time.Duration(1) * time.Second)

	p.ResolvePlayerState(player, resolution)

	err := resolution.AwaitPuzzleIndex(0)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestRecoverRound_PlayerWasAtIndex17_PlayerIsRecoveredAtIndex17(t *testing.T) {
	if testing.Short() {
		return
	}

	seed := int64(42)

	mockCtrl := gomock.NewController(t)
	recoveryHandler := mock.NewMockRecoveryHandler(mockCtrl)
	recoveryHandler.EXPECT().OnRoundRecovered(seed)

	done := make(chan bool, 1)

	// Arrange the persistence to have a current round with seed 42.
	p := persistence.NewPersistence("test_recover_player_state", testingEndpoints)
	p.StoreNewRound(42)
	time.Sleep(time.Duration(200) * time.Millisecond)

	// Store index 17 for player
	player := game.Player("U1")
	p.PlayerSolvedPuzzle(player, 17)
	time.Sleep(time.Duration(200) * time.Millisecond)

	p2 := persistence.NewPersistence("test_recover_player_state", testingEndpoints)
	p2.Recover(recoveryHandler, done)
	success := <-done

	if !success {
		t.Fatalf("Recovery failed")
	}

	// The player state should be resolved as index 17.
	resolution := newAsyncPlayerStateResolution()
	p2.ResolvePlayerState(player, resolution)

	err := resolution.AwaitPuzzleIndex(17)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

//
// Skip puzzles
//

func TestPersistPlayerState_NewPlayer_ZeroPuzzlesSkipped(t *testing.T) {
	if testing.Short() {
		return
	}

	resolution := newAsyncPlayerStateResolution()
	player := game.Player("UNEWPLAYER")

	p := persistence.NewPersistence("konsulatet", testingEndpoints)
	p.ResolvePlayerState(player, resolution)

	err := resolution.AwaitPuzzlesSkipped(0)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPersistPlayerState_PlayerSkipsOnePuzzle_ResolvingPlayerStateSkippedTo1(t *testing.T) {
	if testing.Short() {
		return
	}

	resolution := newAsyncPlayerStateResolution()
	player := game.Player("U1")

	p := persistence.NewPersistence("konsulatet", testingEndpoints)
	p.PlayerSkippedPuzzle(player, 42, 1)
	time.Sleep(time.Duration(1) * time.Second)

	// Act
	p.ResolvePlayerState(player, resolution)

	err := resolution.AwaitPuzzlesSkipped(1)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPersistPlayerState_PlayerSkipsThreePuzzles_ResolvingPlayerStateSkippedTo3(t *testing.T) {
	if testing.Short() {
		return
	}

	resolution := newAsyncPlayerStateResolution()
	player := game.Player("U1")

	p := persistence.NewPersistence("konsulatet", testingEndpoints)
	p.PlayerSkippedPuzzle(player, 42, 1)
	time.Sleep(time.Duration(1) * time.Second)
	p.PlayerSkippedPuzzle(player, 43, 2)
	time.Sleep(time.Duration(1) * time.Second)
	p.PlayerSkippedPuzzle(player, 44, 3)
	time.Sleep(time.Duration(1) * time.Second)

	// Act
	p.ResolvePlayerState(player, resolution)

	err := resolution.AwaitPuzzlesSkipped(3)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

//
// asyncPlayerStateResolution lets us test that a player state is resolved by Persistence.
//

type asyncPlayerStateResolution struct {
	chPuzzleIndex chan game.PlayerState
}

func (p *asyncPlayerStateResolution) PlayerStateResolved(playerState game.PlayerState) {
	p.chPuzzleIndex <- playerState
}

func newAsyncPlayerStateResolution() *asyncPlayerStateResolution {
	return &asyncPlayerStateResolution{make(chan game.PlayerState, 10)}
}

func (p *asyncPlayerStateResolution) AwaitPuzzleIndex(expected int) error {
	select {
	case playerState := <-p.chPuzzleIndex:
		if playerState.PuzzleIndex != expected {
			return fmt.Errorf("Got puzzle index %d, but expected %d", playerState.PuzzleIndex, expected)
		}
	case <-time.After(time.Duration(1) * time.Second):
		return fmt.Errorf("No puzzle index received within 1 second.")
	}
	return nil
}

func (p *asyncPlayerStateResolution) AwaitPuzzlesSkipped(expected int) error {
	select {
	case playerState := <-p.chPuzzleIndex:
		if playerState.Skipped != expected {
			return fmt.Errorf("Got skipped %d, but expected %d", playerState.Skipped, expected)
		}
	case <-time.After(time.Duration(1) * time.Second):
		return fmt.Errorf("No player state received within 1 second.")
	}
	return nil
}
