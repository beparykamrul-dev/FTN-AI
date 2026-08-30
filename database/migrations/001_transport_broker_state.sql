BEGIN;

CREATE TABLE IF NOT EXISTS ftn_services (
    service_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    policy_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ftn_capabilities (
    capability_id TEXT PRIMARY KEY,
    category TEXT NOT NULL,
    risk_class TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    policy_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ftn_transports (
    transport_id TEXT PRIMARY KEY,
    capability_id TEXT NOT NULL REFERENCES ftn_capabilities(capability_id),
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    mode TEXT NOT NULL DEFAULT 'MANUAL',
    policy_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ftn_mesh_peers (
    peer_id TEXT PRIMARY KEY,
    endpoint TEXT NOT NULL,
    state TEXT NOT NULL DEFAULT 'UNKNOWN',
    certificate_ok BOOLEAN NOT NULL DEFAULT FALSE,
    policy_allowed BOOLEAN NOT NULL DEFAULT FALSE,
    last_seen TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ftn_route_decisions (
    decision_id UUID PRIMARY KEY,
    service_id TEXT NOT NULL REFERENCES ftn_services(service_id),
    peer_id TEXT REFERENCES ftn_mesh_peers(peer_id),
    transport_id TEXT REFERENCES ftn_transports(transport_id),
    policy_version TEXT NOT NULL,
    decision TEXT NOT NULL,
    reason TEXT NOT NULL,
    score BIGINT,
    mutation_allowed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ftn_audit_events (
    event_id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    service_id TEXT,
    peer_id TEXT,
    capability_id TEXT,
    transport_id TEXT,
    policy_version TEXT,
    decision TEXT NOT NULL,
    reason TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ftn_route_decisions_service_created
    ON ftn_route_decisions(service_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ftn_audit_events_created
    ON ftn_audit_events(created_at DESC);

COMMIT;
