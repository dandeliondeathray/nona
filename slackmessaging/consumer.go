package slackmessaging

import (
	"log"

	"github.com/Shopify/sarama"
)

// ConsumedMessage contains a message and the topic it was consumed from.
type ConsumedMessage struct {
	Topic string
	Value []byte
}

// Consumer decodes Avro message from Kafka.
type Consumer struct {
	partitionConsumers []sarama.PartitionConsumer
	chanConsumed       chan ConsumedMessage
}

// ConsumedMessages is a channel where all consumed messages are sent.
func (c *Consumer) ConsumedMessages() <-chan ConsumedMessage {
	return c.chanConsumed
}

// Start begins to consume the topics.
func (c *Consumer) Start() error {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		return err
	}

	err = c.consumeTopic(consumer, "nona_PuzzleNotification")
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) consumeTopic(consumer sarama.Consumer, topic string) error {
	log.Println("Creating consumer for topic", topic)
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for i := range partitions {
		p := partitions[i]
		log.Printf("Creating partition consumer for topic %s and partition %d", topic, p)
		partitionConsumer, err := consumer.ConsumePartition(topic, p, sarama.OffsetOldest)
		if err != nil {
			log.Printf("ERROR: Couldn't create partition consumer %d", p)
			return err
		}
		go c.forwardMessage(topic, partitionConsumer.Messages())
		c.partitionConsumers = append(c.partitionConsumers, partitionConsumer)
	}

	return nil

}

func (c *Consumer) forwardMessage(topic string, chMessages <-chan *sarama.ConsumerMessage) {
	for m := range chMessages {
		log.Println("Received message:", m)
		c.chanConsumed <- ConsumedMessage{Topic: m.Topic, Value: m.Value}
	}
}

// NewConsumer creates a Consumer object.
func NewConsumer() *Consumer {
	return &Consumer{make([]sarama.PartitionConsumer, 10), make(chan ConsumedMessage, 100)}
}
