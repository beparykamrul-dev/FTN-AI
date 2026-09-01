package controlplane

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestRemoteEventJournalAppend(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/api/v1/events/append" || r.Header.Get("Authorization") != "Bearer secret" {
            t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusAccepted)
        _, _ = w.Write([]byte(`{"id":"e1","tenant_id":"t1","type":"job.created","sequence":7,"correlation_id":"c1","payload":{"job":"j1"},"created_at":"2026-08-30T00:00:00Z"}`))
    }))
    defer srv.Close()
    j := NewRemoteEventJournal(srv.URL, "secret")
    e, err := j.Append(JournalEvent{TenantID: "t1", Type: "job.created", CorrelationID: "c1", Payload: []byte(`{"job":"j1"}`)})
    if err != nil { t.Fatal(err) }
    if e.ID != "e1" || e.Sequence != 7 || e.TenantID != "t1" { t.Fatalf("unexpected event: %+v", e) }
}

func TestRemoteEventJournalOffsetDefaultsToZero(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) }))
    defer srv.Close()
    j := NewRemoteEventJournal(srv.URL, "secret")
    if got := j.Offset("worker", "t1"); got != 0 { t.Fatalf("expected zero offset, got %d", got) }
}
