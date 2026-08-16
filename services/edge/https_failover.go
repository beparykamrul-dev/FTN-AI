package edge

// HTTPSRoute describes an approved HTTPS edge route.
type HTTPSRoute struct {
	Provider string
	Endpoint string
	Healthy  bool
	LatencyMS uint32
}

// SelectHTTPSRoute chooses a healthy approved route with lowest latency.
func SelectHTTPSRoute(routes []HTTPSRoute) (HTTPSRoute, bool) {
	var best HTTPSRoute
	found := false
	for _, r := range routes {
		if !r.Healthy || r.Provider == "" || r.Endpoint == "" { continue }
		if !found || r.LatencyMS < best.LatencyMS ||
			(r.LatencyMS == best.LatencyMS && r.Endpoint < best.Endpoint) {
			best, found = r, true
		}
	}
	return best, found
}
