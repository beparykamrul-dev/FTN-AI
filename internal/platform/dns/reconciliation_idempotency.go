package dns

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ReconciliationKey provides a stable identity for the exact approved state
// transition. The privileged executor can persist it to reject duplicate work.
func ReconciliationKey(ctx context.Context, audit ReconciliationAudit, plan DriftPlan) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	if audit.Zone != plan.Zone || audit.ExpectedHash != plan.ExpectedHash || audit.ObservedHash != plan.ObservedHash {
		return "", fmt.Errorf("audit and plan mismatch")
	}
	if !audit.CanExecute() || plan.Action != ReconcileSync {
		return "", fmt.Errorf("reconciliation is not approved for execution")
	}
	seed := audit.ID + "|" + plan.Zone + "|" + plan.ExpectedHash + "|" + plan.ObservedHash + "|" + string(plan.Source) + "|" + string(plan.Target)
	h := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(h[:]), nil
}
