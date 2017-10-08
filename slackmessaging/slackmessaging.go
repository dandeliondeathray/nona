package slackmessaging

type UserID string
type Team string

type ChatMessage struct {
	User UserID
	Team Team
	Text string
}

type PuzzleNotification struct {
	User   UserID
	Team   Team
	Puzzle string
}
