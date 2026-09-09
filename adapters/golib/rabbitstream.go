package golib

import (
	"github.com/faustbrian/go-cloudevents"
	cloudrabbitstream "github.com/faustbrian/go-cloudevents/adapters/rabbitstream"
	"github.com/faustbrian/go-rabbitmq-streams"
)

type RabbitStreamTransport = cloudrabbitstream.Transport
type RabbitStreamMessage = cloudrabbitstream.Message
type RabbitStreamState = cloudrabbitstream.State

func EncodeRabbitStream(event cloudevents.Event, transport RabbitStreamTransport) (rabbitstream.Message, error) {
	return cloudrabbitstream.Encode(event, transport)
}

func DecodeRabbitStream(message rabbitstream.Message, limits cloudevents.Limits) (RabbitStreamMessage, RabbitStreamState, error) {
	return cloudrabbitstream.Decode(message, limits)
}
