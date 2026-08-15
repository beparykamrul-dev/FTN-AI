package ftnwan

// Path represents a candidate FTNWAN transport path.
type Path struct {
	ID        string
	Scope     string
	LatencyMS uint32
	LossPPM   uint32
	Healthy   bool
	Cost      uint32
}

// Score returns a deterministic quality score; lower is better.
func Score(p Path) uint64 {
	if !p.Healthy {
		return ^uint64(0)
	}
	return uint64(p.LatencyMS)*1000 + uint64(p.LossPPM)*10 + uint64(p.Cost)
}
