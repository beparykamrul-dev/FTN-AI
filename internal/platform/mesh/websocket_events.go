package mesh

import (
	"encoding/json"
	"net/http"
)

// WebSocketEventEncoder keeps mesh event serialization independent from any
// particular WebSocket library. A gateway can use EncodeEvent as its wire payload.
func EncodeEvent(e Event) ([]byte, error) { return json.Marshal(e) }

// EventStreamHandler provides a safe HTTP capability check for the future
// authenticated WebSocket upgrade layer. Actual upgrades belong in the gateway.
func EventStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte(`{"error":"websocket upgrade is provided by the authenticated FTN gateway"}`))
}
