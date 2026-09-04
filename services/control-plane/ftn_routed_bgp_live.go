package main

import (
	"context"
	"fmt"
	"net/netip"
	"sync"
	"time"

	api "github.com/osrg/gobgp/v4/api"
	"github.com/osrg/gobgp/v4/pkg/server"
	"google.golang.org/protobuf/types/known/anypb"
)

type FTNGoBGPServer struct { mu sync.RWMutex; server *server.BgpServer; started bool }
func NewFTNGoBGPServer() *FTNGoBGPServer { return &FTNGoBGPServer{server: server.NewBgpServer()} }
func (g *FTNGoBGPServer) Start(ctx context.Context, asn uint32, routerID string) error { if ctx==nil{return fmt.Errorf("context is required")};g.mu.Lock();defer g.mu.Unlock();if g.server==nil{return fmt.Errorf("GoBGP server is unavailable")};if g.started{return fmt.Errorf("GoBGP server already started")};if asn==0{return fmt.Errorf("local ASN is required")};addr,err:=netip.ParseAddr(routerID);if err!=nil||!addr.Is4(){return fmt.Errorf("invalid IPv4 router-id")};go g.server.Serve();if err:=g.server.StartBgp(ctx,&api.StartBgpRequest{Global:&api.Global{Asn:asn,RouterId:routerID}});err!=nil{return err};g.started=true;return nil }
func (g *FTNGoBGPServer) Stop() error { g.mu.Lock();defer g.mu.Unlock();if g.server==nil{return fmt.Errorf("GoBGP server is unavailable")};if !g.started{return nil};ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second);defer cancel();if err:=g.server.StopBgp(ctx,&api.StopBgpRequest{});err!=nil{return err};g.started=false;return nil }
func (g *FTNGoBGPServer) AddPeer(ctx context.Context,address string,asn uint32) error {g.mu.RLock();defer g.mu.RUnlock();if g.server==nil{return fmt.Errorf("GoBGP server is unavailable")};if !g.started{return fmt.Errorf("GoBGP server is not started")};if ctx==nil{return fmt.Errorf("context is required")};addr,err:=netip.ParseAddr(address);if err!=nil{return fmt.Errorf("invalid peer address: %w",err)};if asn==0{return fmt.Errorf("peer ASN is required")};return g.server.AddPeer(ctx,&api.AddPeerRequest{Peer:&api.Peer{Conf:&api.PeerConf{NeighborAddress:addr.String(),PeerAsn:asn}}})}
func (g *FTNGoBGPServer) AdvertiseIPv4Prefix(ctx context.Context,prefix string,nextHop string) error {g.mu.RLock();defer g.mu.RUnlock();if g.server==nil{return fmt.Errorf("GoBGP server is unavailable")};if !g.started{return fmt.Errorf("GoBGP server is not started")};if ctx==nil{return fmt.Errorf("context is required")};p,err:=netip.ParsePrefix(prefix);if err!=nil{return fmt.Errorf("invalid prefix: %w",err)};p=p.Masked();nh,err:=netip.ParseAddr(nextHop);if err!=nil{return fmt.Errorf("invalid next-hop: %w",err)};if !p.Addr().Is4()||!nh.Is4(){return fmt.Errorf("IPv4 prefix and next-hop are required")};attrs,err:=anypb.New(&api.NextHopAttribute{NextHop:nh.String()});if err!=nil{return err};_,err=g.server.AddPath(ctx,&api.AddPathRequest{Path:&api.Path{Nlri:&api.NLRI{Nlri:&api.NLRI_Prefix{Prefix:&api.IPAddressPrefix{PrefixLen:uint32(p.Bits()),Prefix:p.Addr().String()}}},Pattrs:[]*anypb.Any{attrs},Family:&api.Family{Afi:api.Family_AFI_IP,Safi:api.Family_SAFI_UNICAST}}});return err}
func (g *FTNGoBGPServer) ListIPv4Routes(ctx context.Context)([]FTNRoute,error){g.mu.RLock();defer g.mu.RUnlock();if g.server==nil{return nil,fmt.Errorf("GoBGP server is unavailable")};if !g.started{return nil,fmt.Errorf("GoBGP server is not started")};if ctx==nil{return nil,fmt.Errorf("context is required")};family:=&api.Family{Afi:api.Family_AFI_IP,Safi:api.Family_SAFI_UNICAST};routes:=make([]FTNRoute,0);err:=g.server.ListPath(ctx,&api.ListPathRequest{TableType:api.TableType_GLOBAL,Family:family},func(destination *api.Destination){for _,path:=range destination.Paths{if path==nil||path.Nlri==nil{continue};p:=path.GetNlri().GetPrefix();if p==nil{continue};routes=append(routes,FTNRoute{Prefix:fmt.Sprintf("%s/%d",p.Prefix,p.PrefixLen),Protocol:"bgp",VRF:"default",Active:!path.IsWithdraw})}});if err!=nil{return nil,err};return routes,nil}
