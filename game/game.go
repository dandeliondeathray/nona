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
	g.persistence.StoreNewRound(seed)
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
	checkSolution := checkWord{
		player:      player,
		word:        word,
		persistence: g.persistence,
		solutions:   g.solutions,
		puzzleChain: g.puzzleChain,
		response:    g.response}
	g.persistence.ResolvePlayerState(player, &checkSolution)
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
	PuzzleIndex int
}

// NewPlayerState creates a state for a player that just started playing.
func NewPlayerState() PlayerState {
	return PlayerState{PuzzleIndex: 0}
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
	puzzle := Puzzle(p.puzzleChain.Get(playerState.PuzzleIndex))
	p.response.OnPuzzleNotification(p.player, puzzle)
}

// checkWord takes a resolved player state and checks the provided word against the current puzzle
// and the dicionary.
type checkWord struct {
	player      Player
	word        Word
	persistence Persistence
	solutions   *Solutions
	puzzleChain *chain.Puzzles
	response    Response
}

func (c *checkWord) PlayerStateResolved(playerState PlayerState) {
	puzzle := Puzzle(c.puzzleChain.Get(playerState.PuzzleIndex))
	if c.solutions.Check(c.word, puzzle) {
		c.persistence.PlayerSolvedPuzzle(c.player, playerState.PuzzleIndex+1)
		c.response.OnCorrectWord(c.player, c.word)
	} else {
		tooMany, tooFew := Diff(normalize(string(c.word)), string(puzzle))
		c.response.OnIncorrectWord(c.player, c.word, tooMany, tooFew)
	}
}
