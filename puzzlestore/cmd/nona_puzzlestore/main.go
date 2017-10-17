package main

import (
	"log"
	"os"

	"github.com/dandeliondeathray/nona/plumber"
	"github.com/dandeliondeathray/nona/puzzlestore"
)

func main() {
	schemasPath := os.Getenv("SCHEMA_PATH")

	codecs, err := plumber.LoadCodecsFromPath(schemasPath)
	if err != nil {
		log.Fatalf("Could not load codecs from path %s", schemasPath)
	}

	service := puzzlestore.NewService()
	service.Start()

	plumber := plumber.NewPlumber(service, codecs)
	plumber.Start([]string{"localhost:9092"})

	go puzzlestore.StartProbes(24689)

	chBlock := make(chan bool)
	<-chBlock
}
