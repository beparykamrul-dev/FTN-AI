package main

import (
	"errors"
	"sync"
)

type FTNRIB interface {
	Install(FTNRoute) error
	Withdraw(FTNRoute) error
	Lookup(string, string) ([]FTNRoute, error)
}

type FTNFIB interface {
	Program(FTNRoute) error
	Remove(FTNRoute) error
}

type FTNBGPAdapter interface {
	Name() string
	Established() bool
	ApplyRoute(FTNRoute) error
	WithdrawRoute(FTNRoute) error
}

type FTNBFDState string

const (
	FTNBFDUp FTNBFDState = "up"
	FTNBFDDown FTNBFDState = "down"
	FTNBFDUnknown FTNBFDState = "unknown"
)

type FTNCoreNode struct {
	ID string `json:"id"`
	Healthy bool `json:"healthy"`
	BGPReady bool `json:"bgp_ready"`
	BFDState FTNBFDState `json:"bfd_state"`
}

type FTNCoreFailoverDecision struct {
	ActiveNode string `json:"active_node"`
	Failover bool `json:"failover"`
	Reason string `json:"reason"`
}

type FTNRoutedEngine struct {
	mu sync.RWMutex
	rib FTNRIB
	fib FTNFIB
	bgp []FTNBGPAdapter
	core []FTNCoreNode
}

func NewFTNRoutedEngine(rib FTNRIB, fib FTNFIB) *FTNRoutedEngine {
	return &FTNRoutedEngine{rib: rib, fib: fib, core: make([]FTNCoreNode, 0, 2)}
}

func (e *FTNRoutedEngine) RegisterBGPAdapter(a FTNBGPAdapter) {
	if a == nil { return }
	e.mu.Lock(); defer e.mu.Unlock()
	e.bgp = append(e.bgp, a)
}

func (e *FTNRoutedEngine) SetCoreNodes(nodes []FTNCoreNode) {
	e.mu.Lock(); defer e.mu.Unlock()
	e.core = append([]FTNCoreNode(nil), nodes...)
}

func (e *FTNRoutedEngine) ApplyApprovedRoute(in FTNRouteIntent, approved bool) error {
	if !approved { return errors.New("route approval required") }
	r, err := NormalizeFTNRoute(in.Route); if err != nil { return err }
	if e.rib == nil || e.fib == nil { return errors.New("RIB and FIB implementations are required") }
	if err := e.rib.Install(r); err != nil { return err }
	if err := e.fib.Program(r); err != nil { _ = e.rib.Withdraw(r); return err }
	for _, a := range e.bgp {
		if a.Established() { if err := a.ApplyRoute(r); err != nil { return err } }
	}
	return nil
}

func (e *FTNRoutedEngine) WithdrawApprovedRoute(in FTNRouteIntent, approved bool) error {
	if !approved { return errors.New("route approval required") }
	r, err := NormalizeFTNRoute(in.Route); if err != nil { return err }
	if e.rib == nil || e.fib == nil { return errors.New("RIB and FIB implementations are required") }
	if err := e.fib.Remove(r); err != nil { return err }
	if err := e.rib.Withdraw(r); err != nil { return err }
	for _, a := range e.bgp { if a.Established() { _ = a.WithdrawRoute(r) } }
	return nil
}

func EvaluateFTNCoreFailover(nodes []FTNCoreNode, current string) FTNCoreFailoverDecision {
	for _, n := range nodes {
		if n.ID == current && n.Healthy && n.BGPReady && n.BFDState != FTNBFDDown { return FTNCoreFailoverDecision{ActiveNode:current, Reason:"current_core_healthy"} }
	}
	for _, n := range nodes {
		if n.Healthy && n.BGPReady && n.BFDState != FTNBFDDown { return FTNCoreFailoverDecision{ActiveNode:n.ID, Failover:n.ID != current, Reason:"alternate_core_healthy"} }
	}
	return FTNCoreFailoverDecision{Reason:"no_healthy_core_available"}
}
