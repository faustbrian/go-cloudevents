package golib

import (
	"context"

	"github.com/faustbrian/go-cloudevents"
	cloudcorrelation "github.com/faustbrian/go-cloudevents/adapters/correlation"
	cloudtelemetry "github.com/faustbrian/go-cloudevents/adapters/telemetry"
	cloudtenancy "github.com/faustbrian/go-cloudevents/adapters/tenancy/v2"
	"github.com/faustbrian/go-correlation"
	telemetrypropagation "github.com/faustbrian/go-telemetry/propagation"
	"github.com/faustbrian/go-tenancy/v2"
)

var (
	// ErrInvalidAdapterInput preserves the target/root sentinel identity for errors.Is.
	ErrInvalidAdapterInput = cloudevents.ErrInvalidAdapterInput
	// ErrMetadataCollision preserves the target/root sentinel identity for errors.Is.
	ErrMetadataCollision = cloudevents.ErrMetadataCollision
	// ErrUntrustedMetadata preserves the target/root sentinel identity for errors.Is.
	ErrUntrustedMetadata = cloudevents.ErrUntrustedMetadata
)

// Loss is the root field/reason omission alias; custom text is not sanitized.
type Loss = cloudevents.AdapterLoss

// Report is the root caller-owned loss-list alias; no loss does not prove domain equivalence.
type Report = cloudevents.AdapterReport

// AddCorrelation delegates equivalent-value collision checks and copied correlation extensions; it does not authenticate identifiers.
func AddCorrelation(event cloudevents.Event, values correlation.Values) (cloudevents.Event, error) {
	return cloudcorrelation.Add(event, values)
}

// ExtractCorrelation delegates explicit trusted adoption under the caller validation policy; absent metadata returns zero values.
func ExtractCorrelation(event cloudevents.Event, trusted bool, policy correlation.Policy) (correlation.Values, error) {
	return cloudcorrelation.Extract(event, trusted, policy)
}

// AddTenant delegates validated tenant/v2 routing metadata addition, not tenant authorization.
func AddTenant(event cloudevents.Event, tenant tenancy.TenantID) (cloudevents.Event, error) {
	return cloudtenancy.Add(event, tenant)
}

// ExtractTenant delegates explicit trusted tenant/v2 parsing; absent metadata and failures return zero tenant plus error.
func ExtractTenant(event cloudevents.Event, trusted bool) (tenancy.TenantID, error) {
	return cloudtenancy.Extract(event, trusted)
}

// InjectTraceContext delegates caller-owned propagation to a new event and reports omitted baggage; SDK lifetime stays external.
func InjectTraceContext(ctx context.Context, event cloudevents.Event, policy *telemetrypropagation.Policy) (cloudevents.Event, Report, error) {
	return cloudtelemetry.InjectTraceContext(ctx, event, policy)
}

// ExtractTraceContext delegates trusted/untrusted filtering; nil context becomes Background and nil policy preserves it.
func ExtractTraceContext(ctx context.Context, event cloudevents.Event, policy *telemetrypropagation.Policy, trusted bool) context.Context {
	return cloudtelemetry.ExtractTraceContext(ctx, event, policy, trusted)
}
