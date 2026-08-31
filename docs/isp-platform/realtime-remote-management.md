# FTN ISP Realtime and Remote Device Management

## Realtime transport

WebSocket endpoint: `wss://api.familytimenet.com/ws/v1`.

Use WebSocket for authorized realtime events only: device health, interface counters, ONU/OLT state, PPPoE state, billing notifications, tickets, incidents, AI updates and job progress. REST remains the source of truth for commands and durable resources.

Every WebSocket connection requires an authenticated application session. Channels are scoped by role, account, tenant and resource permissions. Server-side authorization is rechecked before subscribing to a resource.

## Android API

The Android application uses the versioned REST API for account, billing, support, service status, speed-test and durable actions, and WebSocket for live status/notifications. Android never receives device credentials or direct management-plane addresses.

## Remote management architecture

```text
Android / Admin Panel
        |
     HTTPS/WSS
        |
  API Gateway + Auth/RBAC
        |
  Management Control Plane
        |
  Provider/Device Adapters
   |       |       |       |
MikroTik  OLT     ONU    Switch
   |       |       |       |
   +-------+-------+-------+
           FTN interfaces
```

## FTN device interface abstraction

All supported infrastructure devices are represented through a common adapter contract:

- identity and inventory
- health and capabilities
- interface inventory
- telemetry collection
- configuration read
- configuration change request
- command/job status
- event subscription where supported
- backup/export metadata
- rollback capability where supported

Provider adapters may use documented vendor APIs, SNMP, SSH/NETCONF/RESTCONF or other explicitly supported management interfaces. Credentials are stored server-side and are never returned to clients.

## Device classes

- MikroTik routers
- OLTs
- ONUs/ONTs
- managed switches
- access points
- servers
- UPS/PDU where supported
- FTN service nodes
- monitoring collectors

## Remote action safety

Read operations may be permitted by role. Mutating operations require the appropriate permission and, for high-impact operations, an explicit approval request. The control plane records who requested, approved and executed an operation.

High-impact actions must support preflight validation, timeout, idempotency where applicable, post-change verification and rollback when the adapter supports it.

## WebSocket event envelope

```json
{
  "id": "event-id",
  "type": "device.health.updated",
  "timestamp": "RFC3339",
  "resource": "device",
  "resource_id": "device-id",
  "sequence": 1,
  "payload": {}
}
```

Clients must treat events as hints and refetch authoritative state through REST when required. Sequence numbers allow clients to detect gaps.

## Security requirements

- TLS for HTTPS/WSS
- short-lived access tokens with refresh sessions
- RBAC and resource ownership checks
- tenant/account isolation
- rate limits and connection limits
- per-channel authorization
- audit logging for privileged actions
- server-side secrets only
- fail-closed authorization
- approval gate for high-impact mutations
- no arbitrary command execution from client input
