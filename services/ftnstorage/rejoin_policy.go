package ftnstorage

// RejoinPolicy controls when a repaired/quarantined storage node may re-enter service.
type RejoinPolicy struct {
	RequireHealthy bool
	RequireVerified bool
	RequireCleanScrub bool
}

func (p RejoinPolicy) Valid() bool {
	return p.RequireHealthy && p.RequireVerified && p.RequireCleanScrub
}

func (p RejoinPolicy) Allow(state HealthControllerState, verified, scrubClean bool) bool {
	return p.Valid() && state == HealthHealthy && verified && scrubClean
}
