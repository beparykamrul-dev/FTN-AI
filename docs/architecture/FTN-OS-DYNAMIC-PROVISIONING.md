# FTN OS Dynamic Provisioning

## Scope

The FTN OS Builder uses one small base image and role-specific service profiles. POP, Reseller, User, Client, Cache, DNS, Hosting and Backup profiles are provisioned from the central service registry and entitlement policy.

## Lifecycle

1. Build signed base image.
2. Install as Proxmox VM/LXC or supported ISO/cloud image.
3. First-boot bootstrap detects network and hardware identity.
4. Register node with FTN Control Plane.
5. Obtain short-lived credentials and mTLS identity.
6. Resolve authorized service entitlement.
7. Generate service manifest.
8. Activate only authorized modules and desktop UI groups.
9. Register health/KPI and diagnostics.
10. Snapshot before updates and rollback on failed health checks.

## Identity

Device identity may include MAC address, system UUID and serial number, with hardware inventory for CPU, RAM, motherboard, BIOS/UEFI, storage, NIC, GPU, TPM and temperature where supported.

Sensitive identifiers are protected by the platform's authorization and privacy policies. They are not secrets and must not be exposed through public UI or telemetry.

## Dynamic UI

The FTN desktop shell is manifest-driven. A service appears only after backend entitlement authorization. Removing an entitlement revokes backend access and removes/disables its UI module.

Web, Android and desktop clients consume the same service registry and authorization model.

## Diagnostics

The diagnostic engine correlates identity, hardware health, logs, metrics, dependencies and service state. AI remains advisory: it produces evidence-backed findings and remediation plans. Privileged actions, firmware changes and destructive operations require explicit authorization and approval.

## Deployment safety

- No production secrets in images.
- Signed modules only.
- No blind restart/failover.
- Preflight and post-change health checks.
- Snapshot/rollback before updates.
- Quarantine on security failure.
- Immutable audit trail.
- No external telemetry export by default.
