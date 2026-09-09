// Package kafka maps CloudEvents through the Golib Kafka record boundary. It
// performs no broker I/O and never owns producer or consumer lifecycle.
package kafka

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	golibkafka "github.com/faustbrian/go-kafka"
)

// ErrMetadataCollision classifies transport keys or reserved headers that
// disagree with metadata owned by the CloudEvents Kafka binding.
var ErrMetadataCollision = cloudevents.ErrMetadataCollision

// Transport supplies producer-owned Kafka routing and metadata to Encode.
// Its zero value leaves the topic empty, uses automatic partition selection,
// supplies no key or headers, and preserves an unset timestamp. Encode does
// not retain any caller-owned byte slices.
type Transport struct {
	// Topic is copied to the producer record without validation.
	Topic string
	// Partition is copied to the producer record. Its zero value requests the
	// producer's automatic keyed or unkeyed partition selection.
	Partition golibkafka.PartitionSelection
	// Key is optional transport-owned routing metadata. Encode copies it and
	// rejects a different CloudEvents partitionkey extension.
	Key []byte
	// Headers are ordered transport metadata. Encode copies their values and
	// rejects content-type and ce_ names owned by the CloudEvents binding.
	Headers []golibkafka.Header
	// Timestamp is copied unchanged. Its zero value leaves producer timestamp
	// selection to the downstream Kafka client.
	Timestamp time.Time
}

// State is immutable-by-value Kafka delivery state returned by Decode. Its
// zero value represents an input record whose broker fields were unset.
type State struct {
	// Topic is the consumed Kafka topic.
	Topic string
	// Timestamp is the record timestamp; zero means no timestamp was supplied.
	Timestamp time.Time
	// TimestampType identifies producer, broker, or unknown timestamp origin.
	TimestampType golibkafka.TimestampType
	// Partition is the consumed Kafka partition, including partition zero.
	Partition int32
	// Offset is the consumed Kafka offset, including offset zero.
	Offset int64
	// LeaderEpoch is the consumed record's leader epoch, including epoch zero.
	LeaderEpoch int32
}

// Encode maps event into a producer record without broker I/O. CloudEvents
// metadata owns content-type and ce_ headers, and the partitionkey extension
// owns the Kafka key; conflicting transport metadata returns
// ErrMetadataCollision. All returned byte slices are caller-owned copies.
func Encode(event cloudevents.Event, mode cloudevents.ContentMode, transport Transport) (golibkafka.ProducerRecord, error) {
	key := adapter.CloneBytes(transport.Key)
	if partitionKey, present := cloudevents.KafkaPartitionKey(event); present {
		if key != nil && !bytes.Equal(key, partitionKey) {
			return golibkafka.ProducerRecord{}, fmt.Errorf("%w: kafka key", ErrMetadataCollision)
		}
		key = partitionKey
	}
	for _, header := range transport.Headers {
		if header.Key == "content-type" || strings.HasPrefix(header.Key, "ce_") {
			return golibkafka.ProducerRecord{}, fmt.Errorf("%w: kafka header %s", ErrMetadataCollision, header.Key)
		}
	}
	binding, err := cloudevents.EncodeKafka(event, mode, key)
	if err != nil {
		return golibkafka.ProducerRecord{}, err
	}
	headers := make([]golibkafka.Header, 0, len(binding.Headers)+len(transport.Headers))
	for _, header := range binding.Headers {
		headers = append(headers, golibkafka.Header{Key: header.Key, Value: adapter.CloneBytes(header.Value)})
	}
	for _, header := range transport.Headers {
		headers = append(headers, golibkafka.Header{Key: header.Key, Value: adapter.CloneBytes(header.Value)})
	}
	return golibkafka.ProducerRecord{Topic: transport.Topic, Partition: transport.Partition, Key: adapter.CloneBytes(binding.Key), Value: adapter.CloneBytes(binding.Value), Headers: headers, Timestamp: transport.Timestamp}, nil
}

// Decode maps a borrowed consumed record into an owned CloudEvents message and
// immutable-by-value broker State without acknowledging or retaining the
// record. Limits are applied before copying untrusted metadata. Binding or
// limit failures return zero Message and State values.
func Decode(record golibkafka.ConsumedRecord, limits cloudevents.Limits) (cloudevents.KafkaMessage, State, error) {
	if len(record.Headers) > limits.MaxKafkaHeaders {
		return cloudevents.KafkaMessage{}, State{}, cloudevents.ErrLimitExceeded
	}
	headers := make([]cloudevents.KafkaHeader, len(record.Headers))
	for index, header := range record.Headers {
		headers[index] = cloudevents.KafkaHeader{Key: header.Key, Value: header.Value}
	}
	message, err := cloudevents.DecodeKafka(cloudevents.KafkaRecord{Key: record.Key, Value: record.Value, Headers: headers}, limits)
	if err != nil {
		return cloudevents.KafkaMessage{}, State{}, err
	}
	return message, State{Topic: record.Topic, Timestamp: record.Timestamp, TimestampType: record.TimestampType, Partition: record.Partition, Offset: record.Offset, LeaderEpoch: record.LeaderEpoch}, nil
}
