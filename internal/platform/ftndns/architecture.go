package ftndns

// FTNDNS is the DNS-focused platform boundary. Providers remain adapters;
// FTNDNS owns normalized state, consistency, mesh coordination and policy.
type ProviderID string

type ProviderState struct {
	ID ProviderID `json:"id"`
	Healthy bool `json:"healthy"`
	LatencyMS int64 `json:"latency_ms"`
	SnapshotHash string `json:"snapshot_hash"`
}

type MeshState struct {
	Enabled bool `json:"enabled"`
	Nodes int `json:"nodes"`
	HealthyNodes int `json:"healthy_nodes"`
}

type ConsistencyState struct {
	Consistent bool `json:"consistent"`
	DriftCount int `json:"drift_count"`
}

type State struct {
	Providers []ProviderState `json:"providers"`
	Mesh MeshState `json:"mesh"`
	Consistency ConsistencyState `json:"consistency"`
}
