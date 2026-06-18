package cli

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"
)

func TestNeo4jHostPort(t *testing.T) {
	cases := []struct {
		name     string
		uri      string
		wantHost string
		wantPort string
		wantErr  bool
	}{
		{"explicit port", "bolt://localhost:7687", "localhost", "7687", false},
		{"default port", "bolt://localhost", "localhost", "7687", false},
		{"ipv6 explicit port", "bolt://[::1]:7687", "::1", "7687", false},
		{"with creds", "bolt://user:pass@db:7687", "db", "7687", false},
		{"empty uri", "", "", "", true},
		{"no host", "bolt://", "", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			host, port, err := neo4jHostPort(c.uri)
			if (err != nil) != c.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, c.wantErr)
			}
			if host != c.wantHost {
				t.Errorf("host = %q, want %q", host, c.wantHost)
			}
			if port != c.wantPort {
				t.Errorf("port = %q, want %q", port, c.wantPort)
			}
		})
	}
}

func TestPostgresHostPort(t *testing.T) {
	cases := []struct {
		name     string
		uri      string
		wantHost string
		wantPort string
		wantErr  bool
	}{
		{"full uri", "postgres://u:p@db:5432/agenthound?sslmode=disable", "db", "5432", false},
		{"default port", "postgres://u:p@db/agenthound", "db", "5432", false},
		{"empty uri", "", "", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			host, port, err := postgresHostPort(c.uri)
			if (err != nil) != c.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, c.wantErr)
			}
			if host != c.wantHost {
				t.Errorf("host = %q, want %q", host, c.wantHost)
			}
			if port != c.wantPort {
				t.Errorf("port = %q, want %q", port, c.wantPort)
			}
		})
	}
}

func TestRedactURI(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		want   string
		denied []string // substrings that MUST NOT appear in output
	}{
		{
			name:   "no creds untouched",
			input:  "bolt://localhost:7687",
			want:   "bolt://localhost:7687",
			denied: nil,
		},
		{
			name:   "user+password redacted",
			input:  "postgres://alice:s3cret@db:5432/agenthound?sslmode=disable",
			want:   "postgres://***:***@db:5432/agenthound?sslmode=disable",
			denied: []string{"alice", "s3cret"},
		},
		{
			name:   "user only redacted",
			input:  "bolt://neo4j@db:7687",
			want:   "bolt://***:***@db:7687",
			denied: []string{"neo4j@"},
		},
		{
			name:   "query-string password redacted",
			input:  "postgres://host:5432/db?user=alice&password=hunter2&sslmode=disable",
			want:   "postgres://host:5432/db?password=%2A%2A%2A&sslmode=disable&user=alice",
			denied: []string{"hunter2"},
		},
		{
			name:   "query-string password+userinfo both redacted",
			input:  "postgres://bob:secret@host:5432/db?password=alsosecret&sslmode=disable",
			want:   "postgres://***:***@host:5432/db?password=%2A%2A%2A&sslmode=disable",
			denied: []string{"bob", "secret", "alsosecret"},
		},
		{
			name:   "sslpassword redacted",
			input:  "postgres://host:5432/db?sslpassword=keypw&sslmode=verify-full",
			want:   "postgres://host:5432/db?sslmode=verify-full&sslpassword=%2A%2A%2A",
			denied: []string{"keypw"},
		},
		{
			name:   "garbage left intact",
			input:  "::::not a url::::",
			want:   "::::not a url::::",
			denied: nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := redactURI(c.input)
			if got != c.want {
				t.Fatalf("redactURI = %q, want %q", got, c.want)
			}
			for _, deny := range c.denied {
				if strings.Contains(got, deny) {
					t.Errorf("redactURI(%q) leaked %q in %q", c.input, deny, got)
				}
			}
		})
	}
}

func TestPreflightErrorFormatting(t *testing.T) {
	t.Run("default unreachable header", func(t *testing.T) {
		e := &preflightError{
			Service:  "Neo4j",
			Endpoint: "bolt://localhost:7687",
			Hints:    []string{"hint one", "hint two"},
		}
		got := e.Error()
		for _, want := range []string{
			"Neo4j unreachable at bolt://localhost:7687.",
			"hint one",
			"hint two",
			docsInstallURL,
		} {
			if !strings.Contains(got, want) {
				t.Errorf("Error() output missing %q.\ngot:\n%s", want, got)
			}
		}
	})
	t.Run("custom header for driver-rejected path", func(t *testing.T) {
		e := &preflightError{
			Service:  "Neo4j",
			Endpoint: "bolt://localhost:7687",
			Header:   "Neo4j rejected the connection at bolt://localhost:7687",
			Hints:    []string{"hint"},
		}
		got := e.Error()
		if strings.Contains(got, "unreachable") {
			t.Fatalf("custom-header error must not say 'unreachable'; got:\n%s", got)
		}
		if !strings.Contains(got, "rejected the connection") {
			t.Errorf("expected 'rejected the connection' header; got:\n%s", got)
		}
	})
}

func TestClassifyDriverErrorRedactsAndPreservesMessage(t *testing.T) {
	driverErr := errors.New("Neo4jError: Neo.ClientError.Security.Unauthorized (auth failed)")
	got := classifyDriverError("Neo4j", "bolt://alice:s3cret@db:7687", driverErr)
	msg := got.Error()
	if strings.Contains(msg, "alice") || strings.Contains(msg, "s3cret") {
		t.Fatalf("classifyDriverError leaked credentials: %s", msg)
	}
	if !strings.Contains(msg, "Underlying error:") {
		t.Errorf("classifyDriverError should preserve underlying message; got: %s", msg)
	}
	if !strings.Contains(msg, driverErr.Error()) {
		t.Errorf("classifyDriverError should include the original driver error; got: %s", msg)
	}
	if strings.Contains(msg, "unreachable") {
		t.Errorf("classifyDriverError must not say 'unreachable' (TCP probe already succeeded); got:\n%s", msg)
	}
	if !strings.Contains(msg, "rejected the connection") {
		t.Errorf("classifyDriverError should say 'rejected the connection'; got:\n%s", msg)
	}
}

func TestHostPortErrorsDoNotLeakURI(t *testing.T) {
	// url.Parse error messages include the raw URL ("parse \"...\": ..."),
	// which leaks any credentials the operator embedded. The host/port
	// helpers and their callers must never propagate that.
	leaky := "bolt://alice:s3cret@host:abc/path" // bad port → url.Parse fails
	if _, _, err := neo4jHostPort(leaky); err == nil {
		t.Fatal("expected error for invalid port")
	} else if msg := err.Error(); strings.Contains(msg, "alice") || strings.Contains(msg, "s3cret") {
		t.Fatalf("neo4jHostPort error leaked credentials: %s", msg)
	}

	leakyPG := "postgres://bob:hunter2@host:zzz/db"
	if _, _, err := postgresHostPort(leakyPG); err == nil {
		t.Fatal("expected error for invalid port")
	} else if msg := err.Error(); strings.Contains(msg, "bob") || strings.Contains(msg, "hunter2") {
		t.Fatalf("postgresHostPort error leaked credentials: %s", msg)
	}
}

func TestClassifyDriverErrorPostgresUsesPostgresDefaults(t *testing.T) {
	got := classifyDriverError("PostgreSQL", "postgres://x:y@db:5432/agenthound", errors.New("boom"))
	msg := got.Error()
	// docker-compose.yml line 26-27: POSTGRES_USER/PASSWORD = agenthound.
	// The hint must reflect THAT, not Neo4j's default.
	if !strings.Contains(msg, "user=agenthound password=agenthound") {
		t.Errorf("PostgreSQL hint should reference agenthound/agenthound defaults; got:\n%s", msg)
	}
	if strings.Contains(msg, "user=neo4j") {
		t.Errorf("PostgreSQL hint must not reference neo4j credentials; got:\n%s", msg)
	}
}

func TestClassifyDriverErrorNilPassthrough(t *testing.T) {
	if got := classifyDriverError("Neo4j", "bolt://x:7687", nil); got != nil {
		t.Errorf("classifyDriverError(nil) = %v, want nil", got)
	}
}

func TestReachableUnusedPort(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Port 1 is reserved (tcpmux) and effectively never listening on a
	// dev box; the dial fails fast and reachable() must return false.
	if reachable(ctx, "127.0.0.1", "1") {
		t.Skip("127.0.0.1:1 unexpectedly reachable on this host; cannot exercise the false branch")
	}
}

func TestReachableHappyPath(t *testing.T) {
	// Stand up an ephemeral listener and probe it.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("SplitHostPort: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if !reachable(ctx, "127.0.0.1", port) {
		t.Fatalf("expected reachable on 127.0.0.1:%s", port)
	}
}
