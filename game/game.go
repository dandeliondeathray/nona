package game

import "github.com/dandeliondeathray/nona/game/chain"

// Player uniquely identifies a player or user.
type Player string

// Word represents a possible solution to a puzzle.
type Word string

// Game represents the word puzzle game from the perspective of the user interface.
type Game struct {
	response    Response
	puzzleChain *chain.Puzzles
	dictionary  []string
	persistence Persistence
	solutions   *Solutions
}

// NewRound starts a new round.
func (g *Game) NewRound(seed int64) {
	g.puzzleChain = chain.NewPuzzles(g.dictionary, seed)
}

// GiveMe requests a puzzle notification for a player.
func (g *Game) GiveMe(player Player) {
	notifyPlayerOfPuzzle := puzzleNotification{
		player:      player,
		puzzleChain: g.puzzleChain,
		response:    g.response}

	g.persistence.ResolvePlayerState(player, &notifyPlayerOfPuzzle)
}

// TryWord checks if the supplied word is a correct solution for the current puzzle.
func (g *Game) TryWord(player Player, word Word) {
	// playerState, ok := g.playerState[player]
	// if !ok {
	// 	playerState = 0
	// }
	// puzzle := Puzzle(g.puzzleChain.Get(playerState))
	// if g.solutions.Check(word, puzzle) {
	// 	g.playerState[player] = playerState + 1
	// }
}

// NewGame creates a new game type, given a dictionary.
func NewGame(response Response, persistence Persistence, dictionary []string) *Game {
	return &Game{
		response:    response,
		puzzleChain: nil,
		dictionary:  dictionary,
		persistence: persistence,
		solutions:   NewSolutions(dictionary)}
}

// PlayerState has all state for a given player.
type PlayerState struct {
	NextPuzzle int
}
