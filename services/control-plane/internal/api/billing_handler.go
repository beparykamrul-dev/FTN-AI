package api

import (
	"encoding/json"
	"net/http"
)

func BillingStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(NewResponse(moduleStatus{Module: "billing", Status: "ready"}))
}
