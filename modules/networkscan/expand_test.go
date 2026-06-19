package networkscan

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestExpand_SingleHost(t *testing.T) {
	t.Run("private ipv4 ok", func(t *testing.T) {
		got, err := Expand("10.0.0.5", ExpandOptions{})
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		if len(got) != 1 || got[0] != "10.0.0.5" {
			t.Errorf("got %v, want [10.0.0.5]", got)
		}
	})

	t.Run("public ipv4 refused without flag", func(t *testing.T) {
		_, err := Expand("8.8.8.8", ExpandOptions{})
		if !errors.Is(err, ErrPublicTarget) {
			t.Errorf("err = %v, want ErrPublicTarget", err)
		}
	})

	t.Run("public ipv4 allowed with flag", func(t *testing.T) {
		got, err := Expand("8.8.8.8", ExpandOptions{AllowPublicTargets: true})
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		if len(got) != 1 || got[0] != "8.8.8.8" {
			t.Errorf("got %v, want [8.8.8.8]", got)
		}
	})

	t.Run("link-local ipv4 always refused", func(t *testing.T) {
		_, err := Expand("169.254.0.5", ExpandOptions{AllowPublicTargets: true})
		if !errors.Is(err, ErrLinkLocal) {
			t.Errorf("err = %v, want ErrLinkLocal", err)
		}
	})

	t.Run("link-local ipv6 always refused", func(t *testing.T) {
		_, err := Expand("fe80::1", ExpandOptions{AllowPublicTargets: true})
		if !errors.Is(err, ErrLinkLocal) {
			t.Errorf("err = %v, want ErrLinkLocal", err)
		}
	})

	t.Run("multicast ipv4 always refused", func(t *testing.T) {
		_, err := Expand("224.0.0.1", ExpandOptions{AllowPublicTargets: true})
		if !errors.Is(err, ErrMulticast) {
			t.Errorf("err = %v, want ErrMulticast", err)
		}
	})

	t.Run("multicast ipv6 always refused", func(t *testing.T) {
		_, err := Expand("ff02::1", ExpandOptions{AllowPublicTargets: true})
		if !errors.Is(err, ErrMulticast) {
			t.Errorf("err = %v, want ErrMulticast", err)
		}
	})

	t.Run("hostname passes through", func(t *testing.T) {
		got, err := Expand("internal.example.local", ExpandOptions{AllowPublicTargets: true})
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if len(got) != 1 || got[0] != "internal.example.local" {
			t.Errorf("got %v, want [internal.example.local]", got)
		}
	})

	t.Run("empty spec rejected", func(t *testing.T) {
		_, err := Expand("", ExpandOptions{})
		if !errors.Is(err, ErrInvalidCIDR) {
			t.Errorf("err = %v, want ErrInvalidCIDR", err)
		}
	})

	t.Run("ipv6 ula private ok without flag", func(t *testing.T) {
		got, err := Expand("fc00::1", ExpandOptions{})
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		if len(got) != 1 || got[0] != "fc00::1" {
			t.Errorf("got %v, want [fc00::1]", got)
		}
	})
}

func TestExpand_CIDR(t *testing.T) {
	t.Run("/30 yields 4 ipv4 addresses", func(t *testing.T) {
		got, err := Expand("10.0.0.0/30", ExpandOptions{})
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		want := []string{"10.0.0.0", "10.0.0.1", "10.0.0.2", "10.0.0.3"}
		if !sliceEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("/32 yields 1 address", func(t *testing.T) {
		got, err := Expand("10.0.0.42/32", ExpandOptions{})
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if len(got) != 1 || got[0] != "10.0.0.42" {
			t.Errorf("got %v, want [10.0.0.42]", got)
		}
	})

	t.Run("/128 yields 1 ipv6 address", func(t *testing.T) {
		got, err := Expand("fc00::1/128", ExpandOptions{})
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if len(got) != 1 {
			t.Errorf("got %d hosts, want 1", len(got))
		}
	})

	t.Run("/16 (cap boundary) refused without flag", func(t *testing.T) {
		_, err := Expand("10.0.0.0/8", ExpandOptions{})
		if !errors.Is(err, ErrLargeCIDR) {
			t.Errorf("err = %v, want ErrLargeCIDR", err)
		}
	})

	t.Run("/16 allowed with flag", func(t *testing.T) {
		// Use /20 to keep the test fast — enumerating /16 (~65k) is overkill.
		got, err := Expand("10.0.0.0/20", ExpandOptions{AllowLargeCIDR: true})
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		// /20 = 4096 addresses
		if len(got) != 4096 {
			t.Errorf("got %d hosts, want 4096", len(got))
		}
	})

	t.Run("public CIDR refused without flag", func(t *testing.T) {
		_, err := Expand("8.8.8.0/30", ExpandOptions{})
		if !errors.Is(err, ErrPublicTarget) {
			t.Errorf("err = %v, want ErrPublicTarget", err)
		}
	})

	t.Run("link-local CIDR refused even with public flag", func(t *testing.T) {
		_, err := Expand("169.254.0.0/16", ExpandOptions{AllowPublicTargets: true})
		if !errors.Is(err, ErrLinkLocal) {
			t.Errorf("err = %v, want ErrLinkLocal", err)
		}
	})

	t.Run("invalid CIDR rejected", func(t *testing.T) {
		_, err := Expand("not-a-cidr/24", ExpandOptions{})
		if !errors.Is(err, ErrInvalidCIDR) {
			t.Errorf("err = %v, want ErrInvalidCIDR", err)
		}
	})

	// Straddling CIDR: 192.168.0.0/15 has a private masked base (192.168.0.0)
	// but its range spans into 192.169.0.0/16, which is public. The per-IP
	// gate must catch the public address even though the network base passes.
	// /15 is below the /16 cap, so AllowLargeCIDR is required to reach the
	// per-IP enumeration (the cap check fires first otherwise).
	t.Run("straddling CIDR refused on public address without flag", func(t *testing.T) {
		_, err := Expand("192.168.0.0/15", ExpandOptions{AllowLargeCIDR: true})
		if !errors.Is(err, ErrPublicTarget) {
			t.Errorf("err = %v, want ErrPublicTarget", err)
		}
	})

	t.Run("straddling CIDR allowed with public flag", func(t *testing.T) {
		got, err := Expand("192.168.0.0/15", ExpandOptions{AllowLargeCIDR: true, AllowPublicTargets: true})
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		// /15 = 131072 addresses; just assert the enumeration succeeded and
		// includes a public address from the straddling half.
		if len(got) != 131072 {
			t.Errorf("got %d hosts, want 131072", len(got))
		}
	})
}

func TestExpand_TargetsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "targets.txt")
	content := `
# Lab CIDRs
10.0.0.0/30

# Single host
192.168.1.5

# blank lines and comments are skipped
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write tempfile: %v", err)
	}

	t.Run("file:// prefix", func(t *testing.T) {
		got, err := Expand("file://"+path, ExpandOptions{})
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		// 10.0.0.0/30 → 4 addrs, 192.168.1.5 → 1 addr; total 5.
		if len(got) != 5 {
			t.Errorf("got %d hosts, want 5", len(got))
		}
	})

	t.Run("@ prefix", func(t *testing.T) {
		got, err := Expand("@"+path, ExpandOptions{})
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if len(got) != 5 {
			t.Errorf("got %d hosts, want 5", len(got))
		}
	})

	t.Run("missing file errors", func(t *testing.T) {
		_, err := Expand("@/no/such/file.txt", ExpandOptions{})
		if err == nil {
			t.Error("expected error for missing file")
		}
	})

	t.Run("empty file errors", func(t *testing.T) {
		emptyPath := filepath.Join(dir, "empty.txt")
		if err := os.WriteFile(emptyPath, []byte("# comment only\n"), 0o600); err != nil {
			t.Fatalf("write tempfile: %v", err)
		}
		_, err := Expand("@"+emptyPath, ExpandOptions{})
		if !errors.Is(err, ErrTargetsFileEmpty) {
			t.Errorf("err = %v, want ErrTargetsFileEmpty", err)
		}
	})
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
