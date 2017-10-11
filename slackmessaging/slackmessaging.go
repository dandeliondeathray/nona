package slackmessaging

import (
	"fmt"
	"log"
)

type UserID string
type Team string

type ChatMessage struct {
	User UserID
	Team Team
	Text string
}

func (c ChatMessage) Encode(codecs *Codecs) ([]byte, error) {
	codec, err := codecs.ByName("Chat")
	if err != nil {
		return nil, fmt.Errorf("Could not read schema Chat")
	}

	native := make(map[string]interface{})
	native["user_id"] = c.User
	native["team"] = c.Team
	native["text"] = c.Text

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

// A PuzzleNotification is s
type PuzzleNotification struct {
	User   UserID
	Team   Team
	Puzzle string
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
	codec, err := p.codecs.ByName("PuzzleNotification")
	if err != nil {
		log.Println("Couldn't get codec:", err)
		return
	}

	native, _, err := codec.NativeFromBinary(value)
	if err != nil {
		fmt.Println(err)
	}
	notification, ok := native.(map[string]interface{})
	if !ok {
		fmt.Println("Invalid PuzzleNotification message, after schema decode")
		return
	}

	userID, ok := notification["user_id"].(UserID)
	if !ok {
		fmt.Println("Could not read 'user_id' as string")
		return
	}

	team, ok := notification["team"].(Team)
	if !ok {
		fmt.Println("Could not read 'team' as string")
		return
	}

	puzzle, ok := notification["puzzle"].(string)
	if !ok {
		fmt.Println("Could not read 'puzzle' as string")
		return
	}

	event := PuzzleNotification{User: userID, Team: team, Puzzle: puzzle}
	p.sm.PuzzleNotificationEvent(event)
}
