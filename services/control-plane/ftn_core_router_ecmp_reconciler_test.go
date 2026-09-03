package main

import (
	"errors"
	"testing"
)

type ecmpTestRIB struct{ routes []FTNRoute }
func (r *ecmpTestRIB) Install(x FTNRoute) error { r.routes=append(r.routes,x); return nil }
func (r *ecmpTestRIB) Withdraw(x FTNRoute) error { for i,v:=range r.routes { if v.Prefix==x.Prefix && v.NextHop==x.NextHop { r.routes=append(r.routes[:i],r.routes[i+1:]...); return nil } }; return errors.New("route_not_found") }
func (r *ecmpTestRIB) Lookup(prefix,vrf string)([]FTNRoute,error){ return append([]FTNRoute(nil),r.routes...),nil }

type ecmpTestFIB struct{ routes []FTNRoute; failNext bool }
func (f *ecmpTestFIB) Program(x FTNRoute) error { if f.failNext { f.failNext=false; return errors.New("fib_program_failed") }; f.routes=append(f.routes,x); return nil }
func (f *ecmpTestFIB) Remove(x FTNRoute) error { for i,v:=range f.routes { if v.Prefix==x.Prefix && v.NextHop==x.NextHop { f.routes=append(f.routes[:i],f.routes[i+1:]...); return nil } }; return errors.New("fib_route_not_found") }

func TestBuildFTNECMPReconcilePlanBindsApprovalHash(t *testing.T){
	cur:=FTNECMPSelection{Prefix:"203.0.113.0/24",VRF:"default",Paths:[]FTNRouteCandidate{{PathID:"p1",Route:FTNRoute{Prefix:"203.0.113.0/24",NextHop:"192.0.2.1",Protocol:"bgp",VRF:"default"},Healthy:true,BFDState:FTNBFDUp}}}
	des:=FTNECMPSelection{Prefix:"203.0.113.0/24",VRF:"default",Paths:append([]FTNRouteCandidate{},cur.Paths...)}
	des.Paths=append(des.Paths,FTNRouteCandidate{PathID:"p2",Route:FTNRoute{Prefix:"203.0.113.0/24",NextHop:"192.0.2.2",Protocol:"bgp",VRF:"default"},Healthy:true,BFDState:FTNBFDUp})
	plan,err:=BuildFTNECMPReconcilePlan(cur,des); if err!=nil { t.Fatal(err) }
	if len(plan.Install)!=1 || len(plan.Withdraw)!=0 { t.Fatalf("unexpected diff: %+v",plan) }
	if plan.ApprovalHash=="" || plan.ApprovalHash!=HashFTNECMPReconcilePlan(plan) { t.Fatal("approval hash not bound") }
}

func TestApplyApprovedECMPReconcileRequiresExactApprovalHash(t *testing.T){
	e:=NewFTNRoutedEngine(&ecmpTestRIB{},&ecmpTestFIB{})
	plan,err:=BuildFTNECMPReconcilePlan(FTNECMPSelection{Prefix:"203.0.113.0/24",VRF:"default"},FTNECMPSelection{Prefix:"203.0.113.0/24",VRF:"default",Paths:[]FTNRouteCandidate{{PathID:"p1",Route:FTNRoute{Prefix:"203.0.113.0/24",NextHop:"192.0.2.1",Protocol:"bgp",VRF:"default"},Healthy:true,BFDState:FTNBFDUp}}})
	if err!=nil { t.Fatal(err) }
	if err:=e.ApplyApprovedECMPReconcile(plan,"wrong",true); err==nil { t.Fatal("expected approval binding rejection") }
	if err:=e.ApplyApprovedECMPReconcile(plan,plan.ApprovalHash,false); err==nil { t.Fatal("expected approval rejection") }
}

func TestApplyApprovedECMPReconcileProgramsAllPathsAfterApproval(t *testing.T){
	rib:=&ecmpTestRIB{}; fib:=&ecmpTestFIB{}; e:=NewFTNRoutedEngine(rib,fib)
	plan,err:=BuildFTNECMPReconcilePlan(FTNECMPSelection{Prefix:"203.0.113.0/24",VRF:"default"},FTNECMPSelection{Prefix:"203.0.113.0/24",VRF:"default",Paths:[]FTNRouteCandidate{
		{PathID:"p1",Route:FTNRoute{Prefix:"203.0.113.0/24",NextHop:"192.0.2.1",Protocol:"bgp",VRF:"default"},Healthy:true,BFDState:FTNBFDUp},
		{PathID:"p2",Route:FTNRoute{Prefix:"203.0.113.0/24",NextHop:"192.0.2.2",Protocol:"bgp",VRF:"default"},Healthy:true,BFDState:FTNBFDUp},
	}})
	if err!=nil { t.Fatal(err) }
	if err:=e.ApplyApprovedECMPReconcile(plan,plan.ApprovalHash,true); err!=nil { t.Fatal(err) }
	if len(rib.routes)!=2 || len(fib.routes)!=2 { t.Fatalf("expected two programmed paths, got rib=%d fib=%d",len(rib.routes),len(fib.routes)) }
}

func TestApplyApprovedECMPReconcileRollsBackOnFailure(t *testing.T){
	rib:=&ecmpTestRIB{}; fib:=&ecmpTestFIB{failNext:true}; e:=NewFTNRoutedEngine(rib,fib)
	plan,err:=BuildFTNECMPReconcilePlan(FTNECMPSelection{Prefix:"203.0.113.0/24",VRF:"default"},FTNECMPSelection{Prefix:"203.0.113.0/24",VRF:"default",Paths:[]FTNRouteCandidate{{PathID:"p1",Route:FTNRoute{Prefix:"203.0.113.0/24",NextHop:"192.0.2.1",Protocol:"bgp",VRF:"default"},Healthy:true,BFDState:FTNBFDUp}}})
	if err!=nil { t.Fatal(err) }
	if err:=e.ApplyApprovedECMPReconcile(plan,plan.ApprovalHash,true); err==nil { t.Fatal("expected programming failure") }
	if len(rib.routes)!=0 { t.Fatalf("rollback left RIB state: %+v",rib.routes) }
}
