package main

import (
	"errors"
	"sort"
	"strings"
	"time"
)

type RouterOSQoSState struct {
	ServiceID string       `json:"service_id"`
	Class     TrafficClass `json:"class"`
	DSCP      uint8        `json:"dscp"`
	Priority  uint8        `json:"priority"`
	PathID    string       `json:"path_id"`
}

type RouterOSSnapshot struct {
	DeviceID string             `json:"device_id"`
	Rules    []RouterOSQoSState `json:"rules"`
	CapturedAt time.Time        `json:"captured_at"`
}

type RouterOSDesiredState struct {
	DeviceID string             `json:"device_id"`
	Rules    []RouterOSQoSState `json:"rules"`
}

type RouterOSQoSDiffChange struct {
	Before RouterOSQoSState `json:"before"`
	After  RouterOSQoSState `json:"after"`
}

type RouterOSDiff struct {
	Adds    []RouterOSQoSState       `json:"adds,omitempty"`
	Changes []RouterOSQoSDiffChange  `json:"changes,omitempty"`
	Removes []RouterOSQoSState       `json:"removes,omitempty"`
}

func (d RouterOSDiff) Empty() bool {
	return len(d.Adds) == 0 && len(d.Changes) == 0 && len(d.Removes) == 0
}

func normalizeRouterOSQoSState(rules []RouterOSQoSState) ([]RouterOSQoSState, error) {
	out := make([]RouterOSQoSState, 0, len(rules))
	seen := make(map[string]struct{}, len(rules))
	for _, rule := range rules {
		rule.ServiceID = strings.TrimSpace(rule.ServiceID)
		rule.PathID = strings.TrimSpace(rule.PathID)
		rule.Class = TrafficClass(strings.TrimSpace(string(rule.Class)))
		if rule.ServiceID == "" || rule.PathID == "" {
			return nil, errors.New("routeros_qos_identity_required")
		}
		if rule.DSCP > 63 {
			return nil, errors.New("routeros_qos_dscp_out_of_range")
		}
		if _, ok := trafficPolicyByID(rule.ServiceID); !ok {
			return nil, errors.New("routeros_qos_unknown_service")
		}
		key := rule.ServiceID + "\x00" + rule.PathID
		if _, ok := seen[key]; ok {
			return nil, errors.New("routeros_qos_duplicate_rule")
		}
		seen[key] = struct{}{}
		out = append(out, rule)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ServiceID != out[j].ServiceID { return out[i].ServiceID < out[j].ServiceID }
		return out[i].PathID < out[j].PathID
	})
	return out, nil
}

func NormalizeRouterOSQoSState(rules []RouterOSQoSState) ([]RouterOSQoSState, error) {
	return normalizeRouterOSQoSState(rules)
}

func NormalizeRouterOSSnapshot(snapshot RouterOSSnapshot) (RouterOSSnapshot, error) {
	if strings.TrimSpace(snapshot.DeviceID) == "" { return RouterOSSnapshot{}, errors.New("routeros_snapshot_device_id_required") }
	rules, err := normalizeRouterOSQoSState(snapshot.Rules)
	if err != nil { return RouterOSSnapshot{}, err }
	snapshot.DeviceID = strings.TrimSpace(snapshot.DeviceID)
	snapshot.Rules = rules
	return snapshot, nil
}

func NormalizeRouterOSDesiredState(desired RouterOSDesiredState) (RouterOSDesiredState, error) {
	if strings.TrimSpace(desired.DeviceID) == "" { return RouterOSDesiredState{}, errors.New("routeros_desired_device_id_required") }
	rules, err := normalizeRouterOSQoSState(desired.Rules)
	if err != nil { return RouterOSDesiredState{}, err }
	desired.DeviceID = strings.TrimSpace(desired.DeviceID)
	desired.Rules = rules
	return desired, nil
}

func routerOSQoSKey(r RouterOSQoSState) string { return r.ServiceID + "\x00" + r.PathID }

func DiffRouterOSQoSState(snapshot RouterOSSnapshot, desired RouterOSDesiredState) (RouterOSDiff, error) {
	s, err := NormalizeRouterOSSnapshot(snapshot)
	if err != nil { return RouterOSDiff{}, err }
	d, err := NormalizeRouterOSDesiredState(desired)
	if err != nil { return RouterOSDiff{}, err }
	if s.DeviceID != d.DeviceID { return RouterOSDiff{}, errors.New("routeros_device_mismatch") }

	before := make(map[string]RouterOSQoSState, len(s.Rules))
	after := make(map[string]RouterOSQoSState, len(d.Rules))
	for _, r := range s.Rules { before[routerOSQoSKey(r)] = r }
	for _, r := range d.Rules { after[routerOSQoSKey(r)] = r }

	out := RouterOSDiff{}
	for key, r := range after {
		old, ok := before[key]
		if !ok { out.Adds = append(out.Adds, r); continue }
		if old != r { out.Changes = append(out.Changes, RouterOSQoSDiffChange{Before: old, After: r}) }
	}
	for key, r := range before {
		if _, ok := after[key]; !ok { out.Removes = append(out.Removes, r) }
	}
	sort.Slice(out.Adds, func(i, j int) bool { return routerOSQoSKey(out.Adds[i]) < routerOSQoSKey(out.Adds[j]) })
	sort.Slice(out.Removes, func(i, j int) bool { return routerOSQoSKey(out.Removes[i]) < routerOSQoSKey(out.Removes[j]) })
	sort.Slice(out.Changes, func(i, j int) bool { return routerOSQoSKey(out.Changes[i].After) < routerOSQoSKey(out.Changes[j].After) })
	return out, nil
}
