package ftnservice

// Health captures the minimum control-plane health state for a service.
type Health struct {
	Healthy     bool
	DependenciesHealthy bool
	Ready       bool
	Reason      string
}

func (h Health) CanReceiveTraffic() bool {
	return h.Healthy && h.DependenciesHealthy && h.Ready
}
