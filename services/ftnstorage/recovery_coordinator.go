package ftnstorage

// RecoveryCoordinator builds recovery decisions without directly mutating
// storage. Execution stays behind the repair worker authorization boundary.
type RecoveryCoordinator struct{}

func (RecoveryCoordinator) Plan(ref ChunkRef, replicas []Replica) (RepairCandidate, bool) {
	source, ok := SelectReplica(replicas)
	if !ok {
		return RepairCandidate{}, false
	}
	return RepairCandidate{Ref: ref, TargetNode: source.NodeID, Reason: "verified-replica-recovery"}, true
}
