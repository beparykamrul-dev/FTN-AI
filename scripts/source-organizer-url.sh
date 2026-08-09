#!/usr/bin/env bash
set -Eeuo pipefail

# FTN Source Organizer URL Node
# URL -> temporary clone -> deterministic source classification -> plan
# Default behavior is PLAN ONLY. --apply changes only the temporary workspace.
# No remote repository is modified by this script.
#
# Usage:
#   ./scripts/source-organizer-url.sh <git-url>
#   ./scripts/source-organizer-url.sh <git-url> --apply

URL="${1:-}"
MODE="${2:-plan}"

if [[ -z "$URL" ]]; then
  echo "Usage: $0 <git-url> [--apply]" >&2
  exit 2
fi

if [[ "$MODE" != "plan" && "$MODE" != "--apply" ]]; then
  echo "Mode must be omitted or --apply" >&2
  exit 2
fi

command -v git >/dev/null 2>&1 || { echo "git is required" >&2; exit 127; }
command -v find >/dev/null 2>&1 || { echo "find is required" >&2; exit 127; }

WORKDIR="$(mktemp -d "${TMPDIR:-/tmp}/ftn-source-organizer.XXXXXX")"
trap 'rm -rf "$WORKDIR"' EXIT

REPO="$WORKDIR/repository"
echo "[FTN] Fetching: $URL"
git clone --depth 1 --quiet -- "$URL" "$REPO"

PLAN="$WORKDIR/plan.tsv"
REPORT="$WORKDIR/REPORT.md"
: > "$PLAN"

classify() {
  local f="$1" base ext
  base="$(basename "$f")"
  ext="${base##*.}"

  case "$base" in
    Dockerfile*|docker-compose*.yml|docker-compose*.yaml|compose*.yml|compose*.yaml) echo infrastructure; return;;
    go.mod|go.sum) echo backend; return;;
    package.json|package-lock.json|pnpm-lock.yaml|yarn.lock|bun.lockb) echo frontend; return;;
    *.test.go|*_test.go|*.spec.ts|*.test.ts|*.spec.tsx|*.test.tsx) echo tests; return;;
    *.sql) echo database; return;;
    *.md|*.mdx|*.rst) echo docs; return;;
    *.sh) echo scripts; return;;
    *.env|*.env.*|*.toml|*.ini|*.conf|*.yaml|*.yml|*.json) echo configs; return;;
  esac

  case "$ext" in
    go) echo backend;;
    ts|tsx|js|jsx|css|scss|html) echo frontend;;
    py|rs|java|kt|c|cc|cpp|h|hpp) echo services;;
    *) echo archive;;
  esac
}

while IFS= read -r -d '' f; do
  rel="${f#$REPO/}"
  case "$rel" in
    .git/*|node_modules/*|vendor/*|dist/*|build/*|.next/*|coverage/*) continue;;
  esac

  case "$rel" in
    backend/*|frontend/*|services/*|modules/*|agents/*|database/*|infrastructure/*|scripts/*|configs/*|tests/*|docs/*|archive/*) continue;;
  esac

  category="$(classify "$f")"
  printf '%s\t%s/%s\n' "$rel" "$category" "$rel" >> "$PLAN"
done < <(find "$REPO" -type f -print0)

count="$(wc -l < "$PLAN" | tr -d ' ' )"
{
  echo "# FTN Source Organizer Report"
  echo
  echo "- Source: $URL"
  echo "- Mode: $MODE"
  echo "- Proposed files: $count"
  echo
  echo "## Proposed moves"
  echo
  echo '| Source | Target |'
  echo '|---|---|'
  while IFS=$'\t' read -r src dst; do
    printf '| `%s` | `%s` |\n' "$src" "$dst"
  done < "$PLAN"
} > "$REPORT"

cat "$REPORT"
echo
echo "[FTN] Temporary workspace: $REPO"
echo "[FTN] This script never modifies the remote repository."

echo
if [[ "$MODE" == "plan" ]]; then
  echo "[FTN] PLAN ONLY complete. Review the report before applying locally."
  exit 0
fi

while IFS=$'\t' read -r src dst; do
  [[ -z "$src" ]] && continue
  from="$REPO/$src"
  to="$REPO/$dst"
  [[ -e "$to" ]] && { echo "[FTN] SKIP existing: $dst"; continue; }
  mkdir -p "$(dirname "$to")"
  mv -- "$from" "$to"
done < "$PLAN"

echo "[FTN] Organization applied to temporary workspace only."
