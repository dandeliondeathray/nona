package slack

import (
	"strings"

	"github.com/dandeliondeathray/nona/game"
)

//go:generate mockgen -destination=../mock/mock_handler.go -package=mock github.com/dandeliondeathray/nona/slack Game

// Game is the game interface from Slack towards game.Game.
type Game interface {
	GiveMe(player game.Player)
	TryWord(player game.Player, word game.Word)
	SkipPuzzle(player game.Player)
}

type NonaSlackHandler struct {
	game Game
	self game.Player
}

func NewNonaSlackHandler(game Game, self game.Player) *NonaSlackHandler {
	return &NonaSlackHandler{game: game, self: self}
}

func (c *NonaSlackHandler) OnMessage(player game.Player, text string) {
	if player == c.self {
		return
	}
	trimmedText := strings.TrimSpace(text)
	if trimmedText == "!gemig" || trimmedText == "!" || trimmedText == "! gemig" {
		c.game.GiveMe(player)
	} else if trimmedText == "!skippa" {
		c.game.SkipPuzzle(player)
	} else {
		c.game.TryWord(player, game.Word(trimmedText))
	}
}
