package security

// SecurityBundle groups normalized analysis outputs without retaining raw source.
type SecurityBundle struct {
	Findings   []Finding
	Components []SBOMComponent
}

func (b SecurityBundle) Valid() bool {
	for _, f := range b.Findings { if !f.Valid() { return false } }
	for _, c := range b.Components { if !c.Valid() { return false } }
	return true
}
