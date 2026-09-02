package main

import "testing"

func TestPlanCoreRouteChange(t *testing.T) {
	n := CoreRouterNode{ID:"core-a", Role:"primary", Protocol:"gobgp", Enabled:true}
	p := PlanCoreRouteChange(n, RouteChangeIntent{Action:"announce", Prefix:"203.0.113.0/24"})
	if p.Allowed || len(p.ValidationErrors)==0 || p.ValidationErrors[0]!="approval_required" { t.Fatalf("plan=%+v", p) }
	p = PlanCoreRouteChange(n, RouteChangeIntent{Action:"announce", Prefix:"203.0.113.0/24", ApprovalID:"chg-1"})
	if !p.Allowed || !p.RequiresApproval || !p.PreChangeSnapshot || !p.PostChangeVerify { t.Fatalf("plan=%+v", p) }
	p = PlanCoreRouteChange(n, RouteChangeIntent{Action:"withdraw", Prefix:"203.0.113.0/24", ApprovalID:"chg-2"})
	if p.Risk!="high" { t.Fatalf("risk=%q", p.Risk) }
}
