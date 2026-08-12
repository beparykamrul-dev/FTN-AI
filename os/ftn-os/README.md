# FTN.OS

FTN.OS is the FTN-native operating platform layer. It is designed to run FTN services as managed modules while keeping the underlying host kernel replaceable and production-operable.

## Real-build boundary

FTN.OS does not pretend to replace the Linux kernel in this phase. The initial production architecture is:

```text
Host Linux kernel
      ↓
FTN.OS runtime
      ↓
FTN service manager
      ↓
FTN services / network fabric / control plane
```

## Core responsibilities

- service lifecycle
- module registry
- health and readiness
- dependency ordering
- configuration loading
- local policy enforcement
- metrics/event integration
- controlled shutdown
- audit hooks

## Safety invariants

- no external telemetry export by default
- no secrets in telemetry
- no destructive migrations
- privileged operations remain backend-policy controlled
- AI remains advisory-only
