package edge

// ProviderHealth is the normalized health state of an external edge/CDN/DNS provider.
type ProviderHealth struct {
	Provider string
	EdgeType string
	Healthy bool
	Verified bool
	LatencyMS uint32
	FailureCount uint32
}

func (h ProviderHealth) Eligible() bool {
	return h.Provider != "" && h.EdgeType != "" && h.Healthy && h.Verified
}
