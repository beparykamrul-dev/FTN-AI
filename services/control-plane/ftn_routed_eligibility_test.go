package main

import "testing"

func TestEvaluateFTNRoutedPeerEligibility(t *testing.T) {
	peer := FTNObservedBGPPeer{Address: "192.0.2.2", ASN: 64512, Established: true}
	if got := EvaluateFTNRoutedPeerEligibility(peer, FTNBFDUp); !got.Eligible {
		t.Fatalf("expected eligible peer: %+v", got)
	}
	if got := EvaluateFTNRoutedPeerEligibility(peer, FTNBFDDown); got.Eligible || got.Reason != "bfd_not_up" {
		t.Fatalf("unexpected BFD-down decision: %+v", got)
	}
	peer.Established = false
	if got := EvaluateFTNRoutedPeerEligibility(peer, FTNBFDUp); got.Eligible || got.Reason != "bgp_not_established" {
		t.Fatalf("unexpected BGP-down decision: %+v", got)
	}
}

func TestEvaluateFTNRoutedPeerEligibilityRejectsInvalidPeer(t *testing.T) {
	got := EvaluateFTNRoutedPeerEligibility(FTNObservedBGPPeer{Address: "", ASN: 0}, FTNBFDUp)
	if got.Eligible || got.Reason != "invalid_peer" {
		t.Fatalf("unexpected decision: %+v", got)
	}
}
