package locking_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/githubnext/apm/internal/cache/locking"
)

func TestNewShardLock_LockFilePath_Adjacent(t *testing.T) {
	sl := locking.NewShardLock("/tmp/testdir/shard.bin", 5*time.Second)
	if sl == nil {
		t.Fatal("NewShardLock should not return nil")
	}
}

func TestStagePath_ContainsIncomplete(t *testing.T) {
	sp := locking.StagePath("/some/dir/final")
	if !strings.Contains(sp, ".incomplete.") {
		t.Errorf("StagePath should contain .incomplete. got %q", sp)
	}
}

func TestStagePath_SameDir(t *testing.T) {
	sp := locking.StagePath("/some/dir/final")
	if filepath.Dir(sp) != "/some/dir" {
		t.Errorf("StagePath should be adjacent to final, got dir %q", filepath.Dir(sp))
	}
}

func TestStagePath_HasBasePrefix(t *testing.T) {
	sp := locking.StagePath("/some/dir/final")
	base := filepath.Base(sp)
	if !strings.HasPrefix(base, "final.incomplete.") {
		t.Errorf("staged base should start with 'final.incomplete.' got %q", base)
	}
}

func TestSafeRemoveAll_Directory(t *testing.T) {
	tmp := t.TempDir()
	subdir := filepath.Join(tmp, "toremove")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := locking.SafeRemoveAll(subdir); err != nil {
		t.Fatalf("SafeRemoveAll failed: %v", err)
	}
	if _, err := os.Stat(subdir); !os.IsNotExist(err) {
		t.Error("directory should be removed")
	}
}

func TestSafeRemoveAll_NonexistentNoError(t *testing.T) {
	if err := locking.SafeRemoveAll("/nonexistent/path/xyz"); err != nil {
		t.Errorf("SafeRemoveAll nonexistent should return nil, got %v", err)
	}
}

func TestCleanupIncomplete_NoDir(t *testing.T) {
	n := locking.CleanupIncomplete("/nonexistent/path")
	if n != 0 {
		t.Errorf("expected 0 removals for nonexistent dir, got %d", n)
	}
}

func TestCleanupIncomplete_EmptyDir(t *testing.T) {
	tmp := t.TempDir()
	n := locking.CleanupIncomplete(tmp)
	if n != 0 {
		t.Errorf("expected 0 removals for empty dir, got %d", n)
	}
}

func TestCleanupIncomplete_RemovesIncompleteEntries(t *testing.T) {
	tmp := t.TempDir()
	// Create some .incomplete. dirs and a normal dir.
	for _, d := range []string{"a.incomplete.123.456", "b.incomplete.789.101"} {
		if err := os.MkdirAll(filepath.Join(tmp, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(tmp, "normal"), 0o755); err != nil {
		t.Fatal(err)
	}
	n := locking.CleanupIncomplete(tmp)
	if n != 2 {
		t.Errorf("expected 2 removals, got %d", n)
	}
	if _, err := os.Stat(filepath.Join(tmp, "normal")); err != nil {
		t.Error("normal dir should not be removed")
	}
}

func TestAtomicLand_SuccessMovesStagedToFinal(t *testing.T) {
	tmp := t.TempDir()
	final := filepath.Join(tmp, "final_dir")
	staged := locking.StagePath(final)
	if err := os.MkdirAll(staged, 0o755); err != nil {
		t.Fatal(err)
	}
	sl := locking.NewShardLock(tmp, 5*time.Second)
	ok, err := locking.AtomicLand(staged, final, sl)
	if err != nil {
		t.Fatalf("AtomicLand error: %v", err)
	}
	if !ok {
		t.Error("AtomicLand should return true on success")
	}
	if _, err := os.Stat(final); err != nil {
		t.Error("final should exist after AtomicLand")
	}
	if _, err := os.Stat(staged); !os.IsNotExist(err) {
		t.Error("staged should be removed after AtomicLand")
	}
}

func TestAtomicLand_SkipsIfFinalExists(t *testing.T) {
	tmp := t.TempDir()
	final := filepath.Join(tmp, "existing_dir")
	if err := os.MkdirAll(final, 0o755); err != nil {
		t.Fatal(err)
	}
	staged := locking.StagePath(final)
	if err := os.MkdirAll(staged, 0o755); err != nil {
		t.Fatal(err)
	}
	sl := locking.NewShardLock(tmp, 5*time.Second)
	ok, err := locking.AtomicLand(staged, final, sl)
	if err != nil {
		t.Fatalf("AtomicLand error: %v", err)
	}
	if ok {
		t.Error("AtomicLand should return false when final already exists")
	}
}

func TestShardLock_LockUnlock_Cycle(t *testing.T) {
	tmp := t.TempDir()
	sl := locking.NewShardLock(filepath.Join(tmp, "shard"), 5*time.Second)
	if err := sl.Lock(); err != nil {
		t.Fatalf("Lock failed: %v", err)
	}
	sl.Unlock()
	// Second lock after unlock should succeed.
	if err := sl.Lock(); err != nil {
		t.Fatalf("second Lock failed: %v", err)
	}
	sl.Unlock()
}
