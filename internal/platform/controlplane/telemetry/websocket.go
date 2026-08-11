package telemetry

import (
	"encoding/json"
	"sync"
)

// WebSocketHub is transport-neutral session state for the Control Plane's
// live telemetry channel. Authentication and authorization must be enforced
// by the HTTP/WebSocket gateway before a session is registered.
type WebSocketHub struct {
	mu      sync.RWMutex
	sessions map[string]chan []byte
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{sessions: make(map[string]chan []byte)}
}

func (h *WebSocketHub) Register(sessionID string, buffer int) chan []byte {
	if buffer < 1 { buffer = 1 }
	ch := make(chan []byte, buffer)
	h.mu.Lock()
	defer h.mu.Unlock()
	if old, ok := h.sessions[sessionID]; ok { close(old) }
	h.sessions[sessionID] = ch
	return ch
}

func (h *WebSocketHub) Unregister(sessionID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ch, ok := h.sessions[sessionID]; ok { close(ch); delete(h.sessions, sessionID) }
}

func (h *WebSocketHub) BroadcastHeartbeat(hb Heartbeat) int {
	payload, err := json.Marshal(hb)
	if err != nil { return 0 }
	h.mu.RLock()
	defer h.mu.RUnlock()
	sent := 0
	for _, ch := range h.sessions {
		select {
		case ch <- payload:
			sent++
		default:
			// Backpressure: a slow client does not block the fleet telemetry path.
		}
	}
	return sent
}
