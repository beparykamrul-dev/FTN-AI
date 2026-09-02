package main

import "testing"

func TestSelectMultiTunnelPaths(t *testing.T) {
	p := MultiTunnelPolicy{
		Protocols: []TunnelProtocol{
			{ID: "wireguard", Enabled: true},
			{ID: "openvpn", Enabled: true},
			{ID: "hysteria2", Enabled: true},
		},
		MaxActivePaths: 2, Failover: true, HealthRequired: true, FTNOwnedOnly: true,
	}
	if err := ValidateMultiTunnelPolicy(p); err != nil { t.Fatal(err) }
	got := SelectMultiTunnelPaths(p, []TunnelEndpoint{
		{ID:"wg-1", Protocol:"wireguard", Transport:"udp", Healthy:true, LatencyMS:20},
		{ID:"ovpn-1", Protocol:"openvpn", Transport:"tcp", Healthy:true, LatencyMS:10},
		{ID:"hy-1", Protocol:"hysteria2", Transport:"udp", Healthy:false, LatencyMS:1},
		{ID:"other", Protocol:"pptp", Transport:"tcp", Healthy:true, LatencyMS:2},
	})
	if len(got) != 2 || got[0].ID != "ovpn-1" || got[1].ID != "wg-1" { t.Fatalf("got=%+v", got) }
}

func TestValidateMultiTunnelPolicyRequiresFTNOwnership(t *testing.T) {
	p := MultiTunnelPolicy{MaxActivePaths: 1}
	if err := ValidateMultiTunnelPolicy(p); err == nil { t.Fatal("expected FTN ownership requirement") }
}
