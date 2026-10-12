package golib

import (
	"github.com/faustbrian/go-cloudevents"
	cloudkafka "github.com/faustbrian/go-cloudevents/adapters/kafka"
	"github.com/faustbrian/go-kafka"
)

// KafkaTransport is the target Kafka routing-state alias; broker lifecycle stays external.
type KafkaTransport = cloudkafka.Transport

// KafkaState is the target Kafka decoded broker-state alias, not a settlement handle.
type KafkaState = cloudkafka.State

// EncodeKafka delegates target Kafka strict binding conversion with copied output; it performs no broker I/O.
func EncodeKafka(event cloudevents.Event, mode cloudevents.ContentMode, transport KafkaTransport) (kafka.ProducerRecord, error) {
	return cloudkafka.Encode(event, mode, transport)
}

// DecodeKafka delegates target Kafka bounded owned decoding/state copying without acknowledgement; failures return zero values.
func DecodeKafka(record kafka.ConsumedRecord, limits cloudevents.Limits) (cloudevents.KafkaMessage, KafkaState, error) {
	return cloudkafka.Decode(record, limits)
}
