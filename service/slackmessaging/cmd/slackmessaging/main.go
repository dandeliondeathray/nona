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
	if brokerEnv == "" {
		log.Fatalf("KAFKA_BROKERS not set")
	}
	brokers := strings.Split(brokerEnv, ",")

	codecs, err := plumber.LoadCodecsFromPath(schemasPath)
	if err != nil {
		log.Fatalf("Could not load codecs from path %s", schemasPath)
	}

	service := slackmessaging.NewService()
	service.Start()

	plumber := plumber.NewPlumber(&service, codecs)
	plumber.Start(brokers)

	go slackmessaging.StartProbes(24689)

	chBlock := make(chan bool)
	<-chBlock
}
