package ftndns

import "testing"

func TestDefaultEdgeCachePolicyHasSafePositiveWindows(t *testing.T) {
	p := DefaultEdgeCachePolicy()
	if p.DefaultTTL <= 0 || p.StaleWindow < 0 { t.Fatalf("invalid cache windows: %#v", p) }
	if !p.Enabled { t.Fatal("edge cache must be enabled by default") }
}
