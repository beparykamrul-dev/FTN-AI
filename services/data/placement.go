package data

// Placement selects where service data should live without coupling the
// service to one database engine.
type Placement struct {
	Scope       string // global, regional, local
	Backend     string
	ReadOnly    bool
	Replicated  bool
	Consistency string // strong, eventual
}

func ValidPlacement(p Placement) bool {
	if p.Scope != "global" && p.Scope != "regional" && p.Scope != "local" {
		return false
	}
	if p.Backend == "" {
		return false
	}
	return p.Consistency == "strong" || p.Consistency == "eventual"
}
