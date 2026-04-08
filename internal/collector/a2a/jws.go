package a2a

func VerifySignatures(cardJSON []byte, raw map[string]any) (signed bool, valid bool) {
	sigs, ok := raw["signatures"]
	if !ok {
		return false, false
	}

	sigArr, ok := sigs.([]any)
	if !ok || len(sigArr) == 0 {
		return false, false
	}

	// Phase 2 MVP: signatures field present and non-empty means the card is signed.
	// Full JWS verification (RS256/ES256 with public key discovery) is Phase 5 scope.
	// Conservative default: signed=true, valid=false (unverified).
	return true, false
}
