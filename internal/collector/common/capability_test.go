package common

import (
	"reflect"
	"testing"
)

func TestClassifyCapabilities(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		description string
		inputSchema map[string]any
		want        []string
	}{
		{
			name:        "shell access by name",
			toolName:    "run_shell_command",
			description: "Runs a command",
			want:        []string{"shell_access"},
		},
		{
			name:        "bash keyword",
			toolName:    "bash",
			description: "Execute bash commands in a terminal",
			want:        []string{"code_execution", "shell_access"},
		},
		{
			name:        "file read",
			toolName:    "read_file",
			description: "Reads a file from disk",
			want:        []string{"file_read"},
		},
		{
			name:        "file write",
			toolName:    "write_file",
			description: "Writes content to a file on disk",
			want:        []string{"file_write"},
		},
		{
			name:        "network outbound via fetch",
			toolName:    "web_fetch",
			description: "Fetch a URL and return its content",
			want:        []string{"network_outbound"},
		},
		{
			name:        "database access",
			toolName:    "execute_sql",
			description: "Run a SQL query against the postgres database",
			want:        []string{"code_execution", "database_access", "shell_access"},
		},
		{
			name:        "email send",
			toolName:    "send_email",
			description: "Sends an email via SMTP",
			want:        []string{"email_send"},
		},
		{
			name:        "code execution",
			toolName:    "run_code",
			description: "Execute arbitrary python code",
			want:        []string{"code_execution", "shell_access"},
		},
		{
			name:        "credential access",
			toolName:    "get_secret",
			description: "Retrieve an API key from the vault",
			want:        []string{"credential_access"},
		},
		{
			name:        "multiple capabilities",
			toolName:    "exec_remote",
			description: "Execute a shell command on a remote host via http and read files",
			want:        []string{"code_execution", "file_read", "network_outbound", "shell_access"},
		},
		{
			name:        "no capabilities",
			toolName:    "format_text",
			description: "Formats the input text nicely",
			want:        nil,
		},
		{
			name:        "input schema property detection",
			toolName:    "custom_tool",
			description: "Does something",
			inputSchema: map[string]any{
				"properties": map[string]any{
					"sql_query":   map[string]any{"type": "string"},
					"database_id": map[string]any{"type": "string"},
				},
			},
			want: []string{"database_access"},
		},
		{
			name:        "input schema with password field",
			toolName:    "login",
			description: "Authenticates a user",
			inputSchema: map[string]any{
				"properties": map[string]any{
					"username": map[string]any{"type": "string"},
					"password": map[string]any{"type": "string"},
				},
			},
			want: []string{"credential_access"},
		},
		{
			name:        "case insensitive matching",
			toolName:    "RunBASH",
			description: "Execute SHELL commands in TERMINAL",
			want:        []string{"code_execution", "shell_access"},
		},
		{
			name:        "file scheme url",
			toolName:    "local_reader",
			description: "Reads from file:///etc/config",
			want:        []string{"file_read"},
		},
		{
			name:        "mongodb access",
			toolName:    "mongo_query",
			description: "Query a mongodb collection",
			want:        []string{"database_access"},
		},
		{
			name:        "redis access",
			toolName:    "cache_get",
			description: "Retrieve value from redis",
			want:        []string{"database_access", "network_outbound"},
		},
		{
			name:        "webhook",
			toolName:    "trigger_webhook",
			description: "Send a webhook notification",
			want:        []string{"network_outbound"},
		},
		{
			name:        "empty inputs",
			toolName:    "",
			description: "",
			want:        nil,
		},
		{
			name:        "nil schema is safe",
			toolName:    "tool",
			description: "description",
			inputSchema: nil,
			want:        nil,
		},
		{
			name:        "schema without properties key",
			toolName:    "tool",
			description: "description",
			inputSchema: map[string]any{"type": "object"},
			want:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyCapabilities(tt.toolName, tt.description, tt.inputSchema)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClassifyCapabilities(%q, %q, ...) = %v, want %v",
					tt.toolName, tt.description, got, tt.want)
			}
		})
	}
}

func TestClassifyCapabilitiesDeduplication(t *testing.T) {
	caps := ClassifyCapabilities("shell", "run shell bash terminal command", nil)
	seen := make(map[string]bool)
	for _, c := range caps {
		if seen[c] {
			t.Errorf("duplicate capability: %s", c)
		}
		seen[c] = true
	}
}

func TestClassifyCapabilitiesSorted(t *testing.T) {
	caps := ClassifyCapabilities("exec_remote", "shell http database eval password", nil)
	for i := 1; i < len(caps); i++ {
		if caps[i] < caps[i-1] {
			t.Errorf("capabilities not sorted: %v", caps)
			break
		}
	}
}
