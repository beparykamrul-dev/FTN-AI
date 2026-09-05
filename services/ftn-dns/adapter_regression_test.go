package dns

import (
	"context"
	"testing"
)

type regressionAdapter struct{}
func (regressionAdapter) Name() string { return "regression" }
func (regressionAdapter) Health(context.Context) (Health,error) { return Health{Reachable:true,Secure:true,LatencyMS:1},nil }
func (regressionAdapter) Query(context.Context,string,string) (Response,error) { return Response{},nil }

func TestAdapterHealthRejectsInvalidMetrics(t *testing.T) {
	if (Health{Reachable:true, Secure:true, LatencyMS:-1}).Valid() { t.Fatal("negative latency must be invalid") }
}
