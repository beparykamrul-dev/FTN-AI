package observability

// BackendHealth captures runtime health signals used by FTN routing policy.
type BackendHealth struct {
	Name string
	Healthy bool
	LatencyMS float64
	ErrorRate float64
	FreshnessMS float64
}

func (h BackendHealth) Score() float64 {
	if !h.Healthy { return 0 }
	s := 100.0
	s -= h.LatencyMS / 100.0
	s -= h.ErrorRate * 50.0
	s -= h.FreshnessMS / 1000.0
	if s < 0 { return 0 }
	return s
}
