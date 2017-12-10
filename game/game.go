package game

// Player uniquely identifies a player or user.
type Player string

// Word represents a possible solution to a puzzle.
type Word string

// Game represents the word puzzle game from the perspective of the user interface.
type Game struct {
	response Response
	index    int
}

// NewRound starts a new round.
func (g *Game) NewRound(seed int64) {

}

// GiveMe requests a puzzle notification for a player.
func (g *Game) GiveMe(player Player) {
	if g.index == 0 {
		g.response.OnPuzzleNotification(player, Puzzle("PUSSGURKA"))
		g.index = 1
	} else {
		g.response.OnPuzzleNotification(player, Puzzle("OTHERWORD"))
	}
}

// TryWord checks if the supplied word is a correct solution for the current puzzle.
func (g *Game) TryWord(player Player, word Word) {

}

// NewGame creates a new game type, given a dictionary.
func NewGame(response Response, dictionary []string) *Game {
	return &Game{response: response, index: 0}
}
