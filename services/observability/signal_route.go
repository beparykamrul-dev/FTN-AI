package observability

// SignalRoute selects the preferred backend for one telemetry signal.
type SignalRoute struct {
	Signal string
	Backend string
	Priority uint32
	Healthy bool
}

func (r SignalRoute) Eligible() bool { return r.Signal != "" && r.Backend != "" && r.Healthy }

func SelectSignalRoute(routes []SignalRoute, signal string) (SignalRoute, bool) {
	var best SignalRoute
	found := false
	for _, r := range routes {
		if r.Signal != signal || !r.Eligible() { continue }
		if !found || r.Priority < best.Priority || (r.Priority == best.Priority && r.Backend < best.Backend) {
			best, found = r, true
		}
	}
	return best, found
}
