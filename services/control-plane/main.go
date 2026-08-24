package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"
)

type Service struct { ID string `json:"id"`; Name string `json:"name"`; Status string `json:"status"`; Platform []string `json:"platform"` }
type API struct { Services []Service `json:"services"` }

func main() {
    addr := os.Getenv("FTN_CONTROL_ADDR"); if addr == "" { addr = ":8080" }
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.Header().Set("Content-Type", "application/json"); w.Write([]byte(`{"status":"ok"}`)) })
    mux.HandleFunc("/api/v1/services", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet { http.Error(w, "method not allowed", 405); return }
        _ = json.NewEncoder(w).Encode(API{Services: []Service{
            {"ftn-internet","FTN Internet","available",[]string{"web","android","pc"}},
            {"ftndns","FTNDNS","available",[]string{"web","android","pc"}},
            {"hosting","FTN Hosting","available",[]string{"web","android","pc"}},
            {"drive","FTN Drive","available",[]string{"web","android","pc"}},
            {"cctv","CCTV Cloud","available",[]string{"web","android","pc"}},
            {"fibermap","FTN FiberMap","available",[]string{"web","android"}},
            {"ai","FTN AI Assistant","available",[]string{"web","android","pc"}},
            {"shop","FTN E-Commerce","available",[]string{"web","android"}},
        }})
    })
    mux.HandleFunc("/api/v1/entitlements", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet { http.Error(w, "method not allowed", 405); return }
        raw := r.Header.Get("X-FTN-Services"); out := []string{}
        for _, s := range strings.Split(raw, ",") { s = strings.TrimSpace(s); if s != "" { out = append(out, s) } }
        _ = json.NewEncoder(w).Encode(map[string]any{"services":out})
    })
    log.Printf("FTN control plane listening on %s", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}
