package main

import (
	"errors"
	"sort"
	"strings"
)

type FTNServiceEndpoint struct {
	Service string `json:"service"`
	Protocol string `json:"protocol"`
	Transport string `json:"transport"`
	Address string `json:"address"`
	Healthy bool `json:"healthy"`
	LatencyMS float64 `json:"latency_ms"`
	LossPct float64 `json:"loss_pct"`
	JitterMS float64 `json:"jitter_ms"`
	Capacity float64 `json:"capacity"`
	FTNOwned bool `json:"ftn_owned"`
	Authorized bool `json:"authorized"`
}

func NormalizeFTNServiceEndpoint(e FTNServiceEndpoint) (FTNServiceEndpoint, error) {
	e.Service = strings.ToLower(strings.TrimSpace(e.Service))
	e.Protocol = strings.ToLower(strings.TrimSpace(e.Protocol))
	e.Transport = strings.ToLower(strings.TrimSpace(e.Transport))
	e.Address = strings.TrimSpace(e.Address)
	if e.Service == "" || e.Protocol == "" || e.Transport == "" || e.Address == "" {
		return FTNServiceEndpoint{}, errors.New("ftn_service_endpoint_required")
	}
	if !e.FTNOwned && !e.Authorized {
		return FTNServiceEndpoint{}, errors.New("ftn_service_endpoint_authorization_required")
	}
	if e.LossPct < 0 || e.LossPct > 100 || e.JitterMS < 0 || e.LatencyMS < 0 || e.Capacity < 0 {
		return FTNServiceEndpoint{}, errors.New("ftn_service_endpoint_metric_invalid")
	}
	return e, nil
}

func SelectFTNServiceEndpoints(service string, endpoints []FTNServiceEndpoint, max int) []FTNServiceEndpoint {
	service = strings.ToLower(strings.TrimSpace(service))
	if max < 1 { max = 1 }
	out := make([]FTNServiceEndpoint, 0, len(endpoints))
	for _, e := range endpoints {
		if !e.Healthy || e.Service != service { continue }
		if !e.FTNOwned && !e.Authorized { continue }
		out = append(out, e)
	}
	sort.SliceStable(out, func(i, j int) bool {
		scoreI := out[i].LatencyMS + out[i].LossPct*10 + out[i].JitterMS
		scoreJ := out[j].LatencyMS + out[j].LossPct*10 + out[j].JitterMS
		if scoreI != scoreJ { return scoreI < scoreJ }
		if out[i].Capacity != out[j].Capacity { return out[i].Capacity > out[j].Capacity }
		if out[i].Protocol != out[j].Protocol { return out[i].Protocol < out[j].Protocol }
		return out[i].Address < out[j].Address
	})
	if len(out) > max { out = out[:max] }
	return out
}
