#!/usr/bin/env bash
set -Eeuo pipefail

# FTN Source Organizer
# Production-oriented, approval-first source organization helper.
# Usage:
#   ./scripts/source-organizer.sh /path/to/source
#   ./scripts/source-organizer.sh /path/to/source --apply
#
# Default mode is PLAN ONLY. Nothing is moved or deleted without --apply.

ROOT="${1:-.}"
MODE="plan"
[[ "${2:-}" == "--apply" ]] && MODE="apply"

ROOT="$(cd "$ROOT" && pwd)"
PLAN="$ROOT/.ftn-source-organizer-plan"
mkdir -p "$PLAN"
: > "$PLAN/moves.tsv"

classify() {
  local f="$1" base ext
  base="$(basename "$f")"
  ext="${base##*.}"

  case "$base" in
    Dockerfile*|docker-compose*.yml|docker-compose*.yaml|compose*.yml|compose*.yaml) echo "infrastructure"; return;;
    go.mod|go.sum) echo "backend"; return;;
    package.json|package-lock.json|pnpm-lock.yaml|yarn.lock|bun.lockb) echo "frontend"; return;;
    *.test.go|*_test.go|*.spec.ts|*.test.ts|*.spec.tsx|*.test.tsx) echo "tests"; return;;
    *.sql) echo "database"; return;;
    *.md|*.mdx|*.rst|*.txt) echo "docs"; return;;
    *.sh) echo "scripts"; return;;
    *.env|*.env.*|*.toml|*.ini|*.conf|*.yaml|*.yml|*.json) echo "configs"; return;;
  esac

  case "$ext" in
    go) echo "backend";;
    ts|tsx|js|jsx|css|scss|html) echo "frontend";;
    py|rs|java|kt|c|cc|cpp|h|hpp) echo "services";;
    *) echo "archive";;
  esac
}

while IFS= read -r -d '' f; do
  rel="${f#$ROOT/}"
  case "$rel" in
    .git/*|.ftn-source-organizer-plan/*|node_modules/*|vendor/*|dist/*|build/*|.next/*|coverage/*) continue;;
  esac

  category="$(classify "$f")"
  target="$ROOT/$category/$rel"
  target_dir="$(dirname "$target")"

  # Do not reorganize files already under a canonical category directory.
  case "$rel" in
    backend/*|frontend/*|services/*|modules/*|agents/*|database/*|infrastructure/*|scripts/*|configs/*|tests/*|docs/*|archive/*) continue;;
  esac

  printf '%s\t%s\n' "$rel" "${target#$ROOT/}" >> "$PLAN/moves.tsv"
done < <(find "$ROOT" -type f -print0)

count="$(wc -l < "$PLAN/moves.tsv" | tr -d ' ' )"
echo "FTN Source Organizer"
echo "Root: $ROOT"
echo "Mode: $MODE"
echo "Proposed moves: $count"
echo "Plan: $PLAN/moves.tsv"

if [[ "$MODE" != "apply" ]]; then
  echo
  echo "PLAN ONLY: no files changed. Review the plan, then rerun with --apply."
  exit 0
fi

while IFS=$'\t' read -r rel dest; do
  [[ -z "$rel" ]] && continue
  src="$ROOT/$rel"
  dst="$ROOT/$dest"
  mkdir -p "$(dirname "$dst")"

  if [[ -e "$dst" ]]; then
    echo "SKIP existing: $dest"
    continue
  fi

  mv -- "$src" "$dst"
done < "$PLAN/moves.tsv"

echo "Applied safely. No source file was deleted."
