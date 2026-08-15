package ftnstorage

// RebalanceMove is a declarative movement of a chunk reference between nodes.
type RebalanceMove struct {
	Chunk   ChunkRef
	From    string
	To      string
	Reason  string
}

// RebalancePlan contains storage moves without executing them.
type RebalancePlan struct {
	Moves []RebalanceMove
}

func (p RebalancePlan) Valid() bool {
	for _, m := range p.Moves {
		if m.From == "" || m.To == "" || m.From == m.To || m.Chunk.Hash == "" {
			return false
		}
	}
	return true
}
