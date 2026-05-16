package locking_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/cache/locking"
)

func TestStagePath(t *testing.T) {
	final := "/tmp/some/path/entry"
	staged := locking.StagePath(final)
	if !strings.Contains(staged, ".incomplete.") {
		t.Errorf("StagePath should contain .incomplete. got %q", staged)
	}
	if filepath.Dir(staged) != filepath.Dir(final) {
		t.Errorf("staged dir %q != final dir %q", filepath.Dir(staged), filepath.Dir(final))
	}
}

func TestShardLockLockUnlock(t *testing.T) {
	dir := t.TempDir()
	shardDir := filepath.Join(dir, "shard")
	lock := locking.NewShardLock(shardDir, 0)
	if err := lock.Lock(); err != nil {
		t.Fatalf("Lock() error: %v", err)
	}
	lock.Unlock()
}

func TestAtomicLand(t *testing.T) {
	dir := t.TempDir()
	staged := filepath.Join(dir, "staged")
	final := filepath.Join(dir, "final")

	if err := os.MkdirAll(staged, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(staged, "file.txt"), []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

	lock := locking.NewShardLock(final, 0)
	ok, err := locking.AtomicLand(staged, final, lock)
	if err != nil {
		t.Fatalf("AtomicLand error: %v", err)
	}
	if !ok {
		t.Error("expected AtomicLand to return true")
	}
	if _, err := os.Stat(final); err != nil {
		t.Errorf("final path should exist: %v", err)
	}
}

func TestAtomicLandDestinationAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	staged := filepath.Join(dir, "staged2")
	final := filepath.Join(dir, "final2")

	if err := os.MkdirAll(staged, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(final, 0o700); err != nil {
		t.Fatal(err)
	}

	lock := locking.NewShardLock(final, 0)
	ok, err := locking.AtomicLand(staged, final, lock)
	if err != nil {
		t.Fatalf("AtomicLand error: %v", err)
	}
	if ok {
		t.Error("expected AtomicLand to return false when destination exists")
	}
}

func TestCleanupIncomplete(t *testing.T) {
	parent := t.TempDir()
	// Create stale incomplete dirs
	_ = os.MkdirAll(filepath.Join(parent, "entry.incomplete.1234.5678"), 0o700)
	_ = os.MkdirAll(filepath.Join(parent, "entry.incomplete.9999.0000"), 0o700)
	// Create a normal dir that should not be removed
	_ = os.MkdirAll(filepath.Join(parent, "normal_entry"), 0o700)

	removed := locking.CleanupIncomplete(parent)
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if _, err := os.Stat(filepath.Join(parent, "normal_entry")); err != nil {
		t.Error("normal_entry should still exist")
	}
}

func TestCleanupIncompleteNonexistentParent(t *testing.T) {
	removed := locking.CleanupIncomplete("/nonexistent/path/that/does/not/exist")
	if removed != 0 {
		t.Errorf("expected 0 removed for nonexistent path, got %d", removed)
	}
}

func TestSafeRemoveAll(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir")
	_ = os.MkdirAll(path, 0o700)
	if err := locking.SafeRemoveAll(path); err != nil {
		t.Errorf("SafeRemoveAll error: %v", err)
	}
	// Calling on nonexistent path should not error
	if err := locking.SafeRemoveAll(path); err != nil {
		t.Errorf("SafeRemoveAll on nonexistent should not error: %v", err)
	}
}
