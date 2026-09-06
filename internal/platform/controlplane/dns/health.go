package dns
import("math";"strings";"time")
type ServerHealth struct{NodeID string `json:"node_id"`;Resolver string `json:"resolver"`;Healthy bool `json:"healthy"`;LatencyMs float64 `json:"latency_ms"`;QPS float64 `json:"qps"`;ServfailRate float64 `json:"servfail_rate"`;ObservedAt time.Time `json:"observed_at"`}
func(h ServerHealth)Valid()bool{node,resolver:=strings.TrimSpace(h.NodeID),strings.TrimSpace(h.Resolver);return node!=""&&resolver!=""&&len(node)<=256&&len(resolver)<=256&&!h.ObservedAt.IsZero()&&!h.ObservedAt.After(time.Now().UTC())&&h.LatencyMs>=0&&!math.IsNaN(h.LatencyMs)&&!math.IsInf(h.LatencyMs,0)&&h.QPS>=0&&!math.IsNaN(h.QPS)&&!math.IsInf(h.QPS,0)&&h.ServfailRate>=0&&h.ServfailRate<=100&&!math.IsNaN(h.ServfailRate)&&!math.IsInf(h.ServfailRate,0)}
func Available(h ServerHealth,now time.Time,maxAge time.Duration)bool{if now.IsZero(){now=time.Now().UTC()};return h.Valid()&&maxAge>=0&&!now.Before(h.ObservedAt)&&now.Sub(h.ObservedAt)<=maxAge&&h.Healthy}
