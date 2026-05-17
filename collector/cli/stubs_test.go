package cli

import (
	"testing"
)

// TestStubsRegistered verifies each reserved verb is wired into rootCmd so
// `agenthound --help` lists it. This is the contract of "reserve the verb
// space without implementing anything" — break it and a future PR that adds
// `extract|poison|implant` may collide with an existing alias.
//
// TestNoStubsRemain confirms all offensive verbs have real implementations.
// v0.5 shipped the last stub (extract). No stubs should remain.
func TestNoStubsRemain(t *testing.T) {
	for _, verb := range []string{"scan", "loot", "poison", "implant", "revert", "discover", "extract"} {
		t.Run(verb, func(t *testing.T) {
			cmd, _, err := rootCmd.Find([]string{verb})
			if err != nil {
				t.Fatalf("rootCmd.Find(%q): %v", verb, err)
			}
			if cmd == nil || cmd.RunE == nil {
				t.Fatalf("verb %q not registered or has no RunE", verb)
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
