package ftnstorage

import "context"

// RecoveryExecutor is the authorized execution boundary for storage recovery.
type RecoveryExecutor interface {
	Reconstruct(context.Context, ChunkRef) error
	Rollback(context.Context, RollbackPlan) error
	Rebalance(context.Context, RebalancePlan) error
}
