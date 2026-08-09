package gis

import "time"

type VerificationResult struct {
	PlanID string `json:"plan_id"`
	Status string `json:"status"`
	Recovered []string `json:"recovered,omitempty"`
	StillAffected []string `json:"still_affected,omitempty"`
	VerifiedAt time.Time `json:"verified_at"`
}

// Verification is observation-only. An authorized executor/collector supplies
// the observed asset state; this layer determines whether the planned recovery
// appears to have restored modeled service state.
func VerifyRecovery(plan RecoveryPlan, observed map[string]string) VerificationResult {
	r := VerificationResult{PlanID: plan.ID, Status:"verified", VerifiedAt:time.Now().UTC()}
	for _, step := range plan.Steps {
		if observed[step.Target] == "healthy" { r.Recovered = append(r.Recovered, step.Target) } else { r.StillAffected = append(r.StillAffected, step.Target) }
	}
	if len(r.StillAffected) > 0 { r.Status = "partial" }
	return r
}
