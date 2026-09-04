package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type IdentityRecord struct {
	ID           string
	Username     string
	Email        string
	PasswordHash []byte
	Status       string
}

type Store struct{ DB *sql.DB }

func (s *Store) FindIdentity(ctx context.Context, login string) (IdentityRecord, error) {
	var r IdentityRecord
	login = strings.TrimSpace(login)
	if s == nil || s.DB == nil || login == "" {
		return r, errors.New("invalid identity lookup")
	}
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, status
		FROM ftn_identities
		WHERE username = $1 OR email = $1
		LIMIT 1`, login).Scan(&r.ID, &r.Username, &r.Email, &r.PasswordHash, &r.Status)
	return r, err
}

func (s *Store) HasActiveService(ctx context.Context, identityID, serviceID string) (bool, error) {
	if s == nil || s.DB == nil || strings.TrimSpace(identityID) == "" || strings.TrimSpace(serviceID) == "" {
		return false, errors.New("invalid service authorization lookup")
	}
	var exists bool
	err := s.DB.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM ftn_service_assignments
			WHERE identity_id = $1 AND service_id = $2 AND status = 'active'
		)`, identityID, serviceID).Scan(&exists)
	return exists, err
}
