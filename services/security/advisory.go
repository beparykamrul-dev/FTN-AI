package security

// Advisory is a normalized dependency vulnerability/advisory record.
type Advisory struct {
	ID       string
	Package  string
	Severity string
	FixedIn  string
}

func (a Advisory) Valid() bool {
	return a.ID != "" && a.Package != "" && a.Severity != ""
}
