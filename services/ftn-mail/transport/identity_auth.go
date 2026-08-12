package transport

import (
	"context"
	"errors"
)

type IdentityCredentialStore interface {
	FindCredential(ctx context.Context, login string) (Credential, error)
}

type Credential struct {
	IdentityID    string
	PasswordHash  []byte
	PasswordSalt  []byte
	Status        string
}

type IdentityAuthenticator struct {
	Store IdentityCredentialStore
}

func (a IdentityAuthenticator) Authenticate(ctx context.Context, username, password string) (string, error) {
	if a.Store == nil || username == "" || password == "" {
		return "", errors.New("authentication failed")
	}
	c, err := a.Store.FindCredential(ctx, username)
	if err != nil || c.Status != "active" || !verifyArgon2id(password, c.PasswordHash, c.PasswordSalt) {
		return "", errors.New("authentication failed")
	}
	return c.IdentityID, nil
}
