package observability

// ProportionalChunkAllocator distributes pending chunks according to weighted path capacity.
type ProportionalChunkAllocator struct{}

func (a ProportionalChunkAllocator) Allocate(chunks []MigrationChunk, paths []StoragePath, maxParallel uint32) map[string][]MigrationChunk {
	out := make(map[string][]MigrationChunk)
	if maxParallel == 0 { return out }
	ranked := (WeightedPathScheduler{}).Rank(paths)
	if uint32(len(ranked)) > maxParallel { ranked = ranked[:maxParallel] }
	if len(ranked) == 0 { return out }
	total := 0.0
	for _, p := range ranked { total += p.CapacityMbps }
	if total <= 0 { return out }
	for i, c := range chunks {
		if !c.Valid() || c.Completed { continue }
		cursor := float64(i) / float64(len(chunks))
		acc := 0.0
		for _, p := range ranked {
			acc += p.CapacityMbps / total
			if cursor < acc { out[p.PathID] = append(out[p.PathID], c); break }
		}
	}
	return out
}
