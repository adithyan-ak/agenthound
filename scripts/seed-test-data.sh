#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${AGENTHOUND_URL:-http://localhost:8080}"

# /api/v1/ingest is a mutating endpoint behind the localhost bearer
# token (server/internal/api/middleware/localtoken.go). Resolve the
# token the same way the server does:
#   1. AGENTHOUND_TOKEN explicit override
#   2. AGENTHOUND_TOKEN_PATH file
#   3. $XDG_CONFIG_HOME/agenthound/server.token
#   4. $HOME/.agenthound/server.token  (default)
# Fetch from the running server's /auth/local-token as a last resort —
# that endpoint is intentionally unauthenticated for exactly this case.
resolve_token() {
    if [ -n "${AGENTHOUND_TOKEN:-}" ]; then
        printf '%s' "$AGENTHOUND_TOKEN"
        return
    fi
    local path="${AGENTHOUND_TOKEN_PATH:-}"
    if [ -z "$path" ]; then
        if [ -n "${XDG_CONFIG_HOME:-}" ]; then
            path="${XDG_CONFIG_HOME}/agenthound/server.token"
        else
            path="${HOME}/.agenthound/server.token"
        fi
    fi
    if [ -r "$path" ]; then
        # Strip any trailing newline.
        tr -d '\n' < "$path"
        return
    fi
    # Fall back to the unauthenticated bootstrap endpoint.
    curl -s --max-time 2 "${BASE_URL}/api/v1/auth/local-token" \
        | sed -n 's/.*"token":"\([^"]*\)".*/\1/p'
}

# `|| true` is load-bearing: the curl fallback inside resolve_token can
# exit 7 (connection refused) when the server is down, which under
# `set -e` would kill the script silently before the friendly error
# block below runs. Swallowing the nonzero status here lets the empty
# TOKEN check below print the actionable diagnostic.
TOKEN=$(resolve_token || true)
if [ -z "$TOKEN" ]; then
    echo "ERROR: could not resolve the localhost bearer token."
    echo "  Start the server first (it generates the token on first run),"
    echo "  or set AGENTHOUND_TOKEN explicitly."
    exit 1
fi

for file in testdata/valid_*.json; do
    echo "Ingesting $file..."
    response=$(curl -s -w "\n%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TOKEN}" \
        -d @"$file" \
        "${BASE_URL}/api/v1/ingest")

    http_code=$(printf '%s\n' "$response" | tail -n 1)
    # `head -n -1` is GNU-only (not BSD/macOS). Use sed to drop the
    # last line portably.
    body=$(printf '%s\n' "$response" | sed '$d')

    if [ "$http_code" = "200" ]; then
        echo "  OK: $body"
    else
        echo "  FAILED ($http_code): $body"
        exit 1
    fi
done

echo ""
echo "Seed complete. Checking stats..."
curl -s "${BASE_URL}/api/v1/graph/stats" | python3 -m json.tool 2>/dev/null || \
    curl -s "${BASE_URL}/api/v1/graph/stats"
