package common

import "testing"

func TestClassifyResourceSensitivity(t *testing.T) {
	tests := []struct {
		name string
		uri  string
		want SensitivityLevel
	}{
		// Critical: production databases
		{"postgres prod", "postgres://user:pass@db-prod.internal:5432/mydb", SensitivityCritical},
		{"postgresql prod", "postgresql://user:pass@production-db:5432/app", SensitivityCritical},
		{"mysql prod", "mysql://root@prod-mysql:3306/data", SensitivityCritical},
		{"mongodb prod", "mongodb://prod-cluster.example.com/appdb", SensitivityCritical},
		{"redis prod", "redis://prod-cache:6379", SensitivityCritical},

		// Critical: sensitive system files
		{"etc shadow", "file:///etc/shadow", SensitivityCritical},
		{"etc passwd", "file:///etc/passwd", SensitivityCritical},
		{"root dir", "file:///root/.bashrc", SensitivityCritical},

		// Critical: credential files
		{"dotenv", "file:///app/.env", SensitivityCritical},
		{"private key", "file:///home/user/.ssh/id_rsa.key", SensitivityCritical},
		{"pem cert", "file:///etc/ssl/server.pem", SensitivityCritical},
		{"crt file", "file:///etc/ssl/ca.crt", SensitivityCritical},
		{"p12 keystore", "file:///keys/client.p12", SensitivityCritical},
		{"pfx cert", "file:///certs/server.pfx", SensitivityCritical},
		{"java keystore", "file:///opt/app/keystore.jks", SensitivityCritical},

		// Critical: sensitive paths
		{"credentials dir", "file:///home/user/.aws/credentials", SensitivityCritical},
		{"secrets dir", "file:///app/config/secrets/db.conf", SensitivityCritical},
		{"ssh dir", "file:///home/user/.ssh/config", SensitivityCritical},

		// Critical: prod cloud storage
		{"s3 prod", "s3://prod-data-bucket/exports/users.csv", SensitivityCritical},
		{"gs prod", "gs://production-logs/2024/audit.log", SensitivityCritical},

		// High: non-prod databases
		{"postgres dev", "postgres://user:pass@dev-db:5432/mydb", SensitivityHigh},
		{"mysql staging", "mysql://root@staging-mysql:3306/data", SensitivityHigh},
		{"mongodb test", "mongodb://test-cluster/appdb", SensitivityHigh},
		{"redis dev", "redis://dev-cache:6379", SensitivityHigh},

		// High: log files
		{"var log", "file:///var/log/auth.log", SensitivityHigh},
		{"var log syslog", "file:///var/log/syslog", SensitivityHigh},

		// High: config files with secrets in path
		{"yaml with secret", "file:///app/config/secret-db.yaml", SensitivityHigh},
		{"json with password", "file:///etc/password-store.json", SensitivityHigh},
		{"conf with secret", "file:///app/secret.conf", SensitivityHigh},
		{"yml with password", "file:///config/password-config.yml", SensitivityHigh},
		{"ini with secret", "file:///app/secret-config.ini", SensitivityHigh},
		{"cfg with password", "file:///etc/password.cfg", SensitivityHigh},

		// Medium: general file access
		{"general file", "file:///tmp/data.txt", SensitivityMedium},
		{"file with path", "file:///home/user/documents/report.pdf", SensitivityMedium},

		// Medium: HTTP/HTTPS
		{"http api", "http://api.example.com/v1/data", SensitivityMedium},
		{"https api", "https://api.example.com/v1/users", SensitivityMedium},

		// Medium: non-prod cloud storage
		{"s3 dev", "s3://dev-bucket/test-data.csv", SensitivityMedium},
		{"gs staging", "gs://staging-assets/images/logo.png", SensitivityMedium},

		// Low: everything else
		{"data uri", "data:text/plain;base64,SGVsbG8=", SensitivityLow},
		{"custom scheme", "myprotocol://resource/123", SensitivityLow},
		{"empty string", "", SensitivityLow},
		{"bare path", "/tmp/file.txt", SensitivityLow},
		{"ftp", "ftp://files.example.com/pub/data.zip", SensitivityLow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyResourceSensitivity(tt.uri)
			if got != tt.want {
				t.Errorf("ClassifyResourceSensitivity(%q) = %q, want %q", tt.uri, got, tt.want)
			}
		})
	}
}

func TestSensitivityCaseInsensitive(t *testing.T) {
	got := ClassifyResourceSensitivity("POSTGRES://USER@PROD-DB:5432/DB")
	if got != SensitivityCritical {
		t.Errorf("expected critical for uppercase postgres prod URI, got %q", got)
	}

	got = ClassifyResourceSensitivity("FILE:///ETC/SHADOW")
	if got != SensitivityCritical {
		t.Errorf("expected critical for uppercase /etc/shadow URI, got %q", got)
	}
}
