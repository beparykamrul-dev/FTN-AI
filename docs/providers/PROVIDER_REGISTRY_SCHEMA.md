# FTN Provider Registry Schema

This document defines the normalized metadata required before a provider can be admitted to the FTN registry.

## Required fields

```yaml
name: example-provider
adapter_version: v1
upstream_identity: verified-upstream-id
protocols: []
capabilities: []
read_operations: []
mutation_operations: []
authentication: []
health_probe: true
latency_probe: true
compatibility: []
security_review: pending
license_review: pending
status: research
```

## Status lifecycle

```text
research -> reviewed -> tested -> admitted -> deprecated
```

A provider in `research` or `reviewed` status is never executable. `tested` means the contract suite passed; production use additionally requires security/license review and explicit registry admission.

## Mutation rule

Registry metadata is descriptive and does not grant mutation authority. Any mutation must use the FTN execution chain:

`Plan -> Audit -> Approval -> Guard -> Idempotency -> Executor`
