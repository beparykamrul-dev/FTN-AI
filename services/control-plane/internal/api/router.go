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
	mux.Handle("GET /api/v1/accounts/status", ProtectedAdmin(http.HandlerFunc(AccountsStatus)))
	mux.Handle("GET /api/v1/billing/status", ProtectedBilling(http.HandlerFunc(BillingStatus)))
	mux.Handle("GET /api/v1/noc/status", ProtectedNOC(http.HandlerFunc(NOCStatus)))
	mux.Handle("GET /api/v1/service/status", ProtectedAdmin(http.HandlerFunc(ServiceStatus)))
	return httpx.RequestID(mux)
}
