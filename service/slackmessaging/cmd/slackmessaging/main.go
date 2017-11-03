package main

import (
	"log"
	"os"
	"strings"

	"github.com/dandeliondeathray/nona/service/plumber"
	"github.com/dandeliondeathray/nona/service/slackmessaging"
)

func main() {
	schemasPath := os.Getenv("SCHEMA_PATH")
	brokerEnv := os.Getenv("KAFKA_BROKERS")
	teamsEnv := os.Getenv("TEAMS")
	if brokerEnv == "" {
		log.Fatalf("KAFKA_BROKERS not set")
	}

	if teamsEnv == "" {
		log.Fatalf("TEAMS must be a comma separated list of team names.")
	}

	brokers := strings.Split(brokerEnv, ",")
	teams := strings.Split(teamsEnv, ",")

	codecs, err := plumber.LoadCodecsFromPath(schemasPath)
	if err != nil {
		log.Fatalf("Could not load codecs from path %s", schemasPath)
	}

	service := slackmessaging.NewService(teams)
	service.Start()

	plumber := plumber.NewPlumber(&service, codecs)
	plumber.Start(brokers)

	go slackmessaging.StartProbes(24689)

	chBlock := make(chan bool)
	<-chBlock
}
