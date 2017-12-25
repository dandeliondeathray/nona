package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/dandeliondeathray/nona/game"
	"github.com/dandeliondeathray/nona/persistence"

	"github.com/dandeliondeathray/nona/slack"
)

func main() {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Fatalf("SLACK_TOKEN must be set.")
	}

	dictionaryPath := os.Getenv("DICTIONARY")
	if dictionaryPath == "" {
		log.Fatalf("DICTIONARY must be set to the path of the dictionary file.")
	}

	persistenceEndpointsString := os.Getenv("PERSISTENCE_ENDPOINTS")
	if persistenceEndpointsString == "" {
		log.Fatalf("PERSISTENCE_ENDPOINTS must be set to a comma separated list of etcd instances")
	}
	persistenceEndpoints := strings.Split(persistenceEndpointsString, ",")

	chOutgoing := make(chan slack.OutgoingMessage)
	response := slack.SlackResponse{ChOutgoing: chOutgoing}

	dictionary, err := game.LoadDictionaryFromFile(dictionaryPath)
	if err != nil {
		log.Fatalf("Error when reading dictionary: %v", err)
	}
	etcdPersistence := persistence.NewPersistence("konsulatet", persistenceEndpoints)
	nona := game.NewGame(&response, etcdPersistence, dictionary)
	nona.NewRound(time.Now().Unix())

	slack.RunSlack(token, nona, chOutgoing)
}
