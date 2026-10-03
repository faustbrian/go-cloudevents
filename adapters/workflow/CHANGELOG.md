# Changelog

## Unreleased

- Adopt Workflow v2 through the independently versioned
  `github.com/faustbrian/go-cloudevents/adapters/workflow/v2` module. Update
  adapter and workflow imports together because history events and retained
  definition references now use the workflow-v2 nominal types. Conversion,
  ownership, loss reporting and the CloudEvents root contract are unchanged.

## 1.0.0 - 2026-09-09

- Publish the target-oriented workflow adapter with retained durable state,
  nil-data preservation, canonical types, and explicit losses.
