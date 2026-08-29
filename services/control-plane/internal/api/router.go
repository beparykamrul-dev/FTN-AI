package api

import (
	"net/http"

	"github.com/beparykamrul-dev/FTN-AI/services/control-plane/internal/httpx"
)

func Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", Health)
	return httpx.RequestID(mux)
}
