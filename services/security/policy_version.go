package security

// PolicyVersion identifies the immutable version of a security policy used by a gate decision.
type PolicyVersion struct {
	PolicyID string
	Version  uint64
}

func (p PolicyVersion) Valid() bool { return p.PolicyID != "" && p.Version > 0 }
