package persistence_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dandeliondeathray/nona/game"

	"github.com/dandeliondeathray/nona/persistence"
)

func TestPersistPlayerState_NewPlayer_PlayerStartsAtIndexZero(t *testing.T) {
	if testing.Short() {
		return
	}

	resolution := newAsyncPlayerStateResolution()
	player := game.Player("UNEWPLAYER")

	p := persistence.NewPersistence("konsulatet")
	p.ResolvePlayerState(player, resolution)

	err := resolution.AwaitPuzzleIndex(0)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

// asyncPlayerStateResolution lets us test that a player state is resolved by Persistence.
type asyncPlayerStateResolution struct {
	chPuzzleIndex chan int
}

func (p *asyncPlayerStateResolution) PlayerStateResolved(playerState game.PlayerState) {
	p.chPuzzleIndex <- playerState.PuzzleIndex
}

func newAsyncPlayerStateResolution() *asyncPlayerStateResolution {
	return &asyncPlayerStateResolution{make(chan int, 10)}
}

func (p *asyncPlayerStateResolution) AwaitPuzzleIndex(expected int) error {
	select {
	case index := <-p.chPuzzleIndex:
		if index != expected {
			return fmt.Errorf("Got puzzle index %d, but expected %d", index, expected)
		}
	case <-time.After(time.Duration(1) * time.Second):
		return fmt.Errorf("No puzzle index received within 1 second.")
	}
	return nil
}
