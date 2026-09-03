#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATIONS_DIR="$ROOT/migrations"
START=15
END=23

fail() { echo "migration-validation: ERROR: $*" >&2; exit 1; }
command -v psql >/dev/null || fail "psql is required"

mapfile -t all < <(find "$MIGRATIONS_DIR" -maxdepth 1 -type f -name '*.sql' -printf '%f\n' | sort -V)
((${#all[@]})) || fail "no SQL migrations found"

# The 015..023 production gate must have exactly one migration file per version.
# Legacy duplicate versions outside this range are replayed with migrate.sh's
# existing first-file-wins behavior during the shadow execution below.
for v in $(seq "$START" "$END"); do
  mapfile -t files < <(printf '%s\n' "${all[@]}" | awk -v p="$v" '$0 ~ "^" sprintf("%03d",p) "_" {print}')
  ((${#files[@]} == 1)) || fail "version $v requires exactly one migration file; found ${#files[@]}: ${files[*]-}"
done

echo "migration-validation: unique 015..023 chain OK"

# Replay the same effective migration sequence as migrate.sh. This validates real
# PostgreSQL execution order: missing tables/columns/functions/indexes or trigger
# dependencies fail at the migration that first violates the dependency boundary.
: "${MIGRATION_VALIDATION_DATABASE_URL:?MIGRATION_VALIDATION_DATABASE_URL is required}"
: "${MIGRATION_VALIDATION_ALLOW_CREATE:?set MIGRATION_VALIDATION_ALLOW_CREATE=1 for an isolated validation database}"
[[ "$MIGRATION_VALIDATION_ALLOW_CREATE" == 1 ]] || fail "refusing to create a validation database"

base_db="$(psql "$MIGRATION_VALIDATION_DATABASE_URL" -Atqc 'select current_database()')"
suffix="$(date +%s)-$$"
validation_db="${base_db}_migration_gate_${suffix}"

cleanup() {
  psql "$MIGRATION_VALIDATION_DATABASE_URL" -v ON_ERROR_STOP=1 \
    -c "DROP DATABASE IF EXISTS \"$validation_db\" WITH (FORCE);" >/dev/null || true
}
trap cleanup EXIT

psql "$MIGRATION_VALIDATION_DATABASE_URL" -v ON_ERROR_STOP=1 \
  -c "CREATE DATABASE \"$validation_db\";" >/dev/null

declare -a replay=()
last_version=""
for file in "${all[@]}"; do
  version="${file%%_*}"
  [[ "$version" =~ ^[0-9]+$ ]] || continue
  if [[ "$version" == "$last_version" ]]; then continue; fi
  replay+=("$file")
  last_version="$version"
done

for file in "${replay[@]}"; do
  echo "migration-validation: applying $file"
  psql "$MIGRATION_VALIDATION_DATABASE_URL" -d "$validation_db" -X \
    -v ON_ERROR_STOP=1 -f "$MIGRATIONS_DIR/$file" >/dev/null
done

# Explicit final-state checks for the durable execution integrity surface.
psql "$MIGRATION_VALIDATION_DATABASE_URL" -d "$validation_db" -X \
  -v ON_ERROR_STOP=1 <<'SQL' >/dev/null
DO $$
DECLARE
  required_table text;
  required_column record;
BEGIN
  FOREACH required_table IN ARRAY ARRAY['durable_jobs','execution_attempts','change_approvals'] LOOP
    IF to_regclass('public.' || required_table) IS NULL THEN
      RAISE EXCEPTION 'missing required table: %', required_table;
    END IF;
  END LOOP;

  FOR required_column IN
    SELECT * FROM (VALUES
      ('durable_jobs','approval_id'),
      ('durable_jobs','execution_payload_hash'),
      ('durable_jobs','verification_payload'),
      ('change_approvals','approval_payload'),
      ('change_approvals','approval_payload_hash'),
      ('execution_attempts','attempt_no'),
      ('execution_attempts','worker_id'),
      ('execution_attempts','status')
    ) AS x(table_name,column_name)
  LOOP
    IF NOT EXISTS (
      SELECT 1 FROM information_schema.columns
      WHERE table_schema='public'
        AND table_name=required_column.table_name
        AND column_name=required_column.column_name
    ) THEN
      RAISE EXCEPTION 'missing required column: %.%', required_column.table_name, required_column.column_name;
    END IF;
  END LOOP;
END $$;

DO $$
DECLARE dup record;
BEGIN
  FOR dup IN
    SELECT n.nspname AS schema_name, c.relname AS table_name, t.tgname, count(*) AS trigger_count
      FROM pg_trigger t
      JOIN pg_class c ON c.oid=t.tgrelid
      JOIN pg_namespace n ON n.oid=c.relnamespace
     WHERE NOT t.tgisinternal
     GROUP BY n.nspname,c.relname,t.tgname
    HAVING count(*) > 1
  LOOP
    RAISE EXCEPTION 'duplicate trigger: %.%.% count=%',
      dup.schema_name,dup.table_name,dup.tgname,dup.trigger_count;
  END LOOP;
END $$;
SQL

echo "migration-validation: PostgreSQL execution order, dependency, schema, and final trigger checks PASSED"
