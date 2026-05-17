//go:build !unix

package module

import (
	"errors"
	"io"
	"os"
	"time"
)

const (
	dirLockTimeout  = 5 * time.Second
	dirLockInterval = 50 * time.Millisecond
)

// lockFile acquires a directory-based lock (mkdir is atomic on all
// platforms). Spins with backoff until the lock is acquired or timeout.
// This is the fallback for non-unix platforms (Windows) where flock(2)
// is unavailable. Not a true advisory lock, but prevents the most
// common collision pattern.
func lockFile(path string) (io.Closer, error) {
	lockPath := path + ".lock"
	deadline := time.Now().Add(dirLockTimeout)
	for {
		err := os.Mkdir(lockPath, 0o700)
		if err == nil {
			return &dirLock{path: lockPath}, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}
		if time.Now().After(deadline) {
			return nil, errors.New("lock timeout: could not acquire " + lockPath)
		}
		time.Sleep(dirLockInterval)
	}
}

type dirLock struct {
	path string
}

func (l *dirLock) Close() error {
	return os.Remove(l.path)
}
