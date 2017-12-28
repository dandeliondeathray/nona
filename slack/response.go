package slack

import (
	"fmt"
	"strings"

	"github.com/dandeliondeathray/nona/game"
)

type OutgoingMessage struct {
	Player game.Player
	Text   string
}

type NotificationMessage struct {
	Text string
}

type SlackResponse struct {
	chOutgoing      chan OutgoingMessage
	chNotifications chan NotificationMessage
}

func NewSlackResponse(chOutgoing chan OutgoingMessage, chNotifications chan NotificationMessage) *SlackResponse {
	return &SlackResponse{chOutgoing, chNotifications}
}

func (r *SlackResponse) OnPuzzleNotification(player game.Player, puzzle game.Puzzle, index int) {
	r.chOutgoing <- OutgoingMessage{player, fmt.Sprintf("Pussel %d: %s", index, puzzle)}
}

func (r *SlackResponse) OnCorrectWord(player game.Player, word game.Word) {
	m := fmt.Sprintf("Ordet %s är korrekt!", word)
	r.chOutgoing <- OutgoingMessage{player, m}
}

func (r *SlackResponse) OnIncorrectWord(player game.Player, word game.Word, tooMany, tooFew string) {
	var message string
	if tooMany == "" && tooFew == "" {
		message = fmt.Sprintf("%s finns inte i ordlistan.", word)
	} else {
		message = fmt.Sprintf("%s matchar inte pusslet.", word)
		if tooMany != "" {
			message += fmt.Sprintf(" För många %s.", tooMany)
		}
		if tooFew != "" {
			message += fmt.Sprintf(" För få %s.", tooFew)
		}
	}
	r.chOutgoing <- OutgoingMessage{player, message}
}

func (r *SlackResponse) OnNoRound(player game.Player) {
	r.chOutgoing <- OutgoingMessage{player, "Ingen runda har startats än."}
}

func (r *SlackResponse) OnPerPlayerScores(scoringName string, scores []game.PerPlayerScore) {
	message := []string{fmt.Sprintf("*%s*", scoringName)}

	for _, score := range scores {
		name := channels.getUserName(score.Player)
		scoreText := fmt.Sprintf("%s: %.1f", name, score.Score)
		message = append(message, scoreText)
	}
	if len(scores) == 0 {
		message = append(message, "Inga spelare löste några pussel.")
	}

	r.chNotifications <- NotificationMessage{strings.Join(message, "\n")}
}
