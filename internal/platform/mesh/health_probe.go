package mesh

import("math";"strings";"time")
type ProbeResult struct{LinkID string `json:"link_id"`;Healthy bool `json:"healthy"`;LatencyMS float64 `json:"latency_ms"`;LossPercent float64 `json:"loss_percent"`;ObservedAt time.Time `json:"observed_at"`}
type HealthThresholds struct{MaxLatencyMS float64;MaxLossPercent float64}
func EvaluateProbe(r ProbeResult,t HealthThresholds)LinkState{if strings.TrimSpace(r.LinkID)==""||math.IsNaN(r.LatencyMS)||math.IsInf(r.LatencyMS,0)||math.IsNaN(r.LossPercent)||math.IsInf(r.LossPercent,0)||r.LatencyMS<0||r.LossPercent<0||r.LossPercent>100{return LinkDown};if !r.Healthy||r.LossPercent>=100{return LinkDown};if t.MaxLatencyMS<0||t.MaxLossPercent<0{return LinkDown};if r.LatencyMS>t.MaxLatencyMS||r.LossPercent>t.MaxLossPercent{return LinkDegraded};return LinkUp}
func ProbeAge(now,observed time.Time)time.Duration{if observed.IsZero(){return time.Duration(1<<63-1)};if now.IsZero(){now=time.Now().UTC()}else{now=now.UTC()};observed=observed.UTC();if now.Before(observed){return 0};return now.Sub(observed)}
