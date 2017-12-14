package slack

import (
	"fmt"

	"github.com/dandeliondeathray/nona/game"
)

type OutgoingMessage struct {
	Player game.Player
	Text   string
}

type SlackResponse struct {
	ChOutgoing chan OutgoingMessage
}

func (r *SlackResponse) OnPuzzleNotification(player game.Player, puzzle game.Puzzle) {
	r.ChOutgoing <- OutgoingMessage{player, string(puzzle)}
}

func (r *SlackResponse) OnCorrectWord(player game.Player, word game.Word) {
	m := fmt.Sprintf("Ordet %s är korrekt!", word)
	r.ChOutgoing <- OutgoingMessage{player, m}
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
	r.ChOutgoing <- OutgoingMessage{player, message}
}
