package servicefoundry

import "testing"

func TestRankProviders(t *testing.T) {
	a := ProviderAnalysis{ServiceName: "example", Candidates: []ProviderEvidence{
		{Provider: "A", Score: 70}, {Provider: "B", Score: 95}, {Provider: "C", Score: 80},
	}}
	got := RankProviders(a)
	if got.Best != "B" || got.Candidates[0].Provider != "B" { t.Fatalf("unexpected ranking: %+v", got) }
}

func TestSecurityGate(t *testing.T) {
	g := SecurityGate{}
	if g.Validate() == nil { t.Fatal("incomplete gate unexpectedly passed") }
	g = SecurityGate{ThreatModel:true, AuthN:true, AuthZ:true, InputValidation:true, SecretsIsolation:true, Audit:true, ResourceLimits:true, StaticAnalysis:true, Tests:true, DependencyAudit:true}
	if err := g.Validate(); err != nil { t.Fatalf("complete gate failed: %v", err) }
}
