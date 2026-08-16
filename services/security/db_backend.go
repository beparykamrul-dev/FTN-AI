package security

// DBBackend describes an available persistent security-evidence backend.
type DBBackend struct {
	Name string
	Kind string
	Healthy bool
	Priority uint32
}

func (b DBBackend) Eligible() bool { return b.Name != "" && b.Kind != "" && b.Healthy }

func SelectDBBackend(backends []DBBackend) (DBBackend, bool) {
	var best DBBackend
	found := false
	for _, b := range backends {
		if !b.Eligible() { continue }
		if !found || b.Priority < best.Priority || (b.Priority == best.Priority && b.Name < best.Name) {
			best, found = b, true
		}
	}
	return best, found
}
