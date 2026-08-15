package ftnstorage

import "context"

// RepairCandidate describes a verified chunk that needs recovery.
type RepairCandidate struct {
	Ref        ChunkRef
	TargetNode string
	Reason     string
}

// RepairWorker is an execution boundary; implementations must perform their
// own authorization and integrity checks before writing bytes.
type RepairWorker interface {
	Repair(context.Context, RepairCandidate) error
}
