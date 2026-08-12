package ws

// DependencyPolicy makes the FTN WebSocket runtime explicit: the realtime
// control path must not require a hosted broker, SaaS gateway, or third-party
// realtime service. External systems may be integrated only as optional
// adapters outside the core fabric.
type DependencyPolicy struct {
	CoreTransport       string   `json:"core_transport"`
	CoreState           string   `json:"core_state"`
	CoreEventBus        string   `json:"core_event_bus"`
	AllowedExternalDeps []string `json:"allowed_external_deps,omitempty"`
	ForbiddenCoreDeps   []string `json:"forbidden_core_deps"`
}

func DefaultDependencyPolicy() DependencyPolicy {
	return DependencyPolicy{
		CoreTransport: "self-hosted-rfc6455",
		CoreState: "ftn-controlled-state",
		CoreEventBus: "ftn-event-fabric",
		ForbiddenCoreDeps: []string{
			"hosted-realtime-service",
			"external-pubsub-broker",
			"third-party-websocket-gateway",
			"vendor-managed-session-store",
		},
	}
}
