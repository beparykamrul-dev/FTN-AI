package kernel

import (
	"context"
	"testing"
)

func TestRegistryFailsClosedForInvalidExecution(t *testing.T) {
	r, err := NewRegistry()
	if err != nil { t.Fatal(err) }
	if _, err := r.Execute(context.Background(), ToolRequest{Tool:"", Operation:"run"}); err == nil { t.Fatal("empty tool must fail") }
	if _, err := r.Execute(nil, ToolRequest{Tool:"x", Operation:"run"}); err == nil { t.Fatal("nil context must fail") }
	if got := r.Tools(); len(got) != 0 { t.Fatalf("tools=%v, want empty", got) }
}
