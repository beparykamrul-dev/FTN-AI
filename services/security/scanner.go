package security

// Scanner describes a pluggable static-analysis/security scanner.
type Scanner struct {
	Name    string
	Enabled bool
	Ruleset string
}

func (s Scanner) Valid() bool { return s.Name != "" && s.Enabled && s.Ruleset != "" }
