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
		"PuzzleNotification": puzzleNotificationDecoder}

	consumer := slackmessaging.NewConsumer()
	chConsumed := consumer.ConsumedMessages()

	handler := slackmessaging.NewHandler(chConsumed, decoders)
	handler.Start()

	consumer.Start()

	go slackmessaging.StartProbes(24689)

	chBlock := make(chan bool)
	<-chBlock
}
