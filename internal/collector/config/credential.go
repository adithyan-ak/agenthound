package config

import (
	"strings"

	"github.com/adithyan-ak/agenthound/internal/collector/common"
)

type CredentialInfo struct {
	Type        string // "envVar", "hardcoded", "vaultRef", "inputPrompt"
	Name        string
	Value       string // SHA-256 hash by default, actual value only if includeValues=true
	Source      string
	IsExposed   bool
	HighEntropy bool
	Format      string // "openai", "anthropic", "github", "slack", "aws", "generic"
}

var credentialNamePatterns = []string{
	"KEY", "TOKEN", "SECRET", "PASSWORD", "CREDENTIAL", "AUTH", "API_KEY",
}

var vaultPrefixes = []string{
	"vault://",
	"ssm://",
	"arn:aws:secretsmanager",
}

func ExtractCredentials(env map[string]string, headers map[string]string, source string, includeValues bool) []CredentialInfo {
	var creds []CredentialInfo

	for name, value := range env {
		if !isCredentialName(name) {
			continue
		}
		creds = append(creds, classifyAndBuild(name, value, source, includeValues))
	}

	for name, value := range headers {
		if !isCredentialName(name) {
			continue
		}
		creds = append(creds, classifyAndBuild(name, value, source, includeValues))
	}

	return creds
}

func ClassifyCredentialType(name, value string) string {
	if isVaultRef(value) {
		return "vaultRef"
	}
	if isEnvRef(value) {
		return "envVar"
	}
	return "hardcoded"
}

func classifyAndBuild(name, value, source string, includeValues bool) CredentialInfo {
	ci := CredentialInfo{
		Name:   name,
		Source: source,
		Format: detectFormat(value),
		Type:   ClassifyCredentialType(name, value),
	}

	switch ci.Type {
	case "envVar":
		ci.IsExposed = false
	case "vaultRef":
		ci.IsExposed = false
	default:
		ci.IsExposed = true
		ci.HighEntropy = common.IsLikelySecret(value)
	}

	if includeValues {
		ci.Value = value
	} else {
		ci.Value = common.HashSHA256(value)
	}

	return ci
}

func isCredentialName(name string) bool {
	upper := strings.ToUpper(name)
	for _, pattern := range credentialNamePatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

func isVaultRef(value string) bool {
	lower := strings.ToLower(value)
	for _, prefix := range vaultPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

func isEnvRef(value string) bool {
	return strings.HasPrefix(value, "$") || strings.HasPrefix(value, "${")
}

func detectFormat(value string) string {
	if strings.HasPrefix(value, "sk-ant-") {
		return "anthropic"
	}
	if strings.HasPrefix(value, "sk-") {
		return "openai"
	}
	if strings.HasPrefix(value, "xoxb-") || strings.HasPrefix(value, "xoxp-") || strings.HasPrefix(value, "xoxs-") {
		return "slack"
	}
	if strings.HasPrefix(value, "ghp_") || strings.HasPrefix(value, "gho_") || strings.HasPrefix(value, "ghs_") {
		return "github"
	}
	if strings.HasPrefix(value, "AKIA") {
		return "aws"
	}
	return "generic"
}
