package dns

import (
	"context"
	"testing"
)

type nextDNSAdapter struct{}
func (nextDNSAdapter) Name() string { return "next" }
func (nextDNSAdapter) Health(context.Context) (Health, error) { return Health{Reachable:true, Secure:true, LatencyMS:1, LossRatio:0}, nil }
func (nextDNSAdapter) Query(context.Context, string, string) (Response, error) { return Response{}, nil }

func TestDNSAdapterContractCanBeImplemented(t *testing.T) {
	var a Adapter = nextDNSAdapter{}
	if a.Name() != "next" { t.Fatal("adapter name mismatch") }
}
