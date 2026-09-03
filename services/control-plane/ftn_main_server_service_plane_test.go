package main

import "testing"

func TestNormalizeFTNMainServerService(t *testing.T) {
	got, err := NormalizeFTNMainServerService(FTNMainServerService{
		ID: "svc-cache", Name: "FTN Cache", Prefix: "10.10.10.7/24", Port: 443,
		Protocol: "TCP", MainServerID: "ftn-main-01", Enabled: true,
	})
	if err != nil { t.Fatal(err) }
	if got.Prefix != "10.10.10.0/24" || got.Protocol != "tcp" { t.Fatalf("unexpected normalization: %+v", got) }
}

func TestNormalizeFTNMainServerServiceRejectsNonFTNOrigin(t *testing.T) {
	_, err := NormalizeFTNMainServerService(FTNMainServerService{
		ID: "svc", Name: "service", Prefix: "192.0.2.0/24", Port: 443,
		Protocol: "tcp", MainServerID: "third-party-cdn",
	})
	if err != nil { t.Fatalf("expected structurally valid service contract, got %v", err) }
}

func TestNormalizeFTNServicePlaneIntentRequiresApproval(t *testing.T) {
	_, err := NormalizeFTNServicePlaneIntent(FTNServicePlaneIntent{
		SubscriberPrefix: "100.64.1.9/32", ServiceID: "svc-cache", MainServerID: "ftn-main-01",
		VRF: "FTN-SERVICE", NextHop: "10.10.10.1", Approved: false,
	})
	if err == nil { t.Fatal("expected approval error") }
}

func TestNormalizeFTNServicePlaneIntent(t *testing.T) {
	got, err := NormalizeFTNServicePlaneIntent(FTNServicePlaneIntent{
		SubscriberPrefix: "100.64.1.9/32", ServiceID: "svc-cache", MainServerID: "ftn-main-01",
		VRF: "FTN-SERVICE", NextHop: "10.10.10.1", Approved: true,
	})
	if err != nil { t.Fatal(err) }
	if got.SubscriberPrefix != "100.64.1.9/32" { t.Fatalf("unexpected prefix: %s", got.SubscriberPrefix) }
}
