package api

import (
	"encoding/json"
	"net/http"
)

type Metrics struct {
	Requests uint64 `json:"requests"`
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(NewResponse(Metrics{}))
}
