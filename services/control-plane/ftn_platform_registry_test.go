package main

import "testing"

func TestNormalizeFTNPlatformRegistry(t *testing.T) {
	r, err := NormalizeFTNPlatformRegistry(FTNPlatformRegistry{Services: []FTNPlatformService{{ID:"FTN-DNS", Protocols:[]string{"DoQ","DNS"}, Transports:[]string{"UDP","TCP"}, Enabled:true}}})
	if err != nil || len(r.Services) != 1 || r.Services[0].ID != "ftn-dns" { t.Fatalf("unexpected registry: %+v %v", r, err) }
}

func TestNormalizeFTNPlatformRegistryRejectsMissingCapabilities(t *testing.T) {
	if _, err := NormalizeFTNPlatformRegistry(FTNPlatformRegistry{Services: []FTNPlatformService{{ID:"ftn-api", Enabled:true}}}); err == nil { t.Fatal("expected missing capabilities error") }
}
