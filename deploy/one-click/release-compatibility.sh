#!/usr/bin/env bash
set -Eeuo pipefail

RELEASE_ROOT="${1:-$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)}"
: "${DATABASE_URL:?DATABASE_URL is required}"

fail(){ printf '[FTN][COMPAT][ERROR] %s\n' "$*" >&2; exit 1; }

command -v psql >/dev/null 2>&1 || fail 'psql is required'
[ -d "$RELEASE_ROOT/services/control-plane/migrations" ] || fail 'migration directory is missing'

max_file_version="$(find "$RELEASE_ROOT/services/control-plane/migrations" -maxdepth 1 -type f -name '*.sql' -printf '%f\n' \
  | awk -F_ '$1 ~ /^[0-9]+$/ {print $1+0}' | sort -n | tail -1)"
[ -n "$max_file_version" ] || fail 'no numbered migrations found'

schema_table="$(psql "$DATABASE_URL" -Atqc "SELECT to_regclass('public.schema_migrations')")"
[ "$schema_table" = 'schema_migrations' ] || fail 'schema_migrations table is missing'

applied_version="$(psql "$DATABASE_URL" -Atqc "SELECT COALESCE(MAX(version),0) FROM schema_migrations")"
[[ "$applied_version" =~ ^[0-9]+$ ]] || fail 'invalid applied migration version'

printf '[FTN][COMPAT] applied_schema=%s release_max_schema=%s\n' "$applied_version" "$max_file_version"

# A release may run against an older schema and migrate forward. A release rollback
# must never run against a schema newer than the release knows how to handle.
if (( applied_version > max_file_version )); then
  fail "schema ${applied_version} is newer than release maximum ${max_file_version}; rollback is unsafe"
fi

printf '[FTN][COMPAT] PASS\n'
