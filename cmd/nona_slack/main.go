package main

import (
	"log"
	"os"
	"strings"

	"github.com/dandeliondeathray/nona/score"

	"github.com/dandeliondeathray/nona/control"
	"github.com/dandeliondeathray/nona/game"
	"github.com/dandeliondeathray/nona/persistence"

	"github.com/dandeliondeathray/nona/slack"
)

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

	notificationChannel := os.Getenv("NOTIFICATION_CHANNEL")
	if notificationChannel == "" {
		log.Fatal("NOTIFICATION_CHANNEL must be set to a channel name")
	}

	dictionary, err := game.LoadDictionaryFromFile(dictionaryPath)
	if err != nil {
		log.Fatalf("Error when reading dictionary: %v", err)
	}

	//
	// Arrange game components.
	//
	chOutgoing := make(chan slack.OutgoingMessage)
	chNotifications := make(chan slack.NotificationMessage)
	response := slack.NewSlackResponse(chOutgoing, chNotifications)
	etcdPersistence := persistence.NewPersistence(team, persistenceEndpoints)
	scoring := score.NewScoring(response, etcdPersistence)
	nona := game.NewGame(response, etcdPersistence, dictionary, scoring)

	go control.StartControl(nona)

	//
	// Recover state from database.
	//
	recoveryDone := make(chan bool, 1)
	etcdPersistence.Recover(nona, recoveryDone)
	<-recoveryDone
	// TODO: Don't mark Nona as Kubernetes ready before this is completed.

	slack.RunSlack(token, nona, chOutgoing, chNotifications, notificationChannel)
}
