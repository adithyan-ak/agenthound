package cli

import (
	"context"
	"log/slog"
	"time"

	mcpcollector "github.com/adithyan-ak/agenthound/internal/collector/mcp"
	collector "github.com/adithyan-ak/agenthound/pkg/collector"
	"github.com/spf13/cobra"
)

var collectMCPCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Enumerate MCP servers and collect tool/resource data",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		discover, _ := cmd.Flags().GetBool("discover")
		configPath, _ := cmd.Flags().GetString("config")
		url, _ := cmd.Flags().GetString("url")
		output, _ := cmd.Flags().GetString("output")
		doIngest, _ := cmd.Flags().GetBool("ingest")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		insecure, _ := cmd.Flags().GetBool("insecure")

		if !discover && configPath == "" && url == "" {
			return cmd.Help()
		}

		var mcpOpts []mcpcollector.Option
		if concurrency > 0 {
			mcpOpts = append(mcpOpts, mcpcollector.WithConcurrency(concurrency))
		}
		if timeout > 0 {
			mcpOpts = append(mcpOpts, mcpcollector.WithTimeout(timeout))
		}

		c := mcpcollector.NewMCPCollector(mcpOpts...)
		opts := collector.CollectOptions{
			Discover:   discover,
			ConfigPath: configPath,
			TargetURL:  url,
			Insecure:   insecure,
		}

		slog.Info("collecting MCP data", "discover", discover, "config", configPath, "url", url)
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
	collectMCPCmd.Flags().Bool("discover", false, "Discover all configured MCP servers")
	collectMCPCmd.Flags().String("config", "", "Path to MCP client config file")
	collectMCPCmd.Flags().String("url", "", "URL of a single HTTP MCP server")
	collectMCPCmd.Flags().String("output", "", "Write JSON output to file (default: stdout)")
	collectMCPCmd.Flags().Bool("ingest", false, "Ingest directly into graph database")
	collectMCPCmd.Flags().Int("concurrency", 5, "Max parallel server connections")
	collectMCPCmd.Flags().Duration("timeout", 120*time.Second, "Timeout per server")
	collectMCPCmd.Flags().Bool("insecure", false, "Skip TLS verification for HTTP servers")
	collectCmd.AddCommand(collectMCPCmd)
}
