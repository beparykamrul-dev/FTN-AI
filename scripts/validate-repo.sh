#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

printf '%s\n' '[1/4] validating shell scripts'
while IFS= read -r -d '' file; do
  bash -n "$file"
done < <(git ls-files -z -- '*.sh')

printf '%s\n' '[2/4] validating YAML files'
ruby -ryaml -e '
files = `git ls-files "*.yml" "*.yaml"`.lines.map(&:strip).reject(&:empty?)
errors = []
files.each do |file|
  begin
    YAML.load_stream(File.read(file), aliases: true)
  rescue StandardError => e
    errors << "#{file}: #{e.message}"
  end
end
abort(errors.join("\n")) unless errors.empty?
puts "validated #{files.length} YAML files"
'

printf '%s\n' '[3/4] checking Go formatting'
files=( $(git ls-files '*.go') )
if ((${#files[@]})); then
  unformatted="$(gofmt -l "${files[@]}")"
  if [[ -n "$unformatted" ]]; then
    printf '%s\n' "$unformatted"
    exit 1
  fi
fi

printf '%s\n' '[4/4] testing declared Go module'
if [[ -f modules/account/go.mod ]]; then
  (cd modules/account && go test ./...)
fi

printf '%s\n' 'FTN-AI repository validation passed.'
