package observability

// FailoverRoute provides an ordered fallback path for telemetry backends.
type FailoverRoute struct {
	Signal string
	Backends []string
}

func (r FailoverRoute) Valid() bool {
	return r.Signal != "" && len(r.Backends) > 0
}
