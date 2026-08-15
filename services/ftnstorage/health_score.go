package ftnstorage

// HealthScore is a compact, explainable storage-health indicator.
type HealthScore struct {
	Integrity float64
	Capacity  float64
	Errors    float64
	Recovery  float64
}

func (s HealthScore) Valid() bool {
	return s.Integrity >= 0 && s.Integrity <= 1 && s.Capacity >= 0 && s.Capacity <= 1 && s.Errors >= 0 && s.Errors <= 1 && s.Recovery >= 0 && s.Recovery <= 1
}

func (s HealthScore) Overall() float64 {
	if !s.Valid() { return 0 }
	return 0.4*s.Integrity + 0.2*s.Capacity + 0.2*s.Errors + 0.2*s.Recovery
}
