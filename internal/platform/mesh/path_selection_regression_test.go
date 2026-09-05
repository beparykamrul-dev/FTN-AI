package mesh

import "testing"

func TestDefaultPathPolicyWeightsAreBounded(t *testing.T) {
	p := DefaultPathPolicy()
	if p.LatencyWeight < 0 || p.LossWeight < 0 || p.CapacityWeight < 0 || p.HealthWeight < 0 { t.Fatal("path policy weights must be non-negative") }
	if p.LatencyWeight+p.LossWeight+p.CapacityWeight+p.HealthWeight <= 0 { t.Fatal("path policy must have a positive total weight") }
}
