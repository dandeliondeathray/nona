package game

//go:generate mockgen -destination=../mock/mock_game_persistence.go -package=mock github.com/dandeliondeathray/nona/game Persistence

// PlayerStateResolution is called when a players state has been fetched from the persistence
// layer.
type PlayerStateResolution interface {
	PlayerStateResolved(playerState PlayerState)
}

// AllPlayerStatesResolution is called when a full account of all player states should be resolved.
type AllPlayerStatesResolution interface {
	AllPlayerStatesResolved(playerStates map[Player]PlayerState)
}

// Persistence is the interface towards the database.
type Persistence interface {
	ResolvePlayerState(player Player, resolution PlayerStateResolution)
	ResolveAllPlayerStates(resolution AllPlayerStatesResolution)
	PlayerSolvedPuzzle(player Player, newPuzzleIndex int)
	StoreNewRound(seed int64)
}
