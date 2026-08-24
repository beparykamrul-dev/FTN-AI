package main

import (
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
)

func TestServicesCatalog(t *testing.T) {
  a := &App{services: catalog}
  r := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
  w := httptest.NewRecorder()
  a.services(w, r)
  if w.Code != http.StatusOK { t.Fatalf("status=%d", w.Code) }
  if !strings.Contains(w.Body.String(), "FTN Internet") { t.Fatal("catalog missing FTN Internet") }
  if !strings.Contains(w.Body.String(), "FTN AI Assistant") { t.Fatal("catalog missing AI") }
}

func TestEntitlementsAreServiceScoped(t *testing.T) {
  a := &App{services: catalog}
  r := httptest.NewRequest(http.MethodGet, "/api/v1/entitlements", nil)
  r.Header.Set("X-FTN-Services", "drive,ai")
  w := httptest.NewRecorder()
  a.entitlements(w, r)
  if w.Code != http.StatusOK { t.Fatalf("status=%d", w.Code) }
  body := w.Body.String()
  if !strings.Contains(body, `"service_id":"drive","active":true`) { t.Fatal("drive entitlement missing") }
  if !strings.Contains(body, `"service_id":"hosting","active":false`) { t.Fatal("hosting must remain inactive") }
}

func TestServiceRequestValidation(t *testing.T) {
  a := &App{services: catalog}
  r := httptest.NewRequest(http.MethodPost, "/api/v1/service-requests", strings.NewReader(`{"service_id":"not-a-service"}`))
  r.Header.Set("Content-Type", "application/json")
  w := httptest.NewRecorder()
  a.requests(w, r)
  if w.Code != http.StatusNotFound { t.Fatalf("status=%d", w.Code) }
}
