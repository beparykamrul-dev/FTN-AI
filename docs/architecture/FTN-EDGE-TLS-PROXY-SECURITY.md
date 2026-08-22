# FTN Edge TLS + Proxy Security

FTN services use an FTN-owned edge security boundary in front of public services.

## Certificate lifecycle

```text
FTN service/domain
      ↓
FTN Certificate Controller
      ↓
ACME account + policy
      ↓
Certificate issuance
      ↓
Validation
      ↓
Secure certificate store
      ↓
FTN Edge Proxy
```

The controller automatically renews certificates before expiry, validates the resulting certificate, reloads the proxy safely, and records lifecycle events. Renewal failure must not discard a currently valid certificate.

## FTN proxy security baseline

- TLS 1.2+ with modern TLS 1.3 preference
- automated certificate renewal
- secure private-key storage and restrictive permissions
- HTTP-to-HTTPS policy where appropriate
- HSTS for eligible HTTPS-only domains
- strict security headers
- request size and timeout limits
- rate limiting
- connection limits
- upstream health checks
- origin isolation
- access logging and security audit events
- safe configuration validation before reload
- fail-closed behavior for invalid security configuration

## Fast security path

The proxy should reject obviously invalid or abusive traffic as early as possible, before expensive application/AI processing.

```text
Internet
  ↓
FTN Edge Proxy
  ├─ TLS validation
  ├─ request limits
  ├─ rate limit
  ├─ policy/auth boundary
  └─ abuse filtering
       ↓
   FTN API / Services
       ↓
      AI
```

## Customer certificate UX

Customers should not manually manage ordinary FTN-managed service certificates. Where the service/domain is eligible and authorized, FTN handles issuance and renewal automatically.

Certificate authority/provider credentials remain server-side. FTN's public API and applications remain independent of the underlying CA provider.

## Safety boundary

Certificate issuance must verify domain/service ownership or authorization. The controller must never issue certificates for domains that the FTN account is not authorized to manage.
