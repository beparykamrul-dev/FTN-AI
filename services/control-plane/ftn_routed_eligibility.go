package main

import "strings"

// FTNRoutedEligibility combines independently observed BGP and BFD state for
// a peer. It is advisory state only; it never changes routing configuration.
type FTNRoutedPeerEligibility struct {
	Peer        string `json:"peer"`
	ASN         uint32 `json:"asn"`
	BGPReady    bool   `json:"bgp_ready"`
	BFDState    FTNBFDState `json:"bfd_state"`
	Eligible    bool   `json:"eligible"`
	Reason      string `json:"reason"`
}

func EvaluateFTNRoutedPeerEligibility(peer FTNObservedBGPPeer, bfd FTNBFDState) FTNRoutedPeerEligibility {
	out := FTNRoutedPeerEligibility{
		Peer: strings.TrimSpace(peer.Address), ASN: peer.ASN,
		BGPReady: peer.Established, BFDState: bfd,
	}
	if out.Peer == "" || out.ASN == 0 {
		out.Reason = "invalid_peer"
		return out
	}
	if !peer.Established {
		out.Reason = "bgp_not_established"
		return out
	}
	if bfd != FTNBFDUp {
		out.Reason = "bfd_not_up"
		return out
	}
	out.Eligible = true
	out.Reason = "bgp_established_and_bfd_up"
	return out
}

func EvaluateFTNRoutedCoreReadiness(nodes []FTNCoreNode) FTNCoreFailoverDecision {
	return EvaluateFTNCoreFailover(nodes, "")
}
