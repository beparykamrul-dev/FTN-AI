package api

import (
	"encoding/json"
	"net/http"
)

type readiness struct {
	Status string `json:"status"`
	Ready  bool   `json:"ready"`
}

func Readiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(NewResponse(readiness{Status: "ready", Ready: true}))
}
