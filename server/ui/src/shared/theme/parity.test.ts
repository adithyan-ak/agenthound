import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { describe, it, expect } from "vitest";
import { ACCENT, ACCENT_BRIGHT, SIGNAL_OK } from "./tokens";

// The design system keeps two blessed hex sources in lockstep: tokens.ts
// (TS / SVG / chart consumers) and globals.css (CSS variables for Tailwind /
// the DOM). slop-check bans hex anywhere ELSE but cannot detect when these
// two sources drift apart. This parity test closes that gap: it derives each
// overlapping color from BOTH files and asserts equality. It intentionally
// hardcodes NO hex literal — it imports the token value and reads the CSS
// variable text from globals.css — so the test itself stays within
// slop-check rule #1. (globals.css is read from disk because vitest mocks
// CSS module imports to empty strings.)

// vitest runs from the UI package root (server/ui), so resolve from cwd.
const css = readFileSync(
  resolve(process.cwd(), "src/shared/styles/globals.css"),
  "utf8",
);

function hexToRgb(hex: string): [number, number, number] {
  const body = hex.replace("#", "");
  return [
    parseInt(body.slice(0, 2), 16),
    parseInt(body.slice(2, 4), 16),
    parseInt(body.slice(4, 6), 16),
  ];
}

function cssVarRgb(name: string): [number, number, number] {
  const match = css.match(new RegExp(`--${name}:\\s*(\\d+)\\s+(\\d+)\\s+(\\d+)`));
  if (!match) {
    throw new Error(`CSS variable --${name} not found in globals.css`);
  }
  return [Number(match[1]), Number(match[2]), Number(match[3])];
}

// token (hex export) ↔ globals.css CSS variable. These are the only colors
// defined in BOTH files: node-kind / severity hues live only in tokens.ts,
// and the carbon (mauve) ramp lives only in globals.css, so neither can
// drift. The signature accents (amber phosphor + signal green) appear in
// both and must agree.
const PAIRS: ReadonlyArray<readonly [string, string, string]> = [
  ["ACCENT", ACCENT, "phosphor-raw"],
  ["ACCENT_BRIGHT", ACCENT_BRIGHT, "phosphor-bright-raw"],
  ["SIGNAL_OK", SIGNAL_OK, "signal-raw"],
];

describe("theme token <-> CSS var parity", () => {
  for (const [name, hex, varName] of PAIRS) {
    it(`${name} equals --${varName}`, () => {
      expect(hexToRgb(hex)).toEqual(cssVarRgb(varName));
    });
  }
});
