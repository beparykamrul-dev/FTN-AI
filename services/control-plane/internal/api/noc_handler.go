package api

import (
	"encoding/json"
	"net/http"
)

func NOCStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(NewResponse(moduleStatus{Module: "noc", Status: "ready"}))
}
