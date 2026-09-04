package main

import (
    "context"
    "encoding/json"

    "github.com/jackc/pgx/v5"
)

// appendJobEventTx records a job lifecycle transition in the same transaction
// as the state mutation. This keeps the execution journal authoritative.
func appendJobEventTx(ctx context.Context, tx pgx.Tx, tenantID, eventType, correlationID, aggregateID string, payload any) error {
    b, err := json.Marshal(payload)
    if err != nil { return err }
    // Use PostgreSQL's 64-bit extended hash for the advisory key. The lock
    // serializes per-tenant sequence allocation without relying on a 32-bit key.
    if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1::text,0))`, tenantID); err != nil { return err }
    _, err = tx.Exec(ctx, `insert into event_journal(tenant_id,event_type,sequence,correlation_id,aggregate_id,payload) values(nullif($1,'')::uuid,$2,coalesce((select max(sequence)+1 from event_journal where tenant_id=nullif($1,'')::uuid),1),$3,$4,$5::jsonb)`, tenantID, eventType, correlationID, aggregateID, string(b))
    return err
}

func jobEventPayload(j durableJob, extra map[string]any) map[string]any {
    out := map[string]any{"job_id": j.ID, "job_type": j.JobType, "status": j.Status, "attempt": j.Attempts, "approval_id": j.ApprovalID, "execution_action": j.ExecutionAction}
    for k, v := range extra { out[k] = v }
    return out
}
