package observability

// ReplicationPolicy controls how many independent FTN storage targets should hold critical telemetry.
type ReplicationPolicy struct {
	CriticalReplicas uint32
	StandardReplicas uint32
}

func (p ReplicationPolicy) Valid() bool {
	return p.CriticalReplicas > 0 && p.StandardReplicas > 0 && p.CriticalReplicas >= p.StandardReplicas
}
