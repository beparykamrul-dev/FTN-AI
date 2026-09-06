package mesh

import "testing"

func TestEvaluateProbeFullLossIsDown(t *testing.T) {
 got:=EvaluateProbe(ProbeResult{Healthy:true,LossPercent:100},HealthThresholds{})
 if got!=LinkDown {t.Fatalf("got %s",got)}
}
