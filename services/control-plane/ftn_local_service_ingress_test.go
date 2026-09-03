package main

import "testing"

func TestNormalizeFTNLocalServiceIngress(t *testing.T) {
	v, err := NormalizeFTNLocalServiceIngress(FTNLocalServiceIngress{
		RouterID:"router-1", Tunnel:"WireGuard", SourcePrefix:"192.0.2.7/24", ServicePrefix:"198.51.100.9/24", Service:"ftn-cache", SILKExport:true, FirewallEnforced:true, Approved:true,
	})
	if err != nil { t.Fatal(err) }
	if v.SourcePrefix != "192.0.2.0/24" || v.ServicePrefix != "198.51.100.0/24" || v.Tunnel != "wireguard" { t.Fatalf("ingress=%+v", v) }
	if FTNLocalIngressHash(v) == "" { t.Fatal("expected ingress hash") }
}

func TestNormalizeFTNLocalServiceIngressRequiresFirewall(t *testing.T) {
	_, err := NormalizeFTNLocalServiceIngress(FTNLocalServiceIngress{RouterID:"router-1", Tunnel:"gre", Service:"ftn-cache", Approved:true})
	if err == nil { t.Fatal("expected firewall requirement") }
}
