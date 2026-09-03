package main

import "testing"

func TestNormalizeFTNRoute(t *testing.T) {
	r, err := NormalizeFTNRoute(FTNRoute{Prefix:"10.0.0.1/24", NextHop:"192.0.2.1", Protocol:"BGP"})
	if err != nil { t.Fatal(err) }
	if r.Prefix != "10.0.0.0/24" || r.VRF != "default" || r.Protocol != "bgp" { t.Fatalf("unexpected route: %+v", r) }
}

func TestEvaluateFTNRouteIntent(t *testing.T) {
	in := FTNRouteIntent{Action:"install", Route:FTNRoute{Prefix:"203.0.113.0/24", NextHop:"192.0.2.1", Protocol:"bgp", Active:true}}
	d := EvaluateFTNRouteIntent(in)
	if !d.Allowed || !d.RequiresApproval || d.DecisionHash == "" { t.Fatalf("unexpected decision: %+v", d) }
	in.Action = "delete"
	d = EvaluateFTNRouteIntent(in)
	if d.Allowed || !d.RequiresApproval || d.Risk != "critical" { t.Fatalf("unexpected destructive decision: %+v", d) }
}

func TestFTNRouteHashDeterministic(t *testing.T) {
	in := FTNRouteIntent{Action:"install", Route:FTNRoute{Prefix:"2001:db8::/32", Protocol:"ibgp", Community:[]string{"64512:20", "64512:10"}}}
	if a, b := HashFTNRouteIntent(in), HashFTNRouteIntent(in); a != b { t.Fatalf("hash changed: %s != %s", a, b) }
}
