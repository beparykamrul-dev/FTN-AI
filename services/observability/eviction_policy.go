package observability

// EvictionPolicy protects higher-priority telemetry when the durable buffer is full.
type EvictionPolicy struct {
	MaxEvents uint64
	MinPriority int
}

func (p EvictionPolicy) Allows(priority int) bool {
	return p.MaxEvents > 0 && priority >= p.MinPriority
}
