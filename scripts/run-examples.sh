#!/bin/sh
# Runs the docs' examples (examples/) against a live API, so they are known to run
# unchanged (US-02). The JavaScript, Python and Go ones must print the same lines.
set -eu
: "${VERNACULAR_KEY:?Set VERNACULAR_KEY to an API key from your account page}"
export VERNACULAR_URL="${VERNACULAR_URL:-http://localhost:8080}"
cd "$(dirname "$0")/../examples"
out=$(mktemp -d)
trap 'rm -rf "$out"' EXIT

./search.sh > /dev/null
node search.mjs > "$out/js.txt"
python3 search.py > "$out/py.txt"
go run search.go > "$out/go.txt"

test -s "$out/js.txt" || { echo "examples printed nothing" >&2; exit 1; }
diff "$out/js.txt" "$out/py.txt"
diff "$out/js.txt" "$out/go.txt"
echo "examples ok: curl, JavaScript, Python and Go ($(wc -l < "$out/js.txt" | tr -d ' ') entries)"
