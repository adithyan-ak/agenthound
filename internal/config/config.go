package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

type Config struct {
	Neo4jURI      string
	Neo4jUser     string
	Neo4jPassword string
	PostgresURI   string
	APIPort       int
	LogLevel      string
	JWTSecret     string
	CORSOrigins   []string
	AdminPassword string
	dbExplicit    bool
}

// HasExplicitDBConfig returns true if the user explicitly set Neo4j or PG
// connection config via environment variables or CLI flags (not just defaults).
func (c *Config) HasExplicitDBConfig() bool {
	return c.dbExplicit
}

// LoadWithFlags creates a Config using flag values → env vars → defaults (in priority order).
func LoadWithFlags(flags *pflag.FlagSet) *Config {
	cfg := &Config{
		Neo4jURI:      "bolt://localhost:7687",
		Neo4jUser:     "neo4j",
		Neo4jPassword: "agenthound",
		PostgresURI:   "postgres://agenthound:agenthound@localhost:5432/agenthound?sslmode=disable",
		APIPort:       8080,
		LogLevel:      "info",
	}

	cfg.Neo4jURI = resolve(flags, "neo4j-uri", "AGENTHOUND_NEO4J_URI", cfg.Neo4jURI)
	cfg.Neo4jUser = resolve(flags, "neo4j-user", "AGENTHOUND_NEO4J_USER", cfg.Neo4jUser)
	cfg.Neo4jPassword = resolve(flags, "neo4j-password", "AGENTHOUND_NEO4J_PASSWORD", cfg.Neo4jPassword)
	cfg.PostgresURI = resolve(flags, "pg-uri", "AGENTHOUND_PG_URI", cfg.PostgresURI)

	cfg.dbExplicit = os.Getenv("AGENTHOUND_NEO4J_URI") != "" || os.Getenv("AGENTHOUND_PG_URI") != "" ||
		flagChanged(flags, "neo4j-uri") || flagChanged(flags, "pg-uri")
	cfg.LogLevel = resolve(flags, "log-level", "AGENTHOUND_LOG_LEVEL", cfg.LogLevel)
	cfg.JWTSecret = resolve(flags, "jwt-secret", "AGENTHOUND_JWT_SECRET", "")
	cfg.AdminPassword = resolve(flags, "admin-password", "AGENTHOUND_ADMIN_PASSWORD", "agenthound")

	if portStr := resolve(flags, "port", "AGENTHOUND_API_PORT", ""); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.APIPort = port
		}
	}

	if origins := resolve(flags, "cors-origins", "AGENTHOUND_CORS_ORIGINS", ""); origins != "" {
		for _, o := range strings.Split(origins, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				cfg.CORSOrigins = append(cfg.CORSOrigins, trimmed)
			}
		}
	}
	if len(cfg.CORSOrigins) == 0 {
		cfg.CORSOrigins = []string{"http://localhost:8080"}
	}

	return cfg
}

// Load creates a Config from env vars and defaults (no flags).
func Load() *Config {
	return LoadWithFlags(nil)
}

// Validate checks that all config values are valid.
func (c *Config) Validate() error {
	var errs []string

	if c.APIPort < 1 || c.APIPort > 65535 {
		errs = append(errs, fmt.Sprintf("invalid API port %d: must be 1-65535", c.APIPort))
	}

	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.LogLevel] {
		errs = append(errs, fmt.Sprintf("invalid log level %q: must be debug/info/warn/error", c.LogLevel))
	}

	if c.Neo4jURI == "" {
		errs = append(errs, "neo4j URI must not be empty")
	}
	if c.PostgresURI == "" {
		errs = append(errs, "postgres URI must not be empty")
	}

	if c.JWTSecret == "" {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			errs = append(errs, "failed to generate JWT secret")
		} else {
			c.JWTSecret = hex.EncodeToString(b)
			slog.Warn("no JWT secret configured, using random value (set AGENTHOUND_JWT_SECRET for stable tokens across restarts)")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("config validation: %s", strings.Join(errs, "; "))
	}
	return nil
}

func flagChanged(flags *pflag.FlagSet, name string) bool {
	if flags == nil {
		return false
	}
	f := flags.Lookup(name)
	return f != nil && f.Changed
}

// resolve returns the first non-empty value from: flag, env var, default.
func resolve(flags *pflag.FlagSet, flagName, envName, defaultVal string) string {
	if flags != nil {
		if f := flags.Lookup(flagName); f != nil && f.Changed {
			return f.Value.String()
		}
	}
	if v := os.Getenv(envName); v != "" {
		return v
	}
	return defaultVal
}
