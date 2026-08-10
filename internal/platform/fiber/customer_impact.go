package fiber

// CustomerImpact describes service exposure for a topology change. It is a
// planning result only; it does not disable or modify customer service.
type CustomerImpact struct {
	CustomerID string `json:"customer_id"`
	ServiceID string `json:"service_id,omitempty"`
	ONU string `json:"onu,omitempty"`
	PPPoEUser string `json:"pppoe_user,omitempty"`
	AffectedLinkID string `json:"affected_link_id"`
	Severity string `json:"severity"`
}

func AnalyzeCustomerImpact(change TopologyChange, customers []CustomerProfile) []CustomerImpact {
	if change.Entity.ExternalID == "" { return nil }
	out := make([]CustomerImpact, 0)
	for _, c := range customers {
		if c.ID == "" { continue }
		if c.ONU == change.Entity.ExternalID || c.Router == change.Entity.ExternalID || c.ServiceID == change.Entity.ExternalID {
			severity := "potential"
			if change.Entity.Status == "down" || change.Entity.Status == "cut" { severity = "high" }
			out = append(out, CustomerImpact{CustomerID:c.ID, ServiceID:c.ServiceID, ONU:c.ONU, PPPoEUser:c.PPPoEUser, AffectedLinkID:change.Entity.ExternalID, Severity:severity})
		}
	}
	return out
}
