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

func (r *SlackResponse) OnPuzzleNotification(player game.Player, puzzle game.Puzzle) {
	r.chOutgoing <- OutgoingMessage{player, string(puzzle)}
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
	r.chOutgoing <- OutgoingMessage{player, "No round has been started."}
}

func (r *SlackResponse) OnPerPlayerScores(scoringName string, scores []game.PerPlayerScore) {
	message := []string{fmt.Sprintf("**%s**", scoringName)}

	for i, score := range scores {
		scoreText := fmt.Sprintf("%d: @%s: %f", i+1, score.Player, score.Score)
		message = append(message, scoreText)
	}

	r.chNotifications <- NotificationMessage{strings.Join(message, "\n")}
}
