# Source documentation audit

## Outcome, owner, consumers, risks and smallest checks

Base faa16c8a08c39cb237ad3c27fbc826537fd57573; detached fresh-main task worktree. Tier A comments/docs only.
Core maintainers own CloudEvents validated values, conversion loss and bounded
JSON/HTTP/Kafka mappings; independent adapter modules own explicit conversion
and trust/callback seams. Applications own authenticated identity, transport,
broker delivery, persistence, SDK/client/schema/propagation policy and lifecycle.
Direct consumers are root codecs, thirteen adapter modules (including the
released broad deprecated facade), shared internal adapter helpers and the
separate target-adapter integration test module. Original dirty dependency branch
and all other work remain untouched.

Audit all 96 tracked Go files/15 modules/16 catalog packages; integration is
pure tests, so fifteen actual source-package documentation surfaces. Named risks:
data absence/null/empty and abstract-type loss; duplicate metadata and protocol
precedence; byte/depth/count admission versus allocation; exact source/ID/type
identity versus authentication; copied retained bytes/maps and callback borrowing;
HTTP body ownership/cancellation, Kafka tombstones/partition keys, loss reports,
trusted metadata/schema/telemetry mappings and external operation windows.

Smallest checks: source-derived public/member/package comments and private
invariants, changed formatting, all-source AST/token/literal/implicit-semicolon
and directive preservation, parsed HTML/native package renders, native source/
local docs/offline spec, local links, unchanged structural API/dependencies/
frozen corpus/history, default pinned full/changed/staged secrets, ordinary
complete final-diff truth/usefulness review and explicit Intent commit. Request
the serialized Go slot before native compilation/rendering. Use three owned
caches and remove helpers/renders/caches on success/failure. No unchanged
examples/runtime/race/fuzz/bench/backend/conformance/interoperability/online/
consumer/release campaigns. Native .golib API mode is structural apidiff for
fourteen released modules; no full-prose snapshot is configured. The native structural producer is golib api check using the configured apidiff\nmodule baselines; verification/package.mk/scripts define conformance and interop,\nnot a full Go-doc snapshot generator. Preserve every baseline without execution;
never rewrite signatures or execute an unrelated generator for a structural gate.

## Missing versus investigated

Initially missing: name-led exported fields/constants and some helper/method
contracts. Existing core/package/API comments already protect useful protocol
semantics and must not be expanded mechanically. Investigation needed: pinned
CloudEvents 1.0.2 and selected later clarification scope; loss normalization,
parser/transport size/depth/copy limits and callback/privacy/cancel boundaries;
facade nominal identity versus intentionally separate adapter types.

No runtime or normative selected-expectation change is authorized. Retain stable
decision IDs, typed history prefixes, immutable source pins and frozen official/
SDK fixtures. Preserve useful Deprecated/generation/build/NOTE markers; do not
invent generator, standard completeness, exactly-once or durable delivery claims.

## Reproducible module/package and tracked-file inventory

Use git ls-files '*.go' '*go.mod' '*go.sum' '*AGENTS.md', modules.json/
packages.json, .golib.yaml, Makefile, verification/package.mk, scripts and source
package/declaration/test/fuzz/benchmark/example/marker scans. Every module's Go
minimum remains 1.27.0; native verification pins repository Go 1.27.2.

| Module | Public module path | Source packages |
| --- | --- | --- |
| . | github.com/faustbrian/go-cloudevents | ., internal/adapter |
| adapters/golib | github.com/faustbrian/go-cloudevents/adapters/golib/v4 | . |
| adapters/audit | github.com/faustbrian/go-cloudevents/adapters/audit/v2 | . |
| adapters/correlation | github.com/faustbrian/go-cloudevents/adapters/correlation | . |
| adapters/event-sourcing | github.com/faustbrian/go-cloudevents/adapters/event-sourcing/v2 | . |
| adapters/jsonschema | github.com/faustbrian/go-cloudevents/adapters/jsonschema/v2 | . |
| adapters/kafka | github.com/faustbrian/go-cloudevents/adapters/kafka | . |
| adapters/outbox | github.com/faustbrian/go-cloudevents/adapters/outbox/v2 | . |
| adapters/queue | github.com/faustbrian/go-cloudevents/adapters/queue/v2 | . |
| adapters/rabbitstream | github.com/faustbrian/go-cloudevents/adapters/rabbitstream | . |
| adapters/schema-registry | github.com/faustbrian/go-cloudevents/adapters/schema-registry/v3 | . |
| adapters/telemetry | github.com/faustbrian/go-cloudevents/adapters/telemetry | . |
| adapters/tenancy | github.com/faustbrian/go-cloudevents/adapters/tenancy/v2 | . |
| adapters/workflow | github.com/faustbrian/go-cloudevents/adapters/workflow/v2 | . |
| integration/target-adapters | github.com/faustbrian/go-cloudevents/integration/target-adapters | . |

| Tracked Go file | Package | Audit role/evidence entrypoints |
| --- | --- | --- |
| `adapters/audit/audit.go` | audit | Source API/private invariant |
| `adapters/audit/audit_test.go` | audit_test | Existing evidence: TestAuditMetadataRoundTripSelectsSafeFieldsAndReportsLoss, TestAuditMetadataRoundTripAllowsAbsentContextIdentifiers, TestAuditAddRejectsInvalidRecordTenantAndCollision, TestAuditExtractRejectsUntrustedIncompleteAndMalformedMetadata |
| `adapters/audit/published_v2_test.go` | audit_test | Existing evidence: TestPublishedV2TenantIdentityRoundTrip |
| `adapters/correlation/correlation.go` | correlation | Source API/private invariant |
| `adapters/correlation/correlation_test.go` | correlation_test | Existing evidence: TestCorrelationRoundTripPreservesAllIdentifiers, TestCorrelationAbsentMetadataReturnsZeroWithoutTrust, TestCorrelationRejectsUntrustedCollidingAndMalformedMetadata |
| `adapters/event-sourcing/event_sourcing.go` | eventsourcing | Source API/private invariant |
| `adapters/event-sourcing/event_sourcing_test.go` | eventsourcing | Existing evidence: TestRoundTripRetainsEventStoreStateWithoutAliasing, TestToCloudEventRejectsInvalidInputs, TestFromCloudEventRejectsInvalidStateAndMetadataCollisions, TestFromCloudEventValidatesPayloadTenantAndMessageIdentity, TestUntrustedPortableMetadataCanPopulateMissingRetainedFields, TestRoundTripAllowsAbsentOptionalMetadata |
| `adapters/event-sourcing/published_v2_test.go` | eventsourcing | Existing evidence: TestPublishedV2TenantErrorIdentity |
| `adapters/golib/audit.go` | golib | Source API/private invariant |
| `adapters/golib/audit_schema_test.go` | golib_test | Existing evidence: TestAuditMetadataAdapterSelectsSafeFieldsAndRequiresTrust, TestAuditMetadataRoundTripAllowsAbsentTenant, TestSchemaAdaptersAreExplicitAndRegistryLookupIsCallerMapped |
| `adapters/golib/benchmark_test.go` | golib_test | Existing evidence: BenchmarkKafkaBinaryRoundTrip, BenchmarkRabbitStreamStructuredRoundTrip, BenchmarkQueueRoundTrip |
| `adapters/golib/doc.go` | golib | Source API/private invariant |
| `adapters/golib/envelope_boundaries_test.go` | golib_test | Existing evidence: TestEventSourcingAdapterRejectsEveryUnrepresentableBoundary, TestOutboxAdapterRejectsEveryInvalidBoundary, TestOutboxAdapterPreservesNilPayload, TestOutboxAdapterPreservesNonNilEmptyPayload, TestQueueAdapterRejectsEveryInvalidBoundary, TestQueueAdapterPreservesNilBody, TestQueueAdapterPreservesNonNilEmptyBody, TestCloudEventToQueueRejectsRetainedContentTypeCollision, TestCloudEventToQueueRejectsDeletedRetainedPortableExtensions, TestWorkflowAdapterRejectsEveryInvalidBoundary, TestWorkflowAdapterPreservesNilData, TestWorkflowAdapterPreservesNonNilEmptyData, TestWorkflowAdapterRejectsNonCanonicalEventTypes, TestWorkflowAdapterReportsEveryDiscardedPayloadDeclaration, TestEnvelopeAdaptersReportEveryPayloadKindTheyCannotReconstruct |
| `adapters/golib/envelopes.go` | golib | Source API/private invariant |
| `adapters/golib/envelopes_test.go` | golib_test | Existing evidence: TestEventSourcingMappingPreservesCanonicalEnvelopeOutOfBand, TestGolibOutboxAndQueueMappingsRetainUnrepresentableState, TestWorkflowMappingRequiresCallerOwnedStableID |
| `adapters/golib/example_test.go` | golib_test | Existing evidence: ExampleQueueToCloudEvent |
| `adapters/golib/fuzz_test.go` | golib_test | Existing evidence: FuzzDecodeKafka, FuzzQueuePayloadRoundTrip, FuzzEventSourcingRoundTrip, FuzzOutboxRoundTrip, FuzzWorkflowRoundTrip, FuzzAuditCorrelationTenantMetadata, FuzzDirectJSONSchemaValidation |
| `adapters/golib/kafka.go` | golib | Source API/private invariant |
| `adapters/golib/kafka_test.go` | golib_test | Existing evidence: TestKafkaAdapterPreservesTransportOwnershipAndBindingRoundTrip, TestKafkaAdapterRejectsTransportCollisions |
| `adapters/golib/main_test.go` | golib_test | Existing evidence: TestMain |
| `adapters/golib/metadata.go` | golib | Source API/private invariant |
| `adapters/golib/metadata_test.go` | golib_test | Existing evidence: TestCorrelationAndTenantMetadataRequireExplicitTrustAndRejectCollisions, TestTelemetryAdapterUsesExplicitGolibPropagationPolicy |
| `adapters/golib/public_boundaries_test.go` | golib_test | Existing evidence: TestMetadataAdapterRejectsEveryMalformedOrUntrustedBoundary, TestTelemetryAdapterCoversNilCollisionBaggageAndTrustedExtraction, TestAuditMetadataRejectsIncompleteMalformedAndCollidingInput, TestSchemaValidatorsRejectEveryConfigurationResolutionAndPayloadFailure, TestKafkaAdapterPropagatesBindingDecodeFailuresAndAcceptsEqualKeys, TestKafkaAdapterRejectsOversizedMetadataBeforeCopyingIt |
| `adapters/golib/public_v2_test.go` | golib_test | Existing evidence: TestPublicWorkflowV2FacadeComposition |
| `adapters/golib/public_v3_test.go` | golib_test | Existing evidence: TestPublicTenantV2FacadeRoundTrip |
| `adapters/golib/public_v4_test.go` | golib_test | Existing evidence:  |
| `adapters/golib/rabbitstream.go` | golib | Source API/private invariant |
| `adapters/golib/rabbitstream_test.go` | golib_test | Existing evidence: TestRabbitStreamStructuredAdapterPreservesEventAndTransportState, TestRabbitStreamStructuredAdapterRejectsAmbiguousMappings, TestRabbitStreamStructuredAdapterSupportsUnkeyedTimelessEvents, TestRabbitStreamStructuredAdapterBoundsTransportPropertiesBeforeCopying, TestRabbitStreamStructuredAdapterBoundsTransportPropertiesBeforeEncoding |
| `adapters/golib/schema.go` | golib | Source API/private invariant |
| `adapters/golib/schema_context_test.go` | golib_test | Existing evidence: TestSchemaValidatorsRejectNilContext |
| `adapters/golib/schema_hardening_test.go` | golib_test | Existing evidence: TestRegistryJSONSchemaValidatorRejectsUnmappedEventSchemaWithoutResolution, TestRegistryJSONSchemaValidatorRequiresBoundedExplicitConfiguration, TestRegistryJSONSchemaValidatorUsesCacheIdentityValidation, TestRegistryJSONSchemaValidatorRejectsResolvedNonJSONSchema, TestRegistryJSONSchemaValidatorBoundsResolutionAndPreservesParentCancellation, TestRegistryJSONSchemaValidatorPropagatesCancellationIntoPayloadValidation, TestRegistryJSONSchemaValidatorOwnsImmutableSchemaLookupSnapshot, TestRegistryJSONSchemaValidatorConcurrentValidationDoesNotShareSourceMap |
| `adapters/golib/stress_test.go` | golib_test | Existing evidence: TestConcurrentMetadataConversionsAreIsolated, TestQueueConversionSoakPreservesPayloadAndRetainedState |
| `adapters/jsonschema/jsonschema.go` | jsonschema | Source API/private invariant |
| `adapters/jsonschema/jsonschema_test.go` | jsonschema_test | Existing evidence: TestValidatorEnforcesExplicitSchemaMapping, TestValidatorRequiresContextAndJSONContent, TestValidatorReportsPayloadOutcomes, TestValidatorPreservesV2UnicodeAndCancellationBoundaries |
| `adapters/kafka/kafka.go` | kafka | Source API/private invariant |
| `adapters/kafka/kafka_test.go` | kafka_test | Existing evidence: TestKafkaAdapterPreservesTransportOwnershipAndBindingRoundTrip, TestKafkaAdapterRejectsTransportCollisions, TestKafkaAdapterPropagatesBindingFailuresWithoutTransportState |
| `adapters/outbox/outbox.go` | outbox | Source API/private invariant |
| `adapters/outbox/outbox_test.go` | outbox | Existing evidence: TestRoundTripRetainsOwnedStateWithoutAliasing, TestConversionRejectsInvalidBoundariesAndIdentifierCollision, TestRoundTripPreservesNilAndNonNilEmptyPayloads |
| `adapters/queue/published_v2_test.go` | queue | Existing evidence: TestPublishedV2TenantErrorIdentity |
| `adapters/queue/queue.go` | queue | Source API/private invariant |
| `adapters/queue/queue_test.go` | queue | Existing evidence: TestQueueRoundTripPreservesOwnedStateWithoutAliasing, TestQueueConversionRejectsInvalidInputsAndCollisions, TestQueueExtensionValidationAndLossReporting, TestQueuePreservesNilAndEmptyBodies, TestQueuePrivateHelpersCoverZeroAndCopiedState |
| `adapters/rabbitstream/rabbitstream.go` | rabbitstream | Source API/private invariant |
| `adapters/rabbitstream/rabbitstream_test.go` | rabbitstream_test | Existing evidence: TestRabbitStreamStructuredAdapterPreservesEventAndTransportState, TestRabbitStreamStructuredAdapterRejectsAmbiguousMappings, TestRabbitStreamStructuredAdapterSupportsUnkeyedTimelessEvents, TestRabbitStreamStructuredAdapterBoundsTransportPropertiesBeforeCopying, TestRabbitStreamStructuredAdapterBoundsTransportPropertiesBeforeEncoding |
| `adapters/schema-registry/schema_registry.go` | schemaregistry | Source API/private invariant |
| `adapters/schema-registry/schema_registry_test.go` | schemaregistry_test | Existing evidence: TestNewJSONSchemaValidatorRejectsUnsafeConfiguration, TestJSONSchemaValidatorRejectsUntrustedInputsBeforeResolution, TestJSONSchemaValidatorPreservesResolutionAndFormatFailures, TestJSONSchemaValidatorReportsPayloadOutcomes, TestJSONSchemaValidatorOwnsLookupSnapshot, TestJSONSchemaValidatorPreservesV3BoundsAndPrivatePayloadDiagnostics |
| `adapters/telemetry/telemetry.go` | telemetry | Source API/private invariant |
| `adapters/telemetry/telemetry_test.go` | telemetry | Existing evidence: TestInjectTraceContextPreservesOwnershipAndReportsLoss, TestExtractTraceContextAppliesCallerTrustDecision, TestMetadataCarrierImplementsDeterministicTextMapCarrier |
| `adapters/tenancy/published_v2_test.go` | tenancy | Existing evidence: TestPublishedV2TenantIdentityRoundTrip |
| `adapters/tenancy/tenancy.go` | tenancy | Source API/private invariant |
| `adapters/tenancy/tenancy_test.go` | tenancy | Existing evidence: TestAddValidatesTenantPreservesInputAndRejectsCollisions, TestExtractRequiresPresentTrustedValidStringMetadata |
| `adapters/workflow/public_v2_test.go` | workflow_test | Existing evidence: TestPublicWorkflowV2Composition |
| `adapters/workflow/workflow.go` | workflow | Source API/private invariant |
| `adapters/workflow/workflow_test.go` | workflow | Existing evidence: TestWorkflowRoundTripPreservesOwnedStateWithoutAliasing, TestWorkflowRejectsInvalidBoundaries, TestWorkflowPreservesNilAndEmptyPayloads, TestWorkflowReportsEveryDiscardedCloudEventField |
| `attribute.go` | cloudevents | Source API/private invariant |
| `benchmark_test.go` | cloudevents | Existing evidence: BenchmarkJSONEquivalentEvent |
| `boundary_test.go` | cloudevents | Existing evidence: TestEventValidationRejectsCorruptedInternalDataStates, TestAttributeValidationCoversEveryAbstractTypeAndInvariant, TestAttributeNamesAndCloudStringsEnforceExactCharacterBoundaries, TestValidationReportsSelectedExtensionAndOptionalAttributeFailures, TestTraceStateRejectsWrongAbstractTypeWithValidSyntax, TestTracingExtensionsRejectEveryGrammarBoundary, TestTraceStateAcceptsExactNormativeBoundaries, TestValidateSchemaRejectsEveryPreconditionAndReturnsValidatorError |
| `conformance_test.go` | cloudevents | Existing evidence: TestPinnedOfficialConformanceFixtures, TestOfficialConformanceHTTPBindingScenarios, TestOfficialConformanceKafkaBindingScenarios |
| `conversion.go` | cloudevents | Source API/private invariant |
| `doc.go` | cloudevents | Source API/private invariant |
| `errors.go` | cloudevents | Source API/private invariant |
| `event.go` | cloudevents | Source API/private invariant |
| `event_test.go` | cloudevents | Existing evidence: TestNewEventOwnsInputsAndPreservesContext, TestNewEventReportsCanonicalValidationIssues, TestDataKindsPreserveAbsentNullEmptyAndBinaryValues, TestAttributeTypesEnforceCloudEventsTypeSystem, TestEventsRejectUnassignedExtensionValuesAndZeroValueEncoding |
| `extensions.go` | cloudevents | Source API/private invariant |
| `extensions_test.go` | cloudevents | Existing evidence: TestSelectedExtensionsEnforceNormativeSemantics |
| `fuzz_test.go` | cloudevents | Existing evidence: FuzzDecodeJSON, FuzzDecodeJSONBatch, FuzzDecodeHTTP, FuzzDecodeKafka |
| `http.go` | cloudevents | Source API/private invariant |
| `http_boundary_test.go` | cloudevents | Existing evidence: TestEncodeHTTPRejectsUnsupportedCountsAndMapsOptionalBinaryFields, TestHTTPAttributeEncodingPreservesAndEscapesExactOctetBoundaries, TestDecodeHTTPBoundsContentTypeAndRejectsDuplicateHeadersWithoutCopying, TestDecodeHTTPIgnoresOversizedUnownedHeaderNamesWithoutAllocating, TestHTTPPreflightDistinguishesExactHeaderBoundaries, TestDecodeHTTPRejectsModeBodyAndStructuredMetadataFailures, TestStructuredHTTPMetadataCoversOptionalFieldsAndDirectLimits, TestHTTPBodyLimitsCancellationAndNilOwnership, TestDecodeHTTPBinaryRejectsMalformedAttributesAndData, TestHTTPAttributeQuotedStringGrammar, TestDecodeHTTPBoundsContextAttributeNamesInEveryMode |
| `http_test.go` | cloudevents | Existing evidence: TestEncodeHTTPBinaryMapsAttributesAndExactData, TestDecodeHTTPBinaryDecodesPercentEncodingAndJSONData, TestHTTPStructuredAndBatchModesRoundTrip, TestDecodeHTTPRejectsDuplicateAndConflictingMetadata |
| `integration/target-adapters/audit_schema_test.go` | targetadapters_test | Existing evidence: TestAuditMetadataAdapterSelectsSafeFieldsAndRequiresTrust, TestAuditMetadataRoundTripAllowsAbsentTenant, TestSchemaAdaptersAreExplicitAndRegistryLookupIsCallerMapped |
| `integration/target-adapters/envelope_boundaries_test.go` | targetadapters_test | Existing evidence: TestEventSourcingAdapterRejectsEveryUnrepresentableBoundary, TestOutboxAdapterRejectsEveryInvalidBoundary, TestOutboxAdapterPreservesNilPayload, TestOutboxAdapterPreservesNonNilEmptyPayload, TestQueueAdapterRejectsEveryInvalidBoundary, TestQueueAdapterPreservesNilBody, TestQueueAdapterPreservesNonNilEmptyBody, TestCloudEventToQueueRejectsRetainedContentTypeCollision, TestCloudEventToQueueRejectsDeletedRetainedPortableExtensions, TestWorkflowAdapterRejectsEveryInvalidBoundary, TestWorkflowAdapterPreservesNilData, TestWorkflowAdapterPreservesNonNilEmptyData, TestWorkflowAdapterRejectsNonCanonicalEventTypes, TestWorkflowAdapterReportsEveryDiscardedPayloadDeclaration, TestEnvelopeAdaptersReportEveryPayloadKindTheyCannotReconstruct |
| `integration/target-adapters/envelopes_test.go` | targetadapters_test | Existing evidence: TestEventSourcingMappingPreservesCanonicalEnvelopeOutOfBand, TestGolibOutboxAndQueueMappingsRetainUnrepresentableState, TestWorkflowMappingRequiresCallerOwnedStableID |
| `integration/target-adapters/fuzz_test.go` | targetadapters_test | Existing evidence: FuzzDecodeKafka, FuzzQueuePayloadRoundTrip, FuzzEventSourcingRoundTrip, FuzzOutboxRoundTrip, FuzzWorkflowRoundTrip, FuzzAuditCorrelationTenantMetadata, FuzzDirectJSONSchemaValidation |
| `integration/target-adapters/kafka_schema_cloudevents_example_test.go` | targetadapters_test | Existing evidence: Example_kafkaSchemaCloudEventFlow, TestKafkaSchemaCloudEventRecipeRejectsInvalidEventBeforePublish, TestKafkaSchemaCloudEventRecipeDoesNotAcknowledgeHandlerFailure, TestKafkaSchemaCloudEventRecipeRejectsInvalidConsumedEventBeforeHandling, TestKafkaSchemaCloudEventRecipePreservesCommitFailureWithoutClaimingSettlement, TestKafkaSchemaCloudEventRecipeJoinsOperationAndOrderedShutdownFailures |
| `integration/target-adapters/metadata_test.go` | targetadapters_test | Existing evidence: TestCorrelationAndTenantMetadataRequireExplicitTrustAndRejectCollisions, TestTelemetryAdapterUsesExplicitGolibPropagationPolicy |
| `integration/target-adapters/public_boundaries_test.go` | targetadapters_test | Existing evidence: TestMetadataAdapterRejectsEveryMalformedOrUntrustedBoundary, TestTelemetryAdapterCoversNilCollisionBaggageAndTrustedExtraction, TestAuditMetadataRejectsIncompleteMalformedAndCollidingInput, TestSchemaValidatorsRejectEveryConfigurationResolutionAndPayloadFailure, TestKafkaAdapterPropagatesBindingDecodeFailuresAndAcceptsEqualKeys, TestKafkaAdapterRejectsOversizedMetadataBeforeCopyingIt |
| `integration/target-adapters/schema_context_test.go` | targetadapters_test | Existing evidence: TestSchemaValidatorsRejectNilContext |
| `integration/target-adapters/schema_hardening_test.go` | targetadapters_test | Existing evidence: TestRegistryJSONSchemaValidatorRejectsUnmappedEventSchemaWithoutResolution, TestRegistryJSONSchemaValidatorRequiresBoundedExplicitConfiguration, TestRegistryJSONSchemaValidatorUsesCacheIdentityValidation, TestRegistryJSONSchemaValidatorRejectsResolvedNonJSONSchema, TestRegistryJSONSchemaValidatorBoundsResolutionAndPreservesParentCancellation, TestRegistryJSONSchemaValidatorPropagatesCancellationIntoPayloadValidation, TestRegistryJSONSchemaValidatorOwnsImmutableSchemaLookupSnapshot, TestRegistryJSONSchemaValidatorConcurrentValidationDoesNotShareSourceMap |
| `integration/target-adapters/stress_test.go` | targetadapters_test | Existing evidence: TestConcurrentMetadataConversionsAreIsolated, TestQueueConversionSoakPreservesPayloadAndRetainedState |
| `integration/target-adapters/target_adapters_test.go` | targetadapters_test | Existing evidence: TestTargetAdaptersExposeIndependentIntegrationBoundaries |
| `internal/adapter/adapter.go` | adapter | Source API/private invariant |
| `internal/adapter/adapter_test.go` | adapter | Existing evidence: TestCanonicalAdapterErrorsPreserveLegacyStrings, TestSharedAdapterMechanicsPreserveValuesAndClassifyFailures, TestSharedAdapterDataAndCopyMechanics, TestSharedAdapterMappingAndLossReports |
| `interoperability_test.go` | cloudevents | Existing evidence: TestOfficialGoSDKJSONInteroperabilityInBothDirections, TestOfficialGoSDKJSONBatchInteroperabilityInBothDirections, TestOfficialGoSDKEdgeDataInteroperability, TestOfficialGoSDKKafkaInteroperabilityInBothDirections, TestOfficialGoSDKHTTPInteroperabilityInBothDirections |
| `javascript_interop_test.go` | cloudevents | Existing evidence: TestJavaScriptSDKJSONAndBatchInteroperability, TestJavaScriptSDKEdgeDataInteroperability, TestJavaScriptSDKBinaryTransportInteroperability, TestJavaScriptSDKProducedHTTPInteroperability, TestJavaScriptSDKProducedKafkaInteroperability |
| `json.go` | cloudevents | Source API/private invariant |
| `json_boundary_test.go` | cloudevents | Existing evidence: TestDecodeJSONEnforcesDecodedAndExactResourceBoundaries, TestDecodeJSONRejectsOversizedStringExtensionWithinAllocationBudget, TestDecodeJSONBatchAndDepthAcceptExactLimits, TestEncodeJSONCoversOptionalAttributesAndEveryDataRepresentation, TestDecodeJSONRejectsEveryContextAndExtensionBoundary, TestDecodeJSONRejectsEveryDataAndBatchBoundary, TestJSONInspectionRejectsMalformedAndTrailingValues |
| `json_test.go` | cloudevents | Existing evidence: TestEncodeJSONIsDeterministicAndNormative, TestDecodeJSONRejectsAmbiguousAndOverLimitInput, TestDecodeJSONPreservesUnknownExtensionsAndDataPresence, TestJSONBatchRoundTripAndLimits |
| `kafka.go` | cloudevents | Source API/private invariant |
| `kafka_boundary_test.go` | cloudevents | Existing evidence: TestEncodeKafkaCoversModesOptionalHeadersAndPartitionMapping, TestDecodeKafkaRejectsRecordAndContentTypeBoundaries, TestStructuredKafkaMetadataCoversOptionalAndMalformedHeaders, TestDecodeKafkaBinaryRejectsEveryMalformedHeaderAndDataBoundary, TestKafkaAcceptsExactRecordHeaderAndDataLimits, TestDecodeKafkaBoundsAllCopiedRecordMetadata |
| `kafka_test.go` | cloudevents | Existing evidence: TestEncodeKafkaBinaryPreservesKeyAndTombstoneSemantics, TestDecodeKafkaOfficialBinaryAndStructuredScenarios, TestDecodeKafkaRejectsConflictingStructuredMetadataAndMalformedRecords, TestDecodeKafkaRejectsMissingStructuredValueAsMalformed |
| `limits.go` | cloudevents | Source API/private invariant |
| `lossless_boundary_test.go` | cloudevents | Existing evidence: TestStructuredFormatsPreserveDeclaredJSONPayloadBytes, TestDecodeJSONBoundsUnknownExtensionValues, TestStrictEncodersRejectEveryImplicitConversionLoss, TestLossAwareEncodersRemainLosslessWhenNoReportIsNeeded, TestStrictEncodersReportAmbiguousAndConflictingDataKinds |
| `main_test.go` | cloudevents | Existing evidence: TestMain |
| `schema.go` | cloudevents | Source API/private invariant |
| `schema_test.go` | cloudevents | Existing evidence: TestValidateSchemaIsExplicitAndOwnsValidatorInput, TestValidateSchemaRejectsNilContextPrecisely, TestValidateSchemaRejectsTypedNilValidatorWithoutPanicking |
| `stress_test.go` | cloudevents | Existing evidence: TestConcurrentEncodingAndDecodingIsIsolated, TestCodecSoakRoundTripsStableBytes, TestHTTPDecodePropagatesBodyFailureWithoutClosingCallerResource, TestHTTPDecodeConsumesChunkedShortReads, TestOversizedKafkaKeyIsRejectedWithoutAllocation |
| `validation.go` | cloudevents | Source API/private invariant |

## Investigation and verification disposition

Inventory created before source edits. All 32 production Go sources and all
96 tracked-file declarations/markers/test/fuzz/benchmark/example boundaries were
inspected. No generated, build/platform/tag, TODO/FIXME/HACK/NOTE/WARNING or cgo
path is present; the facade's valid Deprecated header remains unchanged.
Fuzz boundaries are four root decoders, seven deprecated-facade mappings and
seven target-integration mappings. Existing detailed round-trip/trust comments
remain; parser-only targets explicitly do not claim semantic conformance.
The root SDK-comparison benchmark and three facade benchmarks retain their
actual timed operations, with no new performance claim. Examples/bodies/fixtures
and required-service-free metadata remain unchanged and unexecuted locally.

Pinned authorities are CloudEvents 1.0.2 core/JSON/HTTP/Kafka/tracing/partitioning
and selected eight later clarification commits, conformance v0.4.1, SDK Go
v2.16.2/JavaScript v10.0.0 and Node 24.19.0 from the existing specification matrix,
manifest and monitoring register. All twelve stable DEC IDs, normative
requirements, original typed history and official/SDK source/corpus pins remain
frozen. Adjacent code comments link DEC-002 naming, DEC-005 duplicates, DEC-006
URI serialization, DEC-007 determinism, DEC-008/009 HTTP precedence/body ownership,
DEC-010 Kafka tombstones and DEC-012 tracing to their existing protecting tests.
No new standards, authentication, exact-once, persistence or canonical-signing
qualification is claimed.

Investigated source contracts: Event/Data/Attribute retain owned immutable state;
constructor/encoder paths have no Limits. Decoders perform selected input/depth/
retention admission, while token/raw-message work precedes context/batch count
rejection. Optional null context is absent and duplicate envelope context is
rejected, not duplicate payload keys. JSON lexical interior survives structured
formats, but external whitespace/abstract extension-kind/implicit content-type
changes require strict rejection or a report. HTTP borrows header/reader state,
never closes readers and cannot preempt blocked Read; binary body presence differs
from Kafka nil tombstones. Kafka keys and unowned headers copy without broker
delivery/settlement. Adapter retained state is not persistence or authentication.

Most target API comments already described their exact copy/trust/collision
semantics and were retained. Added only entry/lifecycle constraints and repaired
two false current comments: Rabbit Stream State.Timestamp retains transport
Message.Timestamp on Decode, and registry Timeout is one shared cooperative
resolution/decode context. Source pointers/maps/reports must remain stable during
copy/callback use. Direct propagation/schema callbacks are trusted, may block/
panic, and their raw errors are not universally redacted. Event context identifiers
and report field names are not secret-safe merely because payload values are
omitted. No erasure or secure-memory promise is made.

Investigated/unresolved follow-up: DEC-003 maintained prose calls received base64
canonical padded text, while source uses encoding/base64.Strict decoding.
Go 1.27.2 encoding/base64.Strict explicitly still ignores CR/LF; source decoding
therefore does not enforce byte-for-byte canonical received text. This confirmed
current-prose/acceptance mismatch is owned by CloudEvents maintainers at DEC-003,
not repaired or hidden by rewriting the desired selected expectation. The
slice-based batch decoder also accepts top-level JSON null as an empty batch;
no array-only admission repair or normative expectation rewrite is made. Reason
for structured HTTP lacking a separate context-count check is not inferred. New
constructor caps, hard callback deadlines or standards policy would require a
separate behavioral decision; these are accurately documented limits, not
claimed repairs.

## Executed checks and remaining boundary

Go 1.27.2 formatted 36 owned Go files. All 96 files passed exact executable AST,
noncomment token/literal/implicit-semicolon preservation, with no waiver. All
424 parsed comment groups rendered as HTML; 363 public declarations/fields/
interface methods and fifteen canonical source-package headers passed. Retained
Deprecated/directive bytes, all fifteen module pairs, fourteen current structural
API baselines, work/catalog files and specification authorities/history remained
exact. Fifteen native go doc -all source surfaces rendered; the test-only
integration package was parsed, not presented as a new library documentation API.

The first bounded session completed that proof, then stopped at an overly exact
stdlib-comment spelling assertion. It cleaned all resources; the continuation
reused unchanged successful source proof and ran only remaining stdlib source
inspection, native typed digest/history verification, renders and local gates.
The pinned Strict source states CR/LF are ignored, and encoding/json's null
literal branch zeroes slices, supporting the actual accepted-empty outer-null
batch path without an executable behavior probe. Native twelve typed DEC digests
are unchanged and each current digest remains in its original history. Native
source-docs reports zero findings, docs check --local passes, and offline
specification check validates one module/twelve decisions. These are structural
and documentation results, not resolution of the normative/runtime follow-ups.

Default pinned complete-owned and changed/added secret scans passed; staged
content receives the same final scan before explicit Intent commit. Relative
links/retained markers and untouched tracked inputs are checked separately, and
one ordinary complete final-diff truth/usefulness inspection covers all owned
paths. All task-owned GOCACHE/GOMODCACHE/GOTMPDIR/helper/render resources are removed
on failure and success; original dirty checkout and parent tools remain intact.

No executable example/body/dependency/fixture/authority was changed. Unchanged
examples/runtime/test/race/fuzz/benchmark/backend/official-conformance/SDK-interop/
online/consumer/release campaigns were not run. Fourteen structural apidiff
snapshots need no regeneration because declarations are exact. Local committed
source does not establish hosted CI, release or standards qualification.

Remaining questions: maintainers should reconcile DEC-003 canonical received
base64 expectations with Strict's CR/LF acceptance and decide array-only JSON
batch admission versus current null acceptance in a separately scoped behavior
change. Constructor/encoder caps, structured-header admission and hard callback
deadlines remain caller obligations/possible future policy, not newly promised
guarantees. No historical/signed evidence or desired normative behavior was
silently rewritten.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.
