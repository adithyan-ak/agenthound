package action

import "context"

// Reverter is the destructive-action super-interface. Every Poisoner and
// Implanter must compose it so any change made on-target can be undone.
//
// This lives in the SDK from day one because adding it later would be a
// breaking change to every existing destructive-action implementation.
type Reverter interface {
	Revert(ctx context.Context, receipt Receipt) error
}

// Receipt is an empty marker interface. Each destructive action returns a
// concrete receipt type (PoisonReceipt, ImplantReceipt) that satisfies
// Receipt and carries whatever metadata that action needs to undo itself.
type Receipt interface{}
