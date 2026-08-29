package realtime

import (
	"errors"
	"strings"

	"github.com/beparykamrul-dev/FTN-AI/services/isp-api/internal/rbac"
)

var ErrUnauthorized = errors.New("unauthorized realtime connection")

type Principal struct {
	UserID string
	Role   rbac.Role
	Scope  string
}

func AuthorizeConnection(p Principal) error {
	if strings.TrimSpace(p.UserID) == "" || strings.TrimSpace(string(p.Role)) == "" {
		return ErrUnauthorized
	}
	return nil
}

func AuthorizeScopedSubscription(p Principal, channel string) error {
	if err := AuthorizeConnection(p); err != nil {
		return err
	}
	if !AuthorizeSubscription(p.Role, channel) {
		return ErrUnauthorized
	}
	if strings.TrimSpace(p.Scope) == "" {
		return ErrUnauthorized
	}
	return nil
}
