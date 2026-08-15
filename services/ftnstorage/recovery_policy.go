package ftnstorage

// RecoveryPolicy describes automated storage recovery boundaries.
type RecoveryPolicy struct {
	RequireChecksum bool
	AllowAutoRepair bool
	MaxRepairBytes  uint64
	KeepSnapshots   uint32
	VerifyAfterRepair bool
}

func (p RecoveryPolicy) Valid() bool {
	return p.RequireChecksum && p.VerifyAfterRepair && p.KeepSnapshots > 0
}
