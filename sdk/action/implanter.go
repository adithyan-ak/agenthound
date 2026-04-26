package action

import "context"

// Implanter installs a persistent payload (cron entry, hook, modified
// binary, scheduled task) on or adjacent to a Target. Composes Reverter
// for the same reason Poisoner does.
//
// Implementations also implement sdk/module.Module.
type Implanter interface {
	Reverter
	Implant(ctx context.Context, t Target, payload ImplantPayload) (*ImplantReceipt, error)
}

// ImplantPayload is a v0 stub.
type ImplantPayload struct{}

// ImplantReceipt is a v0 stub. It will satisfy the Receipt marker so it can
// be passed back to Revert.
type ImplantReceipt struct{}
