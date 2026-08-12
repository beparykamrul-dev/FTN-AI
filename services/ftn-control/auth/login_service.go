package auth

import (
	"context"
	"errors"
	"time"
)

var ErrInvalidLogin = errors.New("invalid credentials")

type CredentialRecord struct {
	ID           string
	PasswordHash []byte
	PasswordSalt []byte
	Status       string
}

type CredentialStore interface {
	FindCredential(ctx context.Context, login string) (CredentialRecord, error)
}

type LoginService struct {
	Credentials CredentialStore
	Sessions     *SessionStore
	SessionTTL   time.Duration
}

func (s *LoginService) Login(ctx context.Context, login, password string) (Session, error) {
	if login == "" || password == "" || s.Credentials == nil || s.Sessions == nil {
		return Session{}, ErrInvalidLogin
	}
	record, err := s.Credentials.FindCredential(ctx, login)
	if err != nil || record.Status != "active" || !VerifyCredential(password, record.PasswordHash, record.PasswordSalt) {
		return Session{}, ErrInvalidLogin
	}
	ttl := s.SessionTTL
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	session, err := NewSession(record.ID)
	if err != nil {
		return Session{}, err
	}
	if err := s.Sessions.Save(ctx, session, ttl); err != nil {
		return Session{}, err
	}
	return session, nil
}
