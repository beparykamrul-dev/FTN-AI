#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

printf '%s\n' '[1/5] validating shell scripts'
while IFS= read -r -d '' file; do
  bash -n "$file"
done < <(git ls-files -z -- '*.sh')

printf '%s\n' '[2/5] validating YAML files'
ruby -ryaml -e '
files = `git ls-files "*.yml" "*.yaml"`.lines.map(&:strip).reject(&:empty?)
errors = []
files.each do |file|
  begin
    if YAML.respond_to?(:unsafe_load_stream)
      YAML.unsafe_load_stream(File.read(file))
    else
      YAML.load_stream(File.read(file))
    end
  rescue StandardError => e
    errors << "#{file}: #{e.message}"
  end
end
abort(errors.join("\n")) unless errors.empty?
puts "validated #{files.length} YAML files"
'

printf '%s\n' '[3/5] validating JSON files'
if command -v python3 >/dev/null 2>&1; then
  python3 - <<'PY'
import json
import subprocess

files = subprocess.check_output(["git", "ls-files", "*.json"], text=True).splitlines()
for path in files:
    with open(path, encoding="utf-8") as fh:
        json.load(fh)
print(f"validated {len(files)} JSON files")
PY
fi

printf '%s\n' '[4/5] checking Go formatting for changed files'
mapfile -t changed_go < <(git diff-tree --no-commit-id --name-only -r HEAD -- '*.go')
if ((${#changed_go[@]})); then
  unformatted="$(gofmt -l "${changed_go[@]}")"
  if [[ -n "$unformatted" ]]; then
    printf '%s\n' "$unformatted"
    exit 1
  fi
fi

printf '%s\n' '[5/5] testing every declared Go module'
mapfile -t modules < <(git ls-files 'go.mod' | sort -u)
for mod in "${modules[@]}"; do
  dir="$(dirname "$mod")"
  printf 'testing %s\n' "$dir"
  (cd "$dir" && GOPROXY=direct go test ./...)
done

printf '%s\n' 'FTN-AI repository validation passed.'
