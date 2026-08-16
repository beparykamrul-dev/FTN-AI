package observability

// StorageTier identifies a preferred durable-buffer storage class.
type StorageTier struct {
	Name string
	LatencyMS float64
	CapacityGB float64
	Healthy bool
	Priority uint32
}

func (t StorageTier) Eligible() bool {
	return t.Name != "" && t.CapacityGB > 0 && t.Healthy
}

func SelectStorageTier(tiers []StorageTier) (StorageTier, bool) {
	var best StorageTier
	found := false
	for _, t := range tiers {
		if !t.Eligible() { continue }
		if !found || t.Priority < best.Priority || (t.Priority == best.Priority && t.LatencyMS < best.LatencyMS) {
			best, found = t, true
		}
	}
	return best, found
}
