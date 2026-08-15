# FTN Content + Client Transfer Platform

FTN's existing local-service catalog already defines FTP, Drive, Media Server, App Store, Baby Home, Baby Game Store, Learning, AI Assistant and Tools. This layer standardizes how those services are exposed by category across web, Android and desktop clients.

## Client categories

- Storage: FTP + Drive
- Media: Movie/media server
- Family: Baby Home + Baby Game Store
- Learning: learning video, courses and progress
- Tools: file, diagnostics, speed-test and network tools
- Applications: FTN App Store
- Assistant: FTN AI Assistant

## Upload / download

The common transfer contract supports resumable uploads, chunking, SHA-256 checksums, quotas, authorization, resumable downloads and HTTP range requests.

```text
Client
  ↓
FTN API / Edge
  ↓
Authorization + Quota
  ↓
Transfer Service
  ↓
FTN Drive / FTP / Media Storage
```

No client receives storage infrastructure credentials. Transfers are scoped to the authenticated tenant/account and service policy.

## Auto import and server push

`configs/v1/ftn-content-client-platform.yaml` references the existing service catalog rather than creating duplicate services. The existing Control Plane, FTN OS lifecycle, deployment and server-fabric components remain responsible for discovery, startup, health gating and rollback.

A service is enabled only after dependency validation and health checks. Failed startup triggers rollback where supported. Destructive operations remain approval-gated.

## Important implementation boundary

This manifest defines the production contract and wiring. It does not pretend that FTP/media/storage backends are already fully implemented. Backend implementation, API handlers, storage adapters and client UI must be completed and tested against this contract before production enablement.
