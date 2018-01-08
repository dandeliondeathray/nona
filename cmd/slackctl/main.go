package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage: slackctl <token file> <command>")
		os.Exit(1)
	}
	command := os.Args[2]

	tokenBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		println("Could not read token file:", err)
		os.Exit(1)
	}
	token := string(tokenBytes)

	api := slack.New(token)

	if command == "channels" {
		channels, err := api.GetChannels(true)
		if err != nil {
			println("Failed to list channels:", err)
			os.Exit(1)
		}

		println("Channels:")
		for _, c := range channels {
			fmt.Printf("|%20s|%10s|\n", c.Name, c.ID)
		}
	} else {
		println("Unknown command:", command)
		os.Exit(1)
	}
}
