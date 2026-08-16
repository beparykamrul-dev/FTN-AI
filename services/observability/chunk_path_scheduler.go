package observability

// ChunkPathScheduler assigns pending chunks to eligible paths using bounded round-robin selection.
type ChunkPathScheduler struct{}

func (ChunkPathScheduler) Schedule(chunks []MigrationChunk, paths []StoragePath, maxParallel uint32) map[string][]MigrationChunk {
	result := make(map[string][]MigrationChunk)
	if maxParallel == 0 { return result }
	eligible := make([]StoragePath, 0, len(paths))
	for _, p := range paths { if p.Eligible() { eligible = append(eligible, p) } }
	if uint32(len(eligible)) > maxParallel { eligible = eligible[:maxParallel] }
	if len(eligible) == 0 { return result }
	for i, c := range chunks { if c.Valid() && !c.Completed { p := eligible[i%len(eligible)]; result[p.PathID] = append(result[p.PathID], c) } }
	return result
}
