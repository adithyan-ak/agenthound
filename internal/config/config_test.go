package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Clear any env vars that might interfere
	for _, key := range []string{
		"AGENTHOUND_NEO4J_URI", "AGENTHOUND_NEO4J_USER", "AGENTHOUND_NEO4J_PASSWORD",
		"AGENTHOUND_PG_URI", "AGENTHOUND_API_PORT", "AGENTHOUND_LOG_LEVEL",
		"AGENTHOUND_JWT_SECRET", "AGENTHOUND_CORS_ORIGINS", "AGENTHOUND_ADMIN_PASSWORD",
	} {
		t.Setenv(key, "")
		_ = os.Unsetenv(key)
	}

	cfg := Load()

	if cfg.Neo4jURI != "bolt://localhost:7687" {
		t.Errorf("Neo4jURI = %q, want default", cfg.Neo4jURI)
	}
	if cfg.Neo4jUser != "neo4j" {
		t.Errorf("Neo4jUser = %q, want neo4j", cfg.Neo4jUser)
	}
	if cfg.APIPort != 8080 {
		t.Errorf("APIPort = %d, want 8080", cfg.APIPort)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want info", cfg.LogLevel)
	}
	if cfg.AdminPassword != "agenthound" {
		t.Errorf("AdminPassword = %q, want agenthound", cfg.AdminPassword)
	}
	if len(cfg.CORSOrigins) != 1 || cfg.CORSOrigins[0] != "http://localhost:8080" {
		t.Errorf("CORSOrigins = %v, want [http://localhost:8080]", cfg.CORSOrigins)
	}
	if cfg.JWTSecret != "" {
		t.Errorf("JWTSecret should be empty before validation, got %q", cfg.JWTSecret)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("AGENTHOUND_JWT_SECRET", "test-secret-key")
	t.Setenv("AGENTHOUND_CORS_ORIGINS", "http://localhost:3000, https://app.example.com")
	t.Setenv("AGENTHOUND_ADMIN_PASSWORD", "custom-pass")
	t.Setenv("AGENTHOUND_API_PORT", "9090")

	cfg := Load()

	if cfg.JWTSecret != "test-secret-key" {
		t.Errorf("JWTSecret = %q, want test-secret-key", cfg.JWTSecret)
	}
	if cfg.AdminPassword != "custom-pass" {
		t.Errorf("AdminPassword = %q, want custom-pass", cfg.AdminPassword)
	}
	if cfg.APIPort != 9090 {
		t.Errorf("APIPort = %d, want 9090", cfg.APIPort)
	}
	if len(cfg.CORSOrigins) != 2 {
		t.Fatalf("CORSOrigins len = %d, want 2", len(cfg.CORSOrigins))
	}
	if cfg.CORSOrigins[0] != "http://localhost:3000" {
		t.Errorf("CORSOrigins[0] = %q, want http://localhost:3000", cfg.CORSOrigins[0])
	}
	if cfg.CORSOrigins[1] != "https://app.example.com" {
		t.Errorf("CORSOrigins[1] = %q, want https://app.example.com", cfg.CORSOrigins[1])
	}
}

func TestValidateGeneratesJWTSecret(t *testing.T) {
	cfg := &Config{
		Neo4jURI:    "bolt://localhost:7687",
		PostgresURI: "postgres://localhost/test",
		APIPort:     8080,
		LogLevel:    "info",
		CORSOrigins: []string{"http://localhost:8080"},
	}

	err := cfg.Validate()
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if cfg.JWTSecret == "" {
		t.Error("Validate() should generate JWTSecret when empty")
	}
	if len(cfg.JWTSecret) != 64 { // 32 bytes hex-encoded
		t.Errorf("generated JWTSecret len = %d, want 64 hex chars", len(cfg.JWTSecret))
	}
}

func TestValidatePreservesExplicitJWTSecret(t *testing.T) {
	cfg := &Config{
		Neo4jURI:    "bolt://localhost:7687",
		PostgresURI: "postgres://localhost/test",
		APIPort:     8080,
		LogLevel:    "info",
		JWTSecret:   "explicit-secret",
		CORSOrigins: []string{"http://localhost:8080"},
	}

	err := cfg.Validate()
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if cfg.JWTSecret != "explicit-secret" {
		t.Errorf("JWTSecret = %q, want explicit-secret", cfg.JWTSecret)
	}
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want string
	}{
		{
			name: "invalid port",
			cfg:  Config{Neo4jURI: "bolt://x", PostgresURI: "postgres://x", APIPort: 0, LogLevel: "info"},
			want: "invalid API port",
		},
		{
			name: "invalid log level",
			cfg:  Config{Neo4jURI: "bolt://x", PostgresURI: "postgres://x", APIPort: 8080, LogLevel: "verbose"},
			want: "invalid log level",
		},
		{
			name: "empty neo4j URI",
			cfg:  Config{Neo4jURI: "", PostgresURI: "postgres://x", APIPort: 8080, LogLevel: "info"},
			want: "neo4j URI must not be empty",
		},
		{
			name: "empty postgres URI",
			cfg:  Config{Neo4jURI: "bolt://x", PostgresURI: "", APIPort: 8080, LogLevel: "info"},
			want: "postgres URI must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if err == nil {
				t.Fatal("Validate() should return error")
			}
			if got := err.Error(); !contains(got, tt.want) {
				t.Errorf("error = %q, should contain %q", got, tt.want)
			}
		})
	}
}

func TestCORSOriginsEmpty(t *testing.T) {
	t.Setenv("AGENTHOUND_CORS_ORIGINS", "")
	os.Unsetenv("AGENTHOUND_CORS_ORIGINS")

	cfg := Load()
	if len(cfg.CORSOrigins) != 1 || cfg.CORSOrigins[0] != "http://localhost:8080" {
		t.Errorf("empty CORS_ORIGINS should default to [http://localhost:8080], got %v", cfg.CORSOrigins)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
