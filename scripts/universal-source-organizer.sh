#!/usr/bin/env bash
set -Eeuo pipefail

# FTN Universal Source Organizer
# Accepts common Git repository URLs and local paths.
# Default: scan + generate plan only.
# --apply: organize only the temporary fetched workspace; never writes to remote.
# Supported examples:
#   https://github.com/owner/repo.git
#   https://gitlab.com/owner/repo.git
#   https://bitbucket.org/owner/repo.git
#   https://gitea.example.com/owner/repo.git
#   /srv/source/project

INPUT="${1:-}"
MODE="${2:-plan}"

[[ -n "$INPUT" ]] || { echo "Usage: $0 <repo-url-or-local-path> [--apply]" >&2; exit 2; }
[[ "$MODE" == "plan" || "$MODE" == "--apply" ]] || { echo "Mode must be omitted or --apply" >&2; exit 2; }
command -v find >/dev/null || { echo "find is required" >&2; exit 127; }

TMP="$(mktemp -d "${TMPDIR:-/tmp}/ftn-universal-organizer.XXXXXX")"
trap 'rm -rf "$TMP"' EXIT

if [[ -d "$INPUT" ]]; then
  ROOT="$(cd "$INPUT" && pwd)"
elif [[ "$INPUT" =~ ^(https?|ssh|git):// ]] || [[ "$INPUT" =~ ^git@[^:]+:.+ ]]; then
  command -v git >/dev/null || { echo "git is required for repository URLs" >&2; exit 127; }
  ROOT="$TMP/repository"
  echo "[FTN] Fetching repository: $INPUT"
  git clone --depth 1 --quiet -- "$INPUT" "$ROOT"
else
  echo "Unsupported input. Provide a local directory or Git URL." >&2
  exit 2
fi

PLAN="$TMP/source-plan.tsv"
REPORT="$TMP/SOURCE-ORGANIZER-REPORT.md"
: > "$PLAN"

classify() {
  local f="$1" b e
  b="$(basename "$f")"
  e="${b##*.}"
  case "$b" in
    Dockerfile*|compose*.yml|compose*.yaml|docker-compose*.yml|docker-compose*.yaml) echo infrastructure; return;;
    go.mod|go.sum) echo backend; return;;
    package.json|package-lock.json|pnpm-lock.yaml|yarn.lock|bun.lockb) echo frontend; return;;
    *.test.go|*_test.go|*.test.ts|*.spec.ts|*.test.tsx|*.spec.tsx) echo tests; return;;
    *.sql) echo database; return;;
    *.md|*.mdx|*.rst|*.adoc) echo docs; return;;
    *.sh|Makefile) echo scripts; return;;
    *.env|*.env.*|*.toml|*.ini|*.conf|*.json|*.yaml|*.yml) echo configs; return;;
  esac
  case "$e" in
    go) echo backend;;
    ts|tsx|js|jsx|css|scss|html|vue|svelte) echo frontend;;
    py|rs|java|kt|c|cc|cpp|h|hpp|php|rb|dart) echo services;;
    *) echo archive;;
  esac
}

while IFS= read -r -d '' f; do
  rel="${f#$ROOT/}"
  case "$rel" in
    .git/*|node_modules/*|vendor/*|dist/*|build/*|.next/*|coverage/*|target/*) continue;;
    backend/*|frontend/*|services/*|modules/*|agents/*|database/*|infrastructure/*|scripts/*|configs/*|tests/*|docs/*|archive/*) continue;;
  esac
  category="$(classify "$f")"
  printf '%s\t%s/%s\n' "$rel" "$category" "$rel" >> "$PLAN"
done < <(find "$ROOT" -type f -print0)

count="$(wc -l < "$PLAN" | tr -d ' ' )"
{
  echo '# FTN Universal Source Organizer Report'
  echo
  echo "Input: $INPUT"
  echo "Mode: $MODE"
  echo "Files proposed: $count"
  echo
  echo '| Source | Proposed target |'
  echo '|---|---|'
  while IFS=$'\t' read -r src dst; do
    printf '| `%s` | `%s` |\n' "$src" "$dst"
  done < "$PLAN"
} > "$REPORT"

cat "$REPORT"
echo
echo "[FTN] Report: $REPORT"

if [[ "$MODE" == "plan" ]]; then
  echo '[FTN] PLAN ONLY: no files changed.'
  exit 0
fi

while IFS=$'\t' read -r src dst; do
  [[ -z "$src" ]] && continue
  from="$ROOT/$src"
  to="$ROOT/$dst"
  [[ -e "$to" ]] && { echo "[FTN] SKIP existing: $dst"; continue; }
  mkdir -p "$(dirname "$to")"
  mv -- "$from" "$to"
done < "$PLAN"

echo '[FTN] Applied to local/temporary workspace only. Remote repository was not modified.'
