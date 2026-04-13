package cli

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestResolveServerURL_Default(t *testing.T) {
	t.Setenv("AGENTHOUND_SERVER_URL", "")
	cmd := &cobra.Command{}
	cmd.PersistentFlags().String("server-url", "", "")
	root := &cobra.Command{}
	root.PersistentFlags().String("server-url", "", "")
	root.AddCommand(cmd)

	got := resolveServerURL(cmd)
	if got != "http://localhost:8080" {
		t.Errorf("resolveServerURL() = %q, want %q", got, "http://localhost:8080")
	}
}

func TestResolveServerURL_Env(t *testing.T) {
	t.Setenv("AGENTHOUND_SERVER_URL", "http://myhost:9090")
	cmd := &cobra.Command{}
	root := &cobra.Command{}
	root.PersistentFlags().String("server-url", "", "")
	root.AddCommand(cmd)

	got := resolveServerURL(cmd)
	if got != "http://myhost:9090" {
		t.Errorf("resolveServerURL() = %q, want %q", got, "http://myhost:9090")
	}
}

func TestResolveServerURL_Flag(t *testing.T) {
	t.Setenv("AGENTHOUND_SERVER_URL", "http://env-url:9090")
	cmd := &cobra.Command{}
	root := &cobra.Command{}
	root.PersistentFlags().String("server-url", "", "")
	root.AddCommand(cmd)
	_ = root.PersistentFlags().Set("server-url", "http://flag-url:7070")

	got := resolveServerURL(cmd)
	if got != "http://flag-url:7070" {
		t.Errorf("resolveServerURL() = %q, want %q", got, "http://flag-url:7070")
	}
}

func TestPromptInput_Default(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.WriteString("\n")
	_ = w.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	got := promptInput("Username", "admin")
	if got != "admin" {
		t.Errorf("promptInput() = %q, want %q", got, "admin")
	}
}

func TestPromptInput_Custom(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.WriteString("myuser\n")
	_ = w.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	got := promptInput("Username", "admin")
	if got != "myuser" {
		t.Errorf("promptInput() = %q, want %q", got, "myuser")
	}
}

func TestPromptInput_EmptyDefault(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.WriteString("\n")
	_ = w.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	got := promptInput("Token", "")
	if got != "" {
		t.Errorf("promptInput() = %q, want %q", got, "")
	}
}

func TestPromptPassword_Pipe(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.WriteString("secret123\n")
	_ = w.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	got := promptPassword("Password")
	if got != "secret123" {
		t.Errorf("promptPassword() = %q, want %q", got, "secret123")
	}
}
