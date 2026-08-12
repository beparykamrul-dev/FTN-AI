package api

import "net/http"

// RegisterRoutes attaches the SMS API to an existing HTTP mux.
// Authentication is expected to be established by the FTN control-plane middleware;
// the handler performs the operation-level capability check.
func RegisterRoutes(mux *http.ServeMux, h Handler) {
	if mux == nil {
		return
	}
	mux.HandleFunc("POST /v1/sms/messages", h.Send)
}
