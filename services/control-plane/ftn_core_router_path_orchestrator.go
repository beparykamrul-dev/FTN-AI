package main

import (
	"errors"
	"strings"
	"time"
)

// FTNCoreRouterPathOrchestrator is the read/plan boundary between measured
// traffic quality and the canonical routed engine. It never mutates routing
// state without an explicitly approved ECMP reconciliation plan.
type FTNCoreRouterPathOrchestrator struct {
	engine *FTNRoutedEngine
	controllers map[string]*TrafficPathController
}

func NewFTNCoreRouterPathOrchestrator(engine *FTNRoutedEngine) *FTNCoreRouterPathOrchestrator {
	return &FTNCoreRouterPathOrchestrator{engine: engine, controllers: make(map[string]*TrafficPathController)}
}

func (o *FTNCoreRouterPathOrchestrator) controller(serviceID string) *TrafficPathController {
	if c, ok := o.controllers[serviceID]; ok { return c }
	c := &TrafficPathController{}
	o.controllers[serviceID] = c
	return c
}

// SelectPath converts measured service observations into a stable path choice.
// TrafficPathController supplies the existing 5-second hysteresis/failover policy.
func (o *FTNCoreRouterPathOrchestrator) SelectPath(observations []TrafficPathObservation, service TrafficServicePolicy, now time.Time) (TrafficDecision, bool) {
	if o == nil { return TrafficDecision{}, false }
	return o.controller(strings.TrimSpace(service.ID)).Decide(observations, service, now)
}

// BuildECMPPlan filters candidates through health/BFD eligibility and delegates
// state reconciliation to the canonical ECMP planner. This method is plan-only.
func (o *FTNCoreRouterPathOrchestrator) BuildECMPPlan(current FTNECMPSelection, candidates []FTNRouteCandidate, maxPaths int) (FTNECMPReconcilePlan, error) {
	if o == nil || o.engine == nil { return FTNECMPReconcilePlan{}, errors.New("routed_engine_required") }
	selection, err := SelectFTNECMPPaths(candidates, maxPaths)
	if err != nil { return FTNECMPReconcilePlan{}, err }
	if current.Prefix != "" && current.VRF != "" {
		selection.Prefix, selection.VRF = current.Prefix, current.VRF
	}
	return BuildFTNECMPReconcilePlan(current, selection)
}

// ApplyApprovedECMPPlan is the sole orchestration entry point for ECMP mutation.
// The routed engine remains the execution authority and enforces approval binding.
func (o *FTNCoreRouterPathOrchestrator) ApplyApprovedECMPPlan(plan FTNECMPReconcilePlan, approvalHash string, approved bool) error {
	if o == nil || o.engine == nil { return errors.New("routed_engine_required") }
	return o.engine.ApplyApprovedECMPReconcile(plan, approvalHash, approved)
}
