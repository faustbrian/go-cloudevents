package kafka_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/faustbrian/go-cloudevents"
	cloudkafka "github.com/faustbrian/go-cloudevents/adapters/kafka"
	"github.com/faustbrian/go-kafka"
)

func baseEvent(t *testing.T) cloudevents.Event {
	t.Helper()
	data, err := cloudevents.NewJSONData([]byte(`{"ok":true}`))
	if err != nil {
		t.Fatal(err)
	}
	event, err := cloudevents.NewEvent(cloudevents.Attributes{ID: "event-1", Source: "/source", Type: "example.created"}, data)
	if err != nil {
		t.Fatal(err)
	}
	return event
}

func TestKafkaAdapterPreservesTransportOwnershipAndBindingRoundTrip(t *testing.T) {
	t.Parallel()

	partitionKey, err := cloudevents.NewPartitionKeyAttribute("tenant-a")
	if err != nil {
		t.Fatal(err)
	}
	data, err := cloudevents.NewJSONData([]byte(`{"ok":true}`))
	if err != nil {
		t.Fatal(err)
	}
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: "event-1", Source: "/source", Type: "example.created", DataContentType: "application/json",
		Extensions: map[string]cloudevents.Attribute{"partitionkey": partitionKey},
	}, data)
	if err != nil {
		t.Fatal(err)
	}
	timestamp := time.Date(2026, 8, 9, 4, 5, 6, 0, time.UTC)
	transportKey := []byte("tenant-a")
	transportHeader := []byte("1")
	producer, err := cloudkafka.Encode(event, cloudevents.BinaryMode, cloudkafka.Transport{
		Topic: "events", Partition: kafka.ExplicitPartition(3), Timestamp: timestamp,
		Key: transportKey, Headers: []kafka.Header{{Key: "transport-attempt", Value: transportHeader}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if producer.Topic != "events" || producer.Partition != kafka.ExplicitPartition(3) ||
		!producer.Timestamp.Equal(timestamp) || string(producer.Key) != "tenant-a" {
		t.Fatalf("producer transport = %#v", producer)
	}
	transportKey[0] = 'X'
	transportHeader[0] = 'X'
	if string(producer.Key) != "tenant-a" || string(producer.Headers[len(producer.Headers)-1].Value) != "1" {
		t.Fatal("encoded Kafka record aliases transport input")
	}

	consumed := kafka.ConsumedRecord{
		Topic: producer.Topic, Key: producer.Key, Value: producer.Value, Headers: producer.Headers,
		Timestamp: producer.Timestamp, TimestampType: kafka.TimestampCreateTime,
		Partition: 3, Offset: 42, LeaderEpoch: 7,
	}
	message, state, err := cloudkafka.Decode(consumed, cloudevents.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	if message.Event.ID() != event.ID() || message.Mode != cloudevents.BinaryMode ||
		state.Topic != "events" || !state.Timestamp.Equal(timestamp) ||
		state.TimestampType != kafka.TimestampCreateTime || state.Partition != 3 ||
		state.Offset != 42 || state.LeaderEpoch != 7 ||
		len(message.TransportHeaders) != 1 || message.TransportHeaders[0].Key != "transport-attempt" {
		t.Fatalf("decoded Kafka mapping = %#v, %#v", message, state)
	}
	producer.Value[0] = 'X'
	if bytes.Equal(message.Event.Data().Bytes(), producer.Value) {
		t.Fatal("decoded event aliases producer value")
	}
}

func TestKafkaAdapterRejectsTransportCollisions(t *testing.T) {
	t.Parallel()

	event := baseEvent(t)
	if _, err := cloudkafka.Encode(event, cloudevents.BinaryMode, cloudkafka.Transport{
		Headers: []kafka.Header{{Key: "ce_id", Value: []byte("other")}},
	}); !errors.Is(err, cloudkafka.ErrMetadataCollision) {
		t.Fatalf("header collision error = %v", err)
	}
	partitionKey, err := cloudevents.NewPartitionKeyAttribute("event-key")
	if err != nil {
		t.Fatal(err)
	}
	withKey, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: "event-1", Source: "/source", Type: "example.created",
		Extensions: map[string]cloudevents.Attribute{"partitionkey": partitionKey},
	}, cloudevents.NewBinaryData([]byte("body")))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cloudkafka.Encode(withKey, cloudevents.BinaryMode, cloudkafka.Transport{Key: []byte("transport-key")}); !errors.Is(err, cloudkafka.ErrMetadataCollision) {
		t.Fatalf("key collision error = %v", err)
	}
}

func TestKafkaAdapterPropagatesBindingFailuresWithoutTransportState(t *testing.T) {
	t.Parallel()

	if record, err := cloudkafka.Encode(cloudevents.Event{}, cloudevents.BinaryMode, cloudkafka.Transport{Topic: "events"}); !errors.Is(err, cloudevents.ErrInvalidEvent) || record.Topic != "" || record.Key != nil || record.Value != nil || record.Headers != nil {
		t.Fatalf("Encode() = %#v, %v", record, err)
	}

	limits := cloudevents.DefaultLimits()
	tooManyHeaders := kafka.ConsumedRecord{
		Topic:   "events",
		Headers: make([]kafka.Header, limits.MaxKafkaHeaders+1),
	}
	if message, state, err := cloudkafka.Decode(tooManyHeaders, limits); !errors.Is(err, cloudevents.ErrLimitExceeded) || message.Event.ID() != "" || message.Key != nil || message.TransportHeaders != nil || state != (cloudkafka.State{}) {
		t.Fatalf("Decode() oversized metadata = %#v, %#v, %v", message, state, err)
	}

	invalidBinding := kafka.ConsumedRecord{
		Topic: "events",
		Value: []byte(`{"specversion":`),
		Headers: []kafka.Header{{
			Key: "content-type", Value: []byte(cloudevents.JSONMediaType),
		}},
	}
	if message, state, err := cloudkafka.Decode(invalidBinding, limits); !errors.Is(err, cloudevents.ErrInvalidEvent) || message.Event.ID() != "" || message.Key != nil || message.TransportHeaders != nil || state != (cloudkafka.State{}) {
		t.Fatalf("Decode() invalid binding = %#v, %#v, %v", message, state, err)
	}
}
