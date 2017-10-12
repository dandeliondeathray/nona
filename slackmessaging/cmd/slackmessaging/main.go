package main

import (
	"log"
	"os"

	"github.com/dandeliondeathray/nona/slackmessaging"
)

func main() {
	schemasPath := os.Getenv("SCHEMA_PATH")

	codecs, err := slackmessaging.LoadCodecsFromPath(schemasPath)
	if err != nil {
		log.Fatalf("Could not load codecs from path %s: %v", schemasPath, err)
	}
	chChatMessage := make(chan slackmessaging.ChatMessage, 100)
	sm := slackmessaging.NewSlackMessaging(chChatMessage)

	puzzleNotificationDecoder := slackmessaging.NewPuzzleNotificationDecoder(sm, codecs)

	decoders := map[string]slackmessaging.Decoder{
		"nona_PuzzleNotification": puzzleNotificationDecoder}

	consumer := slackmessaging.NewConsumer()
	chConsumed := consumer.ConsumedMessages()

	producer, err := slackmessaging.NewProducer(codecs)
	if err != nil {
		log.Fatalf("Could not create producer: %v", err)
	}
	producer.Start(chChatMessage)

	handler := slackmessaging.NewHandler(chConsumed, decoders)
	handler.Start()

	consumer.Start()

	go slackmessaging.StartProbes(24689)

	chBlock := make(chan bool)
	<-chBlock
}
