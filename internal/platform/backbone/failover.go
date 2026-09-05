package backbone

import("strings";"time")
type HealthState struct{Mode string `json:"mode"`;Healthy bool `json:"healthy"`;LatencyMS uint32 `json:"latency_ms"`;LossPercent float32 `json:"loss_percent"`;CheckedAt time.Time `json:"checked_at"`}
type FailoverDecision struct{From string `json:"from"`;To string `json:"to"`;Reason string `json:"reason"`;RequiresApproval bool `json:"requires_approval"`}
// Evaluate produces a failover proposal only. It never switches traffic.
func Evaluate(primary,secondary HealthState)(FailoverDecision,bool){primary.Mode=strings.TrimSpace(primary.Mode);secondary.Mode=strings.TrimSpace(secondary.Mode);if primary.Healthy{return FailoverDecision{},false};if secondary.Mode==""||!secondary.Healthy{return FailoverDecision{From:primary.Mode,To:secondary.Mode,Reason:"no healthy standby path",RequiresApproval:true},false};return FailoverDecision{From:primary.Mode,To:secondary.Mode,Reason:"primary health check failed",RequiresApproval:true},true}
