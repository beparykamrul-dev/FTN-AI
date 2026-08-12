package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type Session struct {
	ID         string
	IdentityID string
	Revoked    bool
}

func NewSession(identityID string) (Session, error) {
	if identityID == "" {
		return Session{}, errors.New("identity id is required")
	}
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return Session{}, err
	}
	return Session{ID: hex.EncodeToString(buf), IdentityID: identityID}, nil
}

func Revoke(s *Session) {
	if s != nil {
		s.Revoked = true
	}
}
