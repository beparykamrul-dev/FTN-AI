# FTN Source Consolidation Manifest

## Purpose

Preserve and consolidate the existing FTN plans, source implementations, scripts, configurations, and architecture decisions into one production-oriented repository without silently discarding older work.

## Source domains identified

- FTN Enterprise / Sprint foundation
- GoFiber/API and backend services
- React/web and Android clients
- AI assistant, automation, memory and knowledge layers
- ISP/NOC monitoring
- MikroTik/RouterOS and Radius integration
- OLT/ONU/PPPoE discovery and telemetry
- Fiber Map, topology, OTDR and optical metrics
- Device health and recovery
- Alert/SMS/call escalation
- Billing/payment
- TV/IPTV, FTP/file services and app-store concepts
- Edge/CDN/cache and multi-server routing
- DNS/global service architecture

## Existing source evidence

The source library contains prior FTN material covering monitoring with Prometheus/LibreNMS/OTDR/NetFlow/SmokePing, fiber topology and customer mapping, and AI-style alerts. It also contains prior multi-panel ISP architecture and the Sprint-01 foundation structure.

## Consolidation rules

1. Preserve original source before refactoring.
2. Detect duplicates before selecting a canonical implementation.
3. Prefer the newest production-capable implementation when evidence supports it.
4. Keep superseded implementations under an explicit legacy/reference area instead of deleting them silently.
5. Track dependencies and missing interfaces.
6. Separate specification, source, deployment, tests, and generated artifacts.
7. Keep privileged automation policy-controlled and auditable.
8. No prototype/demo implementation should silently replace production code.

## Target layout

```text
SOURCE/
  backend/
  frontend/
  android/
  router/
  fiber/
  monitoring/
  ai/
  alert/
  billing/
  services/
  shared/
DOCS/
  architecture/
  api/
  deployment/
  operations/
  recovery/
TESTS/
DOCKER/
DATABASE/
LEGACY/
  plans/
  superseded/
TOOLS/
```

## Current status

- Source inventory: started
- Domain consolidation map: defined
- Duplicate detection: pending full source extraction
- Canonical implementation selection: pending source comparison
- Production verification: pending
- Repository import: incremental
