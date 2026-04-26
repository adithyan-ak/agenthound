package common

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
