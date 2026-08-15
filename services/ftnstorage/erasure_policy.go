package ftnstorage

// ErasurePolicy defines when FTN may use erasure coding for durable data.
type ErasurePolicy struct {
	DataShards   uint16
	ParityShards uint16
	MinHealthy   uint16
}

func (p ErasurePolicy) Valid() bool {
	return p.DataShards > 0 && p.ParityShards > 0 && p.MinHealthy >= p.DataShards
}
