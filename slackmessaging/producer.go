package slackmessaging

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	codecs       *Codecs
}

func NewProducer(codecs *Codecs) (*Producer, error) {
	p, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		return nil, err
	}

	return &Producer{syncProducer: p, codecs: codecs}, nil
}

func (p *Producer) Start(chChatMessage chan ChatMessage) {
	go p.handleChatMessages(chChatMessage)
}

func (p *Producer) handleChatMessages(chChatMessage chan ChatMessage) {
	for chatMessage := range chChatMessage {
		log.Printf("Sending Chat message: %v", chatMessage)
		topic := fmt.Sprintf("nona_%s_Chat", chatMessage.Team)
		value, err := chatMessage.Encode(p.codecs)
		if err != nil {
			log.Printf("Could not encode ChatMessage: %v", err)
			continue
		}
		err = p.sendMessage(value, topic)
		if err != nil {
			log.Printf("Could not send message: %v", err)
		}
	}
}

func (p *Producer) sendMessage(value []byte, topic string) error {
	msg := &sarama.ProducerMessage{Topic: topic,
		Value: sarama.ByteEncoder(value)}
	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Sent message: Partition %d, Offset %d", partition, offset)
	return nil
}
