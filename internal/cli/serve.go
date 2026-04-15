package cli

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adithyan-ak/agenthound/internal/api"
	"github.com/adithyan-ak/agenthound/internal/auth"
	"github.com/adithyan-ak/agenthound/internal/rules"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the AgentHound API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		infra, cleanup, err := Bootstrap(ctx)
		if err != nil {
			return err
		}
		defer cleanup()

		if err := auth.EnsureAdminUser(ctx, infra.UserStore, cfg.AdminPassword); err != nil {
			return err
		}

		rulesEngine, err := rules.NewEngine(rules.LoadOptions{})
		if err != nil {
			slog.Warn("failed to load rules engine, rules API will return empty", "error", err)
		}

		server := api.NewServer(api.ServerDeps{
			GraphDB:     infra.GraphDB,
			Reader:      infra.Reader,
			PGPool:      infra.PGPool,
			Pipeline:    infra.Pipeline,
			ScanStore:   infra.ScanStore,
			UserStore:   infra.UserStore,
			TokenStore:  infra.TokenStore,
			AuditStore:  infra.AuditStore,
			RulesEngine: rulesEngine,
			JWTSecret:   cfg.JWTSecret,
			CORSOrigins: cfg.CORSOrigins,
		})

		errCh := make(chan error, 1)
		go func() {
			errCh <- server.ListenAndServe(cfg.APIPort)
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
