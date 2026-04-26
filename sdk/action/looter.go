package action

import "context"

// Looter extracts latent secrets, configuration, or state from a Target
// without modifying it. The returned LootResult satisfies a
// ToIngest() *ingest.IngestData contract (added when the first Looter
// implementation lands) so loot folds into the same graph pipeline.
//
// Implementations also implement sdk/module.Module.
type Looter interface {
	Loot(ctx context.Context, t Target, opts LootOptions) (*LootResult, error)
}

// LootOptions is a v0 stub.
type LootOptions struct{}

// LootResult is a v0 stub. Concrete shape (and the ToIngest method) lands
// with the first Looter implementation.
type LootResult struct{}
