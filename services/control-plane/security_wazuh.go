package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type WazuhAlert struct {
	ID        string `json:"id"`
	RuleID    string `json:"rule_id"`
	Severity  string `json:"severity"`
	AgentRef  string `json:"agent_ref"`
	Timestamp string `json:"timestamp"`
	Summary   string `json:"summary"`
}

type WazuhActionClass string

const (
	WazuhObserve       WazuhActionClass = "observe"
	WazuhCorrelate     WazuhActionClass = "correlate"
	WazuhAlertAction   WazuhActionClass = "alert"
	WazuhRecommend     WazuhActionClass = "recommend"
	WazuhBackup        WazuhActionClass = "backup"
	WazuhHealthRecover WazuhActionClass = "health-recover"
	WazuhConfiguration WazuhActionClass = "configuration-change"
	WazuhFirewall      WazuhActionClass = "firewall-change"
	WazuhRouting       WazuhActionClass = "routing-change"
	WazuhCredential    WazuhActionClass = "credential-change"
	WazuhDisable       WazuhActionClass = "service-disable"
	WazuhDelete        WazuhActionClass = "delete"
	WazuhDestructive   WazuhActionClass = "destructive-recovery"
)

func WazuhActionRequiresApproval(action WazuhActionClass) bool {
	switch action {
	case WazuhConfiguration, WazuhFirewall, WazuhRouting, WazuhCredential, WazuhDisable, WazuhDelete, WazuhDestructive:
		return true
	default:
		return false
	}
}

func NormalizeWazuhSeverity(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "critical", "high", "medium", "low", "info":
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return "info"
	}
}

func WazuhAlertHash(a WazuhAlert) string {
	payload := strings.Join([]string{a.ID, a.RuleID, NormalizeWazuhSeverity(a.Severity), a.AgentRef, a.Timestamp, a.Summary}, "|")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
