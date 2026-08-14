# FTN Node Runtime

Shared FTN.OS node-role architecture for Main, POP, Client and Backup nodes.

The same runtime contracts are used across roles; role policy determines which capabilities are enabled.

## Roles

- Main: core gateway and global/local peering control
- POP: local distribution relay and edge bridge
- Client: customer-side dedicated edge and localized mesh
- Backup: failover/state replication and routing convergence assistance

Every role uses FTN identity, authorization, health, reconciliation, signed release and rollback controls. The AI assistant is available on each access surface but cannot bypass authorization or execute arbitrary privileged commands.

CCTV remains intentionally outside this runtime milestone.
