package controlplane

import "encoding/json"

type HealthResponse struct {
	Status string `json:"status"`
	API    string `json:"api"`
}

func HealthJSON() []byte {
	b, _ := json.Marshal(HealthResponse{Status: "ok", API: "v1"})
	return b
}

// EventEnvelope is the transport-neutral shape used by WebSocket gateways.
type EventEnvelope struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}
