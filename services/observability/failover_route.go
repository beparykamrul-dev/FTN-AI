package observability

// FailoverRoute provides an ordered fallback path for telemetry backends.
type FailoverRoute struct {
	Signal string
	Backends []string
}

func (r FailoverRoute) Valid() bool {
	return r.Signal != "" && len(r.Backends) > 0
}

func (r FailoverRoute) Next(current string) (string, bool) {
	for i, name := range r.Backends {
		if name == current && i+1 < len(r.Backends) {
			return r.Backends[i+1], true
		}
	}
	if current == "" && len(r.Backends) > 0 {
		return r.Backends[0], true
	}
	return "", false
}
