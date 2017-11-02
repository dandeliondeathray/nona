package puzzlestore

import (
	"log"
	"reflect"

	"github.com/dandeliondeathray/nona/service/plumber"
)

// Service represents the PuzzleStore micro service.
type Service struct {
	chRequests      chan *PuzzleRequest
	chNotifications chan PuzzleNotification
}

func (s *Service) handleRequests() {
	log.Println("Start handling requests.")
	for req := range s.chRequests {
		log.Println("User requested puzzle:", req)
		s.chNotifications <- PuzzleNotification{User: req.User, Team: req.Team, Puzzle: "PUSSGURKA"}
	}
}

// Start the service by handling all input messages
func (s *Service) Start() {
	go s.handleRequests()
}

// Configuration of the service by specifying which topics to consume from and produce to.
func (s *Service) Configuration() plumber.Configuration {
	return plumber.Configuration{
		ConsumeTopics: []plumber.TopicConfiguration{
			plumber.TopicConfiguration{
				ChMessageType: reflect.TypeOf(PuzzleRequest{}),
				ChMessage:     s.chRequests,
				SchemaName:    "UserRequestsPuzzle",
				Topic:         "nona_UserRequestsPuzzle"}},

		ProduceTopics: []plumber.TopicConfiguration{
			plumber.TopicConfiguration{
				ChMessageType: reflect.TypeOf(PuzzleNotification{}),
				ChMessage:     s.chNotifications,
				SchemaName:    "PuzzleNotification",
				Topic:         "nona_PuzzleNotification"}}}
}

// NewService creates a new PuzzleStore service.
func NewService() *Service {
	return &Service{
		chRequests:      make(chan *PuzzleRequest, 100),
		chNotifications: make(chan PuzzleNotification, 100)}
}

//
// Message types
//

// PuzzleRequest represents the UserRequestsPuzzle message.
type PuzzleRequest struct {
	User      string `avro:"user_id"`
	Team      string `avro:"team"`
	Timestamp int64  `avro:"timestamp"`
}

// PuzzleNotification is the notification that a new puzzle is available
type PuzzleNotification struct {
	User   string `avro:"user_id"`
	Team   string `avro:"team"`
	Puzzle string `avro:"puzzle"`
}
