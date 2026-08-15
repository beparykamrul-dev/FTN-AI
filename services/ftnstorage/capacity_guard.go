package ftnstorage

// CapacityGuard prevents automatic placement from consuming the final
// reserve needed for repair/recovery operations.
type CapacityGuard struct {
	ReservePercent uint8
}

func (g CapacityGuard) Valid() bool { return g.ReservePercent > 0 && g.ReservePercent < 100 }

func (g CapacityGuard) CanAllocate(used, total uint64, requested uint64) bool {
	if total == 0 || !g.Valid() || requested > total-used { return false }
	reserve := total * uint64(g.ReservePercent) / 100
	return used+requested <= total-reserve
}
