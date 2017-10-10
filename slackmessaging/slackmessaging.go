package slackmessaging

type UserID string
type Team string

type ChatMessage struct {
	User UserID
	Team Team
	Text string
}

// A PuzzleNotification is s
type PuzzleNotification struct {
	User   UserID
	Team   Team
	Puzzle string
}

// PuzzleNotificationEvent responds with the puzzle itself.
func PuzzleNotificationEvent(p *PuzzleNotification) string {
	return p.Puzzle
}
