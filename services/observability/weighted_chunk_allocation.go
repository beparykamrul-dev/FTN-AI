package observability

// WeightedChunkAllocation assigns pending chunks proportionally to path quality.
type WeightedChunkAllocation struct{}

func (WeightedChunkAllocation) Allocate(chunks []MigrationChunk, paths []StoragePath, maxParallel uint32) map[string][]MigrationChunk {
	result := make(map[string][]MigrationChunk)
	if maxParallel == 0 { return result }
	s := WeightedPathScheduler{}
	ranked := s.Rank(paths)
	if uint32(len(ranked)) > maxParallel { ranked = ranked[:maxParallel] }
	if len(ranked) == 0 { return result }
	for i, c := range chunks {
		if !c.Valid() || c.Completed { continue }
		p := ranked[i%len(ranked)]
		result[p.PathID] = append(result[p.PathID], c)
	}
	return result
}
