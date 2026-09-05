package model

import (
	"math"
	"testing"
)

func TestEvidenceValidRejectsNonFiniteValues(t *testing.T) {
	base := Evidence{ID:"e1", Kind:"latency", Source:"probe"}
	if !base.Valid() { t.Fatal("baseline evidence should be valid") }
	base.Value = math.NaN()
	if base.Valid() { t.Fatal("NaN evidence must be invalid") }
	base.Value = math.Inf(1)
	if base.Valid() { t.Fatal("Inf evidence must be invalid") }
}

func TestEvidenceNormalizeTrimsText(t *testing.T) {
	e := (Evidence{ID:" e1 ", Kind:" k ", Source:" s ", Summary:" note "}).Normalize()
	if e.ID != "e1" || e.Kind != "k" || e.Source != "s" || e.Summary != "note" { t.Fatalf("unexpected normalization: %#v", e) }
}
