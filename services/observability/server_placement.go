package observability

// ServerStorageCandidate represents a remote FTN storage target.
type ServerStorageCandidate struct {
	ServerID string
	Healthy bool
	FreeGB float64
	LatencyMS float64
	Priority uint32
}

func (c ServerStorageCandidate) Eligible() bool {
	return c.ServerID != "" && c.Healthy && c.FreeGB > 0
}

func SelectServerStorage(candidates []ServerStorageCandidate) (ServerStorageCandidate, bool) {
	var best ServerStorageCandidate
	found := false
	for _, c := range candidates {
		if !c.Eligible() { continue }
		if !found || c.Priority < best.Priority ||
			(c.Priority == best.Priority && c.LatencyMS < best.LatencyMS) ||
			(c.Priority == best.Priority && c.LatencyMS == best.LatencyMS && c.FreeGB > best.FreeGB) {
			best, found = c, true
		}
	}
	return best, found
}
