package ftnservice

// ServiceEndpoint is an advertised FTN service instance.
type ServiceEndpoint struct {
	ServiceID string
	NodeID    string
	Region    string
	Address   string
	Port      uint16
	Healthy   bool
	Weight    uint32
}

func (e ServiceEndpoint) Valid() bool {
	return e.ServiceID != "" && e.NodeID != "" && e.Address != "" && e.Port > 0
}
