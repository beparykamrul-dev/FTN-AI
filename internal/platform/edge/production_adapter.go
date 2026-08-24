package edge

import (
    "context"
    "fmt"
)

// ProductionAdapter is the bridge between the FTN control plane and a
// privileged node agent. It deliberately accepts only approved intents.
type ProductionAdapter interface {
    CoreRouter
    Reconcile(ctx context.Context) error
    ApplyApproved(ctx context.Context, request ChangeRequest) error
}

type AdapterState struct {
    Router RouterState `json:"router"`
    Ready bool `json:"ready"`
    LastError string `json:"last_error,omitempty"`
}

func ValidateApprovedChange(r ChangeRequest) error {
    if r.ID == "" || r.Action == "" || r.Target == "" { return fmt.Errorf("change request is incomplete") }
    if r.ApprovedBy == "" { return fmt.Errorf("approved_by is required") }
    return nil
}

func ApplyApprovedChange(ctx context.Context, a ProductionAdapter, r ChangeRequest) error {
    if err := ValidateApprovedChange(r); err != nil { return err }
    return a.ApplyApproved(ctx, r)
}
