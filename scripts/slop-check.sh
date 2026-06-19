#!/usr/bin/env bash
# slop-check.sh — UI design-system regression gate.
#
# Run via `make slop-check`, also invoked by `make prerelease`.
# Returns non-zero if banned patterns reappear in server/ui/src.
#
# What we ban and why:
#
# 1. Inline hex literals in app code. The full color system lives in
#    server/ui/src/theme/tokens.ts. Any hex reintroduced elsewhere creates
#    a second source of truth and drifts on the next palette change.
#    Allowed: tokens.ts itself, and globals.css (CSS variable definitions).
#
# 2. Hardcoded `grid-cols-N` / `lg:grid-cols-N` Tailwind classes. We rebuilt
#    the layout system around Every Layout primitives (Stack/Switcher/
#    Sidebar/Grid/Cluster/Center). Hardcoded breakpoint grids reintroduce
#    media-query coupling and break container-driven responsiveness.
#    Allowed: ui/dialog.tsx (Radix portal), components/ui/layout/*.
#
# 3. `transition-all`. Rauno Craft rule — animate explicit properties
#    (transform, opacity, color, etc.) only. `transition-all` is a
#    correctness footgun (animates layout properties that thrash) and a
#    performance footgun (paints everything). Use transition-{colors,
#    transform,opacity,[width],...} instead.
#
# 4. `shadow-2xl` / `shadow-xl` on `.glass`. Refactoring UI rule — one
#    light source. Use the elev-1/2/3 scale defined in globals.css; they
#    encode a consistent y-offset+blur ramp from a single direction.
#
# 5. `font-bold` paired with `uppercase tracking-{wider,widest}`. Eyebrow
#    labels already carry 3 emphasis cues (size, case, tracking); piling
#    on a fourth makes everything shout. Use font-semibold for that role.
#
# Override an individual rule by exporting SLOP_SKIP="rule1,rule2".

set -e

cd "$(dirname "$0")/.."

UI_SRC="server/ui/src"
TOKENS_FILE="$UI_SRC/shared/theme/tokens.ts"
GLOBALS_FILE="$UI_SRC/shared/styles/globals.css"

if [ ! -d "$UI_SRC" ]; then
  echo "slop-check: $UI_SRC does not exist; skipping"
  exit 0
fi

failures=0
skip="${SLOP_SKIP:-}"

is_skipped() {
  case ",$skip," in
    *,"$1",*) return 0 ;;
    *)        return 1 ;;
  esac
}

run_rule() {
  local id="$1"; shift
  local label="$1"; shift
  if is_skipped "$id"; then
    echo "=== [$id] $label — SKIPPED via SLOP_SKIP ==="
    return
  fi
  echo "=== [$id] $label ==="
  if "$@"; then
    echo "  ok"
  else
    failures=$((failures + 1))
  fi
}

# 1. Inline hex literals outside the token files.
check_inline_hex() {
  local matches
  matches=$(grep -rn --include='*.ts' --include='*.tsx' \
    -E '#[0-9a-fA-F]{3,8}\b' "$UI_SRC" 2>/dev/null \
    | grep -v -E '^[^:]*tokens\.ts:' \
    | grep -v -E '^[^:]*globals\.css:' \
    | grep -v -E '#[0-9a-fA-F]+/' \
    | grep -v -E '//.*#[0-9a-fA-F]' \
    | grep -v -E '\* .*#[0-9a-fA-F]' || true)
  if [ -n "$matches" ]; then
    echo "  FAIL: inline hex outside tokens.ts. Move to theme/tokens.ts."
    echo "$matches" | head -10 | sed 's/^/    /'
    return 1
  fi
}

# 2. Hardcoded grid-cols-N classes outside the layout primitives + dialog.
check_hardcoded_grids() {
  local matches
  matches=$(grep -rn --include='*.tsx' --include='*.ts' \
    -E '(^|[^a-z-])(sm:|md:|lg:|xl:)?grid-cols-[0-9]' "$UI_SRC" 2>/dev/null \
    | grep -v -E '^[^:]*components/ui/layout/' \
    | grep -v -E '^[^:]*components/ui/dialog\.tsx:' || true)
  if [ -n "$matches" ]; then
    echo "  FAIL: hardcoded grid-cols-N. Use Stack/Switcher/Sidebar/Grid from components/ui/layout."
    echo "$matches" | head -10 | sed 's/^/    /'
    return 1
  fi
}

# 3. transition-all anywhere.
check_transition_all() {
  local matches
  matches=$(grep -rn --include='*.tsx' --include='*.ts' --include='*.css' \
    'transition-all' "$UI_SRC" 2>/dev/null || true)
  if [ -n "$matches" ]; then
    echo "  FAIL: transition-all is banned. Use transition-{colors,transform,opacity,[width],...}."
    echo "$matches" | head -10 | sed 's/^/    /'
    return 1
  fi
}

# 4. shadow-2xl / shadow-xl in app code.
check_heavy_shadows() {
  local matches
  matches=$(grep -rn --include='*.tsx' --include='*.ts' \
    -E 'shadow-2xl|shadow-xl' "$UI_SRC" 2>/dev/null || true)
  if [ -n "$matches" ]; then
    echo "  FAIL: shadow-{xl,2xl} is banned. Use elev-1 / elev-2 / elev-3 from globals.css."
    echo "$matches" | head -10 | sed 's/^/    /'
    return 1
  fi
}

# 5. font-bold combined with uppercase tracking-{wider,widest}.
check_eyebrow_bold() {
  local matches
  matches=$(grep -rn --include='*.tsx' --include='*.ts' \
    -E 'uppercase tracking-(wider|widest)[^"]*font-bold|font-bold[^"]*uppercase tracking-(wider|widest)' \
    "$UI_SRC" 2>/dev/null || true)
  if [ -n "$matches" ]; then
    echo "  FAIL: eyebrow labels (uppercase + tracking) should be font-semibold, not font-bold."
    echo "$matches" | head -10 | sed 's/^/    /'
    return 1
  fi
}

# Sanity: token files must exist — otherwise this gate isn't actually
# protecting anything. Fail loudly rather than silently passing.
if [ ! -f "$TOKENS_FILE" ]; then
  echo "slop-check: $TOKENS_FILE missing — token system has been removed?"
  exit 1
fi
if [ ! -f "$GLOBALS_FILE" ]; then
  echo "slop-check: $GLOBALS_FILE missing — global style system has been removed?"
  exit 1
fi

run_rule inline-hex      "inline hex literals outside tokens.ts" check_inline_hex
run_rule hardcoded-grids "hardcoded grid-cols-N (use layout primitives)" check_hardcoded_grids
run_rule transition-all  "transition-all (use explicit transition-* )" check_transition_all
run_rule heavy-shadows   "shadow-{xl,2xl} (use elev-1/2/3)" check_heavy_shadows
run_rule eyebrow-bold    "font-bold on uppercase+tracking eyebrows" check_eyebrow_bold

echo
if [ "$failures" -gt 0 ]; then
  echo "slop-check: $failures rule(s) failed."
  exit 1
fi
echo "slop-check: all rules passed."
