package slackmessaging

import (
	"fmt"
	"log"
	"reflect"

	"github.com/dandeliondeathray/nona/service/plumber"
)

// Service represents this micro service.
type Service struct {
	chPuzzleNotification chan *PuzzleNotification
	chCorrectSolution    chan *CorrectSolution
	chIncorrectSolution  chan *IncorrectSolution
	chChatMessage        chan ChatMessage
	teams                []string
}

// NewService creates a SlackMessaging service.
func NewService(teams []string) Service {
	return Service{
		make(chan *PuzzleNotification, 100),
		make(chan *CorrectSolution, 100),
		make(chan *IncorrectSolution, 100),
		make(chan ChatMessage, 100), teams}
}

// Start starts the SlackMessaging service.
func (s *Service) Start() {
	go handlePuzzleNotification(s.chPuzzleNotification, s.chChatMessage)
	go handleCorrectSolution(s.chCorrectSolution, s.chChatMessage)
	go handleIncorrectSolution(s.chIncorrectSolution, s.chChatMessage)
}

// Configuration returns information on what topics are consumed and produced, and what types
// are expected from each topic.
func (s Service) Configuration() plumber.Configuration {
	consume := []plumber.TopicConfiguration{
		plumber.TopicConfiguration{
			ChMessageType: reflect.TypeOf(PuzzleNotification{}),
			ChMessage:     s.chPuzzleNotification,
			SchemaName:    "PuzzleNotification",
			Topic:         "nona_PuzzleNotification"},
		plumber.TopicConfiguration{
			ChMessageType: reflect.TypeOf(CorrectSolution{}),
			ChMessage:     s.chCorrectSolution,
			SchemaName:    "CorrectSolution",
			Topic:         "nona_CorrectSolution"},
		plumber.TopicConfiguration{
			ChMessageType: reflect.TypeOf(IncorrectSolution{}),
			ChMessage:     s.chIncorrectSolution,
			SchemaName:    "IncorrectSolution",
			Topic:         "nona_IncorrectSolution"}}

	produce := make([]plumber.TopicConfiguration, len(s.teams))
	for i := range s.teams {
		produce[i] = plumber.TopicConfiguration{
			ChMessageType: reflect.TypeOf(ChatMessage{}),
			ChMessage:     s.chChatMessage,
			SchemaName:    "Chat",
			Topic:         fmt.Sprintf("nona_%s_Chat", s.teams[i])}
	}

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

// IncorrectSolution notifies the user that the word is incorrect.
type IncorrectSolution struct {
	User string `avro:"user_id"`
	Team string `avro:"team"`
}

// CorrectSolution notifies the user that the word is correct.
type CorrectSolution struct {
	User string `avro:"user_id"`
	Team string `avro:"team"`
}

// HandlePuzzleNotification reponds to each puzzle notification with a chat message to the user.
func handlePuzzleNotification(chPuzzleNotification <-chan *PuzzleNotification, chChatMessage chan<- ChatMessage) {
	for p := range chPuzzleNotification {
		log.Printf("Puzzle notification received: %v", p)
		chatMessage := ChatMessage{p.User, p.Team, p.Puzzle}
		chChatMessage <- chatMessage
	}
}

func handleCorrectSolution(chCorrectSolution <-chan *CorrectSolution, chChatMessage chan<- ChatMessage) {
	for p := range chCorrectSolution {
		log.Printf("Correct solution notification received: %v", p)
		chatMessage := ChatMessage{p.User, p.Team, "Korrekt!"}
		chChatMessage <- chatMessage
	}
}

func handleIncorrectSolution(chIncorrectSolution <-chan *IncorrectSolution, chChatMessage chan<- ChatMessage) {
	for p := range chIncorrectSolution {
		log.Printf("Incorrect solution notification received: %v", p)
		chatMessage := ChatMessage{p.User, p.Team, "Inkorrekt."}
		chChatMessage <- chatMessage
	}
}
