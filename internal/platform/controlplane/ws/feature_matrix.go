package ws

// Feature describes a capability selected for FTN WebSocket after comparing
// mature real-time systems and protocol requirements. It is intentionally a
// product-level contract, not a dependency on any third-party server.
type Feature struct {
	Name     string `json:"name"`
	Priority string `json:"priority"`
	Status   string `json:"status"`
	Reason   string `json:"reason"`
}

func FeatureMatrix() []Feature {
	return []Feature{
		{Name: "RFC6455 framing", Priority: "P0", Status: "required", Reason: "Interoperable WebSocket transport foundation."},
		{Name: "TLS/WSS", Priority: "P0", Status: "required", Reason: "Encrypted transport for control and monitoring traffic."},
		{Name: "Heartbeat and liveness", Priority: "P0", Status: "planned", Reason: "Detect dead peers without relying on application traffic."},
		{Name: "Bounded backpressure", Priority: "P0", Status: "implemented", Reason: "Prevent slow consumers from blocking the fabric."},
		{Name: "Topic subscriptions", Priority: "P0", Status: "implemented", Reason: "Multiplex many control and monitoring streams over one connection."},
		{Name: "Ordered sequence IDs", Priority: "P0", Status: "implemented", Reason: "Enable gap detection and resumable streams."},
		{Name: "Reconnect and stream recovery", Priority: "P0", Status: "planned", Reason: "Recover short disconnect gaps without a database thundering herd."},
		{Name: "Presence and connection state", Priority: "P1", Status: "planned", Reason: "Expose node, operator and client liveness to the control plane."},
		{Name: "Request/response correlation", Priority: "P1", Status: "planned", Reason: "Support commands, acknowledgements and deterministic client workflows."},
		{Name: "Idempotency keys", Priority: "P1", Status: "planned", Reason: "Prevent duplicate execution after retries."},
		{Name: "Per-topic authorization", Priority: "P0", Status: "planned", Reason: "Least-privilege access to control, monitoring and audit streams."},
		{Name: "Connection quotas", Priority: "P1", Status: "planned", Reason: "Protect the fabric from resource exhaustion."},
		{Name: "Payload limits", Priority: "P0", Status: "planned", Reason: "Bound memory and parsing risk for untrusted clients."},
		{Name: "Optional per-message compression", Priority: "P2", Status: "planned", Reason: "Use only after measurement; compression has CPU and memory trade-offs."},
		{Name: "Binary codec", Priority: "P2", Status: "planned", Reason: "Reduce overhead for high-volume telemetry after JSON compatibility is stable."},
		{Name: "Multi-node federation", Priority: "P0", Status: "planned", Reason: "Keep WebSocket service available across FTN nodes."},
		{Name: "Event durability boundary", Priority: "P0", Status: "planned", Reason: "Separate realtime delivery from durable event storage."},
		{Name: "Observability", Priority: "P0", Status: "planned", Reason: "Measure connection count, latency, drops, queue depth and recovery."},
	}
}
