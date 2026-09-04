#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

POSTGRES_IMAGE="postgres:17-alpine"
NAME="ftn-db-restore-${GITHUB_RUN_ID:-local}-$$"
DUMP="$(mktemp)"
cleanup() {
  docker rm -f "$NAME" >/dev/null 2>&1 || true
  rm -f "$DUMP"
}
trap cleanup EXIT

command -v docker >/dev/null 2>&1 || { echo 'Docker is required' >&2; exit 1; }
test -f services/control-plane/schema.sql
shopt -s nullglob
migrations=(services/control-plane/migrations/*.sql)
((${#migrations[@]})) || { echo 'No control-plane migrations found' >&2; exit 1; }

declare -A seen_version
ordered=()
while IFS= read -r migration; do
  name="$(basename "$migration")"
  version="${name%%_*}"
  [[ "$version" =~ ^[0-9]+$ ]] || continue
  numeric=$((10#$version))
  if ((numeric >= 15 && numeric <= 24)); then
    if [[ -n "${seen_version[$numeric]:-}" ]]; then
      echo "duplicate critical migration version: $version ($name)" >&2
      exit 1
    fi
    seen_version[$numeric]="$name"
  fi
  # Match migrate.sh semantics: migration versions are applied once; when
  # historical duplicate versions exist, the first file in sort -V order wins.
  if [[ -n "${seen_version[all_$numeric]:-}" ]]; then
    continue
  fi
  seen_version[all_$numeric]="$name"
  ordered+=("$migration")
done < <(printf '%s\n' "${migrations[@]}" | sort -V)

for version in {15..24}; do
  [[ -n "${seen_version[$version]:-}" ]] || {
    echo "missing critical migration version: $(printf '%03d' "$version")" >&2
    exit 1
  }
done

echo "Starting isolated PostgreSQL validation instance"
docker run -d --name "$NAME" \
  -e POSTGRES_DB=ftn \
  -e POSTGRES_USER=ftn \
  -e POSTGRES_PASSWORD=ci-only-password \
  "$POSTGRES_IMAGE" >/dev/null

until docker exec "$NAME" pg_isready -U ftn -d ftn >/dev/null 2>&1 && \
      docker exec "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn -c 'SELECT 1' >/dev/null 2>&1; do
  sleep 1
done

docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn < services/control-plane/schema.sql >/dev/null
for migration in "${ordered[@]}"; do
  echo "Applying $migration"
  docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn < "$migration" >/dev/null
done

echo "Writing deterministic restore fixture"
docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn <<'SQL' >/dev/null
INSERT INTO tenants (slug, name) VALUES ('ci-restore', 'CI Restore Validation') ON CONFLICT (slug) DO NOTHING;
INSERT INTO event_journal (tenant_id, event_type, sequence, correlation_id, causation_id, aggregate_id, payload)
SELECT t.id, 'ci.restore.fixture',
       COALESCE((SELECT MAX(e.sequence) + 1 FROM event_journal e WHERE e.tenant_id=t.id), 1),
       'ci-restore', '', t.id::text, '{"source":"validate-db-restore"}'::jsonb
FROM tenants t WHERE t.slug='ci-restore';
SQL

before="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM tenants WHERE slug='ci-restore';")"
test "$before" = "1"
events_before="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM event_journal WHERE event_type='ci.restore.fixture';")"
test "$events_before" = "1"

echo "Creating custom-format backup"
# Keep the binary dump inside the container and copy it out. This avoids
# docker exec/stdout transport edge cases for custom-format pg_dump output.
docker exec "$NAME" pg_dump -U ftn -d ftn -Fc -f /tmp/ftn-restore.dump
docker cp "$NAME:/tmp/ftn-restore.dump" "$DUMP"
test -s "$DUMP"
pg_restore --list "$DUMP" >/dev/null

docker exec "$NAME" createdb -U ftn ftn_restore
docker cp "$DUMP" "$NAME:/tmp/ftn-restore.dump"
docker exec "$NAME" pg_restore -U ftn -d ftn_restore --exit-on-error /tmp/ftn-restore.dump

after="$(docker exec "$NAME" psql -At -U ftn -d ftn_restore -c "SELECT count(*) FROM tenants WHERE slug='ci-restore';")"
test "$after" = "1"
events_after="$(docker exec "$NAME" psql -At -U ftn -d ftn_restore -c "SELECT count(*) FROM event_journal WHERE event_type='ci.restore.fixture';")"
test "$events_after" = "1"

for table in tenants event_journal durable_jobs execution_attempts change_approvals; do
  docker exec "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn_restore -c "SELECT 1 FROM $table LIMIT 1" >/dev/null
done

trigger_count="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM pg_trigger WHERE tgname='durable_jobs_event_journal' AND NOT tgisinternal;")"
test "$trigger_count" = "1"

schema_ok="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT to_regclass('public.event_journal') IS NOT NULL AND to_regclass('public.durable_jobs') IS NOT NULL AND to_regclass('public.execution_attempts') IS NOT NULL AND to_regclass('public.change_approvals') IS NOT NULL;")"
test "$schema_ok" = "t"

echo "PostgreSQL current migration-chain + backup + restore validation passed."
