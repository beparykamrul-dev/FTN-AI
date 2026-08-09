package observability

import (
	"encoding/json"
	"net/http"
)

type API struct { Traffic *TrafficStore }

func (a *API) TrafficSamples(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(a.Traffic.List())
}
