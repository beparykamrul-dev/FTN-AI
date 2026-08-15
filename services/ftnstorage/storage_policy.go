package ftnstorage

// StoragePolicy defines FTN's data-safety and space-efficiency goals.
type StoragePolicy struct {
	ReplicationFactor uint8
	ErasureCoding    bool
	Compression      bool
	Deduplication    bool
	Checksums        bool
	SnapshotEnabled  bool
	Scope            string // global, regional, local
}

func (p StoragePolicy) Valid() bool {
	if p.ReplicationFactor == 0 || p.ReplicationFactor > 9 {
		return false
	}
	if p.Scope != "global" && p.Scope != "regional" && p.Scope != "local" {
		return false
	}
	return p.Checksums && p.SnapshotEnabled
}
