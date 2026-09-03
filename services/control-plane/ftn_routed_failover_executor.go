package main

import (
	"context"
	"fmt"
	"strings"
)

type FTNCoreFailoverAdapter interface {
	SwitchActiveCore(context.Context, string) error
	VerifyActiveCore(context.Context, string) error
	RollbackActiveCore(context.Context, string) error
}

type FTNFailoverExecutor struct { adapter FTNCoreFailoverAdapter }

func NewFTNFailoverExecutor(adapter FTNCoreFailoverAdapter) *FTNFailoverExecutor {
	return &FTNFailoverExecutor{adapter: adapter}
}

// Execute is intentionally impossible without an injected authorized adapter.
// The adapter is responsible for the actual FTN-owned infrastructure action.
func (e *FTNFailoverExecutor) Execute(ctx context.Context, payload FTNFailoverJobPayload) error {
	if e == nil || e.adapter == nil { return fmt.Errorf("authorized failover adapter required") }
	if strings.TrimSpace(payload.Intent.TargetNode) == "" || payload.Intent.DecisionHash == "" { return fmt.Errorf("validated failover intent required") }
	if !payload.PrechangeSnapshotRequired || !payload.VerificationRequired { return fmt.Errorf("snapshot and verification are required") }
	return e.adapter.SwitchActiveCore(ctx, payload.Intent.TargetNode)
}

func (e *FTNFailoverExecutor) Verify(ctx context.Context, payload FTNFailoverJobPayload) error {
	if e == nil || e.adapter == nil { return fmt.Errorf("authorized failover adapter required") }
	if payload.Intent.TargetNode == "" { return fmt.Errorf("target core required") }
	return e.adapter.VerifyActiveCore(ctx, payload.Intent.TargetNode)
}

func (e *FTNFailoverExecutor) Rollback(ctx context.Context, payload FTNFailoverJobPayload) error {
	if e == nil || e.adapter == nil { return fmt.Errorf("authorized failover adapter required") }
	if !payload.RollbackWhenSafe { return fmt.Errorf("rollback is not enabled") }
	return e.adapter.RollbackActiveCore(ctx, payload.Intent.TargetNode)
}
