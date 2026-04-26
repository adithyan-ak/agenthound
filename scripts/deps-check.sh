#!/usr/bin/env bash
# Verifies that the collector binary doesn't pull server-only deps,
# and the server binary doesn't pull MCP SDK or modules/.
# Run via `make deps-check`.
# Blocking gate: returns non-zero if violations are found.
#
# Note on go-jose: legitimately used by modules/a2a for JWS signature
# verification of A2A Agent Cards (RFC 7515). It belongs in the collector
# and is intentionally NOT on the blocklist.

set -e

cd "$(dirname "$0")/.."

# Collector blocklist — must NOT match
collector_block=(
  '^github.com/neo4j/'
  '^github.com/jackc/pgx'
  '^github.com/go-chi/'
  '^github.com/golang-jwt/'
  '^golang.org/x/crypto/bcrypt'
  '^github.com/adithyan-ak/agenthound/server/internal/'
  '^github.com/adithyan-ak/agenthound/server/cli/'
  '^github.com/adithyan-ak/agenthound/server/ui/'
)

# Server blocklist — must NOT match
server_block=(
  '^github.com/modelcontextprotocol/go-sdk'
  '^github.com/adithyan-ak/agenthound/modules/'
)

fail=0

echo "=== collector deps ==="
collector_deps=$(go list -deps ./collector/cmd/agenthound)
for pat in "${collector_block[@]}"; do
  hits=$(echo "$collector_deps" | grep -E "$pat" || true)
  if [ -n "$hits" ]; then
    echo "ADVISORY: collector links forbidden dep matching $pat:"
    echo "$hits" | sed 's/^/    /'
    fail=1
  fi
done

echo "=== server deps ==="
server_deps=$(go list -deps ./server/cmd/agenthound-server)
for pat in "${server_block[@]}"; do
  hits=$(echo "$server_deps" | grep -E "$pat" || true)
  if [ -n "$hits" ]; then
    echo "ADVISORY: server links forbidden dep matching $pat:"
    echo "$hits" | sed 's/^/    /'
    fail=1
  fi
done

if [ "$fail" -eq 1 ]; then
  echo ""
  echo "deps-check found violations (blocking)."
fi

exit $fail
