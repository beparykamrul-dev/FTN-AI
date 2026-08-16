package observability

// BackendCapability describes which telemetry signals a backend is suitable for.
type BackendCapability struct {
	Backend string
	Signals []string
	Healthy bool
	Priority uint32
}

func (b BackendCapability) Supports(signal string) bool {
	if !b.Healthy || b.Backend == "" { return false }
	for _, s := range b.Signals { if s == signal { return true } }
	return false
}
