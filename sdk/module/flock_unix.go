//go:build unix

package module

import (
	"io"
	"os"

	"golang.org/x/sys/unix"
)

// lockFile acquires an exclusive advisory lock on a sibling .lock file.
// The returned io.Closer releases the lock and removes the lockfile.
// Blocks until the lock is acquired (no timeout — the critical section
// inside WriteReceipt is sub-millisecond, so contention resolves fast).
func lockFile(path string) (io.Closer, error) {
	lockPath := path + ".lock"
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	if err := unix.Flock(int(f.Fd()), unix.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, err
	}
	return &fileLock{f: f, path: lockPath}, nil
}

type fileLock struct {
	f    *os.File
	path string
}

func (l *fileLock) Close() error {
	_ = unix.Flock(int(l.f.Fd()), unix.LOCK_UN)
	_ = l.f.Close()
	_ = os.Remove(l.path)
	return nil
}
