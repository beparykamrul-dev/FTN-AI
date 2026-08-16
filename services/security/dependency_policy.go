package security

// DependencyPolicy defines minimum dependency metadata required by FTN CI.
type DependencyPolicy struct {
	RequireVersion bool
	RequirePURL    bool
	BlockUnknown   bool
}

func (p DependencyPolicy) Accepts(c SBOMComponent) bool {
	if p.RequireVersion && c.Version == "" { return false }
	if p.RequirePURL && c.PURL == "" { return false }
	if p.BlockUnknown && c.Source == "" { return false }
	return c.Valid()
}
