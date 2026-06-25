//go:build unix

package cli

import (
	"context"
	"errors"
	"syscall"
	"testing"
	"time"
)

// TestSignalContext_CancelsOnSignal is the B1 regression: signalContext must
// register a handler so a delivered SIGTERM cancels the returned context
// (enabling the partial-write path) rather than terminating the process.
//
// This test sends a real signal to the test process, which is process-global,
// so it deliberately does NOT call t.Parallel() and is unix-guarded. The
// handler is active between signalContext() and stop(), so the signal is
// caught here rather than killing the run.
func TestSignalContext_CancelsOnSignal(t *testing.T) {
	ctx, stop := signalContext()
	defer stop()

	if err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatalf("kill self with SIGTERM: %v", err)
	}

	select {
	case <-ctx.Done():
		if !errors.Is(ctx.Err(), context.Canceled) {
			t.Errorf("ctx.Err() = %v, want context.Canceled", ctx.Err())
		}
	case <-time.After(2 * time.Second):
		t.Fatal("context not cancelled within 2s after SIGTERM")
	}
}
