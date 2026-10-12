package golib_test

import (
	"testing"
	"time"

	"github.com/faustbrian/go-cloudevents"
	golib "github.com/faustbrian/go-cloudevents/adapters/golib/v4"
	"github.com/faustbrian/go-kafka"
	"github.com/faustbrian/go-queue/job"
	"github.com/faustbrian/go-rabbitmq-streams"
)

// BenchmarkKafkaBinaryRoundTrip measures decode only: event construction and binary encode are setup, with no broker operation.
func BenchmarkKafkaBinaryRoundTrip(b *testing.B) {
	event := benchmarkEvent(b)
	producer, err := golib.EncodeKafka(event, cloudevents.BinaryMode, golib.KafkaTransport{Topic: "events"})
	if err != nil {
		b.Fatal(err)
	}
	record := kafka.ConsumedRecord{Topic: producer.Topic, Key: producer.Key, Value: producer.Value, Headers: producer.Headers}

	b.ReportAllocs()
	for b.Loop() {
		if _, _, err := golib.DecodeKafka(record, cloudevents.DefaultLimits()); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRabbitStreamStructuredRoundTrip measures decode only of one preencoded structured message plus one transport property; it performs no Stream I/O.
func BenchmarkRabbitStreamStructuredRoundTrip(b *testing.B) {
	event := benchmarkEvent(b)
	message, err := golib.EncodeRabbitStream(event, golib.RabbitStreamTransport{Stream: "events"})
	if err != nil {
		b.Fatal(err)
	}
	message.Properties = []rabbitstream.MetadataEntry{{Key: "transport-attempt", Value: []byte("1")}}
	b.SetBytes(int64(len(message.Payload)))
	b.ReportAllocs()
	for b.Loop() {
		if _, _, err := golib.DecodeRabbitStream(message, cloudevents.DefaultLimits()); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkQueueRoundTrip includes both conversions of one fixed payload/retained job state, without enqueueing or execution.
func BenchmarkQueueRoundTrip(b *testing.B) {
	message := job.Message{Timeout: time.Second, Body: []byte(`{"order":"A-123"}`)}
	options := golib.QueueOptions{Source: "/queue", StableID: "job-1", Type: "order.created"}

	b.ReportAllocs()
	for b.Loop() {
		event, state, _, err := golib.QueueToCloudEvent(message, options)
		if err != nil {
			b.Fatal(err)
		}
		if _, _, err := golib.CloudEventToQueue(event, state); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkEvent(b *testing.B) cloudevents.Event {
	b.Helper()
	data, err := cloudevents.NewJSONData([]byte(`{"order":"A-123"}`))
	if err != nil {
		b.Fatal(err)
	}
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: "event-1", Source: "/orders", Type: "order.created", DataContentType: "application/json",
	}, data)
	if err != nil {
		b.Fatal(err)
	}
	return event
}
