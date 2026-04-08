package cli

import (
	"context"
	"log/slog"

	collector "github.com/adithyan-ak/agenthound/internal/collector"
	configcollector "github.com/adithyan-ak/agenthound/internal/collector/config"
	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect data from MCP servers, A2A agents, or config files",
}

var collectConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Parse MCP client configuration files",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		discover, _ := cmd.Flags().GetBool("discover")
		path, _ := cmd.Flags().GetString("path")
		paths, _ := cmd.Flags().GetStringSlice("paths")
		output, _ := cmd.Flags().GetString("output")
		doIngest, _ := cmd.Flags().GetBool("ingest")
		includeCredValues, _ := cmd.Flags().GetBool("include-credential-values")
		projectDir, _ := cmd.Flags().GetString("project-dir")

		if !discover && path == "" && len(paths) == 0 {
			return cmd.Help()
		}

		c := configcollector.NewConfigCollector()
		opts := collector.CollectOptions{
			Discover:                discover,
			ConfigPath:              path,
			ConfigPaths:             paths,
			ProjectDir:              projectDir,
			IncludeCredentialValues: includeCredValues,
		}

		slog.Info("collecting config data", "discover", discover, "path", path)
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
	collectConfigCmd.Flags().Bool("discover", false, "Discover all MCP client config files")
	collectConfigCmd.Flags().String("path", "", "Path to specific config file")
	collectConfigCmd.Flags().StringSlice("paths", nil, "Paths to multiple config files")
	collectConfigCmd.Flags().String("output", "", "Write JSON output to file (default: stdout)")
	collectConfigCmd.Flags().Bool("ingest", false, "Ingest directly into graph database")
	collectConfigCmd.Flags().Bool("include-credential-values", false, "Include raw credential values (default: SHA-256 hash)")
	collectConfigCmd.Flags().String("project-dir", "", "Project directory for instruction file discovery")
	collectCmd.AddCommand(collectConfigCmd)
	rootCmd.AddCommand(collectCmd)
}
