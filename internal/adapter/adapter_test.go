package adapter

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/faustbrian/go-cloudevents"
)

func TestCanonicalAdapterErrorsPreserveLegacyStrings(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		err  error
		want string
	}{
		{cloudevents.ErrInvalidAdapterInput, "cloudevents golib adapter: invalid input"},
		{cloudevents.ErrMetadataCollision, "cloudevents golib adapter: metadata collision"},
		{cloudevents.ErrUntrustedMetadata, "cloudevents golib adapter: metadata is untrusted"},
		{cloudevents.ErrSchemaViolation, "cloudevents golib adapter: schema violation"},
		{cloudevents.ErrSchemaMapping, "cloudevents golib adapter: schema mapping"},
	} {
		if got := test.err.Error(); got != test.want {
			t.Errorf("error string = %q, want %q", got, test.want)
		}
	}
}

func TestSharedAdapterMechanicsPreserveValuesAndClassifyFailures(t *testing.T) {
	t.Parallel()

	if _, err := AddStringExtensions(cloudevents.Event{}, nil); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
		t.Fatalf("invalid event error = %v", err)
	}
	event := helperEvent(t, map[string]cloudevents.Attribute{"same": mustString(t, "value")}, true)
	converted, err := AddStringExtensions(event, map[string]string{"same": "value", "added": "new"})
	if err != nil {
		t.Fatal(err)
	}
	if value, present := converted.Extension("added"); !present || value.String() != "new" {
		t.Fatalf("added extension = %v, %v", value, present)
	}
	if _, err := AddStringExtensions(event, map[string]string{"same": "different"}); !errors.Is(err, cloudevents.ErrMetadataCollision) {
		t.Fatalf("collision error = %v", err)
	}
	if _, err := AddStringExtensions(event, map[string]string{"bad": "\n"}); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
		t.Fatalf("invalid extension error = %v", err)
	}

	if value, present, err := StringExtension(event, "same"); err != nil || !present || value != "value" {
		t.Fatalf("string extension = %q, %v, %v", value, present, err)
	}
	if _, present, err := StringExtension(event, "missing"); err != nil || present {
		t.Fatalf("missing extension present = %v, error = %v", present, err)
	}
	for _, attribute := range []cloudevents.Attribute{cloudevents.NewBooleanAttribute(true), cloudevents.NewBinaryAttribute(nil)} {
		invalid := helperEvent(t, map[string]cloudevents.Attribute{"bad": attribute}, false)
		if _, _, err := StringExtension(invalid, "bad"); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
			t.Fatalf("invalid string extension error = %v", err)
		}
	}
	if !AttributesEqual(mustString(t, "x"), mustString(t, "x")) || AttributesEqual(mustString(t, "x"), mustString(t, "y")) || AttributesEqual(mustString(t, "x"), cloudevents.NewBooleanAttribute(true)) {
		t.Fatal("attribute equality mismatch")
	}
	if _, err := RebuildEvent(cloudevents.Event{}, nil); err == nil {
		t.Fatal("invalid rebuild succeeded")
	}
}

func TestSharedAdapterDataAndCopyMechanics(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		contentType string
		payload     []byte
		kind        cloudevents.DataKind
		wantErr     bool
	}{
		{"", nil, cloudevents.DataBinary, false},
		{"application/problem+json", []byte(`{}`), cloudevents.DataJSON, false},
		{"application/json", []byte("{"), 0, true},
		{"text/plain", []byte("text"), cloudevents.DataText, false},
		{"text/plain", []byte{0xff}, 0, true},
		{"application/octet-stream", []byte{0xff}, cloudevents.DataBinary, false},
		{";", nil, 0, true},
	} {
		data, err := DataFromPayload(test.contentType, test.payload)
		if (err != nil) != test.wantErr {
			t.Fatalf("DataFromPayload(%q) error = %v", test.contentType, err)
		}
		if err == nil && data.Kind() != test.kind {
			t.Fatalf("DataFromPayload(%q) kind = %v", test.contentType, data.Kind())
		}
	}
	if RestoreRetainedNil(cloudevents.NewBinaryData(nil), true) != nil {
		t.Fatal("retained nil became non-nil")
	}
	if value := RestoreRetainedNil(cloudevents.NewBinaryData([]byte{}), false); value == nil {
		t.Fatal("present empty became nil")
	}
	if CloneBytes(nil) != nil || CloneStrings(nil) != nil || CloneTimePointer(nil) != nil {
		t.Fatal("nil clone changed representation")
	}
	empty := CloneBytes([]byte{})
	if empty == nil {
		t.Fatal("empty byte clone became nil")
	}
	bytesValue := []byte("x")
	bytesClone := CloneBytes(bytesValue)
	bytesValue[0] = 'y'
	if string(bytesClone) != "x" {
		t.Fatal("byte clone aliases input")
	}
	stringsValue := map[string]string{"a": "b"}
	stringsClone := CloneStrings(stringsValue)
	stringsValue["a"] = "c"
	if stringsClone["a"] != "b" {
		t.Fatal("string clone aliases input")
	}
	now := time.Now()
	timeClone := CloneTimePointer(&now)
	if timeClone == &now || !timeClone.Equal(now) {
		t.Fatal("time clone mismatch")
	}
}

func TestSharedAdapterMappingAndLossReports(t *testing.T) {
	t.Parallel()

	extensions := map[string]cloudevents.Attribute{}
	if err := PutStringExtension(extensions, "value", "ok"); err != nil {
		t.Fatal(err)
	}
	if err := PutStringExtension(extensions, "bad", "\n"); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
		t.Fatalf("invalid put error = %v", err)
	}
	extensions["z"] = mustString(t, "last")
	extensions["a"] = mustString(t, "first")
	event := helperEvent(t, extensions, true)
	for _, test := range []struct {
		name, retained string
		trusted        bool
		want           string
		wantErr        error
	}{
		{"missing", "", false, "", nil},
		{"missing", "retained", false, "", cloudevents.ErrMetadataCollision},
		{"value", "other", true, "", cloudevents.ErrMetadataCollision},
		{"value", "", false, "", cloudevents.ErrUntrustedMetadata},
		{"value", "", true, "ok", nil},
		{"value", "ok", false, "ok", nil},
	} {
		value, err := MappedString(event, test.name, test.retained, test.trusted)
		if value != test.want || !errors.Is(err, test.wantErr) {
			t.Fatalf("MappedString(%q) = %q, %v", test.name, value, err)
		}
	}
	invalid := helperEvent(t, map[string]cloudevents.Attribute{"value": cloudevents.NewBooleanAttribute(true)}, false)
	if _, err := MappedString(invalid, "value", "", true); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
		t.Fatalf("invalid mapped value error = %v", err)
	}

	report := cloudevents.AdapterReport{}
	AppendOptionalContextLosses(event, &report, "target")
	AppendDataKindLoss(event, &report, "target", false)
	AppendExtensionLosses(event, &report, "target", map[string]struct{}{"value": {}})
	fields := make([]string, len(report.Losses))
	for index, loss := range report.Losses {
		fields[index] = loss.Field
	}
	if !reflect.DeepEqual(fields, []string{"dataschema", "subject", "time", "data.kind", "extensions.a", "extensions.z"}) {
		t.Fatalf("loss fields = %v", fields)
	}
	emptyReport := cloudevents.AdapterReport{}
	withoutData, err := cloudevents.NewEvent(cloudevents.Attributes{ID: "id", Source: "/source", Type: "type"}, cloudevents.Data{})
	if err != nil {
		t.Fatal(err)
	}
	AppendDataKindLoss(withoutData, &emptyReport, "target", true)
	if len(emptyReport.Losses) != 0 {
		t.Fatalf("absent data losses = %v", emptyReport.Losses)
	}
	for _, test := range []struct {
		contentType string
		data        cloudevents.Data
		wantLosses  int
	}{
		{"application/json", mustJSONData(t, `{}`), 0},
		{"text/plain", mustTextData(t, "text"), 0},
		{"application/octet-stream", cloudevents.NewBinaryData([]byte("x")), 0},
	} {
		candidate := eventWithData(t, test.contentType, test.data)
		candidateReport := cloudevents.AdapterReport{}
		AppendDataKindLoss(candidate, &candidateReport, "target", true)
		if len(candidateReport.Losses) != test.wantLosses {
			t.Fatalf("AppendDataKindLoss(%q) = %v", test.contentType, candidateReport.Losses)
		}
	}
	if !IsJSONContentType("application/problem+json") || IsJSONContentType("text/plain") || IsJSONContentType(";") {
		t.Fatal("JSON content-type classification mismatch")
	}
}

func eventWithData(t *testing.T, contentType string, data cloudevents.Data) cloudevents.Event {
	t.Helper()
	event, err := cloudevents.NewEvent(cloudevents.Attributes{ID: "id", Source: "/source", Type: "type", DataContentType: contentType}, data)
	if err != nil {
		t.Fatal(err)
	}
	return event
}

func mustJSONData(t *testing.T, value string) cloudevents.Data {
	t.Helper()
	data, err := cloudevents.NewJSONData([]byte(value))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func mustTextData(t *testing.T, value string) cloudevents.Data {
	t.Helper()
	data, err := cloudevents.NewTextData(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func helperEvent(t *testing.T, extensions map[string]cloudevents.Attribute, withTime bool) cloudevents.Event {
	t.Helper()
	now := time.Unix(1, 2)
	var occurredAt *time.Time
	if withTime {
		occurredAt = &now
	}
	data, err := cloudevents.NewJSONData([]byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	event, err := cloudevents.NewEvent(cloudevents.Attributes{ID: "id", Source: "/source", Type: "type", DataContentType: "application/json", DataSchema: "https://schemas.example/event", Subject: "subject", Time: occurredAt, Extensions: extensions}, data)
	if err != nil {
		t.Fatal(err)
	}
	return event
}

func mustString(t *testing.T, value string) cloudevents.Attribute {
	t.Helper()
	attribute, err := cloudevents.NewStringAttribute(value)
	if err != nil {
		t.Fatal(err)
	}
	return attribute
}
