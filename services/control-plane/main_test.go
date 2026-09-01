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
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected rbac denial without context, status=%d", w.Code)
	}
}

func TestEntitlementsIgnoreClientSuppliedServices(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodGet, "/api/v1/entitlements", nil)
	r.Header.Set("X-FTN-Services", "drive,ai,codec")
	w := httptest.NewRecorder()
	a.entitlements(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected rbac denial without context, status=%d", w.Code)
	}
	if strings.Contains(w.Body.String(), `"active":true`) {
		t.Fatal("client supplied service header must never grant entitlements")
	}
}

func TestServiceRequestValidation(t *testing.T) {
	a := &App{catalog: catalog}
	r := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", strings.NewReader(`{"service_id":"not-a-service"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.requests(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected rbac denial before validation, status=%d", w.Code)
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
