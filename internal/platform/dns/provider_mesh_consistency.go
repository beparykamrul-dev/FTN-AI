package dns

import (
	"context"
	"fmt"
)

// ProviderMeshConsistency coordinates multiple provider implementations behind
// one FTN control-plane operation. Provider-specific SDKs remain implementation
// details; callers do not select provider-specific methods.
type ProviderMeshConsistency struct {
	Registry *SDKRegistry
	AI       *DNSConsistencyAI
}

type MeshConsistencyResult struct {
	Zones map[string]Zone `json:"zones"`
	Report ConsistencyReport `json:"report"`
	Errors map[string]string `json:"errors,omitempty"`
}

func NewProviderMeshConsistency(registry *SDKRegistry, ai *DNSConsistencyAI) *ProviderMeshConsistency {
	return &ProviderMeshConsistency{Registry: registry, AI: ai}
}

func (m *ProviderMeshConsistency) ImportAndAnalyze(ctx context.Context, configs []ProviderConfig, zone string) (MeshConsistencyResult, error) {
	if m == nil || m.Registry == nil { return MeshConsistencyResult{}, fmt.Errorf("provider SDK registry is required") }
	if m.AI == nil { return MeshConsistencyResult{}, fmt.Errorf("consistency AI is required") }
	if ctx == nil { return MeshConsistencyResult{}, fmt.Errorf("context is required") }
	if zone == "" { return MeshConsistencyResult{}, fmt.Errorf("zone is required") }

	result := MeshConsistencyResult{Zones: map[string]Zone{}, Errors: map[string]string{}}
	normalized := map[string][]Record{}
	for _, cfg := range configs {
		sdk, err := m.Registry.Open(cfg)
		if err != nil { result.Errors[cfg.ID] = err.Error(); continue }
		if err := ValidateProviderSDK(ctx, sdk, "import"); err != nil { result.Errors[cfg.ID] = err.Error(); continue }
		z, err := sdk.ImportZone(ctx, zone)
		if err != nil { result.Errors[cfg.ID] = err.Error(); continue }
		result.Zones[cfg.ID] = z
		normalized[cfg.ID] = z.Records
	}

	report, err := m.AI.Analyze(ctx, normalized)
	if err != nil { return MeshConsistencyResult{}, err }
	result.Report = report
	if len(result.Errors) == 0 { result.Errors = nil }
	return result, nil
}
