package ftnstorage

// RecoveryMetrics is the compact telemetry contract for the storage health loop.
type RecoveryMetrics struct {
	Checks       uint64
	Repairs      uint64
	RepairFails  uint64
	Rollbacks    uint64
	Quarantines  uint64
	Verified     uint64
}

func (m *RecoveryMetrics) Record(result HealthLoopResult, verified bool) {
	m.Checks++
	if result.RepairRequired {
		m.Repairs++
	}
	if result.Rollback {
		m.Rollbacks++
	}
	if result.Quarantine {
		m.Quarantines++
	}
	if verified {
		m.Verified++
	}
}
