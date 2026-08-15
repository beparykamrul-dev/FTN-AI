package dynamic

import "errors"

var ErrRepairNotAuthorized = errors.New("FTNDNS repair requires explicit authorization")

// RepairAuthorization is the policy gate for applying a DNS repair plan.
type RepairAuthorization struct {
	Approved   bool
	Actor      string
	PolicyID   string
	PlanVersion uint64
}

// AuthorizeRepair validates the explicit policy gate before execution.
// Backend mutation is intentionally kept outside this package.
func AuthorizeRepair(plan StateRepair, auth RepairAuthorization) error {
	if plan.Version == 0 || !auth.Approved || auth.Actor == "" || auth.PolicyID == "" {
		return ErrRepairNotAuthorized
	}
	if auth.PlanVersion != plan.Version {
		return ErrRepairNotAuthorized
	}
	return nil
}
