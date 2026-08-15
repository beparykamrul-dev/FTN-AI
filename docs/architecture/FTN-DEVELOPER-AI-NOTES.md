# FTN Developer AI Notes

Developer AI has a dedicated note/work-log layer for preserving engineering context without mixing it into customer memory.

## Note types

- Decision: why an architecture or implementation choice was made.
- TODO: unfinished engineering work.
- Finding: observed bug, risk, dependency, or behavior.
- Test: test result and verification evidence.
- Change: concise record of an implementation change.
- Handoff: context another developer/agent needs to continue work.
- Release: release readiness, migration, and rollback notes.

## Required metadata

Each note should carry repository, branch/ref, component, author/agent identity, timestamp, category, and links to relevant issues/commits where available.

## AI behavior

Developer AI can summarize notes, detect stale TODOs, group related findings, identify conflicting decisions, and produce a continuation brief. It must distinguish facts from suggestions and must not silently modify code or repository state.

## Isolation

Developer notes are separate from customer conversations and customer memory. Repository access is permission-scoped.

## Suggested developer panel

```text
FTN Developer AI
├── Important Notes
├── Current Work
├── Decisions
├── Findings
├── Tests
├── TODOs
├── Handoffs
└── Release Readiness
```

This note system is intended to make FTN development continuous across developers and AI agents while retaining an auditable engineering history.
