// Package adapter contains shared mechanics for independently released
// CloudEvents target adapters. It is internal so target modules share error and
// conversion semantics without creating another public abstraction.
package adapter

import (
	"bytes"
	"fmt"
	"mime"
	"sort"
	"strings"
	"time"

	"github.com/faustbrian/go-cloudevents"
)

const (
	CorrelationIDExtension = "correlationid"
	RequestIDExtension     = "requestid"
	CausationIDExtension   = "causationid"
	TenantIDExtension      = "tenantid"
)

func AddStringExtensions(event cloudevents.Event, additions map[string]string) (cloudevents.Event, error) {
	if err := event.Validate(); err != nil {
		return cloudevents.Event{}, fmt.Errorf("%w: event: %w", cloudevents.ErrInvalidAdapterInput, err)
	}
	extensions := event.Extensions()
	for name, value := range additions {
		attribute, err := cloudevents.NewStringAttribute(value)
		if err != nil {
			return cloudevents.Event{}, fmt.Errorf("%w: extension %s: %w", cloudevents.ErrInvalidAdapterInput, name, err)
		}
		if existing, present := extensions[name]; present {
			if !AttributesEqual(existing, attribute) {
				return cloudevents.Event{}, fmt.Errorf("%w: %s", cloudevents.ErrMetadataCollision, name)
			}
		} else {
			extensions[name] = attribute
		}
	}
	return RebuildEvent(event, extensions)
}

func StringExtension(event cloudevents.Event, name string) (string, bool, error) {
	attribute, present := event.Extension(name)
	if !present {
		return "", false, nil
	}
	if attribute.Kind() != cloudevents.AttributeString || attribute.String() == "" {
		return "", false, fmt.Errorf("%w: extension %s", cloudevents.ErrInvalidAdapterInput, name)
	}
	return attribute.String(), true, nil
}

func AttributesEqual(left, right cloudevents.Attribute) bool {
	return left.Kind() == right.Kind() && left.String() == right.String() && bytes.Equal(left.Bytes(), right.Bytes())
}

func RebuildEvent(event cloudevents.Event, extensions map[string]cloudevents.Attribute) (cloudevents.Event, error) {
	dataContentType, _ := event.DataContentType()
	dataSchema, _ := event.DataSchema()
	subject, _ := event.Subject()
	timeValue, hasTime := event.Time()
	var occurredAt *time.Time
	if hasTime {
		occurredAt = &timeValue
	}
	return cloudevents.NewEvent(cloudevents.Attributes{
		ID: event.ID(), Source: event.Source(), Type: event.Type(),
		DataContentType: dataContentType, DataSchema: dataSchema, Subject: subject,
		Time: occurredAt, Extensions: extensions,
	}, event.Data())
}

func DataFromPayload(contentType string, payload []byte) (cloudevents.Data, error) {
	if contentType == "" {
		return cloudevents.NewBinaryData(payload), nil
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return cloudevents.Data{}, fmt.Errorf("%w: content type: %w", cloudevents.ErrInvalidAdapterInput, err)
	}
	mediaType = strings.ToLower(mediaType)
	if mediaType == "application/json" || strings.HasSuffix(mediaType, "+json") {
		return cloudevents.NewJSONData(payload)
	}
	if strings.HasPrefix(mediaType, "text/") {
		return cloudevents.NewTextData(string(payload))
	}
	return cloudevents.NewBinaryData(payload), nil
}

func RestoreRetainedNil(data cloudevents.Data, retainedNil bool) []byte {
	value := data.Bytes()
	if retainedNil && len(value) == 0 {
		return nil
	}
	return value
}

func PutStringExtension(extensions map[string]cloudevents.Attribute, name, value string) error {
	attribute, err := cloudevents.NewStringAttribute(value)
	if err != nil {
		return fmt.Errorf("%w: extension %s: %w", cloudevents.ErrInvalidAdapterInput, name, err)
	}
	extensions[name] = attribute
	return nil
}

func MappedString(event cloudevents.Event, name, retained string, trusted bool) (string, error) {
	value, present, err := StringExtension(event, name)
	if err != nil {
		return "", err
	}
	if present && retained != "" && value != retained {
		return "", fmt.Errorf("%w: %s", cloudevents.ErrMetadataCollision, name)
	}
	if present {
		if retained == "" && !trusted {
			return "", fmt.Errorf("%w: %s", cloudevents.ErrUntrustedMetadata, name)
		}
		return value, nil
	}
	if retained != "" {
		return "", fmt.Errorf("%w: missing %s", cloudevents.ErrMetadataCollision, name)
	}
	return retained, nil
}

func CloneBytes(value []byte) []byte {
	if value == nil {
		return nil
	}
	cloned := make([]byte, len(value))
	copy(cloned, value)
	return cloned
}

func CloneStrings(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

func CloneTimePointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func AppendOptionalContextLosses(event cloudevents.Event, report *cloudevents.AdapterReport, target string) {
	for _, field := range []string{"dataschema", "subject", "time"} {
		present := false
		switch field {
		case "subject":
			_, present = event.Subject()
		case "dataschema":
			_, present = event.DataSchema()
		case "time":
			_, present = event.Time()
		}
		if present {
			report.Losses = append(report.Losses, cloudevents.AdapterLoss{Field: field, Reason: "not represented by " + target})
		}
	}
}

func AppendDataKindLoss(event cloudevents.Event, report *cloudevents.AdapterReport, target string, contentTypePreserved bool) {
	data := event.Data()
	if !data.Present() {
		return
	}
	targetKind := cloudevents.DataBinary
	if contentTypePreserved {
		if contentType, present := event.DataContentType(); present {
			mediaType, _, _ := mime.ParseMediaType(contentType)
			mediaType = strings.ToLower(mediaType)
			switch {
			case mediaType == "application/json" || strings.HasSuffix(mediaType, "+json"):
				targetKind = cloudevents.DataJSON
			case strings.HasPrefix(mediaType, "text/"):
				targetKind = cloudevents.DataText
			}
		}
	}
	if data.Kind() != targetKind {
		report.Losses = append(report.Losses, cloudevents.AdapterLoss{Field: "data.kind", Reason: "not represented by " + target})
	}
}

func AppendExtensionLosses(event cloudevents.Event, report *cloudevents.AdapterReport, target string, selected map[string]struct{}) {
	var names []string
	for name := range event.Extensions() {
		if _, retained := selected[name]; !retained {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		report.Losses = append(report.Losses, cloudevents.AdapterLoss{Field: "extensions." + name, Reason: "not represented by " + target})
	}
}

func IsJSONContentType(value string) bool {
	mediaType, _, err := mime.ParseMediaType(value)
	if err != nil {
		return false
	}
	mediaType = strings.ToLower(mediaType)
	return mediaType == "application/json" || strings.HasSuffix(mediaType, "+json")
}
