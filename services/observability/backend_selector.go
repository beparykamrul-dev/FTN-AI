package observability

// Backend describes one specialized observability backend in the FTN stack.
type Backend struct {
	Name     string
	Signals  []string
	Healthy  bool
	Priority uint32
}

func (b Backend) Eligible() bool { return b.Name != "" && b.Healthy && len(b.Signals) > 0 }

// SelectBackend deterministically selects the healthiest eligible backend with
// the lowest priority value, keeping FTN vendor-neutral.
func SelectBackend(backends []Backend) (Backend, bool) {
	var best Backend
	found := false
	for _, b := range backends {
		if !b.Eligible() { continue }
		if !found || b.Priority < best.Priority || (b.Priority == best.Priority && b.Name < best.Name) {
			best, found = b, true
		}
	}
	return best, found
}
