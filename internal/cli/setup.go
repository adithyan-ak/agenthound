package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/adithyan-ak/agenthound/internal/apiclient"
	"github.com/adithyan-ak/agenthound/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure CLI to connect to an AgentHound server",
	Long: `Set up the CLI to communicate with an AgentHound server.

This command logs in to the server, creates an API token, and saves
the configuration to ~/.config/agenthound/config.yaml.`,
	RunE: runSetup,
}

func init() {
	setupCmd.Flags().String("username", "", "Admin username (default: prompt)")
	setupCmd.Flags().String("password", "", "Admin password (default: prompt)")
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	serverURL := resolveServerURL(cmd)

	client := apiclient.New(serverURL, "")
	if err := client.Health(ctx); err != nil {
		return fmt.Errorf("cannot reach server at %s: %w\nMake sure the AgentHound server is running", serverURL, err)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Server reachable at %s\n", serverURL)

	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")

	if username == "" {
		username = promptInput("Username", "admin")
	}
	if password == "" {
		password = promptPassword("Password")
	}

	jwt, err := client.Login(ctx, username, password)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	_, _ = fmt.Fprintln(os.Stderr, "Authenticated.")

	tokenClient := apiclient.New(serverURL, jwt)
	hostname, _ := os.Hostname()
	tokenName := fmt.Sprintf("cli-%s-%d", hostname, time.Now().Unix())
	token, err := tokenClient.CreateToken(ctx, tokenName)
	if err != nil {
		return fmt.Errorf("create API token: %w", err)
	}

	if err := config.SaveClientConfig(&config.ClientConfig{
		ServerURL: serverURL,
		APIToken:  token,
	}); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	_, _ = fmt.Fprintf(os.Stderr, "\nSetup complete. Config saved to %s\n", config.ClientConfigPath())
	_, _ = fmt.Fprintf(os.Stderr, "Server: %s\n\n", serverURL)
	_, _ = fmt.Fprintf(os.Stderr, "Run 'agenthound scan' to scan your infrastructure.\n")
	return nil
}

func resolveServerURL(cmd *cobra.Command) string {
	if v, _ := cmd.Root().PersistentFlags().GetString("server-url"); v != "" {
		return v
	}
	if v := os.Getenv("AGENTHOUND_SERVER_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

func promptInput(label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Fprintf(os.Stderr, "%s [%s]: ", label, defaultVal)
	} else {
		fmt.Fprintf(os.Stderr, "%s: ", label)
	}
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		val := strings.TrimSpace(scanner.Text())
		if val != "" {
			return val
		}
	}
	return defaultVal
}

func promptPassword(label string) string {
	fmt.Fprintf(os.Stderr, "%s: ", label)
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		pw, err := term.ReadPassword(fd)
		fmt.Fprintln(os.Stderr)
		if err == nil && len(pw) > 0 {
			return string(pw)
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}
