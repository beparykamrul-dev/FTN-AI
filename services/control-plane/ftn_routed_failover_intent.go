package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type FTNFailoverIntent struct {
	CurrentNode string `json:"current_node"`
	TargetNode string `json:"target_node"`
	Reason string `json:"reason"`
	DecisionHash string `json:"decision_hash"`
}

// BuildFTNFailoverIntent converts a health decision into an approval-gated
// durable intent. It never performs the failover itself.
func BuildFTNFailoverIntent(decision FTNCoreFailoverDecision) (FTNFailoverIntent, error) {
	if decision.Failover == false {
		return FTNFailoverIntent{}, fmt.Errorf("failover_not_required")
	}
	if strings.TrimSpace(decision.ActiveNode) == "" {
		return FTNFailoverIntent{}, fmt.Errorf("target_core_required")
	}
	payload := struct{ Target, Reason string }{decision.ActiveNode, strings.TrimSpace(decision.Reason)}
	b, _ := json.Marshal(payload)
	h := sha256.Sum256(b)
	return FTNFailoverIntent{
		TargetNode: decision.ActiveNode,
		Reason: decision.Reason,
		DecisionHash: hex.EncodeToString(h[:]),
	}, nil
}
