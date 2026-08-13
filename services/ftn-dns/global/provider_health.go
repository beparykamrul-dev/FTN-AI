package global

import "time"

// ProviderHealth is the normalized health state used by FTN DNS routing.
type ProviderHealth struct {
	Provider   string
	Reachable  bool
	Secure     bool
	Latency    time.Duration
	LossRatio  float64
	Available  bool
}

// Healthy reports whether a provider is eligible for DNS service selection.
func (p ProviderHealth) Healthy() bool {
	return p.Reachable && p.Secure && p.Available && p.LossRatio >= 0 && p.LossRatio < 1
}

// Score produces a bounded operational score. Lower latency and packet loss
// improve the score; unhealthy providers are excluded by returning zero.
func (p ProviderHealth) Score() float64 {
	if !p.Healthy() {
		return 0
	}
	latencyMS := float64(p.Latency.Milliseconds())
	if latencyMS < 1 {
		latencyMS = 1
	}
	return (1000 / latencyMS) * (1 - p.LossRatio)
}
