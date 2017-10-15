package plumber

import "log"

// Plumber handles the interface between a micro service and the message queue infrastructure.
type Plumber struct {
	service Service
}

// NewPlumber creates a Plumber given the service it will work with.
func NewPlumber(service Service) *Plumber {
	return &Plumber{service}
}

// Start creates producers and consumers against the Kafka message queue, and sets up Avro
// encoders and decoders for the messages.
func (p *Plumber) Start() error {

	serviceConfig := p.service.Configuration()

	for i := range serviceConfig.ConsumeTopics {
		// TODO: Decoder of this message

		// TODO: Set up consumers
		topicConfig := serviceConfig.ConsumeTopics[i]
		log.Println("Consume from ", topicConfig.Topic)
	}

	for i := range serviceConfig.ProduceTopics {
		// TODO: Encoder of this message
		
		// TODO: Set up producers
		topicConfig := serviceConfig.ProduceTopics[i]
		log.Println("Produce to ", topicConfig.Topic)
	}

	return nil
}
