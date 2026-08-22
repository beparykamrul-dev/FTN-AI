# FTN AI Build Factory

A shared build-assistant layer coordinates customer project building and FTN internal project building without mixing their data boundaries.

## Tracks

1. Customer Project Track — requirement to project workspace.
2. FTN Project Track — internal service/module/repository work.
3. Network Operations Track — isolated telemetry and NOC reasoning.

## Shared capabilities

- requirements
- architecture
- task graph
- code/documentation assistance
- testing
- release readiness
- progress summaries

## Isolation

Each track has its own identity, context, memory, permissions and audit scope. Network operations never inherit customer project context automatically.

## Lightweight-first

Use rules/retrieval and small runtime paths first. Escalate only when the requested capability requires a stronger approved layer.
