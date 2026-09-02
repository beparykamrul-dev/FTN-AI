#!/usr/bin/env bash
set -euo pipefail

: "${DATABASE_URL:?DATABASE_URL is required}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Registry must exist before incremental migrations can be tracked.
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS schema_migrations (
  version INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
SQL

shopt -s nullglob
for file in "$ROOT"/migrations/*.sql; do
  name="$(basename "$file")"
  version="${name%%_*}"
  [[ "$version" =~ ^[0-9]+$ ]] || continue
  if psql "$DATABASE_URL" -Atqc "SELECT 1 FROM schema_migrations WHERE version=${version}" | grep -qx 1; then
    continue
  fi
  echo "applying migration ${version}: ${name}"
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f "$file"
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -c "INSERT INTO schema_migrations(version,name) VALUES (${version}, '${name}') ON CONFLICT (version) DO NOTHING;"
done

echo "FTN migrations complete"
