# FTN AI Completion Roadmap

The AI layer is completed in vertical slices before FTN Cloud service expansion.

## Core

- Agent fleet: service, user, developer, assistant scopes
- Specialized category orchestrator
- Lightweight-first shared runtime
- Usage/quota gate
- IAM and approval boundary
- Tenant/user isolation

## Specialized agents

- Studio AI: preview/build analysis
- Call Center AI: conversation assistance and summaries
- Billing AI: billing/payment analysis and explanation
- Network AI: service/network diagnostics
- Developer AI: repository/code/test assistance
- Customer AI: customer-scoped assistant
- Executive AI: concise summaries with drill-down details

## Control-plane requirements

Every agent request resolves scope, category, entitlement, policy and audit context before runtime execution. Side-effecting tools require explicit authorization/approval.

## Completion criterion

The AI layer is considered ready for FTN Cloud integration only when each category has: API adapter, policy, memory boundary, usage metering, health checks, audit events, tests, and a lightweight fallback path.
