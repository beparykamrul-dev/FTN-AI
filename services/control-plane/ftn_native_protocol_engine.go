package main

import (
	"errors"
	"sort"
	"strings"
)

type FTNNativeProtocolEndpoint struct {
	Service string `json:"service"`
	Protocol string `json:"protocol"`
	Transport string `json:"transport"`
	Healthy bool `json:"healthy"`
	Authorized bool `json:"authorized"`
	LatencyMS float64 `json:"latency_ms"`
	LossPct float64 `json:"loss_pct"`
	Capacity float64 `json:"capacity"`
}

func NormalizeFTNNativeProtocolEndpoint(e FTNNativeProtocolEndpoint) (FTNNativeProtocolEndpoint, error) {
	e.Service = strings.ToLower(strings.TrimSpace(e.Service))
	e.Protocol = strings.ToLower(strings.TrimSpace(e.Protocol))
	e.Transport = strings.ToLower(strings.TrimSpace(e.Transport))
	if e.Service == "" || e.Protocol == "" || e.Transport == "" { return FTNNativeProtocolEndpoint{}, errors.New("ftn_endpoint_identity_required") }
	if !e.Authorized { return FTNNativeProtocolEndpoint{}, errors.New("ftn_endpoint_authorization_required") }
	if e.LatencyMS < 0 || e.LossPct < 0 || e.LossPct > 100 || e.Capacity < 0 { return FTNNativeProtocolEndpoint{}, errors.New("ftn_endpoint_metric_invalid") }
	return e, nil
}

func NegotiateFTNNativeProtocol(service string, protocols, transports []string, endpoints []FTNNativeProtocolEndpoint) (FTNNativeProtocolEndpoint, bool) {
	service = strings.ToLower(strings.TrimSpace(service))
	p := make(map[string]bool, len(protocols)); t := make(map[string]bool, len(transports))
	for _, v := range protocols { p[strings.ToLower(strings.TrimSpace(v))] = true }
	for _, v := range transports { t[strings.ToLower(strings.TrimSpace(v))] = true }
	candidates := make([]FTNNativeProtocolEndpoint, 0)
	for _, e := range endpoints {
		if strings.ToLower(e.Service) != service || !e.Healthy || !e.Authorized || !p[strings.ToLower(e.Protocol)] || !t[strings.ToLower(e.Transport)] { continue }
		candidates = append(candidates, e)
	}
	if len(candidates) == 0 { return FTNNativeProtocolEndpoint{}, false }
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].LatencyMS != candidates[j].LatencyMS { return candidates[i].LatencyMS < candidates[j].LatencyMS }
		if candidates[i].LossPct != candidates[j].LossPct { return candidates[i].LossPct < candidates[j].LossPct }
		if candidates[i].Capacity != candidates[j].Capacity { return candidates[i].Capacity > candidates[j].Capacity }
		if candidates[i].Protocol != candidates[j].Protocol { return candidates[i].Protocol < candidates[j].Protocol }
		return candidates[i].Transport < candidates[j].Transport
	})
	return candidates[0], true
}
