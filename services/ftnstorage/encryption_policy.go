package ftnstorage

// EncryptionPolicy defines at-rest encryption requirements without coupling
// storage code to a particular cryptographic provider.
type EncryptionPolicy struct {
	Enabled        bool
	KeyVersion     uint64
	RotateAfterDays uint32
	RequireForCold bool
}

func (p EncryptionPolicy) Valid() bool {
	if !p.Enabled || p.KeyVersion == 0 {
		return false
	}
	return p.RotateAfterDays > 0
}
