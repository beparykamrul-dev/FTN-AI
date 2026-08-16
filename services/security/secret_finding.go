package security

// SecretFinding is a normalized secret-detection result. Secret material itself is never stored.
type SecretFinding struct {
	Scanner   string
	RuleID    string
	Path      string
	Line      uint32
	Kind      string
	Redacted  bool
}

func (f SecretFinding) Valid() bool {
	return f.Scanner != "" && f.RuleID != "" && f.Path != "" && f.Line > 0 && f.Kind != "" && f.Redacted
}
