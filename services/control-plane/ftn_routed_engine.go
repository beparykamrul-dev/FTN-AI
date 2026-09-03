package main

import (
	"errors"
	"fmt"
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

func (e *FTNRoutedEngine) adapters() []FTNBGPAdapter {
	e.mu.RLock(); defer e.mu.RUnlock()
	return append([]FTNBGPAdapter(nil), e.bgp...)
}

func (e *FTNRoutedEngine) ApplyApprovedRoute(in FTNRouteIntent, approved bool) error {
	if !approved { return errors.New("route approval required") }
	r, err := NormalizeFTNRoute(in.Route); if err != nil { return err }
	if e.rib == nil || e.fib == nil { return errors.New("RIB and FIB implementations are required") }

	if err := e.rib.Install(r); err != nil { return err }
	if err := e.fib.Program(r); err != nil {
		_ = e.rib.Withdraw(r)
		return err
	}

	appliedBGP := make([]FTNBGPAdapter, 0)
	rollback := func() {
		for i := len(appliedBGP)-1; i >= 0; i-- {
			_ = appliedBGP[i].WithdrawRoute(r)
		}
		_ = e.fib.Remove(r)
		_ = e.rib.Withdraw(r)
	}
	for _, a := range e.adapters() {
		if !a.Established() { continue }
		if err := a.ApplyRoute(r); err != nil {
			rollback()
			return fmt.Errorf("bgp adapter %s: %w", a.Name(), err)
		}
		appliedBGP = append(appliedBGP, a)
	}
	return nil
}

func (e *FTNRoutedEngine) WithdrawApprovedRoute(in FTNRouteIntent, approved bool) error {
	if !approved { return errors.New("route approval required") }
	r, err := NormalizeFTNRoute(in.Route); if err != nil { return err }
	if e.rib == nil || e.fib == nil { return errors.New("RIB and FIB implementations are required") }

	adapters := e.adapters()
	removedBGP := make([]FTNBGPAdapter, 0)
	rollback := func() {
		_ = e.rib.Install(r)
		_ = e.fib.Program(r)
		for _, a := range removedBGP { _ = a.ApplyRoute(r) }
	}
	for _, a := range adapters {
		if !a.Established() { continue }
		if err := a.WithdrawRoute(r); err != nil {
			rollback()
			return fmt.Errorf("bgp adapter %s withdraw: %w", a.Name(), err)
		}
		removedBGP = append(removedBGP, a)
	}
	if err := e.fib.Remove(r); err != nil {
		rollback()
		return err
	}
	if err := e.rib.Withdraw(r); err != nil {
		rollback()
		return err
	}
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
