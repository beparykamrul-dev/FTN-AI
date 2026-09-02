package main

import (
	"context"
	"testing"
)

type integrationAdapter struct{}
func (integrationAdapter) Protocol() string { return "gobgp" }
func (integrationAdapter) Capabilities(context.Context, CoreRouterNode) ([]string,error) { return []string{"health","peers","route-plan"},nil }
func (integrationAdapter) Health(context.Context, CoreRouterNode) (CoreRouterHealth,error) { return CoreRouterHealth{NodeID:"core-a",Healthy:true},nil }
func (integrationAdapter) Peers(context.Context, CoreRouterNode) ([]CoreRouterPeerState,error) { return []CoreRouterPeerState{{PeerID:"upstream-1",Established:true}},nil }
func (integrationAdapter) PlanRouteChange(_ context.Context, _ CoreRouterNode, i RouteChangeIntent) (RouteChangePlan,error) { return PlanCoreRouteChange(CoreRouterNode{ID:"core-a",Role:"primary",Protocol:"gobgp",Enabled:true},i),nil }

func TestCoreRouterIntegrationInspectAndPlan(t *testing.T) {
	i := CoreRouterIntegration{Adapter:integrationAdapter{}}
	h, peers, err := i.Inspect(context.Background(), CoreRouterNode{ID:"core-a",Role:"primary",Protocol:"gobgp",Enabled:true})
	if err != nil || !h.Healthy || len(peers)!=1 { t.Fatalf("health=%+v peers=%+v err=%v",h,peers,err) }
	p, err := i.Plan(context.Background(), CoreRouterNode{ID:"core-a",Role:"primary",Protocol:"gobgp",Enabled:true}, RouteChangeIntent{Action:"announce",Prefix:"203.0.113.0/24",ApprovalID:"appr-1"})
	if err != nil || !p.Allowed { t.Fatalf("plan=%+v err=%v",p,err) }
	if err := RequireApprovedRoutePlan(p,"appr-1"); err != nil { t.Fatal(err) }
}

func TestRequireApprovedRoutePlanRejectsMissingApproval(t *testing.T) {
	p := RouteChangePlan{Allowed:true,RequiresApproval:true,PreChangeSnapshot:true,PostChangeVerify:true,RollbackWhenSafe:true}
	if err := RequireApprovedRoutePlan(p,""); err == nil { t.Fatal("expected approval requirement") }
}
