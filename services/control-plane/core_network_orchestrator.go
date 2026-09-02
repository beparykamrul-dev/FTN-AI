package main

import (
	"errors"
	"sort"
	"strings"
	"time"
)

type CoreNetworkSnapshot struct {
	CapturedAt time.Time `json:"captured_at"`
	CoreA      CoreRouterHealth `json:"core_a"`
	CoreB      CoreRouterHealth `json:"core_b"`
	Paths      []TunnelEndpoint `json:"paths"`
}

type CoreNetworkDecision struct {
	ActiveCore string `json:"active_core"`
	StandbyCore string `json:"standby_core"`
	Paths []TunnelEndpoint `json:"paths"`
	FailoverRequired bool `json:"failover_required"`
	RequiresApproval bool `json:"requires_approval"`
	SnapshotRequired bool `json:"snapshot_required"`
	VerifyRequired bool `json:"verify_required"`
	RollbackWhenSafe bool `json:"rollback_when_safe"`
	Reason string `json:"reason"`
}

func PlanCoreNetworkFailover(a, b CoreRouterHealth, policy MultiTunnelPolicy, endpoints []TunnelEndpoint) (CoreNetworkDecision, error) {
	if err := ValidateMultiTunnelPolicy(policy); err != nil { return CoreNetworkDecision{}, err }
	if a.NodeID == "" || b.NodeID == "" { return CoreNetworkDecision{}, errors.New("both_core_nodes_required") }
	if a.NodeID == b.NodeID { return CoreNetworkDecision{}, errors.New("core_nodes_must_differ") }
	paths := SelectMultiTunnelPaths(policy, endpoints)
	d := CoreNetworkDecision{Paths: paths, RequiresApproval: true, SnapshotRequired: true, VerifyRequired: true, RollbackWhenSafe: true}
	if a.Healthy { d.ActiveCore, d.StandbyCore, d.Reason = a.NodeID, b.NodeID, "primary-core-healthy"; return d, nil }
	if b.Healthy { d.ActiveCore, d.StandbyCore, d.FailoverRequired, d.Reason = b.NodeID, a.NodeID, true, "primary-core-unhealthy"; return d, nil }
	d.ActiveCore, d.StandbyCore, d.FailoverRequired, d.Reason = "", "", false, "both-core-nodes-unhealthy"
	return d, nil
}

func NormalizeCoreNetworkEndpoints(endpoints []TunnelEndpoint) []TunnelEndpoint {
	out := append([]TunnelEndpoint(nil), endpoints...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Healthy != out[j].Healthy { return out[i].Healthy }
		if out[i].LatencyMS != out[j].LatencyMS { return out[i].LatencyMS < out[j].LatencyMS }
		return strings.TrimSpace(out[i].ID) < strings.TrimSpace(out[j].ID)
	})
	return out
}
