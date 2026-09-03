package main

import "testing"

type testRIB struct { installed, withdrawn int }
func (r *testRIB) Install(FTNRoute) error { r.installed++; return nil }
func (r *testRIB) Withdraw(FTNRoute) error { r.withdrawn++; return nil }
func (r *testRIB) Lookup(string, string) ([]FTNRoute, error) { return nil, nil }

type testFIB struct { programmed, removed int }
func (f *testFIB) Program(FTNRoute) error { f.programmed++; return nil }
func (f *testFIB) Remove(FTNRoute) error { f.removed++; return nil }

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

func TestFTNCoreFailover(t *testing.T) {
	nodes := []FTNCoreNode{{ID:"core-a", Healthy:false, BGPReady:false, BFDState:FTNBFDDown}, {ID:"core-b", Healthy:true, BGPReady:true, BFDState:FTNBFDUp}}
	d := EvaluateFTNCoreFailover(nodes, "core-a")
	if d.ActiveNode != "core-b" || !d.Failover { t.Fatalf("unexpected failover: %+v", d) }
}
