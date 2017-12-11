package game

import "github.com/dandeliondeathray/nona/game/chain"

//go:generate mockgen -destination=../mock/mock_persistence.go -package=mock github.com/dandeliondeathray/nona/game Persistence

// PlayerStateResolution is called when a players state has been fetched from the persistence
// layer.
type PlayerStateResolution interface {
	PlayerStateResolved(playerState PlayerState)
}

// Persistence is the interface towards the database.
type Persistence interface {
	ResolvePlayerState(player Player, resolution PlayerStateResolution)
}

// puzzleNotification tells the player what the current puzzle is, when resolved by the database.
type puzzleNotification struct {
	player      Player
	puzzleChain *chain.Puzzles
	response    Response
}

// PlayerStateResolved for puzzleNotification gets the current puzzle from the chain and notifies
// the player.
func (p *puzzleNotification) PlayerStateResolved(playerState PlayerState) {
	puzzle := Puzzle(p.puzzleChain.Get(playerState.NextPuzzle))
	p.response.OnPuzzleNotification(p.player, puzzle)
}
