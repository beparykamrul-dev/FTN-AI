package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type FTNECMPReconcilePlan struct {
	Prefix string `json:"prefix"`
	VRF string `json:"vrf"`
	Current FTNECMPSelection `json:"current"`
	Desired FTNECMPSelection `json:"desired"`
	Install []FTNRouteIntent `json:"install"`
	Withdraw []FTNRouteIntent `json:"withdraw"`
	ApprovalHash string `json:"approval_hash"`
	RequiresApproval bool `json:"requires_approval"`
	PreChangeSnapshot bool `json:"pre_change_snapshot"`
	PostChangeVerify bool `json:"post_change_verify"`
	RollbackWhenSafe bool `json:"rollback_when_safe"`
}

func HashFTNECMPReconcilePlan(p FTNECMPReconcilePlan) string {
	parts := make([]string, 0, len(p.Install)+len(p.Withdraw)+8)
	parts = append(parts, p.Prefix, p.VRF,
		fmt.Sprintf("requires_approval=%t", p.RequiresApproval),
		fmt.Sprintf("pre_snapshot=%t", p.PreChangeSnapshot),
		fmt.Sprintf("post_verify=%t", p.PostChangeVerify),
		fmt.Sprintf("rollback=%t", p.RollbackWhenSafe))
	appendCandidate := func(kind string, c FTNRouteCandidate) {
		parts = append(parts, kind, c.PathID, c.Route.Prefix, c.Route.NextHop, c.Route.Protocol, c.Route.VRF,
			fmt.Sprintf("score=%.6f", c.Score), string(c.BFDState), fmt.Sprintf("healthy=%t", c.Healthy))
	}
	current := append([]FTNRouteCandidate(nil), p.Current.Paths...)
	desired := append([]FTNRouteCandidate(nil), p.Desired.Paths...)
	sort.SliceStable(current, func(i, j int) bool { return current[i].PathID < current[j].PathID })
	sort.SliceStable(desired, func(i, j int) bool { return desired[i].PathID < desired[j].PathID })
	for _, c := range current { appendCandidate("current", c) }
	for _, c := range desired { appendCandidate("desired", c) }
	for _, in := range p.Install { parts = append(parts, "install", in.Action, in.Reason, in.Route.Prefix, in.Route.NextHop, in.Route.Protocol, in.Route.VRF) }
	for _, in := range p.Withdraw { parts = append(parts, "withdraw", in.Action, in.Reason, in.Route.Prefix, in.Route.NextHop, in.Route.Protocol, in.Route.VRF) }
	s := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(s[:])
}

func BuildFTNECMPReconcilePlan(current, desired FTNECMPSelection) (FTNECMPReconcilePlan, error) {
	if strings.TrimSpace(desired.Prefix) == "" || strings.TrimSpace(desired.VRF) == "" { return FTNECMPReconcilePlan{}, errors.New("ecmp_prefix_vrf_required") }
	if current.Prefix != "" && (current.Prefix != desired.Prefix || current.VRF != desired.VRF) { return FTNECMPReconcilePlan{}, errors.New("ecmp_prefix_vrf_mismatch") }
	cur := append([]FTNRouteCandidate(nil), current.Paths...)
	des := append([]FTNRouteCandidate(nil), desired.Paths...)
	key := func(c FTNRouteCandidate) string { return c.PathID + "|" + c.Route.Prefix + "|" + c.Route.NextHop + "|" + c.Route.Protocol + "|" + c.Route.VRF }
	curMap := make(map[string]FTNRouteCandidate, len(cur)); for _, c := range cur { curMap[key(c)] = c }
	desMap := make(map[string]FTNRouteCandidate, len(des)); for _, c := range des { desMap[key(c)] = c }
	plan := FTNECMPReconcilePlan{Prefix:desired.Prefix, VRF:desired.VRF, Current:current, Desired:desired, Install:make([]FTNRouteIntent,0), Withdraw:make([]FTNRouteIntent,0), RequiresApproval:true, PreChangeSnapshot:true, PostChangeVerify:true, RollbackWhenSafe:true}
	for k, c := range desMap { if _, ok := curMap[k]; !ok { plan.Install = append(plan.Install, FTNRouteIntent{Action:"install", Route:c.Route, Reason:"ecmp_reconcile_install"}) } }
	for k, c := range curMap { if _, ok := desMap[k]; !ok { plan.Withdraw = append(plan.Withdraw, FTNRouteIntent{Action:"withdraw", Route:c.Route, Reason:"ecmp_reconcile_withdraw"}) } }
	sort.Slice(plan.Install, func(i,j int) bool { return plan.Install[i].Route.NextHop < plan.Install[j].Route.NextHop })
	sort.Slice(plan.Withdraw, func(i,j int) bool { return plan.Withdraw[i].Route.NextHop < plan.Withdraw[j].Route.NextHop })
	plan.ApprovalHash = HashFTNECMPReconcilePlan(plan)
	return plan, nil
}

func routesMatchECMPSelection(routes []FTNRoute, desired FTNECMPSelection) bool {
	want := make(map[string]struct{}, len(desired.Paths))
	for _, c := range desired.Paths { want[c.Route.Prefix+"|"+c.Route.NextHop+"|"+c.Route.Protocol+"|"+c.Route.VRF] = struct{}{} }
	got := make(map[string]struct{}, len(routes))
	for _, r := range routes { got[r.Prefix+"|"+r.NextHop+"|"+r.Protocol+"|"+r.VRF] = struct{}{} }
	if len(want) != len(got) { return false }
	for k := range want { if _, ok := got[k]; !ok { return false } }
	return true
}

func (e *FTNRoutedEngine) verifyECMPLiveSelection(selection FTNECMPSelection) error {
	if e.rib == nil { return errors.New("RIB implementation is required") }
	routes, err := e.rib.Lookup(selection.Prefix, selection.VRF)
	if err != nil { return fmt.Errorf("ecmp rib lookup: %w", err) }
	if !routesMatchECMPSelection(routes, selection) { return errors.New("ecmp_post_change_verification_failed") }
	return nil
}

func (e *FTNRoutedEngine) ApplyApprovedECMPReconcile(plan FTNECMPReconcilePlan, approvalHash string, approved bool) error {
	if !approved { return errors.New("ecmp approval required") }
	if !plan.RequiresApproval || !plan.PreChangeSnapshot || !plan.PostChangeVerify || !plan.RollbackWhenSafe { return errors.New("ecmp safety gates required") }
	if strings.TrimSpace(approvalHash) == "" || approvalHash != plan.ApprovalHash || approvalHash != HashFTNECMPReconcilePlan(plan) { return errors.New("ecmp approval binding mismatch") }
	if plan.Current.Prefix != "" {
		if err := e.verifyECMPLiveSelection(plan.Current); err != nil { return fmt.Errorf("ecmp pre-change verification failed: %w", err) }
	}
	applied := make([]FTNRouteIntent, 0, len(plan.Install)+len(plan.Withdraw))
	rollback := func() {
		for i:=len(applied)-1; i>=0; i-- {
			in:=applied[i]
			if in.Action=="install" { _=e.WithdrawApprovedRoute(FTNRouteIntent{Action:"withdraw",Route:in.Route,Reason:"ecmp_safe_rollback"},true) } else { _=e.ApplyApprovedRoute(FTNRouteIntent{Action:"install",Route:in.Route,Reason:"ecmp_safe_rollback"},true) }
		}
	}
	for _, in := range plan.Withdraw { if err:=e.WithdrawApprovedRoute(in,true); err!=nil { rollback(); return fmt.Errorf("ecmp withdraw: %w",err) }; applied=append(applied,in) }
	for _, in := range plan.Install { if err:=e.ApplyApprovedRoute(in,true); err!=nil { rollback(); return fmt.Errorf("ecmp install: %w",err) }; applied=append(applied,in) }
	if err := e.verifyECMPLiveSelection(plan.Desired); err != nil { rollback(); return err }
	return nil
}
