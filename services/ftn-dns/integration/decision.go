package integration

// DecisionPolicy controls how FTN DNS converts route observations into a
// deterministic provider/node decision.
type DecisionPolicy struct {
	RequireHealthy bool
	RequireSecure  bool
}

// Decide returns the first candidate satisfying the policy. Candidates must
// already be ranked by RankRoutes; no network mutation occurs here.
func Decide(policy DecisionPolicy, candidates []RouteCandidate) (RouteCandidate, bool) {
	for _, c := range candidates {
		if policy.RequireHealthy && !c.Healthy { continue }
		if policy.RequireSecure && !c.Secure { continue }
		return c, true
	}
	return RouteCandidate{}, false
}
