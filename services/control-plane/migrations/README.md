# FTN Control Plane migrations

Migrations in this directory are **additive and idempotent**. Existing FTN tables and data must not be removed or rewritten by a foundation migration.

## Rules

1. Never drop an existing table/column as part of a compatibility migration.
2. Use `IF NOT EXISTS` / `ON CONFLICT` where PostgreSQL permits it.
3. Add new constraints only when existing data is compatible; otherwise introduce a staged backfill migration first.
4. Every service keeps its API contract backward-compatible unless a versioned endpoint is introduced.
5. Service registration is centralized in `service_registry`; service-specific state remains owned by the service.
6. Privileged changes require an approval record and an audit event before execution.
7. Credentials are represented by hashes/metadata; plaintext production secrets never belong in Git.

The initial bootstrap currently executes `schema.sql` followed by `002_platform_foundation.sql` for fresh PostgreSQL volumes. Existing databases should receive migrations through the repository's migration runner before the next release.
