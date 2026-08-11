package controlplane

// APIContract describes the stable control-plane boundary used by web, mobile,
// agents and integrations. Authentication, RBAC and approval are enforced by
// the transport/service layer before mutating operations are executed.
type APIContract struct {
	Version string `json:"version"`
	Resources []string `json:"resources"`
	Events []string `json:"events"`
}

func DefaultAPIContract() APIContract {
	return APIContract{
		Version: "v1",
		Resources: []string{"devices", "servers", "networks", "dns", "mesh", "fiber", "customers", "services", "deployments", "changes", "audit"},
		Events: []string{"device.state", "mesh.topology", "deployment.state", "change.state", "alert"},
	}
}
