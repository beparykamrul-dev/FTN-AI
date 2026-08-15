# FTN Security Batch V1

## 1. Certificate Health Monitor

Continuously tracks certificate expiry, issuer, deployment state and renewal status. Alerts before expiry and preserves the last known-good certificate on renewal failure.

## 2. Edge Abuse Protection

Adds an early, inexpensive protection layer for public FTN endpoints: rate limits, connection limits, request-size limits, timeouts and anomaly signals. It is not a substitute for a full WAF/DDoS provider, but establishes FTN-owned baseline controls.

## 3. Secure Secret Boundary

Certificate private keys, API secrets and internal service credentials remain server-side. Customer/mobile clients receive short-lived scoped credentials or application tokens rather than infrastructure secrets.

## 4. Service Security Posture

Every public service should expose health/readiness state, use least-privilege service identity, emit security audit events, validate configuration before reload and support safe rollback.

## Unified path

```text
Internet
  ↓
FTN Edge Proxy
  ├─ TLS / certificate
  ├─ abuse controls
  ├─ request limits
  └─ authentication boundary
       ↓
FTN API
       ↓
Service / AI
       ↓
FTN Database
```
