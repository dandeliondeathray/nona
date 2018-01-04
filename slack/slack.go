package slack

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/dandeliondeathray/nona/game"
	"github.com/nlopes/slack"
)

type slackChannels struct {
	channels map[game.Player]string
	rtmInfo  *slack.Info
	mutex    sync.Mutex
}

func (s *slackChannels) setReplyChannel(player game.Player, channelID string) {
	defer s.mutex.Unlock()
	s.mutex.Lock()

	s.channels[player] = channelID
}

func (s *slackChannels) getReplyChannel(player game.Player) string {
	defer s.mutex.Unlock()
	s.mutex.Lock()

	channelID, ok := s.channels[player]
	if !ok {
		panic(fmt.Sprintf("No reply channel found for player %s", player))
	}
	return channelID
}

func (s *slackChannels) setRTMInfo(rtmInfo *slack.Info) {
	defer s.mutex.Unlock()
	s.mutex.Lock()

	s.rtmInfo = rtmInfo
}

func (s *slackChannels) getUserName(player game.Player) string {
	defer s.mutex.Unlock()
	s.mutex.Lock()

	name := s.rtmInfo.GetUserByID(string(player)).Name
	return name
}

var channels = slackChannels{channels: make(map[game.Player]string), mutex: sync.Mutex{}}

// RunSlack connects to slack and listens to messages.
func RunSlack(token string, nona *game.Game, chOutgoing <-chan OutgoingMessage, chNotifications chan NotificationMessage, notificationChannel string) {
	api := slack.New(token)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()
	go writeOutgoingMessages(chOutgoing, rtm)
	go writeNotificationMessages(chNotifications, rtm, notificationChannel)

	var handler *NonaSlackHandler
	var rtmInfo *slack.Info

	for msg := range rtm.IncomingEvents {
		log.Println("Event Received")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			rtmInfo = rtm.GetInfo()
			channels.setRTMInfo(rtmInfo)
			self := game.Player(rtmInfo.User.ID)
			handler = NewNonaSlackHandler(nona, self)

		case *slack.ConnectedEvent:
			log.Println("Infos:", ev.Info)
			log.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			//rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#konsulatet"))

		case *slack.MessageEvent:
			if handler == nil {
				log.Printf("Event %v was received before handler was initialized!", ev)
				continue
			}
			log.Printf("Message: %v\n", ev)
			msgEvent := msg.Data.(*slack.MessageEvent)
			player := game.Player(msgEvent.User)
			if strings.HasPrefix(msgEvent.Channel, "C") {
				log.Printf("Message '%s' was to a public channel. Ignoring.", msgEvent.Text)
				continue
			}
			channels.setReplyChannel(player, msgEvent.Channel)
			messageText := msgEvent.Text
			if msgEvent.SubType == "message_changed" && msgEvent.SubMessage != nil {
				messageText = msgEvent.SubMessage.Text
			}
			handler.OnMessage(player, messageText)

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
		replyChannel := channels.getReplyChannel(m.Player)
		rtm.SendMessage(rtm.NewOutgoingMessage(m.Text, replyChannel))
	}
}

func writeNotificationMessages(chNotifications <-chan NotificationMessage, rtm *slack.RTM, notificationChannel string) {
	for m := range chNotifications {
		rtm.SendMessage(rtm.NewOutgoingMessage(m.Text, notificationChannel))
	}
}
