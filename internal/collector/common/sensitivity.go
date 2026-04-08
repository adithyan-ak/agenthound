package common

import "strings"

type SensitivityLevel string

const (
	SensitivityCritical SensitivityLevel = "critical"
	SensitivityHigh     SensitivityLevel = "high"
	SensitivityMedium   SensitivityLevel = "medium"
	SensitivityLow      SensitivityLevel = "low"
)

var criticalFileExtensions = []string{".env", ".key", ".pem", ".crt", ".p12", ".pfx", ".jks"}
var criticalFilePaths = []string{"/credentials", "/secrets", "/.ssh/"}
var highConfigExtensions = []string{".conf", ".cfg", ".ini", ".yaml", ".yml", ".json"}

func ClassifyResourceSensitivity(uri string) SensitivityLevel {
	lower := strings.ToLower(uri)

	if isCritical(lower) {
		return SensitivityCritical
	}
	if isHigh(lower) {
		return SensitivityHigh
	}
	if isMedium(lower) {
		return SensitivityMedium
	}
	return SensitivityLow
}

func isCritical(uri string) bool {
	if hasScheme(uri, "postgres://", "postgresql://", "mysql://", "mongodb://") && strings.Contains(uri, "prod") {
		return true
	}
	if hasScheme(uri, "redis://") && strings.Contains(uri, "prod") {
		return true
	}

	if strings.HasPrefix(uri, "file://") {
		path := strings.TrimPrefix(uri, "file://")

		if path == "/etc/shadow" || path == "/etc/passwd" ||
			strings.HasPrefix(path, "/etc/shadow/") || strings.HasPrefix(path, "/etc/passwd/") {
			return true
		}
		if strings.HasPrefix(path, "/root/") || path == "/root" {
			return true
		}
		for _, ext := range criticalFileExtensions {
			if strings.HasSuffix(path, ext) {
				return true
			}
		}
		for _, seg := range criticalFilePaths {
			if strings.Contains(path, seg) {
				return true
			}
		}
	}

	if hasScheme(uri, "s3://", "gs://") && strings.Contains(uri, "prod") {
		return true
	}

	return false
}

func isHigh(uri string) bool {
	if hasScheme(uri, "postgres://", "postgresql://", "mysql://", "mongodb://") {
		return true
	}
	if hasScheme(uri, "redis://") {
		return true
	}

	if strings.HasPrefix(uri, "file://") {
		path := strings.TrimPrefix(uri, "file://")

		if strings.HasPrefix(path, "/var/log/") || path == "/var/log" {
			return true
		}
		for _, ext := range highConfigExtensions {
			if strings.HasSuffix(path, ext) &&
				(strings.Contains(path, "secret") || strings.Contains(path, "password")) {
				return true
			}
		}
	}

	return false
}

func isMedium(uri string) bool {
	if strings.HasPrefix(uri, "file:///") || strings.HasPrefix(uri, "file://localhost/") {
		return true
	}
	if hasScheme(uri, "http://", "https://") {
		return true
	}
	if hasScheme(uri, "s3://", "gs://") {
		return true
	}
	return false
}

func hasScheme(uri string, schemes ...string) bool {
	for _, s := range schemes {
		if strings.HasPrefix(uri, s) {
			return true
		}
	}
	return false
}
