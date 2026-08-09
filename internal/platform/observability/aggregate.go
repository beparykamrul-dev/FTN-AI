package observability

import "sort"

type TrafficAggregate struct {
	Interface string `json:"interface"`
	Protocol  string `json:"protocol,omitempty"`
	App       string `json:"application,omitempty"`
	Bytes     uint64 `json:"bytes"`
	Packets   uint64 `json:"packets"`
}

func Aggregate(samples []TrafficSample) []TrafficAggregate {
	m := make(map[string]*TrafficAggregate)
	for _, s := range samples {
		key := s.Interface + "|" + s.Protocol + "|" + s.App
		v := m[key]
		if v == nil { v = &TrafficAggregate{Interface: s.Interface, Protocol: s.Protocol, App: s.App}; m[key] = v }
		v.Bytes += s.Bytes
		v.Packets += s.Packets
	}
	out := make([]TrafficAggregate, 0, len(m))
	for _, v := range m { out = append(out, *v) }
	sort.Slice(out, func(i, j int) bool { return out[i].Bytes > out[j].Bytes })
	return out
}
