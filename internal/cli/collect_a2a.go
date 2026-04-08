package cli

import (
	"context"
	"log/slog"
	"time"

	a2acollector "github.com/adithyan-ak/agenthound/internal/collector/a2a"
	collector "github.com/adithyan-ak/agenthound/pkg/collector"
	"github.com/spf13/cobra"
)

var collectA2ACmd = &cobra.Command{
	Use:   "a2a",
	Short: "Fetch A2A Agent Cards and collect skill data",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		target, _ := cmd.Flags().GetString("target")
		targets, _ := cmd.Flags().GetStringSlice("targets")
		targetsFile, _ := cmd.Flags().GetString("targets-file")
		output, _ := cmd.Flags().GetString("output")
		doIngest, _ := cmd.Flags().GetBool("ingest")
		authToken, _ := cmd.Flags().GetString("auth-token")
		insecure, _ := cmd.Flags().GetBool("insecure")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		if target == "" && len(targets) == 0 && targetsFile == "" {
			return cmd.Help()
		}

		var a2aOpts []a2acollector.Option
		if concurrency > 0 {
			a2aOpts = append(a2aOpts, a2acollector.WithConcurrency(concurrency))
		}
		if timeout > 0 {
			a2aOpts = append(a2aOpts, a2acollector.WithTimeout(timeout))
		}
		if insecure {
			a2aOpts = append(a2aOpts, a2acollector.WithInsecure(true))
		}

		c := a2acollector.NewA2ACollector(a2aOpts...)
		opts := collector.CollectOptions{
			TargetURL:      target,
			TargetURLs:     targets,
			TargetURLsFile: targetsFile,
			AuthToken:      authToken,
			Insecure:       insecure,
		}

		slog.Info("collecting A2A data", "target", target, "targets", len(targets))
		data, err := c.Collect(ctx, opts)
		if err != nil {
			return err
		}

		slog.Info("collection complete", "nodes", len(data.Graph.Nodes), "edges", len(data.Graph.Edges))

		if doIngest {
			return ingestCollectorOutput(ctx, data)
		}
		return writeCollectorOutput(data, output)
	},
}

func init() {
	collectA2ACmd.Flags().String("target", "", "URL of a single A2A agent")
	collectA2ACmd.Flags().StringSlice("targets", nil, "URLs of multiple A2A agents")
	collectA2ACmd.Flags().String("targets-file", "", "File with A2A agent URLs (one per line)")
	collectA2ACmd.Flags().String("output", "", "Write JSON output to file (default: stdout)")
	collectA2ACmd.Flags().Bool("ingest", false, "Ingest directly into graph database")
	collectA2ACmd.Flags().String("auth-token", "", "Bearer token for authenticated agents")
	collectA2ACmd.Flags().Bool("insecure", false, "Skip TLS verification")
	collectA2ACmd.Flags().Int("concurrency", 5, "Max parallel agent fetches")
	collectA2ACmd.Flags().Duration("timeout", 15*time.Second, "Timeout per agent")
	collectCmd.AddCommand(collectA2ACmd)
}
