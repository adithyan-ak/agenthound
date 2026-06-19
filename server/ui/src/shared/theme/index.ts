// shared/theme barrel — sole TS hex source. Re-exports tokens.ts so callers
// can import design tokens from "@shared/theme" or the explicit
// "@shared/theme/tokens" path (slop-check excludes tokens.ts by basename).
export * from "./tokens";
