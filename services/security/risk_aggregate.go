package security

// RiskAggregate provides a normalized deployment-risk summary.
type RiskAggregate struct {
	Critical uint32
	High     uint32
	Medium   uint32
	Low      uint32
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
	}
	return r
}
