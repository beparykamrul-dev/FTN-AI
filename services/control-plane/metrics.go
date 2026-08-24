package main

import (
 "fmt"
 "net/http"
 "sync/atomic"
)

var requestCount uint64

func metricsHandler(w http.ResponseWriter, _ *http.Request) {
 w.Header().Set("Content-Type", "text/plain; version=0.0.4")
 fmt.Fprintf(w, "ftn_control_requests_total %d\n", atomic.LoadUint64(&requestCount))
}

func metricsMiddleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  atomic.AddUint64(&requestCount, 1)
  next.ServeHTTP(w, r)
 })
}
