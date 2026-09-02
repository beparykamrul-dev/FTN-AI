package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewRouterOSAdapterValidation(t *testing.T) {
	if _, err := NewRouterOSAdapter("", RouterOSCredentials{Username: "u"}, time.Second); err == nil {
		t.Fatal("expected base URL validation error")
	}
	if _, err := NewRouterOSAdapter("http://router", RouterOSCredentials{}, time.Second); err == nil {
		t.Fatal("expected username validation error")
	}
}

func TestRouterOSAdapterCollectsInterfaceAndRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rest/interface" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"name":"ether1","disabled":false,"running":true,"rx-byte":"1234","tx-byte":"5678"}]`))
			return
		}
		if r.URL.Path == "/rest/ip/route" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"dst-address":"0.0.0.0/0","protocol":"static","routing-table":"main","gateway":"192.0.2.1","distance":"10","active":true}]`))
			return
		}
		if r.URL.Path == "/rest/system/resource" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"version":"7.20"}]`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	adapter, err := NewRouterOSAdapter(server.URL, RouterOSCredentials{Username: "u", Password: "p"}, time.Second)
	if err != nil { t.Fatal(err) }
	device := NetworkDevice{ID: "r1", Protocol: "routeros-api"}

	caps, err := adapter.Capabilities(context.Background(), device)
	if err != nil { t.Fatal(err) }
	if len(caps) != 3 { t.Fatalf("capabilities=%v", caps) }

	interfaces, err := adapter.CollectInterfaceState(context.Background(), device)
	if err != nil { t.Fatal(err) }
	if len(interfaces) != 1 || interfaces[0].RxBps != 1234 || !interfaces[0].OperUp {
		t.Fatalf("interfaces=%+v", interfaces)
	}

	routes, err := adapter.CollectRoutingState(context.Background(), device)
	if err != nil { t.Fatal(err) }
	if len(routes) != 1 || routes[0].Prefix != "0.0.0.0/0" || routes[0].Metric != 10 || !routes[0].Active {
		t.Fatalf("routes=%+v", routes)
	}
}
