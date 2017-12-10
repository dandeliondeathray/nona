package game

// Player uniquely identifies a player or user.
type Player string

// Game represents the word puzzle game from the perspective of the user interface.
type Game struct {
	response Response
}

// NewRound starts a new round.
func (g *Game) NewRound(seed int64) {

}

// GiveMe requests a puzzle notification for a player.
func (g *Game) GiveMe(player Player) {
	g.response.OnPuzzleNotification(player, Puzzle("PUSSGURKA"))
}

// NewGame creates a new game type, given a dictionary.
func NewGame(response Response, dictionary []string) *Game {
	return &Game{response: response}
}
