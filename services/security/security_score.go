package security

// SecurityScore summarizes normalized findings into a simple 0-100 score.
func SecurityScore(findings []Finding) int {
	score := 100
	for _, f := range findings {
		switch f.Severity {
		case "critical": score -= 40
		case "high": score -= 25
		case "medium": score -= 10
		case "low": score -= 3
		}
	}
	if score < 0 { return 0 }
	return score
}
