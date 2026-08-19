# FTN AI Autonomous Release Gate

The autonomous service factory can build a service from authorized demand signals, but release is controlled by an explicit FTN release policy.

## Pipeline

```text
Demand Signals
  -> Opportunity Analysis
  -> Service Spec
  -> Build
  -> Tests
  -> Security Checks
  -> Dependency/License Checks
  -> Preview
  -> Release Gate
  -> FTN App Store / FTN Service Catalog
```

## App release checks

- build succeeds
- automated tests pass
- security checks pass
- dependency policy passes
- API compatibility verified
- permissions reviewed
- privacy/data classification verified
- version/changelog generated
- rollback artifact available
- release policy permits automatic publication

## FTN service release

Non-app services use the same gate but additionally verify infrastructure configuration, database migrations, observability, health checks and rollback readiness.

## Autonomous policy

FTN can configure selected low-risk service categories for automatic release. Higher-risk categories remain approval-gated. The AI must never bypass the configured release gate.

## App Store

The FTN App Store is the publication destination for approved FTN applications. The AI prepares the release package and metadata; the release controller enforces the configured publication policy.
