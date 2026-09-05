package graph

import (
	"testing"
	"github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"
)

func TestDependenciesOfIsDeterministicAndDeduplicated(t *testing.T) {
	g := Graph{Dependencies:[]model.Dependency{{From:"a",To:"b",Kind:"rpc"},{From:"a",To:"b",Kind:"rpc"},{From:"c",To:"b",Kind:"db"}}}
	out := g.DependenciesOf("b")
	if len(out) != 2 { t.Fatalf("expected two unique dependencies, got %d", len(out)) }
	if out[0].From != "a" || out[1].From != "c" { t.Fatalf("unexpected ordering: %#v", out) }
}
