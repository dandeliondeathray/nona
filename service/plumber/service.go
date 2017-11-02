package plumber

import (
	"reflect"
)

// TopicConfiguration contains type information for each topic the service will produce or consume.
type TopicConfiguration struct {
	ChMessageType reflect.Type
	ChMessage     interface{}
	SchemaName    string
	Topic         string
}

// Configuration contains type information for each topic to consume from and produce to.
type Configuration struct {
	ConsumeTopics []TopicConfiguration
	ProduceTopics []TopicConfiguration
}

// Service is an interface for micro services.
type Service interface {
	Configuration() Configuration
	Start()
}
