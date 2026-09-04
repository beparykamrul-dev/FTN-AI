package main

import (
    "context"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"
)

func TestClickHouseStoreFromEnvDisabledWithoutURL(t *testing.T) {
    t.Setenv("FTN_CLICKHOUSE_URL", "")
    store, err := NewFTNClickHouseStoreFromEnv()
    if err != nil { t.Fatal(err) }
    if store != nil { t.Fatal("expected ClickHouse store to be disabled without URL") }
}

func TestClickHouseStoreFromEnvRejectsInvalidURL(t *testing.T) {
    t.Setenv("FTN_CLICKHOUSE_URL", "ftp://clickhouse:8123")
    if _, err := NewFTNClickHouseStoreFromEnv(); err == nil { t.Fatal("expected invalid URL error") }
}

func TestClickHouseInsertBatchUsesJSONEachRow(t *testing.T) {
    var gotQuery, gotBody string
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        gotQuery = r.URL.Query().Get("query")
        body, _ := io.ReadAll(r.Body)
        gotBody = string(body)
        w.WriteHeader(http.StatusNoContent)
    }))
    defer server.Close()

    t.Setenv("FTN_CLICKHOUSE_URL", server.URL)
    t.Setenv("FTN_CLICKHOUSE_DATABASE", "ftn_analytics")
    t.Setenv("FTN_CLICKHOUSE_USER", "")
    t.Setenv("FTN_CLICKHOUSE_PASSWORD", "")
    store, err := NewFTNClickHouseStoreFromEnv()
    if err != nil { t.Fatal(err) }

    record := FlowRecord{ExporterID: "r1", Version: 9, SourceIP: "10.0.0.1", DestinationIP: "10.0.0.2", SourcePort: 1234, DestinationPort: 443, Protocol: 6, Bytes: 100, Packets: 2, SamplingRate: 1}
    if err := store.InsertBatch(context.Background(), time.Unix(100, 0).UTC(), "tenant-a", []FlowRecord{record}, "youtube", "", ""); err != nil { t.Fatal(err) }
    if !strings.Contains(gotQuery, "flow_records FORMAT JSONEachRow") { t.Fatalf("unexpected query: %s", gotQuery) }
    if !strings.Contains(gotBody, `"tenant_id":"tenant-a"`) || !strings.Contains(gotBody, `"service_id":"youtube"`) { t.Fatalf("unexpected JSONEachRow body: %s", gotBody) }
}

func TestClickHouseInsertBatchRejectsInvalidRecord(t *testing.T) {
    t.Setenv("FTN_CLICKHOUSE_URL", "http://127.0.0.1:8123")
    store, err := NewFTNClickHouseStoreFromEnv()
    if err != nil { t.Fatal(err) }
    err = store.InsertBatch(context.Background(), time.Now(), "tenant-a", []FlowRecord{{ExporterID: "r1", Version: 7}}, "", "", "")
    if err == nil { t.Fatal("expected unsupported flow version error") }
}
