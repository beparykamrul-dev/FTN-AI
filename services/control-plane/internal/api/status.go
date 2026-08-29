package api

import (
	"encoding/json"
	"net/http"
)

type AggregateStatus struct {
	Status string `json:"status"`
	Services []string `json:"services"`
}

func AggregateServiceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(NewResponse(AggregateStatus{
		Status: "ok",
		Services: []string{"accounts", "billing", "noc", "control-plane"},
	}))
}
