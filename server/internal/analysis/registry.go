package analysis

import "github.com/adithyan-ak/agenthound/server/internal/analysis/processors"

func allProcessors() []PostProcessor {
	return []PostProcessor{
		&processors.HasAccessTo{},
		&processors.CanExecute{},
		&processors.Shadows{},
		&processors.PoisonedDescription{},
		&processors.PoisonedInstructions{},
		&processors.CanReach{},
		// CrossServiceCredentialChain depends on can_reach and emits
		// CAN_REACH edges that span Config Collector + LiteLLM Looter
		// emissions joined on Credential.value_hash. See
		// cross_service_credential_chain.go for the full path
		// description; design rationale in
		// docs/plans/sprint3-offensive-primitives.md 3.7 + 4.5.
		&processors.CrossServiceCredentialChain{},
		&processors.CanExfiltrate{},
		&processors.CanImpersonate{},
		&processors.CrossProtocol{},
		&processors.RiskScore{},
	}
}
