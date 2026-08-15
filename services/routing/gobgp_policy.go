package routing

// GoBGPPolicy defines the FTN control-plane boundary for BGP automation.
type GoBGPPolicy struct {
	NodeID       string
	LocalASN     uint32
	PeerASN      uint32
	Neighbor     string
	RequireAuth  bool
	RequireApproval bool
}

func (p GoBGPPolicy) Valid() bool {
	return p.NodeID != "" && p.LocalASN != 0 && p.PeerASN != 0 && p.Neighbor != ""
}
