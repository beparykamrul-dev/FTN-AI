package ftnwan

// ServiceIdentity describes an authenticated FTNWAN workload/service.
type ServiceIdentity struct {
	ID       string
	Service  string
	NodeID   string
	TenantID string
	Enabled  bool
}

// AccessDecision is the normalized result of a policy evaluation.
type AccessDecision string

const (
	AccessAllow AccessDecision = "allow"
	AccessDeny  AccessDecision = "deny"
)

// EvaluateServiceAccess performs the deterministic baseline checks.
// Cryptographic authentication and policy evaluation are performed by the
// surrounding identity/policy adapters; this helper only handles state gates.
func EvaluateServiceAccess(identity ServiceIdentity) AccessDecision {
	if identity.ID == "" || identity.Service == "" || identity.NodeID == "" || !identity.Enabled {
		return AccessDeny
	}
	return AccessAllow
}
