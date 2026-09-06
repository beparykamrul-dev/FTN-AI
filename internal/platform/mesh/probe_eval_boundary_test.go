package mesh

import "testing"

func TestEvaluateProbeExactThresholdIsUp(t *testing.T) {
 got:=EvaluateProbe(ProbeResult{Healthy:true,LatencyMS:10,LossPercent:1},HealthThresholds{MaxLatencyMS:10,MaxLossPercent:1})
 if got!=LinkUp {t.Fatalf("got %s",got)}
}
