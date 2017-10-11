package slackmessaging

import "log"

type Handler struct {
	chConsumed <-chan ConsumedMessage
	decoders   map[string]Decoder
}

func (h *Handler) Start() {
	go h.handleMessages()
}

func (h *Handler) handleMessages() {
	for message := range h.chConsumed {
		decoder, ok := h.decoders[message.Topic]
		if !ok {
			log.Printf("No decoder for topic %s", message.Topic)
		}
		decoder.Decode(message.Value)
	}
}

type Decoder interface {
	Decode(value []byte)
}

func NewHandler(chConsumed <-chan ConsumedMessage, decoders map[string]Decoder) *Handler {
	return &Handler{chConsumed, decoders}
}
