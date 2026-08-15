package dynamic

import "context"

// Backend is the storage/provider boundary for FTNDNS reconciliation.
// Implementations may target FTN-owned authoritative DNS systems or other
// explicitly approved backends without coupling the policy layer to one DNS engine.
type Backend interface {
	Apply(ctx context.Context, add, remove []Record) error
}

// Reconciler applies only the deterministic delta produced by Reconcile.
type Reconciler struct {
	Backend Backend
}

func (r Reconciler) Apply(ctx context.Context, current, desired []Record) error {
	add, remove := Reconcile(current, desired)
	if len(add) == 0 && len(remove) == 0 {
		return nil
	}
	if r.Backend == nil {
		return ErrBackendUnavailable
	}
	return r.Backend.Apply(ctx, add, remove)
}
