CREATE TABLE IF NOT EXISTS ftn_security_evidence (
    id TEXT PRIMARY KEY,
    kind TEXT NOT NULL,
    run_id TEXT NOT NULL,
    release_id TEXT,
    fingerprint TEXT,
    payload_ref TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ftn_security_evidence_run
    ON ftn_security_evidence (run_id);

CREATE INDEX IF NOT EXISTS idx_ftn_security_evidence_release
    ON ftn_security_evidence (release_id);

CREATE INDEX IF NOT EXISTS idx_ftn_security_evidence_kind_created
    ON ftn_security_evidence (kind, created_at DESC);
