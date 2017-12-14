package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
	p.states[player] = state
}

func main() {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Fatalf("SLACK_TOKEN must be set.")
	}

	dictionaryPath := os.Getenv("DICTIONARY")
	if dictionaryPath == "" {
		log.Fatalf("DICTIONARY must be set to the path of the dictionary file.")
	}

	chOutgoing := make(chan slack.OutgoingMessage)
	response := slack.SlackResponse{ChOutgoing: chOutgoing}

	dictionary, err := game.LoadDictionaryFromFile(dictionaryPath)
	if err != nil {
		log.Fatalf("Error when reading dictionary: %v", err)
	}
	persistence := inMemoryPersistence{make(map[game.Player]game.PlayerState)}
	nona := game.NewGame(&response, &persistence, dictionary)
	nona.NewRound(time.Now().Unix())

	slack.RunSlack(token, nona, chOutgoing)
}
