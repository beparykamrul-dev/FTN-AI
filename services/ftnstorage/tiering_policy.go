package ftnstorage

// TierPolicy describes movement between hot, warm and cold storage tiers.
type TierPolicy struct {
	HotMaxAgeHours  uint32
	WarmMaxAgeHours uint32
	ColdAfterHours  uint32
}

func (p TierPolicy) Valid() bool {
	return p.HotMaxAgeHours <= p.WarmMaxAgeHours && p.WarmMaxAgeHours <= p.ColdAfterHours
}
