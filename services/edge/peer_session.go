package edge

// EdgePeerSession is the normalized state of a server-to-server FTN edge
// peering session. Transport implementation is intentionally separate.
type EdgePeerSession struct {
	LocalNode  string
	RemoteNode string
	MutualTLS  bool
	Healthy    bool
	Verified   bool
	LatencyMS  uint32
}

func (s EdgePeerSession) Ready() bool {
	return s.LocalNode != "" && s.RemoteNode != "" && s.LocalNode != s.RemoteNode &&
		s.MutualTLS && s.Healthy && s.Verified
}
