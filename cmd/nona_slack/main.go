package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dandeliondeathray/nona/game"

	"github.com/dandeliondeathray/nona/slack"
)

type inMemoryPersistence struct {
	states map[game.Player]game.PlayerState
}

func (p *inMemoryPersistence) ResolvePlayerState(player game.Player, resolution game.PlayerStateResolution) {
	state, ok := p.states[player]
	if !ok {
		state = game.NewPlayerState()
		p.states[player] = state
	}

	go resolution.PlayerStateResolved(state)
}

func (p *inMemoryPersistence) PlayerSolvedPuzzle(player game.Player, newPuzzleIndex int) {
	state, ok := p.states[player]
	if !ok {
		panic(fmt.Sprintf("Player %s solved the puzzle, new index is %d, but no state was found", player, newPuzzleIndex))
	}
	state.PuzzleIndex = newPuzzleIndex
}

func main() {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Fatalf("SLACK_TOKEN must be set.")
	}

	chOutgoing := make(chan slack.OutgoingMessage)
	response := slack.SlackResponse{ChOutgoing: chOutgoing}

	dictionary := []string{"PUSSGURKA"}
	persistence := inMemoryPersistence{}
	nona := game.NewGame(&response, &persistence, dictionary)

	slack.RunSlack(token, nona, chOutgoing)
}
