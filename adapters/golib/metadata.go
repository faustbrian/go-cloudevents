package golib

import (
	"context"

	"github.com/faustbrian/go-cloudevents"
	cloudcorrelation "github.com/faustbrian/go-cloudevents/adapters/correlation"
	cloudtelemetry "github.com/faustbrian/go-cloudevents/adapters/telemetry"
	cloudtenancy "github.com/faustbrian/go-cloudevents/adapters/tenancy"
	"github.com/faustbrian/go-correlation"
	telemetrypropagation "github.com/faustbrian/go-telemetry/propagation"
	"github.com/faustbrian/go-tenancy"
)

var (
	ErrInvalidAdapterInput = cloudevents.ErrInvalidAdapterInput
	ErrMetadataCollision   = cloudevents.ErrMetadataCollision
	ErrUntrustedMetadata   = cloudevents.ErrUntrustedMetadata
)

type Loss = cloudevents.AdapterLoss
type Report = cloudevents.AdapterReport

func AddCorrelation(event cloudevents.Event, values correlation.Values) (cloudevents.Event, error) {
	return cloudcorrelation.Add(event, values)
}

func ExtractCorrelation(event cloudevents.Event, trusted bool, policy correlation.Policy) (correlation.Values, error) {
	return cloudcorrelation.Extract(event, trusted, policy)
}

func AddTenant(event cloudevents.Event, tenant tenancy.TenantID) (cloudevents.Event, error) {
	return cloudtenancy.Add(event, tenant)
}

func ExtractTenant(event cloudevents.Event, trusted bool) (tenancy.TenantID, error) {
	return cloudtenancy.Extract(event, trusted)
}

func InjectTraceContext(ctx context.Context, event cloudevents.Event, policy *telemetrypropagation.Policy) (cloudevents.Event, Report, error) {
	return cloudtelemetry.InjectTraceContext(ctx, event, policy)
}

func ExtractTraceContext(ctx context.Context, event cloudevents.Event, policy *telemetrypropagation.Policy, trusted bool) context.Context {
	return cloudtelemetry.ExtractTraceContext(ctx, event, policy, trusted)
}
