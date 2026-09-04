package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"net/netip"
	"sort"
	"strings"
)

type FTNLocalServiceEndpoint struct {
	ServiceID string `json:"service_id"`
	POPID string `json:"pop_id"`
	Address string `json:"address"`
	Healthy bool `json:"healthy"`
	LatencyMS float64 `json:"latency_ms"`
	CapacityPercent float64 `json:"capacity_percent"`
}

type FTNGlobalLocalRouteDecision struct {
	ServiceID string `json:"service_id"`
	POPID string `json:"pop_id"`
	Address string `json:"address"`
	Reason string `json:"reason"`
	DecisionHash string `json:"decision_hash"`
}

// SelectFTNLocalService chooses only from normalized, healthy FTN-owned
// service endpoints. It does not alter routing state.
func SelectFTNLocalService(endpoints []FTNLocalServiceEndpoint) (FTNGlobalLocalRouteDecision, error) {
	candidates := make([]FTNLocalServiceEndpoint, 0, len(endpoints))
	for _, e := range endpoints {
		if strings.TrimSpace(e.ServiceID)=="" || strings.TrimSpace(e.POPID)=="" || !e.Healthy { continue }
		if _, err := netip.ParseAddr(strings.TrimSpace(e.Address)); err != nil { continue }
		if !isFiniteNonNegative(e.LatencyMS) || !isFinitePercent(e.CapacityPercent) { continue }
		candidates = append(candidates, e)
	}
	if len(candidates)==0 { return FTNGlobalLocalRouteDecision{}, fmt.Errorf("no_healthy_local_service") }
	sort.SliceStable(candidates, func(i,j int) bool {
		if candidates[i].LatencyMS != candidates[j].LatencyMS { return candidates[i].LatencyMS < candidates[j].LatencyMS }
		if candidates[i].CapacityPercent != candidates[j].CapacityPercent { return candidates[i].CapacityPercent < candidates[j].CapacityPercent }
		return candidates[i].POPID < candidates[j].POPID
	})
	chosen:=candidates[0]
	v:=fmt.Sprintf("%s|%s|%s|%.3f|%.3f",chosen.ServiceID,chosen.POPID,chosen.Address,chosen.LatencyMS,chosen.CapacityPercent)
	h:=sha256.Sum256([]byte(v))
	return FTNGlobalLocalRouteDecision{ServiceID:chosen.ServiceID,POPID:chosen.POPID,Address:chosen.Address,Reason:"lowest-healthy-latency-capacity",DecisionHash:hex.EncodeToString(h[:])},nil
}

func isFiniteNonNegative(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v,0) && v >= 0 }
func isFinitePercent(v float64) bool { return isFiniteNonNegative(v) && v <= 100 }
