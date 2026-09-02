package router

import (
	"context"
	"testing"
)

type testPlane struct{ applied, deleted []Route; ready bool }
func (p *testPlane) Name() PacketPlane { return PlaneKernel }
func (p *testPlane) Ready(context.Context) bool { return p.ready }
func (p *testPlane) Interfaces(context.Context) ([]Interface,error) { return nil,nil }
func (p *testPlane) Routes(context.Context) ([]Route,error) { return nil,nil }
func (p *testPlane) ApplyRoute(_ context.Context,r Route) error { p.applied=append(p.applied,r); return nil }
func (p *testPlane) DeleteRoute(_ context.Context,r Route) error { p.deleted=append(p.deleted,r); return nil }

func TestRoutePlannerPlanAndApply(t *testing.T) {
	p := &testPlane{ready:true}; rp := RoutePlanner{Plane:p}
	r := Route{Prefix:"10.20.0.0/16",NextHop:"10.20.0.1",Metric:100,Protocol:"static"}
	plan, err := rp.PlanRouteChange(context.Background(),r); if err != nil { t.Fatal(err) }
	if plan == "" { t.Fatal("empty plan") }
	if err := rp.ApplyApprovedPlan(context.Background(),plan); err != nil { t.Fatal(err) }
	if len(p.applied)!=1 || p.applied[0].Prefix!=r.Prefix { t.Fatalf("unexpected applied route: %+v",p.applied) }
	if err := rp.Rollback(context.Background(),plan); err != nil { t.Fatal(err) }
	if len(p.deleted)!=1 { t.Fatalf("expected rollback") }
}

func TestRoutePlannerRejectsInvalidRoute(t *testing.T) {
	p := &testPlane{ready:true}; rp := RoutePlanner{Plane:p}
	if _,err:=rp.PlanRouteChange(context.Background(),Route{Prefix:"not-a-prefix"}); err==nil { t.Fatal("expected invalid prefix error") }
	if _,err:=rp.PlanRouteChange(context.Background(),Route{Prefix:"10.0.0.0/24",NextHop:"2001:db8::1"}); err==nil { t.Fatal("expected address-family error") }
}
