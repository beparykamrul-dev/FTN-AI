package dns

import (
	"context"
	"fmt"
	"net"
	"sort"
	"time"
)

type LatencySample struct { RTT time.Duration `json:"rtt"`; Success bool `json:"success"` }
type LatencySummary struct { Min, Avg, P50, P95, Max time.Duration; LossPercent float64 }

func ProbeLatency(ctx context.Context, address string, count int, timeout time.Duration) (LatencySummary, error) {
	if count <= 0 || timeout <= 0 { return LatencySummary{}, fmt.Errorf("invalid latency probe parameters") }
	var samples []time.Duration; failures := 0
	for i := 0; i < count; i++ {
		select { case <-ctx.Done(): return LatencySummary{}, ctx.Err(); default: }
		start := time.Now(); d := net.Dialer{Timeout: timeout}
		conn, err := d.DialContext(ctx, "udp", address)
		if err != nil { failures++; continue }
		_ = conn.Close(); samples = append(samples, time.Since(start))
	}
	if len(samples) == 0 { return LatencySummary{LossPercent: 100}, fmt.Errorf("all latency probes failed") }
	sort.Slice(samples, func(i,j int) bool { return samples[i] < samples[j] })
	var total time.Duration; for _, v := range samples { total += v }
	p50 := samples[(len(samples)-1)*50/100]; p95 := samples[(len(samples)-1)*95/100]
	return LatencySummary{Min:samples[0], Avg:total/time.Duration(len(samples)), P50:p50, P95:p95, Max:samples[len(samples)-1], LossPercent:float64(failures)*100/float64(count)}, nil
}
