package data

// Region describes a logical database placement region.
type Region struct {
	Name     string
	Primary  string
	Replicas []string
}

// Topology describes a multi-region unified database deployment without
// coupling service code to a particular database engine.
type Topology struct {
	GlobalBackend string
	Regions      []Region
}

func (t Topology) Valid() bool {
	if t.GlobalBackend == "" {
		return false
	}
	for _, r := range t.Regions {
		if r.Name == "" || r.Primary == "" {
			return false
		}
	}
	return true
}
