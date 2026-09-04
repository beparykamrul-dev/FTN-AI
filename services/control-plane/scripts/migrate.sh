#!/usr/bin/env bash
set -euo pipefail
: "${DATABASE_URL:?DATABASE_URL is required}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY,name TEXT NOT NULL,applied_at TIMESTAMPTZ NOT NULL DEFAULT now());
SQL
# Existing databases already containing the FTN Guard baseline must not replay
# historical DDL. Mark the known 001-011 baseline once, then only apply newer files.
if psql "$DATABASE_URL" -Atqc "SELECT to_regclass('public.dns_filter_profiles') IS NOT NULL" | grep -qx t; then
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 <<'SQL'
INSERT INTO schema_migrations(version,name)
SELECT v, 'legacy-baseline-' || lpad(v::text,3,'0') FROM generate_series(1,11) AS v
ON CONFLICT (version) DO NOTHING;
SQL
fi

shopt -s nullglob
files=("$ROOT"/migrations/*.sql)

declare -A seen_versions=()
for file in "${files[@]}"; do
  name="$(basename "$file")"
  version="${name%%_*}"
  [[ "$version" =~ ^[0-9]+$ ]] || continue
  if [[ -n "${seen_versions[$version]:-}" ]]; then
    echo "duplicate migration version ${version}: ${seen_versions[$version]} and ${name}" >&2
    exit 1
  fi
  seen_versions[$version]="$name"
done

mapfile -t files < <(printf '%s\n' "${files[@]}" | sort -V)
for file in "${files[@]}"; do
  name="$(basename "$file")"
  version="${name%%_*}"
  [[ "$version" =~ ^[0-9]+$ ]] || continue

  registered_name="$(psql "$DATABASE_URL" -Atqc "SELECT name FROM schema_migrations WHERE version=${version}" || true)"
  if [[ -n "$registered_name" ]]; then
    if (( 10#$version >= 12 )) && [[ "$registered_name" != "$name" ]]; then
      echo "migration ${version} is registered as ${registered_name}, expected ${name}" >&2
      exit 1
    fi
    continue
  fi

  echo "applying migration ${version}: ${name}"
  # Keep the advisory lock and the migration DDL on the same PostgreSQL
  # session. A second migration runner waits, rechecks the registry, and skips
  # the file once the first runner has committed it.
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 <<SQL
SELECT pg_advisory_lock(hashtextextended('ftn-control-plane-schema-migrations', 0));
SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version=${version}) AS already_applied \gset
\if :already_applied
\else
\i '$file'
INSERT INTO schema_migrations(version, name)
SELECT ${version}, '${name}'
WHERE NOT EXISTS (SELECT 1 FROM schema_migrations WHERE version=${version});
\endif
SELECT pg_advisory_unlock(hashtextextended('ftn-control-plane-schema-migrations', 0));
SQL
done
echo "FTN migrations complete"
