package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidNode(t *testing.T) {
	base := Node{ID: "node-a", Provider: "provider-a", Services: []string{"ftndns"}, Healthy: true}
	if !validNode(base) {
		t.Fatal("expected valid node")
	}
	base.CPUPercent = 101
	if validNode(base) {
		t.Fatal("expected invalid CPU percentage")
	}
}

func TestRegisterNodeUpsertsInMemory(t *testing.T) {
	old := nodes
	nodes = nil
	defer func() { nodes = old }()
	a := &App{catalog: catalog}
	body := `{"id":"node-a","provider":"provider-a","region":"bd","services":["ftndns"],"healthy":true,"latency_ms":2}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/nodes/register", strings.NewReader(body))
	a.registerNode(w, r)
	if w.Code != http.StatusAccepted || len(nodes) != 1 {
		t.Fatalf("status=%d nodes=%d", w.Code, len(nodes))
	}

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/api/v1/nodes/register", strings.NewReader(`{"id":"node-a","provider":"provider-a","region":"bd","services":["ftndns"],"healthy":true,"latency_ms":1}`))
	a.registerNode(w, r)
	if w.Code != http.StatusAccepted || len(nodes) != 1 || nodes[0].LatencyMs != 1 {
		t.Fatalf("upsert failed: status=%d nodes=%d latency=%v", w.Code, len(nodes), nodes[0].LatencyMs)
	}
}

func TestPlacementChoosesHealthyLowLatencyNode(t *testing.T) {
	old := nodes
	nodes = []Node{
		{ID: "slow", Provider: "p1", Region: "BD", Services: []string{"media"}, CPUPercent: 20, RAMPercent: 20, LatencyMs: 80, Healthy: true},
		{ID: "fast", Provider: "p2", Region: "BD", Services: []string{"media"}, CPUPercent: 30, RAMPercent: 30, LatencyMs: 10, Healthy: true},
		{ID: "down", Provider: "p3", Region: "BD", Services: []string{"media"}, LatencyMs: 1, Healthy: false},
	}
	defer func() { nodes = old }()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/placement/preview", strings.NewReader(`{"service_id":"media","region":"BD"}`))
	w := httptest.NewRecorder()
	(&App{catalog: catalog}).placement(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"node_id":"fast"`) {
		t.Fatalf("fast node not selected: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"execution":"approval_required"`) {
		t.Fatal("placement must remain approval-gated")
	}
}

func TestPlacementRejectsUnknownService(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/v1/placement/preview", strings.NewReader(`{"service_id":"does-not-exist"}`))
	w := httptest.NewRecorder()
	(&App{catalog: catalog}).placement(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status=%d", w.Code)
	}
}

func TestPlacementRejectsUnsupportedServiceOnNodes(t *testing.T) {
	old := nodes
	nodes = []Node{{ID: "n1", Provider: "p1", Services: []string{"ftndns"}, Healthy: true}}
	defer func() { nodes = old }()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/placement/preview", strings.NewReader(`{"service_id":"media"}`))
	w := httptest.NewRecorder()
	(&App{catalog: catalog}).placement(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d", w.Code)
	}
}
