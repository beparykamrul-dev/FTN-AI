package ftnstorage

// CompressionPolicy controls when FTN should spend CPU to reduce stored bytes.
type CompressionPolicy struct {
	Enabled       bool
	Algorithm     string
	MinSizeBytes  uint64
	MinSavingsPct uint8
}

func (p CompressionPolicy) Valid() bool {
	if !p.Enabled { return true }
	return p.Algorithm != "" && p.MinSavingsPct <= 100
}
