# FTN Control Plane v2

## Goal

Keep the existing FTN implementation intact while making the Control Plane the stable integration boundary for future modules. New modules must be adapters/plugins around canonical FTN models rather than vendor-specific logic in the core.

## Safety and recovery

- Existing source/configuration must be backed up before replacement or migration.
- Desired state and observed state are separate records.
- Production changes require authenticated identity, RBAC and explicit approval.
- Every deployment/configuration change is versioned and auditable.
- Verification runs after an authorized change; failed verification can produce a rollback plan.

## Managed device model

Server, virtual server, router, OLT, ONU, PC, Android and TV use the same FTN device identity model. Hardware identifiers such as MAC, IP and serial are inventory metadata, not authentication credentials.

## Control Plane domains

- Fleet/server registry and FTN Agent
- Device adapter registry
- Source/import and build pipeline
- Configuration/deployment targets
- DNS, WAN and network adapters
- GIS/fiber/OLT/ONU/IPAM adapters
- Observability and telemetry
- WebSocket live state
- Android/PC/TV client management
- AI detection, recovery planning and verification

## GIS and fiber customer view

The canonical graph may relate fiber segments, joints, splitters, OLTs, ONUs, routers, PPPoE sessions, IPs and customer/service records. Customer-facing applications should receive only the records and actions authorized for that customer/role.

## Firewall / FTN OS direction

The FTN virtual/local router layer is designed around Linux networking primitives such as netlink, nftables and WireGuard through privileged, authenticated agents. The Control Plane expresses desired intent; the dataplane agent performs only authorized operations.

## Module contract

Each module should expose a stable adapter boundary and declare capabilities. The core should reject unsupported capabilities before execution. Modules can therefore evolve independently without rewriting the Control Plane.
