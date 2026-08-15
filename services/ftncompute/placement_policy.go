package ftncompute

// PlacementPolicy defines FTN-native workload placement constraints.
type PlacementPolicy struct {
	Scope          string // global, regional, local
	PreferredNodes []string
	AntiAffinity   []string
	RequireHealthy bool
	MaxLatencyMS   uint32
}

func (p PlacementPolicy) Valid() bool {
	if p.Scope != "global" && p.Scope != "regional" && p.Scope != "local" {
		return false
	}
	return p.MaxLatencyMS >= 0
}
