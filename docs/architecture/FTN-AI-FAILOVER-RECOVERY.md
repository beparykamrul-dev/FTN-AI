# FTN AI Failover and Recovery

FTN AI remains available when an individual approved layer or service becomes unhealthy.

## Recovery path

```text
Request
  ↓
Preferred FTN Native Layer
  ↓ unhealthy
Approved fallback layer
  ↓ unavailable
Service-safe response / retry
```

Recovery controls:

- health-aware routing
- bounded retries
- timeout limits
- circuit breaking
- last-known-good configuration
- audit events for fallback
- operator-visible degraded state

A failed AI layer must not cause unsafe or unauthorized side effects. State-changing operations remain subject to FTN IAM and approval policy.
