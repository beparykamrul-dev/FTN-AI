package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServicesCatalog(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
	w := httptest.NewRecorder()
	a.serviceCatalog(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "FTN Internet") {
		t.Fatal("catalog missing FTN Internet")
	}
	if !strings.Contains(body, "FTN AI Assistant") {
		t.Fatal("catalog missing AI")
	}
	if !strings.Contains(body, "FTN Codec Fabric") {
		t.Fatal("catalog missing codec fabric")
	}
	if !strings.Contains(body, "FTN E2E Transfer") {
		t.Fatal("catalog missing E2E transfer")
	}
}

func TestEntitlementsAreServiceScoped(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodGet, "/api/v1/entitlements", nil)
	r.Header.Set("X-FTN-Services", "drive,ai,codec")
	w := httptest.NewRecorder()
	a.entitlements(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `"service_id":"drive","active":true`) {
		t.Fatal("drive entitlement missing")
	}
	if !strings.Contains(body, `"service_id":"codec","active":true`) {
		t.Fatal("codec entitlement missing")
	}
	if !strings.Contains(body, `"service_id":"hosting","active":false`) {
		t.Fatal("hosting must remain inactive")
	}
}

func TestServiceRequestValidation(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", strings.NewReader(`{"service_id":"not-a-service"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.requests(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status=%d", w.Code)
	}
}

func TestRoutesHealthAndSecurityHeaders(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	a.routes().ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("healthz status=%d", w.Code)
	}
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("security headers missing")
	}
	if w.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatal("frame protection missing")
	}
}
