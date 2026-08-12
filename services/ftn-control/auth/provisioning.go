package auth

import (
	"context"
	"errors"
)

type ProvisioningStore interface {
	CreateIdentity(ctx context.Context, username, email string, passwordHash []byte) (string, error)
	AssignService(ctx context.Context, identityID, serviceID, provisionedBy string) error
}

var ErrInvalidProvisioning = errors.New("invalid provisioning request")

// ProvisionIdentity is intentionally called by the FTN Control Panel only.
func ProvisionIdentity(ctx context.Context, store ProvisioningStore, username, email string, passwordHash []byte) (string, error) {
	if username == "" || email == "" || len(passwordHash) == 0 {
		return "", ErrInvalidProvisioning
	}
	return store.CreateIdentity(ctx, username, email, passwordHash)
}

func AssignService(ctx context.Context, store ProvisioningStore, identityID, serviceID, provisionedBy string) error {
	if identityID == "" || serviceID == "" || provisionedBy == "" {
		return ErrInvalidProvisioning
	}
	return store.AssignService(ctx, identityID, serviceID, provisionedBy)
}
