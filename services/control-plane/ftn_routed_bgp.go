package main

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

type FTNBGPPeerState struct { PeerID string `json:"peer_id"`; RemoteAddress string `json:"remote_address"`; RemoteAS uint32 `json:"remote_as"`; Established bool `json:"established"`; PrefixesReceived uint64 `json:"prefixes_received"`; PrefixesAdvertised uint64 `json:"prefixes_advertised"`; LastError string `json:"last_error,omitempty"` }
type FTNBGPRouteEvent struct { PeerID string `json:"peer_id"`; Route FTNRoute `json:"route"`; Direction string `json:"direction"` }
type FTNGoBGPAdapter struct { mu sync.RWMutex; peers map[string]FTNBGPPeerState; routes map[string]FTNRoute; connected bool }
func NewFTNGoBGPAdapter() *FTNGoBGPAdapter { return &FTNGoBGPAdapter{peers: make(map[string]FTNBGPPeerState), routes: make(map[string]FTNRoute)} }
func (a *FTNGoBGPAdapter) Name() string { return "gobgp-v4" }
func (a *FTNGoBGPAdapter) Established() bool { if a==nil{return false}; a.mu.RLock(); defer a.mu.RUnlock(); return a.connected }
func (a *FTNGoBGPAdapter) SetSessionState(connected bool) { if a==nil{return}; a.mu.Lock(); defer a.mu.Unlock(); a.connected=connected }
func (a *FTNGoBGPAdapter) UpsertPeer(state FTNBGPPeerState) error { if a==nil{return errors.New("BGP adapter is required")}; state.PeerID=strings.TrimSpace(state.PeerID);state.RemoteAddress=strings.TrimSpace(state.RemoteAddress);if state.PeerID==""||state.RemoteAddress==""||state.RemoteAS==0{return errors.New("invalid BGP peer")};a.mu.Lock();defer a.mu.Unlock();if a.peers==nil{a.peers=make(map[string]FTNBGPPeerState)};a.peers[state.PeerID]=state;return nil }
func (a *FTNGoBGPAdapter) Peers() []FTNBGPPeerState { if a==nil{return []FTNBGPPeerState{}};a.mu.RLock();defer a.mu.RUnlock();out:=make([]FTNBGPPeerState,0,len(a.peers));for _,p:=range a.peers{out=append(out,p)};sort.Slice(out,func(i,j int)bool{return out[i].PeerID<out[j].PeerID});return out }
func (a *FTNGoBGPAdapter) ApplyRoute(r FTNRoute) error { if a==nil{return errors.New("BGP adapter is required")};if !a.Established(){return errors.New("BGP session is not established")};key:=r.VRF+"|"+r.Prefix;a.mu.Lock();defer a.mu.Unlock();if a.routes==nil{a.routes=make(map[string]FTNRoute)};a.routes[key]=r;return nil }
func (a *FTNGoBGPAdapter) WithdrawRoute(r FTNRoute) error { if a==nil{return errors.New("BGP adapter is required")};key:=r.VRF+"|"+r.Prefix;a.mu.Lock();defer a.mu.Unlock();delete(a.routes,key);return nil }
func (a *FTNGoBGPAdapter) Routes() []FTNRoute { if a==nil{return []FTNRoute{}};a.mu.RLock();defer a.mu.RUnlock();out:=make([]FTNRoute,0,len(a.routes));for _,r:=range a.routes{out=append(out,r)};sort.Slice(out,func(i,j int)bool{a:=out[i].VRF+"|"+out[i].Prefix;b:=out[j].VRF+"|"+out[j].Prefix;return a<b});return out }
