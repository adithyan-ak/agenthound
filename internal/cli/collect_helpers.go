package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/adithyan-ak/agenthound/internal/appdb"
	"github.com/adithyan-ak/agenthound/internal/graph"
	"github.com/adithyan-ak/agenthound/internal/ingest"
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
	neo4jDriver, err := graph.NewDriver(cfg.Neo4jURI, cfg.Neo4jUser, cfg.Neo4jPassword)
	if err != nil {
		return fmt.Errorf("neo4j: %w", err)
	}
	defer neo4jDriver.Close(ctx)

	pgPool, err := appdb.NewPool(cfg.PostgresURI)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer pgPool.Close()

	if err := graph.InitSchema(ctx, neo4jDriver); err != nil {
		return fmt.Errorf("neo4j schema: %w", err)
	}
	if err := appdb.RunMigrations(ctx, pgPool); err != nil {
		return fmt.Errorf("postgres migrations: %w", err)
	}

	writer := graph.NewWriter(neo4jDriver)
	scanStore := appdb.NewScanStore(pgPool)
	pipeline := ingest.NewPipeline(writer, scanStore)

	result, err := pipeline.Ingest(ctx, data)
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
