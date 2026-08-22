# FTN AI V1 — Repository Index

This index groups the AI work completed on `feat/ftn-service-user-agent-v1` so implementation can continue without losing architectural context.

## Agent foundation

- `internal/platform/agent/fleet.go` — service/user/developer/assistant scopes, runtime and policy boundary.
- `internal/platform/agent/usage.go` — plan-based request/token quota gate.
- `internal/platform/agent/categories.go` — specialized category definitions.
- `internal/platform/agent/orchestrator.go` — category routing/orchestration.
- `internal/platform/agent/summary.go` — important summary + drill-down details.

## AI categories

- Studio AI
- Call Center AI
- Billing AI
- Network AI
- Developer AI
- Customer AI
- Executive AI

## Public API

- `docs/architecture/FTN-AGENT-PUBLIC-API.md`
- Canonical public endpoint: `https://api.familytimenet.com`
- Web/mobile compatible; public reachability does not imply privileged access.

## Customer and organization experience

- `docs/architecture/FTN-AI-SERVICE-EXPERIENCE.md`
- `docs/architecture/FTN-AI-CUSTOMER-SUPPORT-SUITE.md`
- Customer assistant/chart box
- Organization operations assistant
- Inventory/material, technology, tracking, safety and maintenance assistance
- FTN Chat
- FTN Alert Bot
- Call Center AI

## Developer experience

- `docs/architecture/FTN-DEVELOPER-AI-NOTES.md`
- Decision, TODO, Finding, Test, Change, Handoff and Release notes.
- Developer context is isolated from customer memory.

## Completion target

Before connecting the AI fleet deeply into FTN Cloud services, each category should have a real API adapter, scoped memory, IAM/policy, usage metering, health checks, audit events, tests, and lightweight fallback.

## Operating rule

AI is an assistant layer around real FTN services. It may explain, summarize, diagnose, recommend and prepare actions. Side-effecting actions remain behind authorization/approval and are audited.
