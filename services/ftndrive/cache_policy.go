package ftndrive

// CachePolicy controls where a user's frequently accessed object may be cached.
type CachePolicy struct {
	Enabled       bool
	MaxBytes      uint64
	LocalPreferred bool
	AllowPeerCache bool
}

func (p CachePolicy) Valid() bool {
	return p.Enabled && p.MaxBytes > 0
}
