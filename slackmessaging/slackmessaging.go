package slackmessaging

import (
	"log"
	"reflect"

	"github.com/dandeliondeathray/nona/plumber"
)

// Service represents this micro service.
type Service struct {
	chPuzzleNotification chan *PuzzleNotification
	chChatMessage        chan ChatMessage
}

// NewService creates a SlackMessaging service.
func NewService() Service {
	return Service{make(chan *PuzzleNotification, 100), make(chan ChatMessage, 100)}
}

// Start starts the SlackMessaging service.
func (s *Service) Start() {
	go handlePuzzleNotification(s.chPuzzleNotification, s.chChatMessage)
}

// Configuration returns information on what topics are consumed and produced, and what types
// are expected from each topic.
func (s Service) Configuration() plumber.Configuration {
	consume := []plumber.TopicConfiguration{
		plumber.TopicConfiguration{
			ChMessageType: reflect.TypeOf(PuzzleNotification{}),
			ChMessage:     s.chPuzzleNotification,
			SchemaName:    "PuzzleNotification",
			Topic:         "nona_PuzzleNotification"}}

	produce := []plumber.TopicConfiguration{
		plumber.TopicConfiguration{
			ChMessageType: reflect.TypeOf(ChatMessage{}),
			ChMessage:     s.chChatMessage,
			SchemaName:    "Chat",
			Topic:         "nona_konsulatet_Chat"}}

	config := plumber.Configuration{ConsumeTopics: consume, ProduceTopics: produce}
	return config
}

// ChatMessage contains a text message sent to a specific user in a team.
type ChatMessage struct {
	User string `avro:"user_id"`
	Team string `avro:"team"`
	Text string `avro:"text"`
}

// A PuzzleNotification is sent in response to a user requesting the current/next puzzle.
type PuzzleNotification struct {
	User   string `avro:"user_id"`
	Team   string `avro:"team"`
	Puzzle string `avro:"puzzle"`
}

// HandlePuzzleNotification reponds to each puzzle notification with a chat message to the user.
func handlePuzzleNotification(chPuzzleNotification <-chan *PuzzleNotification, chChatMessage chan<- ChatMessage) {
	for p := range chPuzzleNotification {
		log.Printf("Puzzle notification received: %v", p)
		chatMessage := ChatMessage{p.User, p.Team, p.Puzzle}
		chChatMessage <- chatMessage
	}
}
