package dns

import (
	"context"
	"fmt"
)

// ReconciliationAudit binds the exact state transition that was reviewed and
// approved. It contains no provider credentials or secret material.
type ReconciliationAudit struct {
	ID string `json:"id"`
	Zone string `json:"zone"`
	ExpectedHash string `json:"expected_hash"`
	ObservedHash string `json:"observed_hash"`
	Action ReconcileAction `json:"action"`
	Approved bool `json:"approved"`
}

func (a ReconciliationAudit) CanExecute() bool {
	return a.Approved && a.Action == ReconcileSync && a.ID != ""
}

// ReconciliationGuard enforces the final approval boundary. It intentionally
// produces a decision only; privileged provider mutation remains elsewhere.
type ReconciliationGuard struct{}

func NewReconciliationGuard() *ReconciliationGuard { return &ReconciliationGuard{} }

func (g *ReconciliationGuard) Check(ctx context.Context, audit ReconciliationAudit) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if !audit.Approved {
		return fmt.Errorf("DNS reconciliation requires explicit approval")
	}
	if audit.Action != ReconcileSync {
		return fmt.Errorf("reconciliation action %q is not executable", audit.Action)
	}
	if audit.Zone == "" || audit.ExpectedHash == "" || audit.ObservedHash == "" {
		return fmt.Errorf("incomplete reconciliation audit")
	}
	if audit.ExpectedHash == audit.ObservedHash {
		return fmt.Errorf("reconciliation is unnecessary: snapshots already match")
	}
	return nil
}
