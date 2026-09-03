package main

import "testing"

func TestNormalizeFTNNativeProtocolEndpoint(t *testing.T) {
	e, err := NormalizeFTNNativeProtocolEndpoint(FTNNativeProtocolEndpoint{Service:"FTN-API", Protocol:"HTTP3", Transport:"UDP", Authorized:true})
	if err != nil || e.Service != "ftn-api" || e.Protocol != "http3" { t.Fatalf("unexpected normalization: %+v %v", e, err) }
}

func TestNegotiateFTNNativeProtocol(t *testing.T) {
	e, ok := NegotiateFTNNativeProtocol("ftn-api", []string{"http2", "http3"}, []string{"tcp", "udp"}, []FTNNativeProtocolEndpoint{
		{Service:"ftn-api", Protocol:"http3", Transport:"udp", Healthy:true, Authorized:true, LatencyMS:20, LossPct:1, Capacity:80},
		{Service:"ftn-api", Protocol:"http2", Transport:"tcp", Healthy:true, Authorized:true, LatencyMS:10, LossPct:0, Capacity:90},
	})
	if !ok || e.Protocol != "http2" { t.Fatalf("unexpected selection: %+v %v", e, ok) }
}

func TestNegotiateRejectsUnauthorized(t *testing.T) {
	if _, ok := NegotiateFTNNativeProtocol("ftn-dns", []string{"dns"}, []string{"udp"}, []FTNNativeProtocolEndpoint{{Service:"ftn-dns", Protocol:"dns", Transport:"udp", Healthy:true, Authorized:false}}); ok { t.Fatal("expected unauthorized endpoint to be rejected") }
}
