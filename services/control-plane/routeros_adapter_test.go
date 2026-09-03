package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testRouterOSAdapter(t *testing.T, baseURL string) *RouterOSAdapter {
	t.Helper()
	a, err := NewRouterOSAdapter(baseURL, RouterOSCredentialRef("secret/router1"), "", func(context.Context, RouterOSCredentialRef) (RouterOSSecret, error) { return RouterOSSecret{Username: "u", Password: "p"}, nil }, time.Second)
	if err != nil { t.Fatal(err) }
	return a
}

func TestNewRouterOSAdapterValidation(t *testing.T) {
	if _, err := NewRouterOSAdapter("", "secret/router1", "", func(context.Context, RouterOSCredentialRef) (RouterOSSecret, error) { return RouterOSSecret{"u", "p"}, nil }, time.Second); err == nil { t.Fatal("expected base URL validation error") }
	if _, err := NewRouterOSAdapter("http://router", "", "", func(context.Context, RouterOSCredentialRef) (RouterOSSecret, error) { return RouterOSSecret{"u", "p"}, nil }, time.Second); err == nil { t.Fatal("expected credential reference validation error") }
	if _, err := NewRouterOSAdapter("http://router", "secret/router1", "", nil, time.Second); err == nil { t.Fatal("expected secret resolver validation error") }
}

func TestRouterOSAdapterRejectsUnverifiedDevice(t *testing.T) {
	adapter := testRouterOSAdapter(t, "http://router")
	if _, err := adapter.Capabilities(context.Background(), NetworkDevice{ID: "r1", Kind: "unknown", Protocol: "routeros-api", Address: "http://router", Healthy: true}); err == nil { t.Fatal("expected ownership validation error") }
}

func TestRouterOSAdapterRejectsInvalidProtocol(t *testing.T) {
	adapter := testRouterOSAdapter(t, "http://router")
	if _, err := adapter.CollectInterfaceState(context.Background(), NetworkDevice{ID: "r1", Kind: "router", Protocol: "snmp", Address: "http://router", Healthy: true}); err == nil { t.Fatal("expected protocol mismatch") }
}

func TestRouterOSAdapterCollectsInterfaceAndRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet { t.Fatalf("unexpected method %s", r.Method) }
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/rest/interface": _, _ = w.Write([]byte(`[{"name":"ether1","disabled":false,"running":true,"rx-byte":"1234","tx-byte":"5678","rx-error":"2","tx-error":"1"}]`))
		case "/rest/ip/route": _, _ = w.Write([]byte(`[{"dst-address":"0.0.0.0/0","protocol":"static","routing-table":"main","gateway":"192.0.2.1","distance":"10","active":true}]`))
		case "/rest/system/resource": _, _ = w.Write([]byte(`[{"version":"7.20"}]`))
		default: http.NotFound(w, r)
		}
	}))
	defer server.Close()
	adapter := testRouterOSAdapter(t, server.URL)
	device := NetworkDevice{ID: "r1", Kind: "router", Protocol: "routeros-api", Address: server.URL, Healthy: true}
	caps, err := adapter.Capabilities(context.Background(), device); if err != nil { t.Fatal(err) }
	if len(caps) != 4 { t.Fatalf("capabilities=%v", caps) }
	interfaces, err := adapter.CollectInterfaceState(context.Background(), device); if err != nil { t.Fatal(err) }
	if len(interfaces) != 1 || interfaces[0].RxBps != 1234 || interfaces[0].RxErrors != 2 || !interfaces[0].OperUp { t.Fatalf("interfaces=%+v", interfaces) }
	routes, err := adapter.CollectRoutingState(context.Background(), device); if err != nil { t.Fatal(err) }
	if len(routes) != 1 || routes[0].Prefix != "0.0.0.0/0" || routes[0].Metric != 10 || !routes[0].Active { t.Fatalf("routes=%+v", routes) }
}

func TestRouterOSAdapterReadQoSSnapshotIsReadOnly(t *testing.T) {
	var mutations int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet { mutations++; t.Fatalf("unexpected mutation method %s", r.Method) }
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/rest/queue/simple" { t.Fatalf("unexpected path %s", r.URL.Path) }
		_, _ = w.Write([]byte(`[{"name":"ftn-qos-whatsapp","comment":"FTN-QOS service=whatsapp class=voice path=path-a dscp=46 priority=90"},{"name":"unmanaged","comment":"operator queue"}]`))
	}))
	defer server.Close()
	adapter := testRouterOSAdapter(t, server.URL)
	device := NetworkDevice{ID:"r1", Kind:"router", Protocol:"routeros-api", Address:server.URL, Healthy:true}
	snapshot, err := adapter.ReadQoSSnapshot(context.Background(), device, true); if err != nil { t.Fatal(err) }
	if snapshot.DeviceID != "r1" || len(snapshot.Rules) != 1 { t.Fatalf("snapshot=%+v", snapshot) }
	if snapshot.Rules[0].ServiceID != "whatsapp" || snapshot.Rules[0].PathID != "path-a" || snapshot.Rules[0].DSCP != 46 || snapshot.Rules[0].Priority != 90 { t.Fatalf("rule=%+v", snapshot.Rules[0]) }
	if mutations != 0 { t.Fatalf("mutation calls=%d", mutations) }
}

func TestRouterOSAdapterReadQoSSnapshotRequiresOwnership(t *testing.T) {
	adapter := testRouterOSAdapter(t, "http://router")
	device := NetworkDevice{ID:"r1", Kind:"router", Protocol:"routeros-api", Address:"http://router", Healthy:true}
	if _, err := adapter.ReadQoSSnapshot(context.Background(), device, false); err == nil { t.Fatal("expected ownership error") }
}

func TestRouterOSAdapterReadQoSSnapshotFailsClosedOnMalformedFTNComment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "application/json"); _, _ = w.Write([]byte(`[{"comment":"FTN-QOS service=whatsapp class=voice path=broken"}]`)) }))
	defer server.Close()
	adapter := testRouterOSAdapter(t, server.URL)
	device := NetworkDevice{ID:"r1", Kind:"router", Protocol:"routeros-api", Address:server.URL, Healthy:true}
	if _, err := adapter.ReadQoSSnapshot(context.Background(), device, true); err == nil { t.Fatal("expected malformed FTN-QOS comment error") }
}
