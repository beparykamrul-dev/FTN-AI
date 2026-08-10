package dns

import (
	"context"
	"fmt"
)

type DNSReconciliationExecutor interface {
	Sync(ctx context.Context, plan DriftPlan) error
}

// ExecuteReconciliation is the final orchestration boundary. Provider mutation
// is delegated to an injected executor and is blocked unless the audit guard
// explicitly authorizes the operation.
func ExecuteReconciliation(ctx context.Context, guard *ReconciliationGuard, audit ReconciliationAudit, executor DNSReconciliationExecutor, plan DriftPlan) error {
	if guard == nil { return fmt.Errorf("reconciliation guard is required") }
	if executor == nil { return fmt.Errorf("reconciliation executor is required") }
	if err := guard.Check(ctx, audit); err != nil { return err }
	if plan.Action != ReconcileSync { return fmt.Errorf("plan action %q is not executable", plan.Action) }
	if plan.Zone != audit.Zone || plan.ExpectedHash != audit.ExpectedHash || plan.ObservedHash != audit.ObservedHash { return fmt.Errorf("audit and reconciliation plan do not match") }
	return executor.Sync(ctx, plan)
}
