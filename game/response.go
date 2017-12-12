package game

// Puzzle is a word with the letters shuffled.
type Puzzle string

//go:generate mockgen -destination=../mock/mock_response.go -package=mock github.com/dandeliondeathray/nona/game Response

// Response is the interface from Nona to the player.
type Response interface {
	OnPuzzleNotification(player Player, puzzle Puzzle)
	OnCorrectWord(player Player, word Word)
}
