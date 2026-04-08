package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/adithyan-ak/agenthound/internal/model"
)

func writeCollectorOutput(data *model.IngestData, outputPath string) error {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	if outputPath == "" {
		_, err = os.Stdout.Write(encoded)
		if err != nil {
			return fmt.Errorf("write stdout: %w", err)
		}
		fmt.Println()
		return nil
	}

	if err := os.WriteFile(outputPath, encoded, 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	slog.Info("output written", "path", outputPath, "nodes", len(data.Graph.Nodes), "edges", len(data.Graph.Edges))
	return nil
}

func ingestCollectorOutput(ctx context.Context, data *model.IngestData) error {
	infra, cleanup, err := Bootstrap(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	result, err := infra.Pipeline.Ingest(ctx, data)
	if err != nil {
		return fmt.Errorf("ingest: %w", err)
	}

	fmt.Printf("Ingest complete:\n")
	fmt.Printf("  Scan ID:       %s\n", result.ScanID)
	fmt.Printf("  Nodes written: %d\n", result.NodesWritten)
	fmt.Printf("  Edges written: %d\n", result.EdgesWritten)
	fmt.Printf("  Duration:      %s\n", result.Duration)
	if len(result.Warnings) > 0 {
		fmt.Printf("  Warnings:      %d\n", len(result.Warnings))
		for _, w := range result.Warnings {
			slog.Warn(w)
		}
	}
	return nil
}
