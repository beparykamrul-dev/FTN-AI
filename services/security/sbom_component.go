package security

// SBOMComponent is a normalized software component record for FTN security analysis.
type SBOMComponent struct {
	Name      string
	Version   string
	PURL      string
	License   string
	Source    string
}

func (c SBOMComponent) Valid() bool {
	return c.Name != "" && c.Version != ""
}
