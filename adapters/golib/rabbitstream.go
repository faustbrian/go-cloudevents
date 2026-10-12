package golib

import (
	"github.com/faustbrian/go-cloudevents"
	cloudrabbitstream "github.com/faustbrian/go-cloudevents/adapters/rabbitstream"
	"github.com/faustbrian/go-rabbitmq-streams"
)

// RabbitStreamTransport is the structured target adapter transport alias; it defines no binary AMQP mapping.
type RabbitStreamTransport = cloudrabbitstream.Transport

// RabbitStreamMessage is the target decoded event/properties alias with owned returned bytes.
type RabbitStreamMessage = cloudrabbitstream.Message

// RabbitStreamState is the target retained delivery-state alias, not an acknowledgement handle.
type RabbitStreamState = cloudrabbitstream.State

// EncodeRabbitStream delegates target structured JSON mapping/transport validation without broker I/O.
func EncodeRabbitStream(event cloudevents.Event, transport RabbitStreamTransport) (rabbitstream.Message, error) {
	return cloudrabbitstream.Encode(event, transport)
}

// DecodeRabbitStream delegates target default transport limits and supplied CloudEvent limits; it does not acknowledge delivery.
func DecodeRabbitStream(message rabbitstream.Message, limits cloudevents.Limits) (RabbitStreamMessage, RabbitStreamState, error) {
	return cloudrabbitstream.Decode(message, limits)
}
