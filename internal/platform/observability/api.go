package observability

import("encoding/json";"net/http")
type API struct{Traffic *TrafficStore}
func(a *API)TrafficSamples(w http.ResponseWriter,r *http.Request){if w==nil||r==nil{return};if r.Method!=http.MethodGet{http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return};if a==nil||a.Traffic==nil{http.Error(w,"observability unavailable",http.StatusServiceUnavailable);return};if err:=r.Context().Err();err!=nil{http.Error(w,"request cancelled",http.StatusRequestTimeout);return};w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");if err:=json.NewEncoder(w).Encode(a.Traffic.List());err!=nil{return}}
