package edge

// Provider identifies an approved external edge/CDN/DNS ingress source.
type Provider struct {
	Name      string
	Kind      string // cdn, edge, dns, proxy
	Enabled   bool
	Verified  bool
	Endpoints []string
}

func (p Provider) Valid() bool {
	return p.Name != "" && p.Enabled && p.Verified && len(p.Endpoints) > 0
}
