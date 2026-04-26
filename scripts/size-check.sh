#!/usr/bin/env bash
# Measures the linux/amd64 stripped collector binary size and reports
# whether it has grown beyond the recorded baseline + 10%.
# Run via `make size-check`.
# Advisory in this commit; flipped to blocking in Step 7.

set -e

cd "$(dirname "$0")/.."

# BASELINE_BYTES recorded at Step 5 from a fresh prototype build below.
# Update this number consciously when a dep change intentionally raises the bar.
BASELINE_BYTES=9412792

OUT=$(mktemp)
trap 'rm -f "$OUT"' EXIT

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -trimpath \
  -ldflags='-s -w' \
  -o "$OUT" \
  ./collector/cmd/agenthound

BYTES=$(wc -c < "$OUT" | tr -d ' ')
LIMIT=$(( BASELINE_BYTES * 110 / 100 ))

echo "collector linux/amd64 size: $BYTES bytes ($((BYTES/1024/1024)) MiB)"
echo "baseline:                   $BASELINE_BYTES bytes"
echo "limit (baseline +10%):      $LIMIT bytes"

fail=0
if [ "$BYTES" -gt "$LIMIT" ]; then
  echo "ADVISORY: collector binary exceeds baseline + 10%."
  fail=1
fi

# Step 5: exit 0 always (advisory). Step 7 will change to: exit $fail
exit 0
