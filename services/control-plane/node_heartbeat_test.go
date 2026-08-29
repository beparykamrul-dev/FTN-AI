package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterNodeUpsertsInMemory(t *testing.T) {
	nodes = nil
	a := &App{catalog: catalog}
	body := `{"id":"n1","provider":"ftn","region":"ctg","services":["ftndns"],"cpu_percent":20,"ram_percent":30,"ssd_percent":10,"hdd_percent":10,"net_mbps":1000,"latency_ms":5,"packet_loss_percent":0,"healthy":true}`
	for i := 0; i < 2; i++ {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/nodes/register", strings.NewReader(body))
		w := httptest.NewRecorder()
		a.registerNode(w, r)
		if w.Code != http.StatusAccepted { t.Fatalf("status=%d", w.Code) }
	}
	if len(nodes) != 1 { t.Fatalf("expected one node after upsert, got %d", len(nodes)) }
}

func TestRegisterNodeRejectsInvalidMetrics(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodPost, "/api/v1/nodes/register", strings.NewReader(`{"id":"n1","provider":"ftn","cpu_percent":101,"healthy":true}`))
	w := httptest.NewRecorder()
	a.registerNode(w, r)
	if w.Code != http.StatusBadRequest { t.Fatalf("status=%d", w.Code) }
}

func TestPlacementRejectsUnknownService(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodPost, "/api/v1/placement/preview", strings.NewReader(`{"service_id":"unknown"}`))
	w := httptest.NewRecorder()
	a.placement(w, r)
	if w.Code != http.StatusNotFound { t.Fatalf("status=%d", w.Code) }
}
