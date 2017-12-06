package main

import (
	"log"
	"os"
	"strings"

	"github.com/dandeliondeathray/nona/service/chain"
	"github.com/dandeliondeathray/nona/service/plumber"
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

	teams := chain.NewTeams([]string{"BEBIS"})

	service := chain.NewService(teams)
	service.Start()

	plumber := plumber.NewPlumber(service, codecs)
	if err = plumber.Start(brokers); err != nil {
		log.Fatalf("Could not start plumber: %s", err)
	}

	service.ListenAndServe(8080)

	chBlock := make(chan bool)
	<-chBlock
}
