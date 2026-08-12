package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrSessionRevoked = errors.New("session revoked or expired")

type SessionStore struct{ DB *sql.DB }

func (s *SessionStore) Save(ctx context.Context, session Session, ttl time.Duration) error {
	if session.ID == "" || session.IdentityID == "" || ttl <= 0 {
		return errors.New("invalid session")
	}
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO ftn_sessions (id, identity_id, expires_at)
		VALUES ($1, $2, $3)`, session.ID, session.IdentityID, time.Now().Add(ttl))
	return err
}

func (s *SessionStore) Validate(ctx context.Context, sessionID string) (string, error) {
	var identityID string
	var revokedAt sql.NullTime
	var expiresAt time.Time
	err := s.DB.QueryRowContext(ctx, `
		SELECT identity_id, expires_at, revoked_at
		FROM ftn_sessions WHERE id = $1`, sessionID).Scan(&identityID, &expiresAt, &revokedAt)
	if err != nil || revokedAt.Valid || !expiresAt.After(time.Now()) {
		return "", ErrSessionRevoked
	}
	return identityID, nil
}

func (s *SessionStore) Revoke(ctx context.Context, sessionID string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE ftn_sessions SET revoked_at = now()
		WHERE id = $1 AND revoked_at IS NULL`, sessionID)
	return err
}
