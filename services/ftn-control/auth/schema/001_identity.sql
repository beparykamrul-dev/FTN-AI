CREATE TABLE IF NOT EXISTS ftn_identities (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','disabled','locked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ftn_service_assignments (
    id UUID PRIMARY KEY,
    identity_id UUID NOT NULL REFERENCES ftn_identities(id) ON DELETE CASCADE,
    service_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','revoked','suspended')),
    provisioned_by UUID REFERENCES ftn_identities(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at TIMESTAMPTZ,
    UNIQUE(identity_id, service_id)
);

CREATE INDEX IF NOT EXISTS idx_ftn_service_assignments_identity
    ON ftn_service_assignments(identity_id, status);

CREATE TABLE IF NOT EXISTS ftn_sessions (
    id TEXT PRIMARY KEY,
    identity_id UUID NOT NULL REFERENCES ftn_identities(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ftn_sessions_identity
    ON ftn_sessions(identity_id);
