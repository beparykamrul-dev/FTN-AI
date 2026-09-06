package mesh

import "testing"

func TestEvaluateProbeThresholdBreachIsDegraded(t *testing.T) {
 got:=EvaluateProbe(ProbeResult{Healthy:true,LatencyMS:20,LossPercent:0},HealthThresholds{MaxLatencyMS:10,MaxLossPercent:1})
 if got!=LinkDegraded {t.Fatalf("got %s",got)}
}
