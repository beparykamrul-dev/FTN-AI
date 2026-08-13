package proxy

import "time"

// FailoverPolicy controls when FTN Proxy may move a request to another
// upstream. It is intentionally bounded to avoid retry storms.
type FailoverPolicy struct {
	MaxAttempts    int
	RetryDelay     time.Duration
	RequireHealthy bool
	RequireSecure  bool
}

func DefaultFailoverPolicy() FailoverPolicy {
	return FailoverPolicy{MaxAttempts: 2, RetryDelay: 25 * time.Millisecond, RequireHealthy: true, RequireSecure: true}
}

// SelectFailover returns the next ranked upstream after an unsuccessful
// attempt. The caller owns network I/O and circuit-breaker state.
func SelectFailover(policy FailoverPolicy, ranked []Upstream, failedID string, attempt int) (Upstream, bool) {
	max := policy.MaxAttempts
	if max <= 0 { max = 2 }
	if attempt >= max { return Upstream{}, false }
	for _, u := range ranked {
		if u.ID == failedID { continue }
		if policy.RequireHealthy && !u.Healthy { continue }
		if policy.RequireSecure && !u.Secure { continue }
		return u, true
	}
	return Upstream{}, false
}
