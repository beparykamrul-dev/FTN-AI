package realtime

import (
	"strings"

	"github.com/beparykamrul-dev/FTN-AI/services/isp-api/internal/rbac"
)

func AuthorizeSubscription(role rbac.Role, channel string) bool {
	channel = strings.TrimSpace(channel)
	if !IsAllowedChannel(channel) {
		return false
	}
	switch channel {
	case "account":
		return rbac.Allows(role, rbac.ReadOwnAccount)
	case "billing":
		return rbac.Allows(role, rbac.ReadBilling)
	case "notifications", "tickets":
		return rbac.Allows(role, rbac.CreateTicket)
	case "topology", "device_telemetry", "discovery", "incidents":
		return rbac.Allows(role, rbac.ReadNetwork)
	case "recovery":
		return rbac.Allows(role, rbac.ProposeRecovery) || rbac.Allows(role, rbac.ApproveRecovery)
	default:
		return false
	}
}
