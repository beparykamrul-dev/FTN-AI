package main

import (
	"errors"
	"testing"
)

type testRIB struct { installed, withdrawn int }
func (r *testRIB) Install(FTNRoute) error { r.installed++; return nil }
func (r *testRIB) Withdraw(FTNRoute) error { r.withdrawn++; return nil }
func (r *testRIB) Lookup(string, string) ([]FTNRoute, error) { return nil, nil }

type testFIB struct { programmed, removed int }
func (f *testFIB) Program(FTNRoute) error { f.programmed++; return nil }
func (f *testFIB) Remove(FTNRoute) error { f.removed++; return nil }

type testBGP struct { apply, withdraw int; applyErr error; withdrawErr error }
func (a *testBGP) Name() string { return "test-bgp" }
func (a *testBGP) Established() bool { return true }
func (a *testBGP) ApplyRoute(FTNRoute) error { a.apply++; return a.applyErr }
func (a *testBGP) WithdrawRoute(FTNRoute) error { a.withdraw++; return a.withdrawErr }

func TestFTNRoutedEngineRequiresApproval(t *testing.T) {
	r, f := &testRIB{}, &testFIB{}
	e := NewFTNRoutedEngine(r, f)
	in := FTNRouteIntent{Action:"install", Route:FTNRoute{Prefix:"203.0.113.0/24", Protocol:"bgp", Active:true}}
	if err := e.ApplyApprovedRoute(in, false); err == nil { t.Fatal("expected approval error") }
	if r.installed != 0 || f.programmed != 0 { t.Fatal("route changed without approval") }
}

func TestFTNRoutedEngineProgramsRIBAndFIB(t *testing.T) {
	r, f := &testRIB{}, &testFIB{}
	e := NewFTNRoutedEngine(r, f)
	in := FTNRouteIntent{Action:"install", Route:FTNRoute{Prefix:"203.0.113.0/24", Protocol:"bgp", Active:true}}
	if err := e.ApplyApprovedRoute(in, true); err != nil { t.Fatal(err) }
	if r.installed != 1 || f.programmed != 1 { t.Fatalf("rib=%d fib=%d", r.installed, f.programmed) }
}

func TestFTNRoutedEngineRollsBackWhenBGPApplyFails(t *testing.T) {
	r, f := &testRIB{}, &testFIB{}
	bgp := &testBGP{applyErr: errors.New("bgp failure")}
	e := NewFTNRoutedEngine(r, f); e.RegisterBGPAdapter(bgp)
	in := FTNRouteIntent{Action:"install", Route:FTNRoute{Prefix:"203.0.113.0/24", Protocol:"bgp", Active:true}}
	if err := e.ApplyApprovedRoute(in, true); err == nil { t.Fatal("expected BGP failure") }
	if r.installed != 1 || r.withdrawn != 1 { t.Fatalf("rib install=%d withdraw=%d", r.installed, r.withdrawn) }
	if f.programmed != 1 || f.removed != 1 { t.Fatalf("fib program=%d remove=%d", f.programmed, f.removed) }
}

func TestFTNRoutedEngineRollsBackWhenBGPWithdrawFails(t *testing.T) {
	r, f := &testRIB{}, &testFIB{}
	bgp := &testBGP{withdrawErr: errors.New("bgp withdraw failure")}
	e := NewFTNRoutedEngine(r, f); e.RegisterBGPAdapter(bgp)
	in := FTNRouteIntent{Action:"withdraw", Route:FTNRoute{Prefix:"203.0.113.0/24", Protocol:"bgp", Active:true}}
	if err := e.WithdrawApprovedRoute(in, true); err == nil { t.Fatal("expected BGP withdraw failure") }
	if r.installed != 1 || f.programmed != 1 { t.Fatalf("rollback rib=%d fib=%d", r.installed, f.programmed) }
}

func TestFTNCoreFailover(t *testing.T) {
	nodes := []FTNCoreNode{{ID:"core-a", Healthy:false, BGPReady:false, BFDState:FTNBFDDown}, {ID:"core-b", Healthy:true, BGPReady:true, BFDState:FTNBFDUp}}
	d := EvaluateFTNCoreFailover(nodes, "core-a")
	if d.ActiveNode != "core-b" || !d.Failover { t.Fatalf("unexpected failover: %+v", d) }
}
