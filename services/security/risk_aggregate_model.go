package security

// RiskAggregate is the normalized security risk summary used by policy evaluation.
type RiskAggregate struct {
	Critical uint32
	High uint32
	Medium uint32
	Low uint32
	Total uint32
}

func AggregateRisk(findings []Finding) RiskAggregate {
	var r RiskAggregate
	for _, f := range findings {
		switch f.Severity {
		case "critical": r.Critical++
		case "high": r.High++
		case "medium": r.Medium++
		case "low": r.Low++
		}
		r.Total++
	}
	return r
}
