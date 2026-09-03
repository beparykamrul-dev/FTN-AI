package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type FTNFailoverJobPayload struct {
	Intent FTNFailoverIntent `json:"intent"`
	PrechangeSnapshotRequired bool `json:"prechange_snapshot_required"`
	VerificationRequired bool `json:"verification_required"`
	RollbackWhenSafe bool `json:"rollback_when_safe"`
}

// BuildFTNFailoverJobPayload binds a validated failover intent to the durable
// job contract. Submission remains separate so the caller must supply an
// approved change_approval and idempotency key.
func BuildFTNFailoverJobPayload(intent FTNFailoverIntent) (json.RawMessage, error) {
	if intent.TargetNode == "" || intent.DecisionHash == "" {
		return nil, fmt.Errorf("validated failover intent required")
	}
	payload := FTNFailoverJobPayload{Intent:intent, PrechangeSnapshotRequired:true, VerificationRequired:true, RollbackWhenSafe:true}
	b, err := json.Marshal(payload)
	if err != nil { return nil, err }
	return b, nil
}

// FTNFailoverExecutionBoundary is deliberately an adapter boundary: the
// control plane can submit/claim the job, but actual router failover remains
// owned by a separately authorized infrastructure adapter.
type FTNFailoverExecutionBoundary interface {
	Execute(context.Context, FTNFailoverJobPayload) error
	Verify(context.Context, FTNFailoverJobPayload) error
	Rollback(context.Context, FTNFailoverJobPayload) error
}
