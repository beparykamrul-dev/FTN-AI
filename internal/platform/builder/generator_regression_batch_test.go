package builder

import "testing"

func TestGeneratorRejectsNilReceiver(t *testing.T) {
	var g *Generator
	if g.Generate(Manifest{}, "/tmp/ftn-test") == nil { t.Fatal("nil generator must fail closed") }
}
