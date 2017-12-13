package slack

import (
	"log"
	"os"

	"github.com/dandeliondeathray/nona/game"
	"github.com/nlopes/slack"
)

// RunSlack connects to slack and listens to messages.
func RunSlack(token string, nona *game.Game, chOutgoing <-chan OutgoingMessage) {
	api := slack.New(token)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()
	go writeOutgoingMessages(chOutgoing, rtm)

	for msg := range rtm.IncomingEvents {
		log.Println("Event Received")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			log.Println("Infos:", ev.Info)
			log.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			//rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#konsulatet"))

		case *slack.MessageEvent:
			log.Printf("Message: %v\n", ev)
			msgEvent := msg.Data.(*slack.MessageEvent)
			player := game.Player(msgEvent.User)
			if msgEvent.Text == "!gemig" {
				nona.GiveMe(player)
			} else {
				nona.TryWord(player, game.Word(msgEvent.Text))
			}

		case *slack.PresenceChangeEvent:
			log.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			log.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			log.Fatalf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Println("Invalid credentials")
			return

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func writeOutgoingMessages(chOutgoing <-chan OutgoingMessage, rtm *slack.RTM) {
	for m := range chOutgoing {
		rtm.SendMessage(rtm.NewOutgoingMessage(m.User, m.Text))
	}
}
