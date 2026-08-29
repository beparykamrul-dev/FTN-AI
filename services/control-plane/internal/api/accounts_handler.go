package api

import (
	"encoding/json"
	"net/http"
)

type moduleStatus struct {
	Module string `json:"module"`
	Status string `json:"status"`
}

func AccountsStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(NewResponse(moduleStatus{Module: "accounts", Status: "ready"}))
}
