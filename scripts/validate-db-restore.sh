#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

POSTGRES_IMAGE="postgres:17-alpine"
NAME="ftn-db-restore-${GITHUB_RUN_ID:-local}-$$"
DUMP="$(mktemp)"
RESTORE_DUMP="$(mktemp)"
cleanup() {
  docker rm -f "$NAME" >/dev/null 2>&1 || true
  rm -f "$DUMP" "$RESTORE_DUMP"
}
trap cleanup EXIT

for required in \
  services/control-plane/schema.sql \
  services/control-plane/migrations/002_platform_foundation.sql \
  services/control-plane/migrations/003_durable_execution.sql \
  services/control-plane/migrations/004_approval_execution.sql \
  services/control-plane/migrations/005_job_integrity_leases.sql \
  services/control-plane/migrations/006_job_event_automation.sql; do
  test -f "$required"
done

echo "Starting isolated PostgreSQL validation instance"
docker run -d --name "$NAME" \
  -e POSTGRES_DB=ftn \
  -e POSTGRES_USER=ftn \
  -e POSTGRES_PASSWORD=ci-only-password \
  "$POSTGRES_IMAGE" >/dev/null

ready_deadline=$((SECONDS + 90))
until docker exec "$NAME" pg_isready -U ftn -d ftn >/dev/null 2>&1; do
  if (( SECONDS >= ready_deadline )); then
    echo "PostgreSQL readiness timeout" >&2
    exit 1
  fi
  sleep 1
done

for migration in \
  services/control-plane/schema.sql \
  services/control-plane/migrations/002_platform_foundation.sql \
  services/control-plane/migrations/003_durable_execution.sql \
  services/control-plane/migrations/004_approval_execution.sql \
  services/control-plane/migrations/005_job_integrity_leases.sql \
  services/control-plane/migrations/006_job_event_automation.sql; do
  echo "Applying $migration"
  docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn < "$migration" >/dev/null
done

echo "Writing deterministic integrity fixture"
docker exec -i "$NAME" psql -v ON_ERROR_STOP=1 -U ftn -d ftn <<'SQL' >/dev/null
INSERT INTO tenants (slug, name) VALUES ('ci-restore', 'CI Restore Validation') ON CONFLICT (slug) DO NOTHING;
SQL

before="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM tenants WHERE slug='ci-restore';")"
test "$before" = "1"

echo "Creating custom-format backup"
docker exec "$NAME" pg_dump -U ftn -d ftn -Fc > "$DUMP"
test -s "$DUMP"

docker exec "$NAME" createdb -U ftn ftn_restore
docker exec -i "$NAME" pg_restore -U ftn -d ftn_restore --exit-on-error < "$DUMP"

after="$(docker exec "$NAME" psql -At -U ftn -d ftn_restore -c "SELECT count(*) FROM tenants WHERE slug='ci-restore';")"
test "$after" = "1"

event_trigger="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT count(*) FROM pg_trigger WHERE tgname='durable_jobs_event_journal';")"
test "$event_trigger" = "1"

schema_version="$(docker exec "$NAME" psql -At -U ftn -d ftn -c "SELECT to_regclass('public.event_journal') IS NOT NULL AND to_regclass('public.durable_jobs') IS NOT NULL;")"
test "$schema_version" = "t"

echo "PostgreSQL migration + backup + restore validation passed."
