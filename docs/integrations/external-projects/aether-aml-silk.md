# External Project Integration Registry

This document records external projects considered for FTN-AI integration. Projects are referenced by capability; upstream source is not copied into the core runtime unless an explicit adapter is implemented and its license/dependency requirements are satisfied.

## Aether Manager

Upstream: https://github.com/aetherdev01/aether-manager

Role: Android/device-management reference and optional client-side telemetry/profile integration boundary. The upstream project targets rooted Android devices and provides performance profiles, kernel/system tweaks and real-time monitoring. FTN should consume only an explicit, permissioned telemetry/control API; it must not assume root access or silently modify client devices.

## Agent Memory Leaderboard

Upstream: https://github.com/AML-memory/agent-memory-leaderboard

Role: AI memory evaluation/benchmark reference. Use its Add/Search evaluation concepts to validate FTN-AI long-term memory quality. Do not copy benchmark secrets, private evaluation data, credentials, or production infrastructure.

## silk-codec

Upstream: https://github.com/KasukuSakura/silk-codec

Role: optional audio codec adapter for FTN communication/media components. Keep it isolated behind a codec interface so the core control plane does not depend on native codec libraries. Prefer a pinned upstream release and platform-specific build artifacts when integration is actually enabled.

## Integration policy

- External repositories are capability references/adapters, not automatic dependencies of every FTN service.
- Pin versions/commits before production use.
- Preserve upstream license and attribution information.
- Keep credentials, private benchmark data, device secrets and production keys outside Git.
- Add health checks, compatibility checks and rollback metadata before enabling an adapter in production.
- High-impact operations remain approval-gated and auditable.
