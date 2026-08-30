BEGIN;

CREATE TABLE IF NOT EXISTS ftn_services (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    policy_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ftn_capabilities (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    risk_class TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    policy_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ftn_transports (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    policy_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ftn_mesh_peers (
    id TEXT PRIMARY KEY,
    endpoint TEXT NOT NULL,
    state TEXT NOT NULL,
    certificate_ok BOOLEAN NOT NULL DEFAULT FALSE,
    policy_allowed BOOLEAN NOT NULL DEFAULT FALSE,
    last_seen TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ftn_path_candidates (
    id BIGSERIAL PRIMARY KEY,
    service_id TEXT NOT NULL REFERENCES ftn_services(id),
    peer_id TEXT NOT NULL REFERENCES ftn_mesh_peers(id),
    transport_id TEXT NOT NULL REFERENCES ftn_transports(id),
    healthy BOOLEAN NOT NULL DEFAULT FALSE,
    allowed BOOLEAN NOT NULL DEFAULT FALSE,
    score BIGINT NOT NULL DEFAULT 0,
    observed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ftn_policy_versions (
    version TEXT PRIMARY KEY,
    policy_hash TEXT NOT NULL,
    approved BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ftn_approvals (
    id BIGSERIAL PRIMARY KEY,
    policy_version TEXT NOT NULL REFERENCES ftn_policy_versions(version),
    operation TEXT NOT NULL,
    approved BOOLEAN NOT NULL DEFAULT FALSE,
    approved_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS ftn_route_decisions (
    id BIGSERIAL PRIMARY KEY,
    service_id TEXT NOT NULL REFERENCES ftn_services(id),
    candidate_id BIGINT REFERENCES ftn_path_candidates(id),
    decision TEXT NOT NULL,
    reason TEXT NOT NULL,
    policy_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ftn_audit_events (
    id BIGSERIAL PRIMARY KEY,
    event_type TEXT NOT NULL,
    service_id TEXT REFERENCES ftn_services(id),
    peer_id TEXT REFERENCES ftn_mesh_peers(id),
    transport_id TEXT REFERENCES ftn_transports(id),
    policy_version TEXT,
    decision TEXT,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ftn_path_candidates_service ON ftn_path_candidates(service_id, observed_at DESC);
CREATE INDEX IF NOT EXISTS idx_ftn_route_decisions_service ON ftn_route_decisions(service_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ftn_audit_events_created ON ftn_audit_events(created_at DESC);

COMMIT;
