package dns

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

type ConsistencyIssue struct {
	Zone string `json:"zone"`
	Record string `json:"record"`
	Kind string `json:"kind"`
	Severity string `json:"severity"`
	Details string `json:"details"`
}

type ConsistencyReport struct {
	Consistent bool `json:"consistent"`
	Issues []ConsistencyIssue `json:"issues,omitempty"`
}

// DNSConsistencyAI provides deterministic evidence-based reconciliation
// recommendations. It deliberately does not mutate providers automatically.
type DNSConsistencyAI struct{}

func NewDNSConsistencyAI() *DNSConsistencyAI { return &DNSConsistencyAI{} }

func (a *DNSConsistencyAI) Analyze(ctx context.Context, zones map[string][]Record) (ConsistencyReport, error) {
	if zones == nil { return ConsistencyReport{Consistent: true}, nil }
	select { case <-ctx.Done(): return ConsistencyReport{}, ctx.Err(); default: }
	report := ConsistencyReport{Consistent: true}
	for zone, records := range zones {
		seen := make(map[string]string)
		for _, r := range records {
			name := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(r.Name), "."))
			typ := strings.ToUpper(strings.TrimSpace(r.Type))
			key := name + "|" + typ
			value := strings.TrimSpace(r.Value)
			if previous, ok := seen[key]; ok && previous != value {
				report.Consistent = false
				report.Issues = append(report.Issues, ConsistencyIssue{Zone: zone, Record: key, Kind: "conflicting_values", Severity: "high", Details: "multiple values observed for the same normalized record key"})
			} else { seen[key] = value }
		}
	}
	sort.Slice(report.Issues, func(i, j int) bool { return report.Issues[i].Zone+report.Issues[i].Record < report.Issues[j].Zone+report.Issues[j].Record })
	return report, nil
}

func (a *DNSConsistencyAI) ValidateReport(report ConsistencyReport) error {
	if report.Consistent && len(report.Issues) > 0 { return fmt.Errorf("consistent report cannot contain issues") }
	return nil
}
