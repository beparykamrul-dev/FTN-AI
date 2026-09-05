package agent

import "sort"

type SummaryItem struct {
	Priority int
	Title    string
	Text     string
	Details  string
}

type Summary struct {
	Important []SummaryItem
	Details   []SummaryItem
}

// BuildSummary separates the high-value operator/customer view from drill-down details.
func BuildSummary(items []SummaryItem) Summary {
	important := append([]SummaryItem(nil), items...)
	details := append([]SummaryItem(nil), items...)
	sort.SliceStable(important, func(i, j int) bool {
		return important[i].Priority > important[j].Priority
	})
	sort.SliceStable(details, func(i, j int) bool {
		return details[i].Priority > details[j].Priority
	})
	if len(important) > 5 {
		important = important[:5]
	}
	return Summary{Important: important, Details: details}
}
