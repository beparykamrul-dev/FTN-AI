-- FTN migration registry. This migration is itself idempotent and is used by
-- the production migration runner for existing PostgreSQL volumes.
CREATE TABLE IF NOT EXISTS schema_migrations (
  version INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
INSERT INTO schema_migrations(version,name) VALUES (11,'011_migration_registry.sql') ON CONFLICT (version) DO NOTHING;
