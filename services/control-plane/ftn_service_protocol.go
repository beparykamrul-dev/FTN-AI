package main

import (
	"errors"
	"sort"
	"strings"
)

type FTNServiceProtocolEndpoint struct {
	ServiceID string `json:"service_id"`
	Protocol string `json:"protocol"`
	Transport string `json:"transport"`
	Healthy bool `json:"healthy"`
	LatencyMS float64 `json:"latency_ms"`
	LossPct float64 `json:"loss_pct"`
	JitterMS float64 `json:"jitter_ms"`
	Capacity float64 `json:"capacity"`
}

func NormalizeFTNServiceProtocolEndpoint(e FTNServiceProtocolEndpoint) (FTNServiceProtocolEndpoint, error) {
	e.ServiceID = strings.ToLower(strings.TrimSpace(e.ServiceID))
	e.Protocol = strings.ToLower(strings.TrimSpace(e.Protocol))
	e.Transport = strings.ToLower(strings.TrimSpace(e.Transport))
	if e.ServiceID == "" || e.Protocol == "" || e.Transport == "" {
		return FTNServiceProtocolEndpoint{}, errors.New("ftn_service_protocol_identity_required")
	}
	if e.LatencyMS < 0 || e.LossPct < 0 || e.JitterMS < 0 || e.Capacity < 0 {
		return FTNServiceProtocolEndpoint{}, errors.New("ftn_service_protocol_metrics_invalid")
	}
	return e, nil
}

func SelectFTNServiceProtocols(endpoints []FTNServiceProtocolEndpoint, max int) []FTNServiceProtocolEndpoint {
	if max < 1 { max = 1 }
	out := make([]FTNServiceProtocolEndpoint, 0, len(endpoints))
	for _, e := range endpoints {
		n, err := NormalizeFTNServiceProtocolEndpoint(e)
		if err == nil && n.Healthy { out = append(out, n) }
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].LatencyMS != out[j].LatencyMS { return out[i].LatencyMS < out[j].LatencyMS }
		if out[i].LossPct != out[j].LossPct { return out[i].LossPct < out[j].LossPct }
		if out[i].JitterMS != out[j].JitterMS { return out[i].JitterMS < out[j].JitterMS }
		if out[i].Capacity != out[j].Capacity { return out[i].Capacity > out[j].Capacity }
		if out[i].Protocol != out[j].Protocol { return out[i].Protocol < out[j].Protocol }
		return out[i].Transport < out[j].Transport
	})
	if len(out) > max { out = out[:max] }
	return out
}
