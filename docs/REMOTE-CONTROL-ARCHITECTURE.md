# FTN Remote Control Architecture

FTN-AI can manage authorized servers, network systems, containers, and applications from one web control plane.

## Flow

Web UI -> API Gateway -> RBAC/Approval -> Command Router -> mTLS Agent -> Server/App/Network Adapter

## Domains

- Server: health, services, logs, metrics, deployment, rollback
- Network: interfaces, routes, VLANs, DNS, BGP status, device adapters
- Applications: status, logs, config read, deploy, rollback
- Containers: list, logs, start, stop, restart
- Infrastructure: backups, health checks, controlled reloads

## Safety boundary

The control plane MUST NOT expose arbitrary shell execution. Every operation is an explicit allowlisted capability. High-risk operations require authenticated user approval, second confirmation, and an immutable audit event.

Agents should use mutually authenticated TLS for remote transport. Local Unix sockets may be used between an agent and a local gateway. Credentials and private keys must never be stored in source code.

## Server registry

Each server should have:

- stable server ID
- hostname
- management address
- agent version
- capabilities
- environment (production/staging)
- health state
- last-seen timestamp
- certificate identity

## Web experience

The dashboard should provide a server/network/app tree, live status, metrics, logs, operation history, approval prompts, deployment state, and audit trail.

## Remote execution model

1. User authenticates.
2. RBAC checks capability.
3. Router validates the requested operation against the allowlist.
4. High-risk operations pause for explicit approval.
5. Agent executes only the typed operation.
6. Result and audit event are returned to the UI.
7. Failed operations are surfaced without automatic destructive retries.
