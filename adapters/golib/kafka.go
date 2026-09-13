package golib

import (
	"github.com/faustbrian/go-cloudevents"
	cloudkafka "github.com/faustbrian/go-cloudevents/adapters/kafka"
	"github.com/faustbrian/go-kafka"
)

type KafkaTransport = cloudkafka.Transport
type KafkaState = cloudkafka.State

func EncodeKafka(event cloudevents.Event, mode cloudevents.ContentMode, transport KafkaTransport) (kafka.ProducerRecord, error) {
	return cloudkafka.Encode(event, mode, transport)
}

func DecodeKafka(record kafka.ConsumedRecord, limits cloudevents.Limits) (cloudevents.KafkaMessage, KafkaState, error) {
	return cloudkafka.Decode(record, limits)
}
