package security

// Finding is a normalized security-analysis result. Raw source is not stored here.
type Finding struct {
	Scanner  string
	RuleID   string
	Severity string
	Path     string
	Line     uint32
	Message  string
}

func (f Finding) Valid() bool {
	return f.Scanner != "" && f.RuleID != "" && f.Severity != "" && f.Path != ""
}
