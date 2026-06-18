package cli

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/adithyan-ak/agenthound/server/internal/servercfg"
)

const (
	preflightDialTimeout = 1500 * time.Millisecond
	docsInstallURL       = "https://docs.agenthound.io/getting-started/install/"
)

// runtimePreflight performs cheap TCP probes against Neo4j and Postgres
// before the real driver init in Bootstrap. The goal is to translate
// generic driver errors ("dial tcp [::1]:7687: connect: connection
// refused") into a friendly, actionable diagnostic that tells a
// first-time user exactly what to do.
//
// The probe is cheap (~ms on the happy path); the full driver init
// still happens afterward and remains the source of truth for auth /
// schema errors.
func runtimePreflight(ctx context.Context, cfg *servercfg.Config) error {
	if err := probeNeo4j(ctx, cfg); err != nil {
		return err
	}
	if err := probePostgres(ctx, cfg); err != nil {
		return err
	}
	return nil
}

func probeNeo4j(ctx context.Context, cfg *servercfg.Config) error {
	host, port, err := neo4jHostPort(cfg.Neo4jURI)
	if err != nil {
		// `url.Parse` errors include the raw URL ("parse \"bolt://user:secret@..\":
		// ..."), which would leak embedded credentials. Return a sanitized
		// shape: don't surface the underlying url error string at all.
		return fmt.Errorf("invalid AGENTHOUND_NEO4J_URI: %s (set a valid URI like bolt://host:7687)", err)
	}
	if reachable(ctx, host, port) {
		return nil
	}
	return &preflightError{
		Service:  "Neo4j",
		Endpoint: redactURI(cfg.Neo4jURI),
		Hints: []string{
			"Most common cause: the database stack isn't running.",
			"  Start it with:    docker compose -f docker/docker-compose.yml up -d graph-db app-db",
			"",
			"  Override URI:     export AGENTHOUND_NEO4J_URI=bolt://your-host:7687",
		},
	}
}

func probePostgres(ctx context.Context, cfg *servercfg.Config) error {
	host, port, err := postgresHostPort(cfg.PostgresURI)
	if err != nil {
		return fmt.Errorf("invalid AGENTHOUND_PG_URI: %s (set a valid URI like postgres://user:pass@host:5432/db?sslmode=disable)", err)
	}
	if reachable(ctx, host, port) {
		return nil
	}
	return &preflightError{
		Service:  "PostgreSQL",
		Endpoint: fmt.Sprintf("postgres://%s:%s", host, port),
		Hints: []string{
			"Most common cause: the database stack isn't running.",
			"  Start it with:    docker compose -f docker/docker-compose.yml up -d graph-db app-db",
			"",
			"  Override URI:     export AGENTHOUND_PG_URI=postgres://user:pass@host:5432/db?sslmode=disable",
		},
	}
}

// redactURI replaces credentials in the URI with ***. Covers both the
// userinfo authority (user:password@host) AND libpq-style query-string
// credentials (?user=...&password=...) — pgx accepts both shapes via
// pgxpool.New, so an operator-supplied AGENTHOUND_PG_URI could embed
// credentials in either place. Returns the input unchanged if it doesn't
// parse — better to show *something* than to drop the operator's URI
// entirely.
//
// We intentionally rebuild the string by hand instead of using
// url.UserPassword: the latter percent-encodes the literal "***" into
// %2A%2A%2A which is unreadable in operator-facing diagnostics.
func redactURI(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	hasUser := u.User != nil
	// Redact known credential-bearing query parameters (libpq / pgconn).
	// Anything we don't recognize is preserved verbatim — sslmode, host,
	// dbname, port etc. are not secrets.
	rawQuery := u.RawQuery
	queryRedacted := false
	if rawQuery != "" {
		q := u.Query()
		for _, key := range []string{"password", "passfile", "sslpassword"} {
			if q.Has(key) {
				q.Set(key, "***")
				queryRedacted = true
			}
		}
		if queryRedacted {
			rawQuery = q.Encode()
		}
	}
	if !hasUser && !queryRedacted {
		return raw
	}
	host := u.Host
	tail := ""
	if u.Path != "" {
		tail += u.Path
	}
	if rawQuery != "" {
		tail += "?" + rawQuery
	}
	if u.Fragment != "" {
		tail += "#" + u.Fragment
	}
	authority := host
	if hasUser {
		authority = "***:***@" + host
	}
	return fmt.Sprintf("%s://%s%s", u.Scheme, authority, tail)
}

// neo4jHostPort parses a bolt:// URI into host/port. Defaults to 7687
// if the port is omitted.
//
// Errors are *sanitized* — `url.Parse` includes the raw URL in its error
// message, which would leak credentials embedded by the operator. We
// return only error classes ("malformed URI", "missing host"), never
// the underlying url.Error.
func neo4jHostPort(rawURI string) (string, string, error) {
	u, err := url.Parse(rawURI)
	if err != nil {
		return "", "", fmt.Errorf("malformed URI")
	}
	host := u.Hostname()
	port := u.Port()
	if host == "" {
		return "", "", fmt.Errorf("missing host")
	}
	if port == "" {
		port = "7687"
	}
	return host, port, nil
}

// postgresHostPort parses a postgres:// URI into host/port. Defaults
// to 5432. See neo4jHostPort for the rationale on sanitized errors.
func postgresHostPort(rawURI string) (string, string, error) {
	u, err := url.Parse(rawURI)
	if err != nil {
		return "", "", fmt.Errorf("malformed URI")
	}
	host := u.Hostname()
	port := u.Port()
	if host == "" {
		return "", "", fmt.Errorf("missing host")
	}
	if port == "" {
		port = "5432"
	}
	return host, port, nil
}

// reachable opens a short-timeout TCP connection. Returns true if the
// dial succeeded; the caller decides whether to escalate to a real
// driver handshake.
func reachable(ctx context.Context, host, port string) bool {
	dialer := net.Dialer{Timeout: preflightDialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(host, port))
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// preflightError carries a structured diagnostic that prints as a
// formatted block instead of a one-line driver error.
//
// Header lets callers pick the right phrasing: "X unreachable at Y" for
// the TCP-probe-failed path, "X rejected the connection at Y" for the
// driver-failed path. Mixing them up turns into a confusing diagnostic
// (e.g. printing "unreachable" right after a successful TCP probe).
type preflightError struct {
	Service  string
	Endpoint string
	Header   string // optional override; defaults to "<Service> unreachable at <Endpoint>"
	Hints    []string
}

func (e *preflightError) Error() string {
	var b strings.Builder
	header := e.Header
	if header == "" {
		header = fmt.Sprintf("%s unreachable at %s", e.Service, e.Endpoint)
	}
	fmt.Fprintf(&b, "%s.\n\n", header)
	for _, line := range e.Hints {
		fmt.Fprintf(&b, "  %s\n", line)
	}
	fmt.Fprintf(&b, "\n  See %s for full setup.", docsInstallURL)
	return b.String()
}

// classifyDriverError takes the generic error returned by graph.NewDriver
// or appdb.NewPool *after* the TCP probe succeeded. The TCP probe alone
// already proves the port is open, so the only thing we add here is a
// human-readable wrapper around the driver's own message — we deliberately
// do NOT try to confidently classify "auth failure" vs "schema failure"
// vs "handshake failure" by string-matching, because that's how false
// positives end up in operator-facing diagnostics. The driver's original
// error is preserved verbatim for the operator to read.
func classifyDriverError(service, endpoint string, err error) error {
	if err == nil {
		return nil
	}
	credHint := "AGENTHOUND_NEO4J_USER / AGENTHOUND_NEO4J_PASSWORD"
	defaultsHint := "Defaults from docker compose: user=neo4j password=agenthound"
	if service == "PostgreSQL" {
		credHint = "the credentials embedded in AGENTHOUND_PG_URI"
		defaultsHint = "Defaults from docker compose: user=agenthound password=agenthound database=agenthound"
	}
	redacted := redactURI(endpoint)
	return &preflightError{
		Service:  service,
		Endpoint: redacted,
		// Backend is reachable here — the TCP probe already passed —
		// so DON'T say "unreachable". Use a header that matches
		// reality.
		Header: fmt.Sprintf("%s rejected the connection at %s", service, redacted),
		Hints: []string{
			"The port is open but the driver could not complete the handshake.",
			fmt.Sprintf("  Most often this is a credential mismatch — check %s.", credHint),
			"  " + defaultsHint,
			"",
			fmt.Sprintf("  Underlying error: %s", err.Error()),
		},
	}
}
