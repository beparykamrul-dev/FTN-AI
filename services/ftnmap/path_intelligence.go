package ftnmap

// PathScore is an explainable score for a service-to-service path.
type PathScore struct {
	LatencyMS  uint32
	LossPct    float64
	Healthy    bool
	Available  bool
}

func (p PathScore) Score() float64 {
	if !p.Healthy || !p.Available { return 0 }
	s := 1.0
	if p.LatencyMS > 0 { s -= float64(p.LatencyMS)/10000 }
	s -= p.LossPct/100
	if s < 0 { return 0 }
	return s
}
