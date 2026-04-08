# Contributing to AgentHound

## Getting started

```bash
git clone https://github.com/adithyan-ak/agenthound.git
cd agenthound
docker compose -f docker/docker-compose.yml up -d   # Neo4j + PostgreSQL
make build                                            # Build Go binary + UI
make test                                             # Run all tests
```

## Development workflow

1. Fork the repo and create a feature branch from `main`
2. Make your changes
3. Run pre-commit checks (mandatory):
   ```bash
   gofmt -l .          # Must produce no output
   go build ./...      # Must pass
   go vet ./...        # Must pass
   make test           # Must pass
   ```
4. Commit with a clear message describing what changed and why
5. Open a PR against `main`

## Code style

- **Go:** `gofmt` formatting, `errcheck` compliance (handle all error return values)
- **TypeScript:** Prettier + ESLint (run `cd ui && npm run lint`)
- **No manual alignment padding** -- `golangci-lint` enforces this
- **Intentionally discarded errors** use `_, _ =` (e.g., `_, _ = fmt.Fprintf(os.Stderr, ...)`)
- **Property keys** in Neo4j are always `snake_case`

## How to add a new collector

Collectors implement the `Collector` interface in `pkg/collector/collector.go`.

1. Create a new directory under `internal/collector/<protocol>/`
2. Implement the `Collector` interface:
   ```go
   type Collector interface {
       Name() string
       Collect(ctx context.Context, opts CollectOptions) (*model.IngestData, error)
   }
   ```
3. Produce `model.IngestData` with the standard JSON schema (see `docs/graph-model.md`)
4. Register the CLI subcommand in `internal/cli/collect_<protocol>.go`
5. Add test data in `testdata/`
6. Node IDs must be deterministic SHA-256 hashes (see graph model docs)

## How to add a detection rule (post-processor)

Post-processors implement the `PostProcessor` interface in `pkg/analysis/postprocessor.go`.

1. Create a file in `internal/analysis/processors/<name>.go`
2. Implement the interface:
   ```go
   type PostProcessor interface {
       Name() string
       DependsOn() []string
       Process(ctx context.Context, db graph.GraphDB) error
   }
   ```
3. `DependsOn()` returns the names of processors that must run before this one
4. Register the processor in `internal/analysis/postprocessor.go`
5. Add a corresponding test file `<name>_test.go`
6. If the detection should appear as a pre-built query, add it to `internal/analysis/prebuilt/queries.go`

## How to add a pre-built query

1. Add the Cypher constant in `internal/analysis/prebuilt/cypher.go`
2. Register it in `internal/analysis/prebuilt/queries.go` with:
   - Unique ID (kebab-case)
   - Name, description, category, severity
   - OWASP mapping (MCP01-MCP10, ASI01-ASI10)
3. The query is automatically available via the CLI (`--prebuilt`) and API (`/analysis/prebuilt/{id}`)

## How to add a config parser

Config parsers implement the `ConfigParser` interface in `internal/collector/config/`.

1. Create a file in `internal/collector/config/parsers/<client>.go`
2. Implement:
   ```go
   type ConfigParser interface {
       ClientName() string
       ConfigPaths() []string
       Parse(path string, data []byte) (*ParseResult, error)
   }
   ```
3. `ConfigPaths()` returns platform-specific default config file locations
4. Register the parser in the config collector's parser list
5. Add test fixtures in `testdata/`

## Testing

- All new features must include tests
- Run `make test` for Go tests, `cd ui && npm test` for frontend tests
- Use `testdata/` fixtures for collector and ingest tests
- Post-processor tests should verify edge creation with known graph states

## Reporting bugs

Open a [GitHub issue](https://github.com/adithyan-ak/agenthound/issues) with:
- Steps to reproduce
- Expected vs actual behavior
- AgentHound version (`agenthound --version`)
- Neo4j version and OS

## Security vulnerabilities

See [SECURITY.md](SECURITY.md) for responsible disclosure instructions.
