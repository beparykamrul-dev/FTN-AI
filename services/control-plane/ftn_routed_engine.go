package main

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

type FTNRIB interface { Install(FTNRoute) error; Withdraw(FTNRoute) error; Lookup(string,string)([]FTNRoute,error) }
type FTNFIB interface { Program(FTNRoute) error; Remove(FTNRoute) error }
type FTNBGPAdapter interface { Name() string; Established() bool; ApplyRoute(FTNRoute) error; WithdrawRoute(FTNRoute) error }
type FTNBFDState string
const ( FTNBFDUp FTNBFDState="up"; FTNBFDDown FTNBFDState="down"; FTNBFDUnknown FTNBFDState="unknown" )
type FTNCoreNode struct { ID string `json:"id"`; Healthy bool `json:"healthy"`; BGPReady bool `json:"bgp_ready"`; BFDState FTNBFDState `json:"bfd_state"` }
type FTNCoreFailoverDecision struct { ActiveNode string `json:"active_node"`; Failover bool `json:"failover"`; Reason string `json:"reason"` }
type FTNRoutedEngine struct { mu sync.RWMutex; rib FTNRIB; fib FTNFIB; bgp []FTNBGPAdapter; core []FTNCoreNode }
func NewFTNRoutedEngine(rib FTNRIB,fib FTNFIB)*FTNRoutedEngine{return &FTNRoutedEngine{rib:rib,fib:fib,core:make([]FTNCoreNode,0,2)}}
func(e *FTNRoutedEngine)RegisterBGPAdapter(a FTNBGPAdapter){if e==nil||a==nil{return};e.mu.Lock();defer e.mu.Unlock();e.bgp=append(e.bgp,a)}
func(e *FTNRoutedEngine)SetCoreNodes(nodes []FTNCoreNode){if e==nil{return};clean:=make([]FTNCoreNode,0,len(nodes));seen:=map[string]struct{}{};for _,n:=range nodes{n.ID=strings.TrimSpace(n.ID);if n.ID==""{continue};if _,ok:=seen[n.ID];ok{continue};seen[n.ID]=struct{}{};clean=append(clean,n)};sort.Slice(clean,func(i,j int)bool{return clean[i].ID<clean[j].ID});e.mu.Lock();e.core=clean;e.mu.Unlock()}
func(e *FTNRoutedEngine)ApplyApprovedRoute(in FTNRouteIntent,approved bool)error{if e==nil{return errors.New("routed engine is required")};if !approved{return errors.New("route approval required")};r,err:=NormalizeFTNRoute(in.Route);if err!=nil{return err};if e.rib==nil||e.fib==nil{return errors.New("RIB and FIB implementations are required")};if err:=e.rib.Install(r);err!=nil{return err};if err:=e.fib.Program(r);err!=nil{_ = e.rib.Withdraw(r);return err};e.mu.RLock();adapters:=append([]FTNBGPAdapter(nil),e.bgp...);e.mu.RUnlock();for _,a:=range adapters{if a!=nil&&a.Established(){if err:=a.ApplyRoute(r);err!=nil{return err}}};return nil}
func(e *FTNRoutedEngine)WithdrawApprovedRoute(in FTNRouteIntent,approved bool)error{if e==nil{return errors.New("routed engine is required")};if !approved{return errors.New("route approval required")};r,err:=NormalizeFTNRoute(in.Route);if err!=nil{return err};if e.rib==nil||e.fib==nil{return errors.New("RIB and FIB implementations are required")};if err:=e.fib.Remove(r);err!=nil{return err};if err:=e.rib.Withdraw(r);err!=nil{return err};e.mu.RLock();adapters:=append([]FTNBGPAdapter(nil),e.bgp...);e.mu.RUnlock();for _,a:=range adapters{if a!=nil&&a.Established(){_ = a.WithdrawRoute(r)}};return nil}
func EvaluateFTNCoreFailover(nodes []FTNCoreNode,current string)FTNCoreFailoverDecision{current=strings.TrimSpace(current);clean:=append([]FTNCoreNode(nil),nodes...);sort.SliceStable(clean,func(i,j int)bool{return clean[i].ID<clean[j].ID});for _,n:=range clean{if n.ID==current&&n.Healthy&&n.BGPReady&&n.BFDState!=FTNBFDDown{return FTNCoreFailoverDecision{ActiveNode:current,Reason:"current_core_healthy"}}};for _,n:=range clean{if n.Healthy&&n.BGPReady&&n.BFDState!=FTNBFDDown&&strings.TrimSpace(n.ID)!=""{return FTNCoreFailoverDecision{ActiveNode:n.ID,Failover:n.ID!=current,Reason:"alternate_core_healthy"}}};return FTNCoreFailoverDecision{Reason:"no_healthy_core_available"}}
