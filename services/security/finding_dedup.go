package security

// DeduplicateFindings removes exact duplicate normalized findings deterministically.
func DeduplicateFindings(findings []Finding) []Finding {
	seen := make(map[string]struct{}, len(findings))
	out := make([]Finding, 0, len(findings))
	for _, f := range findings {
		key := f.Scanner + "|" + f.RuleID + "|" + f.Path + "|" + f.Message
		if _, ok := seen[key]; ok { continue }
		seen[key] = struct{}{}
		out = append(out, f)
	}
	return out
}
