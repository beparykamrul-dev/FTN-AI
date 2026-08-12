package auth

import "errors"

var ErrServiceAccessDenied = errors.New("service access denied")

type ServiceAssignment struct {
	IdentityID string
	ServiceID  string
	Active     bool
}

func Authorize(a ServiceAssignment, identityID, serviceID string) error {
	if !a.Active || a.IdentityID != identityID || a.ServiceID != serviceID {
		return ErrServiceAccessDenied
	}
	return nil
}
