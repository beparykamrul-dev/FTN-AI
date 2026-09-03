package main

import (
	"context"
	"fmt"
	"net/netip"

	api "github.com/osrg/gobgp/v4/api"
)

type FTNObservedBGPPeer struct {
	Address     string
	ASN         uint32
	Established bool
}

type FTNBGPPeerSource interface {
	ListPeers(context.Context) ([]FTNObservedBGPPeer, error)
}

// FTNGoBGPPeerSource exposes only normalized peer state from the live GoBGP
// server. It has no configuration or route mutation capability.
type FTNGoBGPPeerSource struct {
	server *FTNGoBGPServer
}

func NewFTNGoBGPPeerSource(server *FTNGoBGPServer) *FTNGoBGPPeerSource {
	return &FTNGoBGPPeerSource{server: server}
}

func (s *FTNGoBGPPeerSource) ListPeers(ctx context.Context) ([]FTNObservedBGPPeer, error) {
	if s == nil || s.server == nil {
		return nil, fmt.Errorf("GoBGP server is required")
	}
	s.server.mu.RLock()
	started, gobgp := s.server.started, s.server.server
	s.server.mu.RUnlock()
	if !started || gobgp == nil {
		return nil, fmt.Errorf("GoBGP server is not started")
	}
	peers := make([]FTNObservedBGPPeer, 0)
	err := gobgp.ListPeer(ctx, &api.ListPeerRequest{}, func(peer *api.Peer) {
		if peer == nil || peer.Conf == nil || peer.Conf.NeighborAddress == "" || peer.Conf.PeerAsn == 0 {
			return
		}
		if _, err := netip.ParseAddr(peer.Conf.NeighborAddress); err != nil {
			return
		}
		peers = append(peers, FTNObservedBGPPeer{
			Address: peer.Conf.NeighborAddress,
			ASN: peer.Conf.PeerAsn,
			Established: peer.State != nil && peer.State.SessionState == api.PeerState_ESTABLISHED,
		})
	})
	return peers, err
}
