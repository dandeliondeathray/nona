package slackmessaging

import (
	"fmt"
	"log"
)

// TopicConfiguration contains type information for each topic
type TopicConfiguration struct {
	ChMessageType Type
	ChMessage     interface{}
	SchemaName    string
	Topic         string
}

// Configuration contains type information for each topic to consume from and produce to.
type Configuration struct {
	ConsumeTopics []TopicConfiguration
	ProduceTopics []TopicConfiguration
}

// Service represents this micro service.
type Service struct {
	chPuzzleNotification chan PuzzleNotification
	chChatMessage        chan ChatMessage
}

// NewService creates a SlackMessaging service.
func NewService() Service {
	return Service{make(chan PuzzleNotification, 100), make(chan ChatMessage, 100)}
}

// Start starts the SlackMessaging service.
func (s *Service) Start() {
	go handlePuzzleNotification(s.chPuzzleNotification, s.chChatMessage)
}

// Configuration returns information on what topics are consumed and produced, and what types
// are expected from each topic.
func (s Service) Configuration() Configuration {
	consume := []TopicConfiguration{
		TopicConfiguration{PuzzleNotification, s.chPuzzleNotification, "PuzzleNotification", "nona_PuzzleNotification"}}

	produce := []TopicConfiguration{
		TopicConfiguration{ChatMessage, s.chChatMessage, "Chat", "nona_konsulatet_Chat"}}

	config := Configuration{ConsumeTopics: consume, ProduceTopics: produce}
	return config
}

// UserID is a unique and stable identifier for a user.
type UserID string

// Team is an identifier for a Slack team.
type Team string

// ChatMessage contains a text message sent to a specific user in a team.
type ChatMessage struct {
	User UserID
	Team Team
	Text string
}

// A PuzzleNotification is sent in response to a user requesting the current/next puzzle.
type PuzzleNotification struct {
	User   UserID
	Team   Team
	Puzzle string
}

// HandlePuzzleNotification reponds to each puzzle notification with a chat message to the user.
func handlePuzzleNotification(chPuzzleNotification <-chan PuzzleNotification, chChatMessage chan<- ChatMessage) {
	for p := range chPuzzleNotification {
		log.Printf("Puzzle notification received: %v", p)
		chatMessage := ChatMessage{p.User, p.Team, p.Puzzle}
		chChatMessage <- chatMessage
	}
}

//
//
// Old code that I'm trying to replace
//
//

func (c ChatMessage) Encode(codecs *Codecs) ([]byte, error) {
	codec, err := codecs.ByName("Chat")
	if err != nil {
		return nil, fmt.Errorf("Could not read schema Chat")
	}

	native := make(map[string]interface{})
	native["user_id"] = string(c.User)
	native["team"] = string(c.Team)
	native["text"] = string(c.Text)

	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		return nil, fmt.Errorf("Could encode Chat message %v", native)
	}

	return binary, nil
}

type SlackMessaging struct {
	chChatMessage chan ChatMessage
}

func NewSlackMessaging(chChatMessage chan ChatMessage) *SlackMessaging {
	return &SlackMessaging{chChatMessage}
}

// PuzzleNotificationEvent responds with the puzzle itself.
func (s *SlackMessaging) PuzzleNotificationEvent(p PuzzleNotification) {
	chatMessage := ChatMessage{p.User, p.Team, p.Puzzle}
	s.chChatMessage <- chatMessage
}

type PuzzleNotificationDecoder struct {
	sm     *SlackMessaging
	codecs *Codecs
}

func NewPuzzleNotificationDecoder(sm *SlackMessaging, codecs *Codecs) *PuzzleNotificationDecoder {
	return &PuzzleNotificationDecoder{sm, codecs}
}

func (p *PuzzleNotificationDecoder) Decode(value []byte) {
	codec, err := p.codecs.ByName("nona_PuzzleNotification")
	if err != nil {
		log.Println("Couldn't get codec:", err)
		return
	}

	log.Println("Trying to decode binary value")
	native, _, err := codec.NativeFromBinary(value)
	if err != nil {
		fmt.Printf("Could not decode Avro message as PuzzleNotification: %s", err)
		return
	}
	log.Printf("Native: %v", native)
	notification, ok := native.(map[string]interface{})
	if !ok {
		fmt.Println("Invalid PuzzleNotification message, after schema decode")
		return
	}

	userID, ok := notification["user_id"].(string)
	if !ok {
		fmt.Println("Could not read 'user_id' as string")
		return
	}

	team, ok := notification["team"].(string)
	if !ok {
		fmt.Println("Could not read 'team' as string")
		return
	}

	puzzle, ok := notification["puzzle"].(string)
	if !ok {
		fmt.Println("Could not read 'puzzle' as string")
		return
	}

	event := PuzzleNotification{User: UserID(userID), Team: Team(team), Puzzle: puzzle}
	p.sm.PuzzleNotificationEvent(event)
}
