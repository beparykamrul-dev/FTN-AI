package ftncompute

// ServerIdentity gives every FTN main server one stable control-plane identity.
// Multiple physical servers can therefore appear as one logical FTN fabric.
type ServerIdentity struct {
	ID        string
	ClusterID string
	Region    string
	Role      string
	Epoch     uint64
}

func (i ServerIdentity) Valid() bool {
	return i.ID != "" && i.ClusterID != "" && i.Region != "" && i.Role != "" && i.Epoch > 0
}
