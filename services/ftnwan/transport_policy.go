package ftnwan

// TransportPolicy defines the preferred FTNWAN transport properties.
type TransportPolicy struct {
	RequireEncryption bool
	AllowDirect       bool
	AllowMeshRelay    bool
	MaxLossPPM        uint32
	MaxLatencyMS      uint32
}

func (p TransportPolicy) Accept(path Path) bool {
	if !path.Healthy || path.ID == "" {
		return false
	}
	if p.MaxLossPPM > 0 && path.LossPPM > p.MaxLossPPM {
		return false
	}
	if p.MaxLatencyMS > 0 && path.LatencyMS > p.MaxLatencyMS {
		return false
	}
	return true
}
