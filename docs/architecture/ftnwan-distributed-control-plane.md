# FTNWAN Distributed Control Plane

## Purpose

Define the multi-server control-plane model for FTNWAN without coupling the control plane to any single router, VPN, proxy, or vendor implementation.

## Model

```text
FTN Global Control Plane
        |
   +----+----+
   |    |    |
  CP-01 CP-02 CP-03
   |    |    |
   +----+----+
        |
  Consistent State
        |
  FTNWAN Device Fabric
```

## State

Each control-plane node maintains:

- node identity and authorization
- topology snapshot
- device registry
- network intent state
- policy version
- route/path state
- configuration version
- audit sequence
- health/lease information

## Reconciliation

```text
Observe -> Validate -> Replicate -> Consistency Check -> Commit
```

Network changes follow:

```text
Intent -> Policy -> Plan -> Approval -> Apply -> Verify -> Rollback
```

## Failover

A node becomes eligible to take control only after identity, authorization, lease, epoch and state consistency are verified. Existing data-plane paths are not automatically changed merely because a control-plane node fails.

## Security boundary

The control plane is separate from the FTNWAN data plane. Replication must be authenticated and integrity-protected. Secrets and private keys are never stored in topology/state snapshots.
