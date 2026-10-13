package cloudevents_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/faustbrian/go-cloudevents"
)

const admissionEventPrefix = `{"specversion":"1.0","id":"r13","source":"/source","type":"example",`

func TestReceivedBase64RejectsLineBreaksAcrossBindings(t *testing.T) {
	for _, text := range []string{`YQ\r==`, `YQ\n==`, `YQ==\r\n`, `\r\nYQ==`, `Y\u000dQ==`} {
		input := []byte(admissionEventPrefix + `"data_base64":"` + text + `"}`)
		t.Run(text, func(t *testing.T) {
			t.Run("JSON", func(t *testing.T) {
				event, err := cloudevents.DecodeJSON(input, cloudevents.DefaultLimits())
				if !errors.Is(err, cloudevents.ErrInvalidEvent) || event.ID() != "" || event.Data().Present() {
					t.Fatalf("JSON admitted noncanonical received base64: event=%q, present=%v, err=%v", event.ID(), event.Data().Present(), err)
				}
			})
			t.Run("HTTP structured", func(t *testing.T) {
				body := &admissionBody{Reader: bytes.NewReader(input)}
				message, err := cloudevents.DecodeHTTP(context.Background(), http.Header{"Content-Type": {cloudevents.JSONMediaType}}, body, cloudevents.DefaultLimits())
				if !errors.Is(err, cloudevents.ErrInvalidEvent) || len(message.Events) != 0 || body.closed {
					t.Fatalf("structured HTTP admission/borrowed ownership: events=%d, closed=%v, err=%v", len(message.Events), body.closed, err)
				}
			})
			t.Run("Kafka structured", func(t *testing.T) {
				record := cloudevents.KafkaRecord{Value: input, Headers: []cloudevents.KafkaHeader{{Key: "content-type", Value: []byte(cloudevents.JSONMediaType)}}}
				kafka, err := cloudevents.DecodeKafka(record, cloudevents.DefaultLimits())
				if !errors.Is(err, cloudevents.ErrInvalidEvent) || kafka.Event.ID() != "" || kafka.Event.Data().Present() {
					t.Fatalf("structured Kafka admitted noncanonical base64: event=%q, err=%v", kafka.Event.ID(), err)
				}
			})
		})
	}
}

// data_base64 is a binary JSON-string member, not a nullable context attribute.
// In particular, null must not silently become present empty binary data.
func TestReceivedBase64RejectsNullInsteadOfEmptyBinary(t *testing.T) {
	for _, raw := range []string{"null", " \nnull\t "} {
		input := []byte(admissionEventPrefix + `"data_base64":` + raw + `}`)
		event, err := cloudevents.DecodeJSON(input, cloudevents.DefaultLimits())
		if !errors.Is(err, cloudevents.ErrInvalidEvent) || event.ID() != "" || event.Data().Present() {
			t.Fatalf("null base64 became binary data: id=%q, present=%v, err=%v", event.ID(), event.Data().Present(), err)
		}
	}
}

func TestReceivedBase64PreservesCanonicalAndEscapedStrings(t *testing.T) {
	for _, test := range []struct {
		name string
		text string
		want string
	}{
		{name: "canonical", text: "YQ==", want: "a"},
		{name: "escaped JSON spelling", text: `\u0059\u0051==`, want: "a"},
		{name: "empty binary", text: "", want: ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			input := []byte(admissionEventPrefix + `"data_base64":"` + test.text + `"}`)
			event, err := cloudevents.DecodeJSON(input, cloudevents.DefaultLimits())
			if err != nil || !event.Data().Present() || event.Data().Kind() != cloudevents.DataBinary || string(event.Data().Bytes()) != test.want {
				t.Fatalf("binary control: kind=%v, present=%v, bytes=%q, err=%v", event.Data().Kind(), event.Data().Present(), event.Data().Bytes(), err)
			}
		})
	}
	limits := cloudevents.DefaultLimits()
	limits.MaxDataBytes = 1
	if _, err := cloudevents.DecodeJSON([]byte(admissionEventPrefix+`"data_base64":"YQ=="}`), limits); err != nil {
		t.Fatalf("exact decoded byte limit: %v", err)
	}
	if _, err := cloudevents.DecodeJSON([]byte(admissionEventPrefix+`"data_base64":"YWI="}`), limits); !errors.Is(err, cloudevents.ErrLimitExceeded) {
		t.Fatalf("decoded byte overflow: %v", err)
	}
	if _, err := cloudevents.DecodeJSON([]byte(admissionEventPrefix+`"data_base64":"Zh=="}`), cloudevents.DefaultLimits()); !errors.Is(err, cloudevents.ErrInvalidEvent) {
		t.Fatalf("nonzero trailing bits admitted: %v", err)
	}
}

func TestJSONBatchRequiresOuterArrayAcrossHTTP(t *testing.T) {
	for _, input := range []string{"null", " \r\nnull\t ", "[null]", "[" + admissionEventPrefix + `"data":null},null]`} {
		t.Run(input, func(t *testing.T) {
			t.Run("JSON", func(t *testing.T) {
				events, err := cloudevents.DecodeJSONBatch([]byte(input), cloudevents.DefaultLimits())
				if !errors.Is(err, cloudevents.ErrInvalidEvent) || events != nil {
					t.Fatalf("batch admitted non-event/outer null: events=%v, err=%v", events, err)
				}
			})
			t.Run("HTTP batch", func(t *testing.T) {
				body := &admissionBody{Reader: bytes.NewReader([]byte(input))}
				message, err := cloudevents.DecodeHTTP(context.Background(), http.Header{"Content-Type": {cloudevents.JSONBatchMediaType}}, body, cloudevents.DefaultLimits())
				if !errors.Is(err, cloudevents.ErrInvalidEvent) || len(message.Events) != 0 || body.closed {
					t.Fatalf("HTTP batch admission/borrowed ownership: events=%d, closed=%v, err=%v", len(message.Events), body.closed, err)
				}
			})
		})
	}
}

func TestJSONBatchPreservesEmptyArraysAndEventNullFields(t *testing.T) {
	for _, input := range []string{"[]", " \r\n[ \t ]\n "} {
		events, err := cloudevents.DecodeJSONBatch([]byte(input), cloudevents.DefaultLimits())
		if err != nil || events == nil || len(events) != 0 {
			t.Fatalf("empty array control: events=%v, err=%v", events, err)
		}
	}
	input := []byte("[ \n" + admissionEventPrefix + `"subject":null,"time":null,"data":null}` + "\n ]")
	events, err := cloudevents.DecodeJSONBatch(input, cloudevents.DefaultLimits())
	if err != nil || len(events) != 1 {
		t.Fatalf("event-local null control: events=%d, err=%v", len(events), err)
	}
	event := events[0]
	if !event.Data().Present() || event.Data().Kind() != cloudevents.DataJSON || string(event.Data().Bytes()) != "null" {
		t.Fatalf("null data lost: kind=%v, present=%v, bytes=%q", event.Data().Kind(), event.Data().Present(), event.Data().Bytes())
	}
	if _, present := event.Subject(); present {
		t.Fatal("null optional subject became present")
	}
	if _, present := event.Time(); present {
		t.Fatal("null optional time became present")
	}
	limits := cloudevents.DefaultLimits()
	limits.MaxBatchEvents = 0
	if _, err := cloudevents.DecodeJSONBatch([]byte("[]"), limits); !errors.Is(err, cloudevents.ErrLimitExceeded) {
		t.Fatalf("zero batch admission limit changed: %v", err)
	}
}

type admissionBody struct {
	*bytes.Reader
	closed bool
}

func (b *admissionBody) Close() error {
	b.closed = true
	return nil
}
