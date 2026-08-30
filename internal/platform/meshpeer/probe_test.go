package meshpeer

import (
	"context"
	"testing"
	"time"
)

func TestProbeRejectsEmptyEndpoint(t *testing.T) {
	result := Probe(context.Background(), "", "", nil, time.Second)
	if result.Reachable || result.ReadOK || result.Error == "" {
		t.Fatalf("expected rejected empty endpoint, got %+v", result)
	}
}

func TestProbeDoesNotClaimUnreachablePeerHealthy(t *testing.T) {
	result := Probe(context.Background(), "127.0.0.1:1", "", nil, 100*time.Millisecond)
	if result.Reachable || result.ReadOK {
		t.Fatalf("expected unreachable peer, got %+v", result)
	}
}
