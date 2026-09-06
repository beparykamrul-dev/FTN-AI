package mesh

import "testing"

func TestEvaluateProbeUnhealthyIsDown(t *testing.T) {
 got:=EvaluateProbe(ProbeResult{Healthy:false,LossPercent:0},HealthThresholds{MaxLatencyMS:10,MaxLossPercent:1})
 if got!=LinkDown {t.Fatalf("got %s",got)}
}
