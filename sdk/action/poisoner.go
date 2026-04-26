package action

import "context"

// Poisoner injects content into an upstream artifact a Target consumes
// (config file, instruction file, tool description, etc.). Composes
// Reverter so every poison can be undone — this is non-negotiable for
// red-team safety.
//
// Implementations also implement sdk/module.Module.
type Poisoner interface {
	Reverter
	Poison(ctx context.Context, t Target, payload PoisonPayload) (*PoisonReceipt, error)
}

// PoisonPayload is a v0 stub.
type PoisonPayload struct{}

// PoisonReceipt is a v0 stub. It will satisfy the Receipt marker so it can
// be passed back to Revert.
type PoisonReceipt struct{}
