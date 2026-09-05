package dns

import (
	"testing"
	"time"
)

func TestServerHealthRejectsInvalidRates(t *testing.T) {
	h := ServerHealth{NodeID:"n1", Resolver:"r1", ObservedAt:time.Now(), LatencyMs:1, QPS:1, ServfailRate:101}
	if h.Valid() { t.Fatal("servfail rate above 100 must be invalid") }
	h.ServfailRate = 0
	if !h.Valid() { t.Fatal("valid DNS health must pass") }
}
