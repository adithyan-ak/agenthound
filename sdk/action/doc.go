// Package action declares the AgentHound action interface contracts.
//
// Each Action represents a distinct phase of an AI infrastructure red-team
// engagement. Modules implement one or more of these interfaces to plug into
// the framework via the sdk/module registry.
//
//	Scan         — CIDR/range expansion → []Target
//	Fingerprint  — single Target → service identification
//	Enumerate    — single Target → graph patch (today's Collector behaviour)
//	Loot         — extract latent secrets / state from a service
//	Extract      — pull a specific resource by reference
//	Poison       — inject content into an upstream artifact (composes Reverter)
//	Implant      — install a persistent payload (composes Reverter)
//
// Stability: v0 — UNSTABLE. Action signatures, the Target struct shape, and
// the Module.IsDestructive() bool contract may change before v1. The Reverter
// super-interface lives here from day one because adding it later would be a
// breaking change to every Poisoner / Implanter implementation.
//
// A real module satisfies BOTH this package's action interface AND the
// sdk/module.Module interface. The action interfaces here deliberately do
// NOT embed module.Module to avoid an import cycle (sdk/module depends on
// sdk/action for the Action enum). Implementations declare both contracts
// explicitly.
package action
