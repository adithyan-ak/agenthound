package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type ClientConfig struct {
	ServerURL string `yaml:"server_url"`
	APIToken  string `yaml:"api_token"`
}

func ClientConfigPath() string {
	if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
		return filepath.Join(v, "agenthound", "config.yaml")
	}
	return filepath.Join(os.Getenv("HOME"), ".config", "agenthound", "config.yaml")
}

func LoadClientConfig(flags *pflag.FlagSet) (*ClientConfig, error) {
	cfg := &ClientConfig{}

	cfg.ServerURL = resolve(flags, "server-url", "AGENTHOUND_SERVER_URL", "")
	cfg.APIToken = resolve(flags, "api-token", "AGENTHOUND_API_TOKEN", "")

	if cfg.ServerURL == "" || cfg.APIToken == "" {
		fileCfg, err := loadClientConfigFile()
		if err != nil {
			return nil, fmt.Errorf("reading client config: %w", err)
		}
		if fileCfg != nil {
			if cfg.ServerURL == "" {
				cfg.ServerURL = fileCfg.ServerURL
			}
			if cfg.APIToken == "" {
				cfg.APIToken = fileCfg.APIToken
			}
		}
	}

	if cfg.ServerURL == "" && cfg.APIToken == "" {
		return nil, nil
	}

	return cfg, nil
}

func loadClientConfigFile() (*ClientConfig, error) {
	data, err := os.ReadFile(ClientConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var cfg ClientConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", ClientConfigPath(), err)
	}
	return &cfg, nil
}

func SaveClientConfig(cfg *ClientConfig) error {
	path := ClientConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
