package data

// ReplicationPolicy describes a declarative replication requirement.
type ReplicationPolicy struct {
	Backend string
	Scope   string
	Copies  uint16
	Mode    ConsistencyMode
}

func (p ReplicationPolicy) Valid() bool {
	if p.Backend == "" || p.Copies == 0 {
		return false
	}
	if p.Scope != "global" && p.Scope != "regional" && p.Scope != "local" {
		return false
	}
	return ValidConsistency(p.Mode)
}
