package observability

// StoragePressure describes available local durable-buffer capacity.
type StoragePressure struct {
	UsedPercent float64
	ReservedPercent float64
}

func (s StoragePressure) Safe() bool {
	return s.UsedPercent >= 0 && s.UsedPercent < (100-s.ReservedPercent)
}
