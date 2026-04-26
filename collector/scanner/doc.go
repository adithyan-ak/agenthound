// Package scanner is a placeholder for the future network discovery engine.
//
// The scanner takes a CIDR range or list of targets and returns []action.Target
// (defined in sdk/action). When a real implementation lands, it will reuse
// sdk/rules/MatcherSpec for fingerprinting matches — there is intentionally no
// parallel matcher framework in this package.
//
// This package is collector-only; the server binary does not import it.
package scanner
