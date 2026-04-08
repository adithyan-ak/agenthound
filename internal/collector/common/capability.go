package common

import (
	"sort"
	"strings"
)

type Capability string

const (
	CapShellAccess      Capability = "shell_access"
	CapFileRead         Capability = "file_read"
	CapFileWrite        Capability = "file_write"
	CapNetworkOutbound  Capability = "network_outbound"
	CapDatabaseAccess   Capability = "database_access"
	CapEmailSend        Capability = "email_send"
	CapCodeExecution    Capability = "code_execution"
	CapCredentialAccess Capability = "credential_access"
)

var capabilityKeywords = []struct {
	cap      Capability
	keywords []string
}{
	{CapShellAccess, []string{"shell", "bash", "terminal", "command", "exec", "subprocess", "sh -c", "powershell"}},
	{CapFileRead, []string{"read file", "file_read", "get_file", "list_directory", "readdir", "cat ", "file://", "open file", "read_file"}},
	{CapFileWrite, []string{"write file", "file_write", "save_file", "create_file", "write_file", "overwrite", "append_file", "put_file"}},
	{CapNetworkOutbound, []string{"http", "fetch", "curl", "request", "webhook", "api call", "post ", "get ", "download", "upload", "send_request"}},
	{CapDatabaseAccess, []string{"sql", "query", "database", "postgres", "mysql", "mongodb", "redis", "execute_sql", "db_", "collection"}},
	{CapEmailSend, []string{"email", "smtp", "send_mail", "sendmail", "mail_send", "send_email", "notify"}},
	{CapCodeExecution, []string{"eval", "execute", "run_code", "python", "javascript", "compile", "interpret", "lambda", "exec_code"}},
	{CapCredentialAccess, []string{"password", "credential", "secret", "api_key", "token", "ssh_key", "private_key", "auth", "vault"}},
}

func ClassifyCapabilities(toolName, description string, inputSchema map[string]any) []string {
	combined := strings.ToLower(toolName + " " + description)

	if inputSchema != nil {
		if props, ok := inputSchema["properties"].(map[string]any); ok {
			for key := range props {
				combined += " " + strings.ToLower(key)
			}
		}
	}

	seen := make(map[string]bool)
	var caps []string

	for _, entry := range capabilityKeywords {
		if seen[string(entry.cap)] {
			continue
		}
		for _, kw := range entry.keywords {
			if strings.Contains(combined, kw) {
				caps = append(caps, string(entry.cap))
				seen[string(entry.cap)] = true
				break
			}
		}
	}

	sort.Strings(caps)
	return caps
}
