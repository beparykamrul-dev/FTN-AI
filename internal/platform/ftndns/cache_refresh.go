package ftndns

import "time"

type RefreshDecision struct {
	Refresh bool
	StaleAllowed bool
}

// DecideRefresh provides stale-while-revalidate semantics for hot DNS data.
func DecideRefresh(expiresAt, now time.Time, refreshBefore, staleWindow time.Duration) RefreshDecision {
	if !expiresAt.After(now) {
		return RefreshDecision{Refresh:true, StaleAllowed: now.Before(expiresAt.Add(staleWindow))}
	}
	return RefreshDecision{Refresh: !expiresAt.After(now.Add(refreshBefore))}
}
