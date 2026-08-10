package dns

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

type ReconcileAction string

const (
	ReconcileNoop ReconcileAction = "noop"
	ReconcileReview ReconcileAction = "review"
	ReconcileSync ReconcileAction = "sync"
)

type DriftPlan struct {
	Zone string `json:"zone"`
	Source ProviderType `json:"source"`
	Target ProviderType `json:"target"`
	ExpectedHash string `json:"expected_hash"`
	ObservedHash string `json:"observed_hash"`
	Action ReconcileAction `json:"action"`
	Reason string `json:"reason"`
}

// BuildDriftPlan creates an approval-first reconciliation plan. It never
// mutates a DNS provider; execution belongs to a separate privileged layer.
func BuildDriftPlan(ctx context.Context, expected, observed ZoneSnapshot) (DriftPlan, error) {
	select { case <-ctx.Done(): return DriftPlan{}, ctx.Err(); default: }
	if strings.TrimSpace(expected.Zone) == "" || strings.TrimSpace(observed.Zone) == "" { return DriftPlan{}, fmt.Errorf("zone is required") }
	if !strings.EqualFold(strings.TrimSuffix(expected.Zone, "."), strings.TrimSuffix(observed.Zone, ".")) { return DriftPlan{}, fmt.Errorf("zone mismatch") }
	plan := DriftPlan{Zone: observed.Zone, Source: expected.Provider, Target: observed.Provider, ExpectedHash: expected.Hash, ObservedHash: observed.Hash, Action: ReconcileNoop, Reason: "provider state matches expected snapshot"}
	if expected.Hash != observed.Hash {
		plan.Action = ReconcileReview
		plan.Reason = "provider state drifted from expected snapshot; manual approval required"
	}
	return plan, nil
}

func SortDriftPlans(plans []DriftPlan) {
	sort.Slice(plans, func(i, j int) bool {
		return strings.ToLower(plans[i].Zone)+string(plans[i].Target) < strings.ToLower(plans[j].Zone)+string(plans[j].Target)
	})
}
