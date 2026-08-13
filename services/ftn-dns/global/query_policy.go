package global

import (
	"sort"
	"strings"
)

// QueryObservation is a provider-neutral result produced by a DNS adapter.
type QueryObservation struct {
	Server    string
	Name      string
	Type      string
	Reachable bool
	Secure    bool
	LatencyMS int64
	Error     string
}

// SelectBestQueryPath returns the healthiest, lowest-latency observation.
// It does not force traffic through any provider; it only ranks supplied paths.
func SelectBestQueryPath(observations []QueryObservation) (QueryObservation, bool) {
	candidates := append([]QueryObservation(nil), observations...)
	sort.SliceStable(candidates, func(i, j int) bool {
		return betterQueryPath(candidates[i], candidates[j])
	})
	for _, item := range candidates {
		if item.Reachable && item.Error == "" {
			return item, true
		}
	}
	return QueryObservation{}, false
}

func betterQueryPath(a, b QueryObservation) bool {
	if a.Reachable != b.Reachable {
		return a.Reachable
	}
	if a.Error == "" != (b.Error == "") {
		return a.Error == ""
	}
	if a.Secure != b.Secure {
		return a.Secure
	}
	if a.LatencyMS != b.LatencyMS {
		return a.LatencyMS < b.LatencyMS
	}
	return strings.ToLower(a.Server) < strings.ToLower(b.Server)
}
