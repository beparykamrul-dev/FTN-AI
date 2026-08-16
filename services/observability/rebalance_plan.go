package observability

// RebalancePlan describes a controlled movement of replicas from overloaded or unhealthy nodes.
type RebalancePlan struct {
	SourceNode string
	TargetNode string
	Reason string
	Priority uint32
}

func (p RebalancePlan) Valid() bool {
	return p.SourceNode != "" && p.TargetNode != "" && p.SourceNode != p.TargetNode && p.Reason != ""
}
