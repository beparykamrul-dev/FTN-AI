package main

import "testing"

func TestSelectFTNCoreRouterHAPrefersCurrentHealthyNode(t *testing.T) {
	nodes := []FTNCoreNode{
		{ID:"core-b", Healthy:true, BGPReady:true, BFDState:FTNBFDUp},
		{ID:"core-a", Healthy:true, BGPReady:true, BFDState:FTNBFDUp},
	}
	state, err := SelectFTNCoreRouterHA(nodes, "core-a")
	if err != nil { t.Fatal(err) }
	if state.ActiveNode != "core-a" || state.Failover { t.Fatalf("unexpected state: %+v", state) }
	if state.StandbyNode != "core-b" { t.Fatalf("expected standby core-b: %+v", state) }
}

func TestSelectFTNCoreRouterHAFailsOverWhenCurrentDown(t *testing.T) {
	nodes := []FTNCoreNode{
		{ID:"core-a", Healthy:false, BGPReady:false, BFDState:FTNBFDDown},
		{ID:"core-b", Healthy:true, BGPReady:true, BFDState:FTNBFDUp},
	}
	state, err := SelectFTNCoreRouterHA(nodes, "core-a")
	if err != nil { t.Fatal(err) }
	if state.ActiveNode != "core-b" || !state.Failover { t.Fatalf("expected failover: %+v", state) }
}

func TestSelectFTNCoreRouterHARejectsNoEligibleCore(t *testing.T) {
	_, err := SelectFTNCoreRouterHA([]FTNCoreNode{{ID:"core-a", Healthy:true, BGPReady:false, BFDState:FTNBFDDown}}, "core-a")
	if err == nil { t.Fatal("expected no eligible core error") }
}
