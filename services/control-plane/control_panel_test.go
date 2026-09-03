package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestControlPanelRoute(t *testing.T) {
    app := &App{catalog: catalog}
    req := httptest.NewRequest(http.MethodGet, "/control-panel", nil)
    rec := httptest.NewRecorder()
    app.controlPanel(rec, req)
    if rec.Code != http.StatusOK { t.Fatalf("status=%d", rec.Code) }
    body := rec.Body.String()
    for _, want := range []string{"FTN Control Panel", "/api/v1/dashboard", "/api/v1/ai/analyze", "approval-first"} {
        if !strings.Contains(body, want) { t.Fatalf("missing %q", want) }
    }
}

func TestControlPanelRejectsNonGET(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/control-panel", nil)
    rec := httptest.NewRecorder()
    (&App{}).controlPanel(rec, req)
    if rec.Code != http.StatusMethodNotAllowed { t.Fatalf("status=%d", rec.Code) }
}
