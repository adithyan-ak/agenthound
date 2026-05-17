package cli

import (
	"errors"
	"testing"
)

// TestStubsRegistered verifies each reserved verb is wired into rootCmd so
// `agenthound --help` lists it. This is the contract of "reserve the verb
// space without implementing anything" — break it and a future PR that adds
// `extract|poison|implant` may collide with an existing alias.
//
// loot landed in v0.2 (collector/cli/loot.go) and is no longer a stub.
func TestStubsRegistered(t *testing.T) {
	for _, verb := range []string{"extract"} {
		t.Run(verb, func(t *testing.T) {
			cmd, _, err := rootCmd.Find([]string{verb})
			if err != nil {
				t.Fatalf("rootCmd.Find(%q): %v", verb, err)
			}
			if cmd == nil || cmd.Use != verb {
				t.Fatalf("verb %q not registered (got %v)", verb, cmd)
			}
			if cmd.RunE == nil {
				t.Fatalf("verb %q has no RunE (cobra would print help instead of returning the stub error)", verb)
			}
		})
	}
}

// TestStubsReturnError confirms the stub RunE returns errStubNotImplemented.
// main.go translates a non-nil RunE error into exit-code 1, so this is the
// load-bearing assertion that the still-stubbed verbs fail with non-zero
// status. (loot is no longer in this list; it's a real command in v0.2.)
func TestStubsReturnError(t *testing.T) {
	for _, verb := range []string{"extract"} {
		t.Run(verb, func(t *testing.T) {
			cmd, _, err := rootCmd.Find([]string{verb})
			if err != nil {
				t.Fatalf("rootCmd.Find(%q): %v", verb, err)
			}
			runErr := cmd.RunE(cmd, nil)
			if !errors.Is(runErr, errStubNotImplemented) {
				t.Errorf("RunE for %q returned %v, want errStubNotImplemented", verb, runErr)
			}
		})
	}
}

// TestLootRegistered confirms `loot` is wired but is NOT a stub anymore.
// Mirror of the contract above for the stubs that remain.
func TestLootRegistered(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"loot"})
	if err != nil {
		t.Fatalf("rootCmd.Find(loot): %v", err)
	}
	if cmd == nil || cmd.Use == "" {
		t.Fatalf("loot not registered (got %v)", cmd)
	}
	if cmd.RunE == nil {
		t.Fatal("loot has no RunE")
	}
}

// TestPoisonRegistered guards against an accidental revert to the v0.3
// stub state: `poison` MUST be wired and have a real RunE. v0.4 ships
// the destructive primitive; downgrading it back to a stub silently is
// a regression we want to fail loudly on.
func TestPoisonRegistered(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"poison"})
	if err != nil {
		t.Fatalf("rootCmd.Find(poison): %v", err)
	}
	if cmd == nil || cmd.Use == "" {
		t.Fatalf("poison not registered (got %v)", cmd)
	}
	if cmd.RunE == nil {
		t.Fatal("poison has no RunE")
	}
}

// TestRevertRegistered guards the v0.4 revert verb the same way.
func TestRevertRegistered(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"revert"})
	if err != nil {
		t.Fatalf("rootCmd.Find(revert): %v", err)
	}
	if cmd == nil || cmd.Use == "" {
		t.Fatalf("revert not registered (got %v)", cmd)
	}
	if cmd.RunE == nil {
		t.Fatal("revert has no RunE")
	}
}

// TestImplantRegistered guards the v0.4 implant verb the same way.
func TestImplantRegistered(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"implant"})
	if err != nil {
		t.Fatalf("rootCmd.Find(implant): %v", err)
	}
	if cmd == nil || cmd.Use == "" {
		t.Fatalf("implant not registered (got %v)", cmd)
	}
	if cmd.RunE == nil {
		t.Fatal("implant has no RunE")
	}
}
