# FTN Source Wrapper / Middleware Contract

## Purpose

Source wrappers provide a stable boundary between external discovery/input systems and FTN provider adapters. They normalize metadata without coupling the control plane to a specific source implementation.

## Pipeline

```text
Source
  -> Wrapper
  -> Validation
  -> Normalization
  -> Provider Adapter
  -> Snapshot
  -> Consistency Engine
```

## Required properties

- deterministic normalization
- context-aware cancellation
- explicit source identity
- schema/version validation
- bounded input handling
- structured error classification
- no provider mutation during import
- trace/audit correlation ID support

## Middleware rules

Middleware may enrich, filter, validate, or normalize an input, but must not silently change authoritative DNS state. Any mutation belongs to the approved reconciliation execution path.

## Compatibility

Wrappers should expose capability metadata and tolerate provider-specific fields by preserving them as extensions where practical. Breaking schema changes require an explicit version transition.
