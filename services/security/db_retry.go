package security

// DBRetryPolicy captures bounded retry behavior for transient distributed-DB failures.
type DBRetryPolicy struct {
	MaxAttempts uint32
	BaseDelayMS uint32
}

func (p DBRetryPolicy) Valid() bool { return p.MaxAttempts > 0 && p.BaseDelayMS > 0 }
