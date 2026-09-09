// Package rabbitstream maps structured CloudEvents through Golib RabbitMQ
// Stream messages without claiming a non-existent binary AMQP binding.
package rabbitstream

import (
	"bytes"
	"fmt"
	"mime"
	"strings"
	"time"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	golibrabbitstream "github.com/faustbrian/go-rabbitmq-streams"
)

var (
	// ErrInvalidInput classifies malformed optional-adapter input while
	// preserving any more specific wrapped cause.
	ErrInvalidInput = cloudevents.ErrInvalidAdapterInput
	// ErrMetadataCollision classifies ambiguous targets or transport metadata
	// that disagrees with canonical CloudEvents metadata.
	ErrMetadataCollision = cloudevents.ErrMetadataCollision
)

// Transport supplies RabbitMQ Stream routing and application properties to
// Encode. Exactly one of Stream and SuperStream must be non-empty. The zero
// value is therefore invalid. Encode does not retain Properties or their byte
// slices.
type Transport struct {
	// Stream is the direct stream target and is mutually exclusive with
	// SuperStream.
	Stream string
	// SuperStream is the partitioned logical target and is mutually exclusive
	// with Stream.
	SuperStream string
	// RoutingKey selects a Super Stream partition. A non-empty CloudEvents
	// partitionkey extension owns this value and must agree with it.
	RoutingKey string
	// CorrelationID is opaque transport metadata copied unchanged.
	CorrelationID string
	// Properties preserves ordered application metadata, including duplicate
	// keys. Encode returns a deep copy after validation.
	Properties []golibrabbitstream.MetadataEntry
}

// Message is the decoded CloudEvent and transport-owned application metadata.
// Event and TransportProperties own their returned data. The zero value means
// no message has been decoded.
type Message struct {
	// Event is the decoded structured CloudEvent.
	Event cloudevents.Event
	// TransportProperties preserves ordered RabbitMQ application properties,
	// including duplicate keys, as caller-owned copies.
	TransportProperties []golibrabbitstream.MetadataEntry
}

// State is immutable-by-value RabbitMQ Stream delivery state returned by
// Decode. Presence flags distinguish explicit zero publishing IDs and offsets
// from values that were not supplied.
type State struct {
	// Stream is the direct target or consumed backing stream.
	Stream string
	// SuperStream is the logical partitioned target, if supplied.
	SuperStream string
	// Partition is the selected backing stream after delivery, if supplied.
	Partition string
	// RoutingKey is the transport routing key used for partition selection.
	RoutingKey string
	// CorrelationID is opaque transport correlation metadata.
	CorrelationID string
	// PublishingID is meaningful only when HasPublishingID is true.
	PublishingID uint64
	// HasPublishingID distinguishes an explicit zero ID from an unset ID.
	HasPublishingID bool
	// Offset is meaningful only when HasOffset is true.
	Offset uint64
	// HasOffset distinguishes explicit offset zero from an unset offset.
	HasOffset bool
	// Timestamp is application event time, not broker receipt time. Its zero
	// value means the CloudEvent did not supply a time.
	Timestamp time.Time
}

// Encode maps event into one structured RabbitMQ Stream message without broker
// I/O. It requires exactly one target, preserves ordered Properties by deep
// copy, and rejects routing metadata that conflicts with the CloudEvents
// partitionkey extension. The returned message owns its payload and metadata.
func Encode(event cloudevents.Event, transport Transport) (golibrabbitstream.Message, error) {
	if (transport.Stream == "") == (transport.SuperStream == "") {
		return golibrabbitstream.Message{}, fmt.Errorf("%w: RabbitStream target", ErrMetadataCollision)
	}
	routingKey := transport.RoutingKey
	if partitionKey, present := cloudevents.KafkaPartitionKey(event); present {
		if routingKey != "" && routingKey != string(partitionKey) {
			return golibrabbitstream.Message{}, fmt.Errorf("%w: RabbitStream routing key", ErrMetadataCollision)
		}
		routingKey = string(partitionKey)
	}
	payload, err := cloudevents.EncodeJSON(event)
	if err != nil {
		return golibrabbitstream.Message{}, err
	}
	timestamp, _ := event.Time()
	message := golibrabbitstream.Message{Stream: transport.Stream, SuperStream: transport.SuperStream, RoutingKey: routingKey, MessageID: event.ID(), CorrelationID: transport.CorrelationID, Timestamp: timestamp, ContentType: cloudevents.JSONMediaType, Payload: payload, Properties: transport.Properties}
	if err := message.Validate(golibrabbitstream.DefaultLimits()); err != nil {
		return golibrabbitstream.Message{}, err
	}
	message.Properties = cloneMetadata(transport.Properties)
	return message, nil
}

// Decode maps a borrowed structured RabbitMQ Stream message into an owned
// CloudEvent Message and immutable-by-value State without acknowledging the
// delivery. It rejects non-JSON modes and transport identifiers that conflict
// with canonical event metadata. Limits apply to the CloudEvent payload;
// RabbitMQ transport metadata is bounded by the dependency's default limits.
func Decode(message golibrabbitstream.Message, limits cloudevents.Limits) (Message, State, error) {
	if err := validateMessage(message); err != nil {
		return Message{}, State{}, err
	}
	mediaType, _, err := mime.ParseMediaType(message.ContentType)
	if err != nil || !strings.EqualFold(mediaType, cloudevents.JSONMediaType) {
		return Message{}, State{}, cloudevents.ErrUnsupportedMode
	}
	event, err := cloudevents.DecodeJSON(message.Payload, limits)
	if err != nil {
		return Message{}, State{}, err
	}
	if message.MessageID != "" && message.MessageID != event.ID() {
		return Message{}, State{}, fmt.Errorf("%w: RabbitStream message ID", ErrMetadataCollision)
	}
	if partitionKey, present := cloudevents.KafkaPartitionKey(event); present && message.RoutingKey != "" && !bytes.Equal(partitionKey, []byte(message.RoutingKey)) {
		return Message{}, State{}, fmt.Errorf("%w: RabbitStream routing key", ErrMetadataCollision)
	}
	return Message{Event: event, TransportProperties: cloneMetadata(message.Properties)}, State{Stream: message.Stream, SuperStream: message.SuperStream, Partition: message.Partition, RoutingKey: message.RoutingKey, CorrelationID: message.CorrelationID, PublishingID: message.PublishingID, HasPublishingID: message.HasPublishingID, Offset: message.Offset, HasOffset: message.HasOffset, Timestamp: message.Timestamp}, nil
}

func validateMessage(message golibrabbitstream.Message) error {
	limits := golibrabbitstream.DefaultLimits()
	if message.Partition != "" || message.HasOffset || message.Offset != 0 {
		return message.ValidateDelivery(limits)
	}
	return message.Validate(limits)
}

func cloneMetadata(entries []golibrabbitstream.MetadataEntry) []golibrabbitstream.MetadataEntry {
	if entries == nil {
		return nil
	}
	cloned := make([]golibrabbitstream.MetadataEntry, len(entries))
	for index, entry := range entries {
		cloned[index] = golibrabbitstream.MetadataEntry{Key: entry.Key, Value: adapter.CloneBytes(entry.Value)}
	}
	return cloned
}
