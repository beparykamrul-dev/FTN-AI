package callcenter

type CustomerContext struct {
	CustomerID string `json:"customer_id"`
	Name string `json:"name,omitempty"`
	ServiceID string `json:"service_id,omitempty"`
	PackageName string `json:"package_name,omitempty"`
	PPPoEUser string `json:"pppoe_user,omitempty"`
	ONUID string `json:"onu_id,omitempty"`
	RouterID string `json:"router_id,omitempty"`
	IP string `json:"ip,omitempty"`
	ServiceStatus string `json:"service_status,omitempty"`
	ActiveIncident string `json:"active_incident,omitempty"`
}

type ContextProvider interface {
	Load(customerID string) (CustomerContext, error)
}

// BuildSessionContext loads customer context through an injected provider so
// billing, CRM, network and GIS databases remain replaceable integrations.
func BuildSessionContext(provider ContextProvider, customerID string) (CustomerContext, error) {
	if provider == nil { return CustomerContext{CustomerID: customerID}, nil }
	return provider.Load(customerID)
}
