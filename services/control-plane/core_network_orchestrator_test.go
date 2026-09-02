package main

import (
	"testing"
)

func TestPlanCoreNetworkFailoverUsesStandbyWhenPrimaryUnhealthy(t *testing.T) {
	policy := MultiTunnelPolicy{MaxActivePaths: 2, FTNOwnedOnly: true, HealthRequired: true, Failover: true, Protocols: []TunnelProtocol{{ID:"wireguard", Enabled:true},{ID:"hysteria2", Enabled:true}}}
	paths := []TunnelEndpoint{{ID:"wg-a", Protocol:"wireguard", Healthy:true, LatencyMS:10},{ID:"hy-a", Protocol:"hysteria2", Healthy:true, LatencyMS:20}}
	d, err := PlanCoreNetworkFailover(CoreRouterHealth{NodeID:"core-a", Healthy:false}, CoreRouterHealth{NodeID:"core-b", Healthy:true}, policy, paths)
	if err != nil { t.Fatal(err) }
	if d.ActiveCore != "core-b" || !d.FailoverRequired || len(d.Paths) != 2 { t.Fatalf("unexpected decision: %+v", d) }
	if !d.RequiresApproval || !d.SnapshotRequired || !d.VerifyRequired || !d.RollbackWhenSafe { t.Fatalf("safety gates missing: %+v", d) }
}

func TestPlanCoreNetworkFailoverBlocksWhenBothCoresUnhealthy(t *testing.T) {
	policy := MultiTunnelPolicy{MaxActivePaths: 1, FTNOwnedOnly: true, Protocols: []TunnelProtocol{{ID:"wireguard", Enabled:true}}}
	d, err := PlanCoreNetworkFailover(CoreRouterHealth{NodeID:"core-a", Healthy:false}, CoreRouterHealth{NodeID:"core-b", Healthy:false}, policy, nil)
	if err != nil { t.Fatal(err) }
	if d.ActiveCore != "" || d.FailoverRequired || d.Reason != "both-core-nodes-unhealthy" { t.Fatalf("unexpected decision: %+v", d) }
}

func TestNormalizeCoreNetworkEndpoints(t *testing.T) {
	in := []TunnelEndpoint{{ID:"z", Healthy:false, LatencyMS:1},{ID:"b", Healthy:true, LatencyMS:20},{ID:"a", Healthy:true, LatencyMS:10}}
	out := NormalizeCoreNetworkEndpoints(in)
	if out[0].ID != "a" || out[1].ID != "b" || out[2].ID != "z" { t.Fatalf("unexpected order: %+v", out) }
}
