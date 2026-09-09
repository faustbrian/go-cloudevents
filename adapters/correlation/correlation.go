// Package correlation maps trusted Golib correlation values to and from
// CloudEvents extensions without hidden I/O or mutable aliases.
package correlation

import (
	"fmt"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	golibcorrelation "github.com/faustbrian/go-correlation"
)

var (
	// ErrInvalidInput reports a malformed event, extension, or identifier. It
	// aliases the root adapter error so callers can use errors.Is across module
	// boundaries.
	ErrInvalidInput = cloudevents.ErrInvalidAdapterInput
	// ErrMetadataCollision reports that Add would overwrite a non-equivalent
	// caller-owned extension. Existing metadata is never replaced.
	ErrMetadataCollision = cloudevents.ErrMetadataCollision
	// ErrUntrustedMetadata reports that inbound correlation extensions were
	// present but the caller had not established their trustworthiness.
	ErrUntrustedMetadata = cloudevents.ErrUntrustedMetadata
)

// Add returns a new event with the non-empty correlation identifiers. The
// input event and values remain caller-owned and unmodified. Equivalent
// existing extensions are accepted; conflicting or invalid extensions fail
// with ErrMetadataCollision or ErrInvalidInput.
func Add(event cloudevents.Event, values golibcorrelation.Values) (cloudevents.Event, error) {
	additions := map[string]string{}
	if values.CorrelationID != "" {
		additions[adapter.CorrelationIDExtension] = values.CorrelationID.String()
	}
	if values.RequestID != "" {
		additions[adapter.RequestIDExtension] = values.RequestID.String()
	}
	if values.CausationID != "" {
		additions[adapter.CausationIDExtension] = values.CausationID.String()
	}
	return adapter.AddStringExtensions(event, additions)
}

// Extract parses correlation extensions only after the caller marks the
// inbound metadata trusted. The zero Values result denotes an event with no
// correlation extensions. Policy remains caller-owned and controls identifier
// validation; present metadata is never returned when trusted is false.
func Extract(event cloudevents.Event, trusted bool, policy golibcorrelation.Policy) (golibcorrelation.Values, error) {
	correlationValue, hasCorrelation, err := adapter.StringExtension(event, adapter.CorrelationIDExtension)
	if err != nil {
		return golibcorrelation.Values{}, err
	}
	requestValue, hasRequest, err := adapter.StringExtension(event, adapter.RequestIDExtension)
	if err != nil {
		return golibcorrelation.Values{}, err
	}
	causationValue, hasCausation, err := adapter.StringExtension(event, adapter.CausationIDExtension)
	if err != nil {
		return golibcorrelation.Values{}, err
	}
	if !hasCorrelation && !hasRequest && !hasCausation {
		return golibcorrelation.Values{}, nil
	}
	if !trusted {
		return golibcorrelation.Values{}, ErrUntrustedMetadata
	}
	var values golibcorrelation.Values
	if hasCorrelation {
		values.CorrelationID, err = golibcorrelation.ParseCorrelationID(correlationValue, policy)
		if err != nil {
			return golibcorrelation.Values{}, fmt.Errorf("%w: correlationid: %w", ErrInvalidInput, err)
		}
	}
	if hasRequest {
		values.RequestID, err = golibcorrelation.ParseRequestID(requestValue, policy)
		if err != nil {
			return golibcorrelation.Values{}, fmt.Errorf("%w: requestid: %w", ErrInvalidInput, err)
		}
	}
	if hasCausation {
		values.CausationID, err = golibcorrelation.ParseCausationID(causationValue, policy)
		if err != nil {
			return golibcorrelation.Values{}, fmt.Errorf("%w: causationid: %w", ErrInvalidInput, err)
		}
	}
	return values, nil
}
