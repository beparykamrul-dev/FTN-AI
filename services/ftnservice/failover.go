package ftnservice

// FailoverDecision describes a safe service endpoint switch.
type FailoverDecision struct {
	ServiceID  string
	FromNode   string
	ToNode     string
	Reason     string
	Approved   bool
}

func BuildFailover(serviceID, fromNode string, candidates []ServiceEndpoint) (FailoverDecision, bool) {
	var best ServiceEndpoint
	found := false
	for _, c := range candidates {
		if !c.Valid() || !c.Healthy || c.ServiceID != serviceID || c.NodeID == fromNode {
			continue
		}
		if !found || c.Weight > best.Weight || (c.Weight == best.Weight && c.NodeID < best.NodeID) {
			best, found = c, true
		}
	}
	if !found {
		return FailoverDecision{}, false
	}
	return FailoverDecision{ServiceID: serviceID, FromNode: fromNode, ToNode: best.NodeID, Reason: "unified-service-failover"}, true
}
