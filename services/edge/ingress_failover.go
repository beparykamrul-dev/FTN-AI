package edge

// IngressCandidate is an approved external provider ingress endpoint.
type IngressCandidate struct {
	Provider   string
	Endpoint   string
	TLSHealthy bool
	EdgeHealthy bool
	LatencyMS  uint32
}

// SelectIngress chooses a healthy TLS-valid ingress endpoint with the lowest
// observed latency. The caller remains responsible for policy authorization.
func SelectIngress(candidates []IngressCandidate) (IngressCandidate, bool) {
	var best IngressCandidate
	found := false
	for _, c := range candidates {
		if c.Provider == "" || c.Endpoint == "" || !c.TLSHealthy || !c.EdgeHealthy {
			continue
		}
		if !found || c.LatencyMS < best.LatencyMS ||
			(c.LatencyMS == best.LatencyMS && c.Endpoint < best.Endpoint) {
			best, found = c, true
		}
	}
	return best, found
}
