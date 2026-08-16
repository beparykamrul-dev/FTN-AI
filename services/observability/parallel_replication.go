package observability

// ParallelReplicationPolicy bounds concurrent replication streams.
type ParallelReplicationPolicy struct {
	MaxStreams uint32
	RequireHealthyPaths bool
}

func (p ParallelReplicationPolicy) Valid() bool { return p.MaxStreams > 0 }
