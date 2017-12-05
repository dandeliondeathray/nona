package chain

import (
	"log"
	"reflect"

	"github.com/dandeliondeathray/nona/service/plumber"
)

// Service represents the PuzzleStore micro service.
type Service struct {
	chNewRounds chan *NewRound
	teams       *Teams
}

func (s *Service) handleNewRounds() {
	log.Println("Start handling requests.")
	for newRound := range s.chNewRounds {
		log.Printf("New round for team %s with seed %d", newRound.Team, newRound.Seed)
		// TODO: Implement me
	}
}

// Start the service by handling all input messages
func (s *Service) Start() {
	go s.handleNewRounds()
}

// Configuration of the service by specifying which topics to consume from and produce to.
func (s *Service) Configuration() plumber.Configuration {
	return plumber.Configuration{
		ConsumeTopics: []plumber.TopicConfiguration{
			plumber.TopicConfiguration{
				ChMessageType: reflect.TypeOf(NewRound{}),
				ChMessage:     s.chNewRounds,
				SchemaName:    "NewRound",
				Topic:         "nona_NewRound"}},

		ProduceTopics: []plumber.TopicConfiguration{}}
}

// NewService creates a new Chain service.
func NewService(teams *Teams) *Service {
	return &Service{chNewRounds: make(chan *NewRound, 100), teams: teams}
}

//
// Message types
//

// NewRound is sent when a team starts a new round.
type NewRound struct {
	Team string `avro:"team"`
	Seed int64  `avro:"seed"`
}
