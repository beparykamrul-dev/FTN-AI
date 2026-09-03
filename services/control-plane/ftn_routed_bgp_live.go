package main

import (
	"context"
	"fmt"
	"net/netip"
	"sync"

	api "github.com/osrg/gobgp/v4/api/gobgpapi"
	"github.com/osrg/gobgp/v4/pkg/server"
	"google.golang.org/protobuf/types/known/anypb"
)

type FTNGoBGPServer struct {
	mu sync.RWMutex
	server *server.BgpServer
	started bool
}

func NewFTNGoBGPServer() *FTNGoBGPServer { return &FTNGoBGPServer{server: server.NewBgpServer()} }

func (g *FTNGoBGPServer) Start(ctx context.Context, asn uint32, routerID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.started { return fmt.Errorf("GoBGP server already started") }
	if asn == 0 { return fmt.Errorf("local ASN is required") }
	if _, err := netip.ParseAddr(routerID); err != nil { return fmt.Errorf("invalid router-id: %w", err) }
	g.server.StartBgp(ctx)
	if err := g.server.StartGrpc(ctx, ":50051"); err != nil { return err }
	if err := g.server.AddPath(g.server, &api.AddPathRequest{Path: &api.Path{Family: &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST}}}); err != nil { _ = err }
	g.started = true
	return nil
}

func (g *FTNGoBGPServer) Stop() error {
	g.mu.Lock(); defer g.mu.Unlock()
	if !g.started { return nil }
	g.server.StopBgp()
	g.started = false
	return nil
}

func (g *FTNGoBGPServer) AddPeer(ctx context.Context, address string, asn uint32) error {
	g.mu.RLock(); defer g.mu.RUnlock()
	if !g.started { return fmt.Errorf("GoBGP server is not started") }
	if _, err := netip.ParseAddr(address); err != nil { return fmt.Errorf("invalid peer address: %w", err) }
	if asn == 0 { return fmt.Errorf("peer ASN is required") }
	return g.server.AddPeer(ctx, &api.AddPeerRequest{Peer: &api.Peer{Conf: &api.PeerConf{NeighborAddress: address, PeerAsn: asn}}})
}

func (g *FTNGoBGPServer) AdvertiseIPv4Prefix(ctx context.Context, prefix string, nextHop string) error {
	g.mu.RLock(); defer g.mu.RUnlock()
	if !g.started { return fmt.Errorf("GoBGP server is not started") }
	p, err := netip.ParsePrefix(prefix); if err != nil { return fmt.Errorf("invalid prefix: %w", err) }
	p = p.Masked()
	nh, err := netip.ParseAddr(nextHop); if err != nil { return fmt.Errorf("invalid next-hop: %w", err) }
	if !p.Addr().Is4() || !nh.Is4() { return fmt.Errorf("IPv4 prefix and next-hop are required") }
	nlri, err := anypb.New(&api.IPAddressPrefix{PrefixLen: uint32(p.Bits()), Prefix: p.Addr().String()}); if err != nil { return err }
	attrs, err := anypb.New(&api.NextHopAttribute{NextHop: nh.String()}); if err != nil { return err }
	_, err = g.server.AddPath(ctx, &api.AddPathRequest{Path: &api.Path{Nlri: nlri, Pattrs: []*anypb.Any{attrs}, Family: &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST}}})
	return err
}
