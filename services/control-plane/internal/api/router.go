package api

import (
	"net/http"

	"github.com/beparykamrul-dev/FTN-AI/services/control-plane/internal/httpx"
)

func Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", Health)
	mux.HandleFunc("GET /ready", Readiness)
	mux.HandleFunc("GET /version", Version)
	mux.HandleFunc("GET /api/v1/accounts/status", AccountsStatus)
	mux.HandleFunc("GET /api/v1/billing/status", BillingStatus)
	mux.HandleFunc("GET /api/v1/noc/status", NOCStatus)
	return httpx.RequestID(mux)
}
