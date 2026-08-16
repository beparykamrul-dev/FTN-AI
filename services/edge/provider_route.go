package edge

// ProviderRoute is a normalized route advertised/selected for external provider traffic.
type ProviderRoute struct {
	Provider string
	EdgeType string
	Endpoint string
	Healthy bool
	Verified bool
	LatencyMS uint32
}

func (r ProviderRoute) Eligible() bool {
	return r.Provider != "" && r.EdgeType != "" && r.Endpoint != "" && r.Healthy && r.Verified
}

func SelectProviderRoute(routes []ProviderRoute) (ProviderRoute, bool) {
	var best ProviderRoute
	found := false
	for _, r := range routes {
		if !r.Eligible() { continue }
		if !found || r.LatencyMS < best.LatencyMS ||
			(r.LatencyMS == best.LatencyMS && r.Provider < best.Provider) {
			best, found = r, true
		}
	}
	return best, found
}
