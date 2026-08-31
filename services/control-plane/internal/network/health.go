package network

type Health struct {
	DeviceID string `json:"device_id"`
	Status   string `json:"status"`
	LatencyMS int64 `json:"latency_ms"`
	PacketLoss float64 `json:"packet_loss"`
}
