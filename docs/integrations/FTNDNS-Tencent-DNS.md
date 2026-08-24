# FTNDNS + Tencent Cloud DNS Integration

## FTNDNS

`ftndns.net` is the planned FTN Dynamic DNS (DDNS) service.

Core model:

`FTN Number / Account -> authenticated DDNS client -> FTNDNS API -> DNS provider adapter -> authoritative DNS`

Requirements:

- authenticated updates only;
- per-record ownership and authorization;
- IPv4/IPv6 support;
- TTL controls within policy;
- update rate limits;
- signed/API-key based client authentication;
- audit trail for every record change;
- automatic health/status reporting;
- provider-independent adapter interface;
- safe rollback of accidental record changes;
- no exposure of provider credentials to end users.

## Tencent Cloud DNS

Tencent Cloud DNS/DNSPod is an optional provider adapter in the FTN DNS control panel. It is not the authoritative source for FTN identity or policy.

Panel flow:

`DNS Provider -> Tencent Cloud -> Credentials/Secret Reference -> Zone Discovery -> Record Management -> Health/Audit`

The panel should support:

- provider connection status;
- zone discovery;
- record listing;
- A/AAAA/CNAME/TXT/MX/NS records where supported;
- DDNS record selection;
- create/update/delete through an explicit policy-controlled API;
- TTL management;
- synchronization status;
- provider error reporting;
- audit history;
- credential rotation/revocation.

Credentials must be stored as secret references, never in source code, logs, frontend bundles or ordinary configuration files. Provider API operations must be scoped to the minimum required permissions.

## Provider-neutral DNS architecture

`FTNDNS -> DNS Adapter Interface -> {PowerDNS, Knot, CoreDNS, Tencent DNS/DNSPod, Cloudflare DNS, other authorized providers}`

This keeps FTNDNS portable and prevents a single external provider from becoming an FTN core dependency.

## Operational rule

A DNS provider failure must not corrupt FTN's internal DNS state. Changes should use idempotency, validation, reconciliation and audit records before/after provider synchronization.
