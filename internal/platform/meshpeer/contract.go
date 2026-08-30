package meshpeer

import "time"

type PeerState string

const (
	PeerUnknown     PeerState = "UNKNOWN"
	PeerDiscovered  PeerState = "DISCOVERED"
	PeerVerified    PeerState = "VERIFIED"
	PeerHealthy     PeerState = "HEALTHY"
	PeerUnhealthy   PeerState = "UNHEALTHY"
	PeerQuarantined PeerState = "QUARANTINED"
)

type Peer struct {
	ID            string
	Endpoint      string
	State         PeerState
	CertificateOK bool
	PolicyAllowed bool
	LastSeen      time.Time
}

type ProbeResult struct {
	Reachable bool
	ReadOK    bool
	LatencyMS int64
	Error     string
}

type PathCandidate struct {
	PeerID    string
	Transport string
	Healthy   bool
	Allowed   bool
	Score     int64
}

// SelectCandidates only evaluates already-approved candidates. It does not
// mutate routes or activate transports.
func SelectCandidates(peers []Peer, candidates []PathCandidate) []PathCandidate {
	allowedPeers := make(map[string]bool, len(peers))
	for _, p := range peers {
		allowedPeers[p.ID] = p.PolicyAllowed && p.CertificateOK && p.State == PeerHealthy
	}

	out := make([]PathCandidate, 0, len(candidates))
	for _, c := range candidates {
		if c.Allowed && c.Healthy && allowedPeers[c.PeerID] {
			out = append(out, c)
		}
	}
	return out
}
