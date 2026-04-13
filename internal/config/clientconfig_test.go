package config

import (
	"os"
	"path/filepath"
	"testing"
)

func configDir(home string) string {
	return filepath.Join(home, ".config", "agenthound")
}

func TestLoadClientConfigFromEnv(t *testing.T) {
	t.Setenv("AGENTHOUND_SERVER_URL", "http://remote:9090")
	t.Setenv("AGENTHOUND_API_TOKEN", "ah_envtoken")

	cfg, err := LoadClientConfig(nil)
	if err != nil {
		t.Fatalf("LoadClientConfig() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadClientConfig() returned nil")
	}
	if cfg.ServerURL != "http://remote:9090" {
		t.Errorf("ServerURL = %q, want http://remote:9090", cfg.ServerURL)
	}
	if cfg.APIToken != "ah_envtoken" {
		t.Errorf("APIToken = %q, want ah_envtoken", cfg.APIToken)
	}
}

func TestLoadClientConfigFromFile(t *testing.T) {
	tmp := t.TempDir()
	content := []byte("server_url: http://file-server:8080\napi_token: ah_filetoken\n")

	t.Setenv("AGENTHOUND_SERVER_URL", "")
	t.Setenv("AGENTHOUND_API_TOKEN", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", tmp)

	dir := configDir(tmp)
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), content, 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadClientConfig(nil)
	if err != nil {
		t.Fatalf("LoadClientConfig() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadClientConfig() returned nil, want config from file")
	}
	if cfg.ServerURL != "http://file-server:8080" {
		t.Errorf("ServerURL = %q, want http://file-server:8080", cfg.ServerURL)
	}
	if cfg.APIToken != "ah_filetoken" {
		t.Errorf("APIToken = %q, want ah_filetoken", cfg.APIToken)
	}
}

func TestLoadClientConfigNothingConfigured(t *testing.T) {
	t.Setenv("AGENTHOUND_SERVER_URL", "")
	t.Setenv("AGENTHOUND_API_TOKEN", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", t.TempDir())

	cfg, err := LoadClientConfig(nil)
	if err != nil {
		t.Fatalf("LoadClientConfig() error = %v", err)
	}
	if cfg != nil {
		t.Errorf("LoadClientConfig() = %+v, want nil when nothing configured", cfg)
	}
}

func TestEnvOverridesFile(t *testing.T) {
	tmp := t.TempDir()
	dir := configDir(tmp)
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	content := []byte("server_url: http://file-server:8080\napi_token: ah_filetoken\n")
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), content, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", tmp)

	t.Setenv("AGENTHOUND_SERVER_URL", "http://env-server:9090")
	t.Setenv("AGENTHOUND_API_TOKEN", "ah_envtoken")

	cfg, err := LoadClientConfig(nil)
	if err != nil {
		t.Fatalf("LoadClientConfig() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadClientConfig() returned nil")
	}
	if cfg.ServerURL != "http://env-server:9090" {
		t.Errorf("ServerURL = %q, want env override http://env-server:9090", cfg.ServerURL)
	}
	if cfg.APIToken != "ah_envtoken" {
		t.Errorf("APIToken = %q, want env override ah_envtoken", cfg.APIToken)
	}
}

func TestSaveClientConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", tmp)

	cfg := &ClientConfig{
		ServerURL: "http://saved:8080",
		APIToken:  "ah_savedtoken",
	}
	if err := SaveClientConfig(cfg); err != nil {
		t.Fatalf("SaveClientConfig() error = %v", err)
	}

	path := ClientConfigPath()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("config file not created: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions = %o, want 0600", perm)
	}

	loaded, err := LoadClientConfig(nil)
	if err != nil {
		t.Fatalf("LoadClientConfig() after save error = %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadClientConfig() returned nil after save")
	}
	if loaded.ServerURL != cfg.ServerURL {
		t.Errorf("ServerURL = %q, want %q", loaded.ServerURL, cfg.ServerURL)
	}
	if loaded.APIToken != cfg.APIToken {
		t.Errorf("APIToken = %q, want %q", loaded.APIToken, cfg.APIToken)
	}
}

func TestClientConfigPathExpected(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", tmp)

	path := ClientConfigPath()
	want := filepath.Join(tmp, ".config", "agenthound", "config.yaml")
	if path != want {
		t.Errorf("ClientConfigPath() = %q, want %q", path, want)
	}
}
