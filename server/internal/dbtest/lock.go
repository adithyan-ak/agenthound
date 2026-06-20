// Package dbtest provides cross-process serialization for integration tests
// that share a single Neo4j instance.
//
// The post-processors run whole-graph Cypher with no scan_id scoping (by
// production design: e.g. risk_score lists every node of a kind, shadows
// matches across all servers). `go test ./...` runs each package's test binary
// in parallel, so two DB-touching packages mutating the same Neo4j at once race
// — a DETACH DELETE in one binary can vanish a node mid-traversal in another,
// surfacing as Neo.ClientError.Statement.EntityNotFound or phantom zero-count
// assertions. Each such package's TestMain holds this exclusive advisory lock
// for the duration of its run, so no two integration binaries touch the DB
// concurrently. Within a package, tests already run sequentially.
package dbtest

import (
	"os"
	"path/filepath"
	"syscall"
)

// lockPath is the well-known advisory-lock file shared by every integration
// test binary in this module.
func lockPath() string {
	return filepath.Join(os.TempDir(), "agenthound-neo4j-itest.lock")
}

// Lock acquires an exclusive, blocking flock on the shared lock file and
// returns a release function. It is a no-op (returning a no-op release) when
// AGENTHOUND_NEO4J_URI is unset, so unit-only runs (`go test -short`, no DB)
// keep their full package-level parallelism.
func Lock() (func(), error) {
	if os.Getenv("AGENTHOUND_NEO4J_URI") == "" {
		return func() {}, nil
	}
	f, err := os.OpenFile(lockPath(), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return func() {}, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return func() {}, err
	}
	return func() {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
	}, nil
}
