package main

import (
	"errors"
	"sort"
	"strings"
)

// RouterOSTrafficQoSState is the adapter-neutral snapshot of FTN-managed QoS.
// It intentionally contains only state owned by FTN's QoS namespace.
type RouterOSTrafficQoSState struct {
	Rules []RouterOSTrafficQoSRule `json:"rules"`
}

type RouterOSTrafficQoSDiff struct {
	Add    []RouterOSTrafficQoSRule `json:"add,omitempty"`
	Remove []RouterOSTrafficQoSRule `json:"remove,omitempty"`
}

func qosRuleKey(r RouterOSTrafficQoSRule) string {
	return strings.TrimSpace(r.ServiceID) + "\x00" + strings.TrimSpace(r.PathID)
}

func normalizeQoSRules(rules []RouterOSTrafficQoSRule) []RouterOSTrafficQoSRule {
	seen := make(map[string]struct{}, len(rules))
	out := make([]RouterOSTrafficQoSRule, 0, len(rules))
	for _, r := range rules {
		r.ServiceID = strings.TrimSpace(r.ServiceID)
		r.PathID = strings.TrimSpace(r.PathID)
		if r.ServiceID == "" || r.PathID == "" || r.DSCP > 63 {
			continue
		}
		if _, ok := trafficPolicyByID(r.ServiceID); !ok {
			continue
		}
		key := qosRuleKey(r)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return qosRuleKey(out[i]) < qosRuleKey(out[j]) })
	return out
}

// DiffRouterOSTrafficQoS computes an idempotent desired-state diff. No device
// mutation occurs here.
func DiffRouterOSTrafficQoS(current RouterOSTrafficQoSState, desired RouterOSTrafficQoSPlan) RouterOSTrafficQoSDiff {
	currentRules := normalizeQoSRules(current.Rules)
	desiredRules := normalizeQoSRules(desired.Rules)
	currentByKey := make(map[string]RouterOSTrafficQoSRule, len(currentRules))
	desiredByKey := make(map[string]RouterOSTrafficQoSRule, len(desiredRules))
	for _, r := range currentRules { currentByKey[qosRuleKey(r)] = r }
	for _, r := range desiredRules { desiredByKey[qosRuleKey(r)] = r }
	out := RouterOSTrafficQoSDiff{}
	for key, r := range desiredByKey {
		if cur, ok := currentByKey[key]; !ok || cur.Class != r.Class || cur.DSCP != r.DSCP || cur.Priority != r.Priority {
			out.Add = append(out.Add, r)
		}
	}
	for key, r := range currentByKey {
		if _, ok := desiredByKey[key]; !ok { out.Remove = append(out.Remove, r) }
	}
	sort.Slice(out.Add, func(i, j int) bool { return qosRuleKey(out.Add[i]) < qosRuleKey(out.Add[j]) })
	sort.Slice(out.Remove, func(i, j int) bool { return qosRuleKey(out.Remove[i]) < qosRuleKey(out.Remove[j]) })
	return out
}

func (d RouterOSTrafficQoSDiff) Empty() bool { return len(d.Add) == 0 && len(d.Remove) == 0 }

// ValidateRouterOSTrafficQoSApply is the last local gate before a transport
// executor is allowed to mutate a router. The transport executor must still
// perform the pre-change snapshot and post-change verification.
func ValidateRouterOSTrafficQoSApply(device NetworkDevice, intent NetworkExecutionIntent, diff RouterOSTrafficQoSDiff) error {
	if err := ValidateNetworkExecutionIntent(intent); err != nil { return err }
	if intent.Action != NetworkConfiguration && intent.Action != NetworkApply {
		return errors.New("routeros_qos_requires_configuration_action")
	}
	if !intent.RollbackSafe {
		return errors.New("routeros_qos_requires_rollback_safe")
	}
	if strings.TrimSpace(intent.Device.ID) != strings.TrimSpace(device.ID) {
		return errors.New("routeros_qos_device_mismatch")
	}
	if diff.Empty() {
		return errors.New("routeros_qos_no_change")
	}
	return nil
}
