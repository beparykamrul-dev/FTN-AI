# FTN Certificate + Edge Security Automation

FTN provides a centralized security automation layer for eligible FTN-managed domains and services.

## Certificate automation

```text
User / FTN Service
      ↓
Domain ownership / authorization check
      ↓
FTN Certificate Controller
      ↓
ACME issuance / renewal
      ↓
Certificate validation
      ↓
Atomic certificate deployment
      ↓
FTN Edge Proxy
```

Certificates are renewed before expiry. A failed renewal keeps the last known-good certificate active. Private keys are never exposed to clients.

## Fast security path

The FTN edge proxy performs inexpensive controls before requests reach application and AI workloads:

- TLS termination and protocol policy
- host/SNI validation
- request-size limits
- connection and timeout limits
- rate limiting
- abuse/anomaly signals
- authentication/policy handoff
- upstream health checks
- security headers

Expensive AI processing happens only after the edge security boundary accepts the request.

## Strong proxy posture

The proxy is treated as a security boundary, not merely a reverse proxy. Configuration is validated before reload, invalid security configuration fails closed, and certificate/key changes are deployed atomically.

## FTN ownership

`api.familytimenet.com` remains an FTN-owned API boundary. CA vendors and internal proxy implementations can change behind FTN's control plane without changing customer applications.

## Authorization rule

Automatic certificate issuance is limited to domains that FTN can verify as owned or explicitly authorized by the account/service. The automation must not be used to obtain certificates for unrelated domains.
