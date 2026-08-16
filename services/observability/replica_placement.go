package observability

// ReplicaTarget is a candidate server/storage domain for a telemetry replica.
type ReplicaTarget struct {
	NodeID string
	Domain FailureDomain
	Healthy bool
	FreeGB float64
}

func (t ReplicaTarget) Eligible() bool {
	return t.NodeID != "" && t.Domain.Eligible() && t.Healthy && t.FreeGB > 0
}

// SelectReplicaTargets selects up to count eligible targets while preferring distinct failure domains.
func SelectReplicaTargets(targets []ReplicaTarget, count uint32) []ReplicaTarget {
	if count == 0 { return nil }
	selected := make([]ReplicaTarget, 0, count)
	seen := make(map[string]bool)
	for _, t := range targets {
		if !t.Eligible() || seen[t.Domain.ID] { continue }
		selected = append(selected, t)
		seen[t.Domain.ID] = true
		if uint32(len(selected)) == count { return selected }
	}
	for _, t := range targets {
		if !t.Eligible() { continue }
		duplicate := false
		for _, s := range selected { if s.NodeID == t.NodeID { duplicate = true; break } }
		if duplicate { continue }
		selected = append(selected, t)
		if uint32(len(selected)) == count { break }
	}
	return selected
}
