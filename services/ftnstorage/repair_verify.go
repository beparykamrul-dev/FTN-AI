package ftnstorage

// RepairVerification records the post-repair integrity decision.
type RepairVerification struct {
	ChunkHash string
	NodeID    string
	Expected  uint64
	Observed  uint64
	Verified  bool
	Reason    string
}

func VerifyRepair(ref ChunkRef, nodeID string, expected, observed uint64) RepairVerification {
	return RepairVerification{
		ChunkHash: ref.Hash,
		NodeID:    nodeID,
		Expected:  expected,
		Observed:  observed,
		Verified:  expected == observed,
		Reason:    "post-repair-integrity-check",
	}
}
