package observability

import("encoding/json";"net/http")
type API struct{Traffic *TrafficStore}
func(a *API)TrafficSamples(w http.ResponseWriter,r *http.Request){if r==nil||r.Method!=http.MethodGet{http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return};if a==nil||a.Traffic==nil{http.Error(w,"observability unavailable",http.StatusServiceUnavailable);return};w.Header().Set("Content-Type","application/json");if err:=json.NewEncoder(w).Encode(a.Traffic.List());err!=nil{http.Error(w,"encode failed",http.StatusInternalServerError)}}
