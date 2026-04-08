package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adithyan-ak/agenthound/internal/api"
	"github.com/spf13/cobra"

	"log/slog"
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

		server := api.NewServer(infra.GraphDB, infra.Reader, infra.PGPool, infra.Pipeline, infra.ScanStore)

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
	serveCmd.Flags().Int("port", 8080, "API server port")
	rootCmd.AddCommand(serveCmd)
}
