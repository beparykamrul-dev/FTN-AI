package main

import (
	"errors"
	"strings"
	"time"
)

// RouterOSQoSVerification is a side-effect-free post-change verification result.
// Verification never applies, repairs, or mutates RouterOS state.
type RouterOSQoSVerification struct {
	DeviceID       string       `json:"device_id"`
	Verified       bool         `json:"verified"`
	Drift          bool         `json:"drift"`
	Diff           RouterOSDiff `json:"diff"`
	ExpectedRules  int          `json:"expected_rules"`
	ObservedRules  int          `json:"observed_rules"`
	VerifiedAt     time.Time    `json:"verified_at"`
}

// VerifyRouterOSQoSState compares an observed post-change snapshot against the
// exact desired state. Invalid input and device mismatches fail closed.
func VerifyRouterOSQoSState(snapshot RouterOSSnapshot, desired RouterOSDesiredState) (RouterOSQoSVerification, error) {
	s, err := NormalizeRouterOSSnapshot(snapshot)
	if err != nil {
		return RouterOSQoSVerification{}, err
	}
	d, err := NormalizeRouterOSDesiredState(desired)
	if err != nil {
		return RouterOSQoSVerification{}, err
	}
	if s.DeviceID != d.DeviceID {
		return RouterOSQoSVerification{}, errors.New("routeros_device_mismatch")
	}
	diff, err := DiffRouterOSQoSState(s, d)
	if err != nil {
		return RouterOSQoSVerification{}, err
	}
	return RouterOSQoSVerification{
		DeviceID:      s.DeviceID,
		Verified:      diff.Empty(),
		Drift:         !diff.Empty(),
		Diff:          diff,
		ExpectedRules: len(d.Rules),
		ObservedRules: len(s.Rules),
		VerifiedAt:    time.Now().UTC(),
	}, nil
}

// RouterOSQoSRollbackPlan identifies the exact pre-change state as the rollback
// target. It is still plan-only: applying it requires a separate explicit
// approval and execution boundary.
type RouterOSQoSRollbackPlan struct {
	DeviceID         string             `json:"device_id"`
	Target           RouterOSSnapshot   `json:"target"`
	Current          RouterOSSnapshot   `json:"current"`
	Diff             RouterOSDiff       `json:"diff"`
	RequiresApproval bool               `json:"requires_approval"`
	ApplyAllowed     bool               `json:"apply_allowed"`
	CreatedAt        time.Time          `json:"created_at"`
}

// BuildRouterOSQoSRollbackPlan computes the changes required to restore the
// captured pre-change snapshot. It never performs a network write.
func BuildRouterOSQoSRollbackPlan(preChange RouterOSSnapshot, current RouterOSSnapshot) (RouterOSQoSRollbackPlan, error) {
	pre, err := NormalizeRouterOSSnapshot(preChange)
	if err != nil {
		return RouterOSQoSRollbackPlan{}, err
	}
	cur, err := NormalizeRouterOSSnapshot(current)
	if err != nil {
		return RouterOSQoSRollbackPlan{}, err
	}
	if strings.TrimSpace(pre.DeviceID) == "" || pre.DeviceID != cur.DeviceID {
		return RouterOSQoSRollbackPlan{}, errors.New("routeros_rollback_device_mismatch")
	}
	diff, err := DiffRouterOSQoSState(cur, RouterOSDesiredState{DeviceID: pre.DeviceID, Rules: pre.Rules})
	if err != nil {
		return RouterOSQoSRollbackPlan{}, err
	}
	return RouterOSQoSRollbackPlan{
		DeviceID:         pre.DeviceID,
		Target:           pre,
		Current:          cur,
		Diff:             diff,
		RequiresApproval: !diff.Empty(),
		ApplyAllowed:     false,
		CreatedAt:        time.Now().UTC(),
	}, nil
}
