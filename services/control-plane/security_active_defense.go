package main

import "strings"

type ActiveDefenseDecision struct {
	Action            WazuhActionClass `json:"action"`
	TargetScope       string            `json:"target_scope"`
	Automatic         bool              `json:"automatic"`
	RequiresApproval  bool              `json:"requires_approval"`
	Reason            string            `json:"reason"`
}

// BuildActiveDefenseDecision selects bounded defensive containment only.
// It never authorizes retaliation or actions against external systems.
func BuildActiveDefenseDecision(alert WazuhAlert, ftwnOwned bool) ActiveDefenseDecision {
	if !ftwnOwned {
		return ActiveDefenseDecision{Action: WazuhObserve, TargetScope: "none", Automatic: false, RequiresApproval: true, Reason: "target_ownership_not_verified"}
	}
	severity := NormalizeWazuhSeverity(alert.Severity)
	switch severity {
	case "critical", "high":
		return ActiveDefenseDecision{Action: WazuhHealthRecover, TargetScope: "ftn-owned-asset", Automatic: true, Reason: "bounded_defensive_containment"}
	case "medium":
		return ActiveDefenseDecision{Action: WazuhAlertAction, TargetScope: "ftn-owned-asset", Automatic: true, Reason: "alert_and_rate_limit_candidate"}
	default:
		return ActiveDefenseDecision{Action: WazuhObserve, TargetScope: "ftn-owned-asset", Automatic: true, Reason: "observe_only"}
	}
}

func NormalizeDefenseTargetScope(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "ftn-owned-asset" || v == "none" { return v }
	return "none"
}
