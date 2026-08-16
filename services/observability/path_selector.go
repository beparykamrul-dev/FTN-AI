package observability

// PathSelector chooses the healthiest eligible path with the highest usable capacity.
type PathSelector struct{}

func (PathSelector) Select(paths []StoragePath) (StoragePath, bool) {
	var best StoragePath
	found := false
	for _, p := range paths {
		if !p.Eligible() { continue }
		if !found || p.LatencyMS < best.LatencyMS || (p.LatencyMS == best.LatencyMS && p.LossPercent < best.LossPercent) || (p.LatencyMS == best.LatencyMS && p.LossPercent == best.LossPercent && p.CapacityMbps > best.CapacityMbps) {
			best, found = p, true
		}
	}
	return best, found
}
