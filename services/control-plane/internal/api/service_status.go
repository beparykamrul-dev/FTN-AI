package api

import (
	"encoding/json"
	"net/http"

	"github.com/beparykamrul-dev/FTN-AI/services/control-plane/internal/services"
)

func ServiceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(NewResponse(services.ServiceStatus{ID: "control-plane", Name: "control-plane", Status: "ok"}))
}
