package cloudevents

import (
	"fmt"
	"mime"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

// KafkaHeader is a Kafka record header. Value is always caller-owned on public
// input and output boundaries.
type KafkaHeader struct {
	// Key is the exact case-sensitive header name; content-type and ce_ belong to the binding.
	Key string
	// Value is borrowed on input and copied when retained; nil headers are not record tombstones.
	Value []byte
}

// KafkaRecord is the transport-neutral portion of a Kafka record owned by the
// CloudEvents binding. Topic, partition, offset, timestamp, retries, and broker
// settlement remain owned by the Kafka caller.
type KafkaRecord struct {
	// Key is borrowed on input/copied on output; nil remains distinct from non-nil empty bytes.
	Key []byte
	// Value is borrowed on input/copied on output; nil binary record data is a tombstone, not an empty payload.
	Value []byte
	// Headers is ordered input metadata; decoder-owned copies do not imply broker delivery.
	Headers []KafkaHeader
}

// KafkaMessage is a decoded binding result. Transport headers are headers not
// owned by the CloudEvents binding and are returned without interpretation.
type KafkaMessage struct {
	// Mode records the detected structured/binary binding.
	Mode ContentMode
	// Event is the decoded immutable value, not an acknowledged message.
	Event Event
	// Key is borrowed on input/copied on output; nil remains distinct from non-nil empty bytes.
	Key []byte
	// TransportHeaders contains copied unowned headers in received order; values are not interpreted/redacted.
	TransportHeaders []KafkaHeader
}

// EncodeKafka maps one Event without implicit representation loss. Use
// EncodeKafkaWithReport to accept and inspect target-binding changes.
func EncodeKafka(event Event, mode ContentMode, key []byte) (KafkaRecord, error) {
	record, report, err := EncodeKafkaWithReport(event, mode, key)
	if err != nil {
		return KafkaRecord{}, err
	}
	if err := rejectConversionLoss(report); err != nil {
		return KafkaRecord{}, err
	}
	return record, nil
}

// EncodeKafkaWithReport maps one Event to the stable Kafka protocol binding
// while reporting representation changes. Only binary/structured modes exist;
// batch fails with ErrUnsupportedMode. Key/header/value bytes are copied with
// no encode Limits or broker I/O. Absent binary data becomes nil (a tombstone);
// present empty data stays non-nil empty (Kafka 1.0.2, CLOUDEVENTS-DEC-010).
func EncodeKafkaWithReport(event Event, mode ContentMode, key []byte) (KafkaRecord, ConversionReport, error) {
	switch mode {
	case BinaryMode:
		if err := event.Validate(); err != nil {
			return KafkaRecord{}, ConversionReport{}, err
		}
		record := KafkaRecord{Key: cloneBytesPreservingNil(key), Headers: encodeKafkaBinaryHeaders(event)}
		if event.data.present {
			record.Value = cloneBytesPreservingEmpty(event.data.bytes)
		}
		return record, binaryConversionReport(event), nil
	case StructuredMode:
		value, report, err := EncodeJSONWithReport(event)
		if err != nil {
			return KafkaRecord{}, ConversionReport{}, err
		}
		return KafkaRecord{
			Key: cloneBytesPreservingNil(key), Value: value,
			Headers: []KafkaHeader{{Key: "content-type", Value: []byte(JSONMediaType)}},
		}, report, nil
	default:
		return KafkaRecord{}, ConversionReport{}, ErrUnsupportedMode
	}
}

// KafkaPartitionKey implements the official binding's opt-in partitionkey
// mapper. It does not modify the Event.
func KafkaPartitionKey(event Event) ([]byte, bool) {
	attribute, present := event.extensions["partitionkey"]
	if !present || attribute.kind != AttributeString || attribute.text == "" {
		return nil, false
	}
	return []byte(attribute.text), true
}

func encodeKafkaBinaryHeaders(event Event) []KafkaHeader {
	headers := make([]KafkaHeader, 0)
	appendAttribute := func(name, value string) {
		headers = append(headers, KafkaHeader{Key: "ce_" + name, Value: []byte(value)})
	}
	appendAttribute("specversion", event.SpecVersion())
	appendAttribute("id", event.id)
	appendAttribute("source", event.source)
	appendAttribute("type", event.eventType)
	if event.dataSchema != "" {
		appendAttribute("dataschema", event.dataSchema)
	}
	if event.subject != "" {
		appendAttribute("subject", event.subject)
	}
	if event.hasTime {
		appendAttribute("time", event.time.UTC().Format(time.RFC3339Nano))
	}
	for name, attribute := range event.extensions {
		appendAttribute(name, attribute.String())
	}
	if event.dataContentType != "" {
		headers = append(headers, KafkaHeader{Key: "content-type", Value: []byte(event.dataContentType)})
	} else if event.data.present {
		headers = append(headers, KafkaHeader{Key: "content-type", Value: []byte(implicitDataContentType(event.data.kind))})
	}
	slices.SortFunc(headers, func(left, right KafkaHeader) int {
		return strings.Compare(left.Key, right.Key)
	})
	return headers
}

// Kafka tombstones and adapter retained-state reversal depend on nil versus
// present-empty bytes; append-to-nil would collapse that distinction (DEC-010).
func cloneBytesPreservingNil(value []byte) []byte {
	if value == nil {
		return nil
	}
	return cloneBytesPreservingEmpty(value)
}

func cloneBytesPreservingEmpty(value []byte) []byte {
	owned := make([]byte, len(value))
	copy(owned, value)
	return owned
}

// DecodeKafka decodes Kafka 1.0.2 without broker I/O. Supplied limits are not
// defaulted; copied key/unowned headers receive their own admission. Content-Type
// selects structured JSON or binary mode; singleton metadata/conflicts fail.
// Nil binary Value means absent data/tombstone; non-nil empty means present.
// Received key is retained independently of partitionkey, without enforcing
// equality. All failures return zero message; settlement/retry remain external.
func DecodeKafka(record KafkaRecord, limits Limits) (KafkaMessage, error) {
	if limits.MaxEventBytes <= 0 || limits.MaxDataBytes <= 0 ||
		limits.MaxKafkaKeyBytes <= 0 || len(record.Key) > limits.MaxKafkaKeyBytes ||
		limits.MaxKafkaHeaders <= 0 || len(record.Headers) > limits.MaxKafkaHeaders ||
		limits.MaxKafkaHeaderNameBytes <= 0 || limits.MaxKafkaHeaderValueBytes <= 0 {
		return KafkaMessage{}, ErrLimitExceeded
	}
	contentTypes := make([][]byte, 0, 1)
	cloudEventHeaders := make(map[string][][]byte)
	transportHeaders := make([]KafkaHeader, 0)
	for _, header := range record.Headers {
		if len(header.Key) > limits.MaxKafkaHeaderNameBytes || len(header.Value) > limits.MaxKafkaHeaderValueBytes {
			return KafkaMessage{}, ErrLimitExceeded
		}
		if !utf8.ValidString(header.Key) {
			return KafkaMessage{}, fmt.Errorf("%w: kafka header name", ErrInvalidEvent)
		}
		if header.Key == "content-type" {
			contentTypes = append(contentTypes, header.Value)
		} else if strings.HasPrefix(header.Key, "ce_") {
			name := strings.TrimPrefix(header.Key, "ce_")
			cloudEventHeaders[name] = append(cloudEventHeaders[name], header.Value)
		} else {
			transportHeaders = append(transportHeaders, KafkaHeader{
				Key:   header.Key,
				Value: cloneBytesPreservingNil(header.Value),
			})
		}
	}
	if len(contentTypes) > 1 {
		return KafkaMessage{}, fmt.Errorf("%w: duplicate content-type", ErrInvalidEvent)
	}
	contentType := ""
	mode := BinaryMode
	if len(contentTypes) == 1 {
		if !utf8.Valid(contentTypes[0]) {
			return KafkaMessage{}, fmt.Errorf("%w: content-type", ErrInvalidEvent)
		}
		contentType = string(contentTypes[0])
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return KafkaMessage{}, fmt.Errorf("%w: content-type", ErrInvalidEvent)
		}
		lowerMediaType := strings.ToLower(mediaType)
		if lowerMediaType == JSONMediaType {
			mode = StructuredMode
		} else if strings.HasPrefix(lowerMediaType, "application/cloudevents") {
			return KafkaMessage{}, ErrUnsupportedMode
		}
	}

	var event Event
	var err error
	if mode == StructuredMode {
		if record.Value == nil {
			return KafkaMessage{}, fmt.Errorf("%w: structured value", ErrInvalidEvent)
		}
		if int64(len(record.Value)) > limits.MaxEventBytes {
			return KafkaMessage{}, ErrLimitExceeded
		}
		event, err = DecodeJSON(record.Value, limits)
		if err == nil {
			err = validateStructuredKafkaMetadata(cloudEventHeaders, event, limits)
		}
	} else {
		event, err = decodeKafkaBinary(record.Value, contentType, cloudEventHeaders, limits)
	}
	if err != nil {
		return KafkaMessage{}, err
	}
	return KafkaMessage{
		Mode:             mode,
		Event:            event,
		Key:              cloneBytesPreservingNil(record.Key),
		TransportHeaders: transportHeaders,
	}, nil
}

func validateStructuredKafkaMetadata(headers map[string][][]byte, event Event, limits Limits) error {
	expected := map[string]string{
		"specversion": event.SpecVersion(),
		"id":          event.id,
		"source":      event.source,
		"type":        event.eventType,
	}
	if event.dataContentType != "" {
		expected["datacontenttype"] = event.dataContentType
	}
	if event.dataSchema != "" {
		expected["dataschema"] = event.dataSchema
	}
	if event.subject != "" {
		expected["subject"] = event.subject
	}
	if event.hasTime {
		expected["time"] = event.time.UTC().Format(time.RFC3339Nano)
	}
	for name, attribute := range event.extensions {
		expected[name] = attribute.String()
	}
	if len(headers) > limits.MaxAttributes {
		return ErrLimitExceeded
	}
	for name, values := range headers {
		if limits.MaxAttributeNameBytes <= 0 || len(name) > limits.MaxAttributeNameBytes {
			return ErrLimitExceeded
		}
		if len(values) != 1 {
			return fmt.Errorf("%w: duplicate ce_%s", ErrInvalidEvent, name)
		}
		if limits.MaxAttributeValueBytes <= 0 || len(values[0]) > limits.MaxAttributeValueBytes {
			return ErrLimitExceeded
		}
		if !utf8.Valid(values[0]) {
			return fmt.Errorf("%w: ce_%s", ErrInvalidEvent, name)
		}
		want, present := expected[name]
		if !present || string(values[0]) != want {
			return fmt.Errorf("%w: conflicting ce_%s", ErrInvalidEvent, name)
		}
	}
	return nil
}

func decodeKafkaBinary(value []byte, contentType string, headers map[string][][]byte, limits Limits) (Event, error) {
	if _, present := headers["datacontenttype"]; present {
		return Event{}, fmt.Errorf("%w: ce_datacontenttype", ErrInvalidEvent)
	}
	if len(headers) > limits.MaxAttributes {
		return Event{}, ErrLimitExceeded
	}
	decoded := make(map[string]string, len(headers))
	for name, values := range headers {
		if limits.MaxAttributeNameBytes <= 0 || len(name) > limits.MaxAttributeNameBytes {
			return Event{}, ErrLimitExceeded
		}
		if len(values) != 1 {
			return Event{}, fmt.Errorf("%w: duplicate ce_%s", ErrInvalidEvent, name)
		}
		if limits.MaxAttributeValueBytes <= 0 || len(values[0]) > limits.MaxAttributeValueBytes {
			return Event{}, ErrLimitExceeded
		}
		if !utf8.Valid(values[0]) {
			return Event{}, fmt.Errorf("%w: ce_%s", ErrInvalidEvent, name)
		}
		decoded[name] = string(values[0])
	}
	if decoded["specversion"] != specVersion {
		return Event{}, fmt.Errorf("%w: specversion", ErrInvalidEvent)
	}
	var occurredAt *time.Time
	if timeValue, present := decoded["time"]; present {
		parsed, err := time.Parse(time.RFC3339Nano, timeValue)
		if err != nil {
			return Event{}, fmt.Errorf("%w: time", ErrInvalidEvent)
		}
		occurredAt = &parsed
	}
	extensions := make(map[string]Attribute)
	for name, attributeValue := range decoded {
		if _, reserved := reservedAttributeNames[name]; reserved {
			continue
		}
		attribute, err := NewStringAttribute(attributeValue)
		if err != nil {
			return Event{}, fmt.Errorf("%w: ce_%s", ErrInvalidEvent, name)
		}
		extensions[name] = attribute
	}
	data := Data{}
	if value != nil {
		if int64(len(value)) > limits.MaxDataBytes {
			return Event{}, ErrLimitExceeded
		}
		switch {
		case isJSONMediaType(contentType):
			jsonData, err := NewJSONData(value)
			if err != nil {
				return Event{}, err
			}
			data = jsonData
		case isTextMediaType(contentType):
			textData, err := NewTextData(string(value))
			if err != nil {
				return Event{}, err
			}
			data = textData
		default:
			data = NewBinaryData(value)
		}
	}
	return NewEvent(Attributes{
		ID:              decoded["id"],
		Source:          decoded["source"],
		Type:            decoded["type"],
		DataContentType: contentType,
		DataSchema:      decoded["dataschema"],
		Subject:         decoded["subject"],
		Time:            occurredAt,
		Extensions:      extensions,
	}, data)
}
