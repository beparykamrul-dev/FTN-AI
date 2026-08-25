package main

import (
 "fmt"
 "net/http"
 "sync/atomic"
)

func metricsHandler(w http.ResponseWriter, _ *http.Request) {
 w.Header().Set("Content-Type", "text/plain; version=0.0.4")
 fmt.Fprintf(w, "ftn_control_requests_total %d\n", atomic.LoadUint64(&requestCount))
}
