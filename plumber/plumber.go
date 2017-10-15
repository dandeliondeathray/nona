package plumber

import (
	"fmt"
	"log"
	"reflect"

	"github.com/Shopify/sarama"
)

// Plumber handles the interface between a micro service and the message queue infrastructure.
type Plumber struct {
	service            Service
	codecs             *Codecs
	partitionConsumers []sarama.PartitionConsumer
}

// NewPlumber creates a Plumber given the service it will work with.
func NewPlumber(service Service, codecs *Codecs) *Plumber {
	return &Plumber{service, codecs, make([]sarama.PartitionConsumer, 100)}
}

// Start creates producers and consumers against the Kafka message queue, and sets up Avro
// encoders and decoders for the messages.
func (p *Plumber) Start(brokers []string) error {

	kafkaConfig := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, kafkaConfig)
	if err != nil {
		return err
	}

	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return fmt.Errorf("Could not create producer: %v", err)
	}
	chToProducer := make(chan *sarama.ProducerMessage, 100)
	go produceMessages(chToProducer, producer)

	serviceConfig := p.service.Configuration()

	for i := range serviceConfig.ConsumeTopics {
		topicConfig := serviceConfig.ConsumeTopics[i]
		chConsumed := make(chan []byte, 100)

		// Set up decoding of messages from this topic.
		codec, err := p.codecs.ByName(topicConfig.SchemaName)
		if err != nil {
			return fmt.Errorf("Could not find schema name %s", topicConfig.SchemaName)
		}
		recordEncoding := NewRecordEncoding(codec)

		go decodeKafkaMessagesFromTopic(chConsumed, recordEncoding, topicConfig)

		// Set up consumers for this topic, and send all messages via a channel to the decoder.
		log.Println("Consume from", topicConfig.Topic)
		p.consumeTopic(consumer, topicConfig.Topic, chConsumed)
	}

	for i := range serviceConfig.ProduceTopics {
		topicConfig := serviceConfig.ProduceTopics[i]

		// Set up encoding of message to this topic
		codec, err := p.codecs.ByName(topicConfig.SchemaName)
		if err != nil {
			return fmt.Errorf("Could not find schema name %s", topicConfig.SchemaName)
		}
		recordEncoding := NewRecordEncoding(codec)
		log.Println("RecordEncoding", recordEncoding)

		go encodeMessagesToTopic(chToProducer, recordEncoding, topicConfig)
	}

	return nil
}

func decodeKafkaMessagesFromTopic(chConsumed <-chan []byte, recordEncoding *RecordEncoding, topicConfig TopicConfiguration) {
	chToService := reflect.ValueOf(topicConfig.ChMessage)
	for consumed := range chConsumed {
		decoded, err := recordEncoding.DecodeFromType(consumed, topicConfig.ChMessageType)
		if err != nil {
			log.Printf("Could not decode message as %v: %v", topicConfig.ChMessageType.Name(), consumed)
			continue
		}
		log.Printf("Decoded message: %v", decoded)

		chToService.Send(reflect.ValueOf(decoded))
	}
}

func encodeMessagesToTopic(chToProducer chan<- *sarama.ProducerMessage, recordEncoding *RecordEncoding, topicConfig TopicConfiguration) {
	chMessages := reflect.ValueOf(topicConfig.ChMessage)
	for {
		message, ok := chMessages.Recv()
		if !ok {
			break
		}
		bytes, err := recordEncoding.Encode(message.Interface())
		if err != nil {
			log.Printf("Could not encode message %v", message.Interface())
			continue
		}
		chToProducer <- &sarama.ProducerMessage{Topic: topicConfig.Topic, Value: sarama.ByteEncoder(bytes)}
	}
}

func (p *Plumber) consumeTopic(consumer sarama.Consumer, topic string, chToDecoder chan<- []byte) error {
	log.Println("Creating consumer for topic", topic)
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for i := range partitions {
		partition := partitions[i]
		log.Printf("Creating partition consumer for topic %s and partition %d", topic, partition)
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			log.Printf("ERROR: Couldn't create partition consumer %d", partition)
			return err
		}
		go forwardMessagesToDecoder(partitionConsumer.Messages(), chToDecoder)
		p.partitionConsumers = append(p.partitionConsumers, partitionConsumer)
	}

	return nil

}

func forwardMessagesToDecoder(chMessages <-chan *sarama.ConsumerMessage, chToDecoder chan<- []byte) {
	for m := range chMessages {
		log.Println("Received message on partition", m.Partition)
		chToDecoder <- m.Value
	}
}

func produceMessages(chProducerMessages <-chan *sarama.ProducerMessage, producer sarama.SyncProducer) {
	for m := range chProducerMessages {
		partition, offset, err := producer.SendMessage(m)
		if err != nil {
			log.Printf("Could not produce message %v: %v", m, err)
			continue
		}

		log.Printf("Produced message %v to partition %d and offset %d", m, partition, offset)
	}
}
