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
	playerState map[Player]int
	solutions   *Solutions
}

// NewRound starts a new round.
func (g *Game) NewRound(seed int64) {
	g.puzzleChain = chain.NewPuzzles(g.dictionary, seed)
}

// GiveMe requests a puzzle notification for a player.
func (g *Game) GiveMe(player Player) {
	playerState, ok := g.playerState[player]
	if !ok {
		playerState = 0
	}
	puzzle := Puzzle(g.puzzleChain.Get(playerState))
	g.response.OnPuzzleNotification(player, puzzle)
}

// TryWord checks if the supplied word is a correct solution for the current puzzle.
func (g *Game) TryWord(player Player, word Word) {
	playerState, ok := g.playerState[player]
	if !ok {
		playerState = 0
	}
	puzzle := Puzzle(g.puzzleChain.Get(playerState))
	if g.solutions.Check(word, puzzle) {
		g.playerState[player] = playerState + 1
	}
}

// NewGame creates a new game type, given a dictionary.
func NewGame(response Response, dictionary []string) *Game {
	return &Game{
		response:    response,
		puzzleChain: nil,
		dictionary:  dictionary,
		playerState: make(map[Player]int),
		solutions:   NewSolutions(dictionary)}
}
