package auth

import (
	"context"
	"crypto/subtle"
	"errors"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginResult struct {
	IdentityID string
	Session    Session
}

type Service struct{ Store *Store }

// Authenticate verifies the identity and creates a server-side session.
// Password verification must use the same Argon2id parameters used at provisioning.
func (s *Service) Authenticate(ctx context.Context, login string, password, salt []byte) (LoginResult, error) {
	r, err := s.Store.FindIdentity(ctx, login)
	if err != nil || r.Status != "active" {
		return LoginResult{}, ErrInvalidCredentials
	}
	candidate := HashPassword(password, salt)
	if len(candidate) != len(r.PasswordHash) || subtle.ConstantTimeCompare(candidate, r.PasswordHash) != 1 {
		return LoginResult{}, ErrInvalidCredentials
	}
	session, err := NewSession(r.ID)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{IdentityID: r.ID, Session: session}, nil
}
