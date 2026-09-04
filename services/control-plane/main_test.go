package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServicesCatalogRequiresAuthorization(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
	w := httptest.NewRecorder()
	a.serviceCatalog(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected protected catalog to reject missing authorization context, status=%d", w.Code)
	}
}

func TestEntitlementsIgnoreClientSuppliedServices(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodGet, "/api/v1/entitlements", nil)
	r.Header.Set("X-FTN-Services", "drive,ai,codec")
	w := httptest.NewRecorder()
	a.entitlements(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected protected entitlement endpoint to reject missing authorization context, status=%d", w.Code)
	}
	if strings.Contains(w.Body.String(), `"active":true`) {
		t.Fatal("client supplied service header must never grant entitlements")
	}
}

func TestServiceRequestValidationRequiresAuthorization(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", strings.NewReader(`{"service_id":"not-a-service"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.requests(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected protected service request endpoint to reject missing authorization context, status=%d", w.Code)
	}
}

func TestJSONResponseAlwaysEmitsValidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	jsonResponse(w, http.StatusAccepted, map[string]any{
		"status":        "accepted",
		"service_id":    "internet",
		"firmware_push": false,
		"message":       "request accepted; device changes require authorization",
	})
	if w.Code != http.StatusAccepted {
		t.Fatalf("status=%d", w.Code)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type=%q", got)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not valid JSON: %v; body=%q", err, w.Body.String())
	}
	if body["status"] != "accepted" || body["service_id"] != "internet" {
		t.Fatalf("unexpected response body: %#v", body)
	}
}

func TestRequestContextMiddlewareRequiresDatabaseIdentity(t *testing.T) {
	a := &App{}
	r := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
	w := httptest.NewRecorder()
	h := requestContextMiddleware(a, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("protected handler must not execute without DB identity")
	}))
	h.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", w.Code)
	}
	if w.Header().Get("X-Request-ID") == "" || w.Header().Get("X-Correlation-ID") == "" {
		t.Fatal("request/correlation ids must be emitted")
	}
}
