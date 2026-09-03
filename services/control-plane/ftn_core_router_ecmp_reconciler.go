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
	parts := make([]string, 0, len(p.Install)+len(p.Withdraw)+2)
	parts = append(parts, p.Prefix, p.VRF)
	for _, in := range p.Install { parts = append(parts, "install", in.Route.Prefix, in.Route.NextHop, in.Route.Protocol, in.Route.VRF) }
	for _, in := range p.Withdraw { parts = append(parts, "withdraw", in.Route.Prefix, in.Route.NextHop, in.Route.Protocol, in.Route.VRF) }
	s := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(s[:])
}

func BuildFTNECMPReconcilePlan(current, desired FTNECMPSelection) (FTNECMPReconcilePlan, error) {
	if strings.TrimSpace(desired.Prefix) == "" || strings.TrimSpace(desired.VRF) == "" { return FTNECMPReconcilePlan{}, errors.New("ecmp_prefix_vrf_required") }
	if current.Prefix != "" && (current.Prefix != desired.Prefix || current.VRF != desired.VRF) { return FTNECMPReconcilePlan{}, errors.New("ecmp_prefix_vrf_mismatch") }
	cur := append([]FTNRouteCandidate(nil), current.Paths...)
	des := append([]FTNRouteCandidate(nil), desired.Paths...)
	key := func(c FTNRouteCandidate) string { return c.PathID + "|" + c.Route.Prefix + "|" + c.Route.NextHop }
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

func (e *FTNRoutedEngine) ApplyApprovedECMPReconcile(plan FTNECMPReconcilePlan, approvalHash string, approved bool) error {
	if !approved { return errors.New("ecmp approval required") }
	if !plan.RequiresApproval || !plan.PreChangeSnapshot || !plan.PostChangeVerify || !plan.RollbackWhenSafe { return errors.New("ecmp safety gates required") }
	if strings.TrimSpace(approvalHash) == "" || approvalHash != plan.ApprovalHash || approvalHash != HashFTNECMPReconcilePlan(plan) { return errors.New("ecmp approval binding mismatch") }
	applied := make([]FTNRouteIntent, 0, len(plan.Install)+len(plan.Withdraw))
	rollback := func() { for i:=len(applied)-1; i>=0; i-- { in:=applied[i]; if in.Action=="install" { _=e.WithdrawApprovedRoute(FTNRouteIntent{Action:"withdraw",Route:in.Route,Reason:"ecmp_safe_rollback"},true) } else { _=e.ApplyApprovedRoute(FTNRouteIntent{Action:"install",Route:in.Route,Reason:"ecmp_safe_rollback"},true) } } }
	for _, in := range plan.Withdraw { if err:=e.WithdrawApprovedRoute(in,true); err!=nil { rollback(); return fmt.Errorf("ecmp withdraw: %w",err) }; applied=append(applied,in) }
	for _, in := range plan.Install { if err:=e.ApplyApprovedRoute(in,true); err!=nil { rollback(); return fmt.Errorf("ecmp install: %w",err) }; applied=append(applied,in) }
	return nil
}
