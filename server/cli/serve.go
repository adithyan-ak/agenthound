package cli

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/adithyan-ak/agenthound/sdk/rules"
	"github.com/adithyan-ak/agenthound/server/internal/api"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the AgentHound API server",
	// SilenceUsage / SilenceErrors are set on rootCmd; inherited here.
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		infra, cleanup, err := Bootstrap(ctx)
		if err != nil {
			return err
		}
		defer cleanup()

		rulesEngine, err := rules.NewEngine(rules.LoadOptions{})
		if err != nil {
			slog.Warn("failed to load rules engine, rules API will return empty", "error", err)
		}

		// Warn loudly if the bind address isn't loopback. The server has
		// no application-layer auth; OriginGuard's CSRF protection
		// assumes the only callers reaching loopback are the operator.
		// Binding 0.0.0.0 exposes mutating endpoints to anyone who can
		// spoof an Origin header (trivial for a LAN attacker).
		warnIfNonLoopbackBind(cfg.Bind)

		server := api.NewServer(api.ServerDeps{
			GraphDB:      infra.GraphDB,
			Reader:       infra.Reader,
			PGPool:       infra.PGPool,
			Pipeline:     infra.Pipeline,
			ScanStore:    infra.ScanStore,
			FindingStore: infra.FindingStore,
			RulesEngine:  rulesEngine,
			CORSOrigins:  cfg.CORSOrigins,
		})

		errCh := make(chan error, 1)
		go func() {
			errCh <- server.ListenAndServe(cfg.Bind)
		}()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case err := <-errCh:
			return err
		case sig := <-sigCh:
			slog.Info("shutting down", "signal", sig)
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			return server.Shutdown(shutdownCtx)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// warnIfNonLoopbackBind emits a warning when the configured bind host
// is not in 127.0.0.0/8 (or the literal "localhost"). Loopback binds
// are the only deployment the threat model defends; remote access
// should go through VPN / SSH tunnel / reverse proxy with mTLS.
func warnIfNonLoopbackBind(bind string) {
	host, _, err := net.SplitHostPort(bind)
	if err != nil {
		return // invalid bind is handled elsewhere
	}
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "localhost" || host == "" {
		return
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return
	}
	slog.Warn("non-loopback bind: OriginGuard alone is insufficient against LAN attackers",
		"bind", bind,
		"guidance", "place behind VPN, SSH tunnel, or reverse proxy with mTLS")
}
