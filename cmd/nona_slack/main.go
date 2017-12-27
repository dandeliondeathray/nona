package main

import (
	"log"
	"os"
	"strings"

	"github.com/dandeliondeathray/nona/control"
	"github.com/dandeliondeathray/nona/game"
	"github.com/dandeliondeathray/nona/persistence"

	"github.com/dandeliondeathray/nona/slack"
)

type fakeScoring struct{}

func (f *fakeScoring) ProduceScores(seed int64) {

}

func main() {
	//
	// Read configuration from environment.
	//
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

	team := os.Getenv("TEAM")
	if team == "" {
		log.Fatal("TEAM must be set to the team name.")
	}

	dictionary, err := game.LoadDictionaryFromFile(dictionaryPath)
	if err != nil {
		log.Fatalf("Error when reading dictionary: %v", err)
	}

	//
	// Arrange game components.
	//
	chOutgoing := make(chan slack.OutgoingMessage)
	response := slack.SlackResponse{ChOutgoing: chOutgoing}
	etcdPersistence := persistence.NewPersistence(team, persistenceEndpoints)
	scoring := fakeScoring{}
	nona := game.NewGame(&response, etcdPersistence, dictionary, &scoring)

	go control.StartControl(nona)

	//
	// Recover state from database.
	//
	recoveryDone := make(chan bool, 1)
	etcdPersistence.Recover(nona, recoveryDone)
	<-recoveryDone
	// TODO: Don't mark Nona as Kubernetes ready before this is completed.

	slack.RunSlack(token, nona, chOutgoing)
}
