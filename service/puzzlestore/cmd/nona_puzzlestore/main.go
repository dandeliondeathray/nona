package main

import (
	"log"
	"os"
	"strings"

	"github.com/dandeliondeathray/nona/service/plumber"
	"github.com/dandeliondeathray/nona/service/puzzlestore"
)

func main() {
	schemasPath := os.Getenv("SCHEMA_PATH")
	brokerEnv := os.Getenv("KAFKA_BROKERS")
	if brokerEnv == "" {
		log.Fatalf("No KAFKA_BROKERS set!")
	}
	brokers := strings.Split(brokerEnv, ",")
	log.Println("Kafka brokers:", brokers)

	codecs, err := plumber.LoadCodecsFromPath(schemasPath)
	if err != nil {
		log.Fatalf("Could not load codecs from path %s", schemasPath)
	}

	service := puzzlestore.NewService()
	service.Start()

	plumber := plumber.NewPlumber(service, codecs)
	if err = plumber.Start(brokers); err != nil {
		log.Fatalf("Could not start plumber: %s", err)
	}

	go puzzlestore.StartProbes(24689)

	chBlock := make(chan bool)
	<-chBlock
}
