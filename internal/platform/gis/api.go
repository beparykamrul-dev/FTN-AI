package gis

import (
	"encoding/json"
	"net/http"
)

type API struct {
	IPAM *IPAM
	Map  *MapStore
}

func (a *API) IPAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if a == nil || a.IPAM == nil {
		http.Error(w, "GIS IPAM not configured", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(a.IPAM.List())
}

func (a *API) MapSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if a == nil || a.Map == nil {
		http.Error(w, "GIS map not configured", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(a.Map.Snapshot())
}
