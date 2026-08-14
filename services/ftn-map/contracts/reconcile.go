package fiber

import "context"

type ReconcileState string
const (
	ReconcileHealthy ReconcileState = "healthy"
	ReconcilePending ReconcileState = "pending"
	ReconcileConflict ReconcileState = "conflict"
)

type ReconcileResult struct {
	EntityID string
	State ReconcileState
	Applied int
	Conflicts int
	Message string
}

// Reconciler converges local POP/client topology with the FTN control plane
// without blindly overwriting newer or conflicting field state.
type Reconciler interface {
	Reconcile(context.Context, FiberTopology) (ReconcileResult, error)
}
