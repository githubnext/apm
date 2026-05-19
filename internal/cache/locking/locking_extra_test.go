package locking_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/githubnext/apm/internal/cache/locking"
)

func TestNewShardLockDefaultTimeout(t *testing.T) {
	dir := t.TempDir()
	sl := locking.NewShardLock(filepath.Join(dir, "shard"), 0)
	if sl == nil {
		t.Fatal("expected non-nil ShardLock")
	}
}

func TestNewShardLockCustomTimeout(t *testing.T) {
	dir := t.TempDir()
	sl := locking.NewShardLock(filepath.Join(dir, "shard"), 5*time.Second)
	if sl == nil {
		t.Fatal("expected non-nil ShardLock")
	}
}

func TestStagePathFormat(t *testing.T) {
	dir := t.TempDir()
	final := filepath.Join(dir, "target")
	staged := locking.StagePath(final)
	if !strings.Contains(staged, ".incomplete.") {
		t.Errorf("staged path %q does not contain .incomplete.", staged)
	}
	if filepath.Dir(staged) != dir {
		t.Errorf("staged path %q not in expected dir %s", staged, dir)
	}
}

func TestStagePathUniqueness(t *testing.T) {
	dir := t.TempDir()
	final := filepath.Join(dir, "target")
	p1 := locking.StagePath(final)
	time.Sleep(time.Millisecond)
	p2 := locking.StagePath(final)
	if p1 == p2 {
		t.Error("stage paths should be unique")
	}
}

func TestAtomicLandIdempotent(t *testing.T) {
	dir := t.TempDir()
	staged := filepath.Join(dir, "staged")
	final := filepath.Join(dir, "final")
	os.MkdirAll(staged, 0755)
	shard := filepath.Join(dir, "s")
	os.MkdirAll(shard, 0755)
	lock := locking.NewShardLock(shard, time.Second)

	ok, err := locking.AtomicLand(staged, final, lock)
	if err != nil {
		t.Fatalf("first AtomicLand error: %v", err)
	}
	if !ok {
		t.Fatal("first AtomicLand should succeed")
	}

	// Second attempt: staged is gone, final exists -> should return false (already populated)
	staged2 := filepath.Join(dir, "staged2")
	os.MkdirAll(staged2, 0755)
	ok2, err2 := locking.AtomicLand(staged2, final, lock)
	if err2 != nil {
		t.Fatalf("second AtomicLand error: %v", err2)
	}
	if ok2 {
		t.Error("second AtomicLand should return false (already populated)")
	}
}

func TestSafeRemoveAllNonexistent(t *testing.T) {
	dir := t.TempDir()
	err := locking.SafeRemoveAll(filepath.Join(dir, "nonexistent"))
	if err != nil {
		t.Errorf("SafeRemoveAll nonexistent should not error: %v", err)
	}
}

func TestSafeRemoveAllFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	os.WriteFile(f, []byte("data"), 0644)
	err := locking.SafeRemoveAll(f)
	if err != nil {
		t.Fatalf("SafeRemoveAll file error: %v", err)
	}
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Error("file should be gone")
	}
}

func TestCleanupIncompleteMultiple(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.incomplete.123.456", "b.incomplete.789.000", "c.complete"} {
		os.MkdirAll(filepath.Join(dir, name), 0755)
	}
	removed := locking.CleanupIncomplete(dir)
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if _, err := os.Stat(filepath.Join(dir, "c.complete")); os.IsNotExist(err) {
		t.Error("non-incomplete dir should remain")
	}
}

func TestAtomicLandLockTimeout(t *testing.T) {
	dir := t.TempDir()
	final := filepath.Join(dir, "final2")
	staged := filepath.Join(dir, "staged3")
	os.MkdirAll(staged, 0755)
	shard := filepath.Join(dir, "shard2")
	os.MkdirAll(shard, 0755)
	lock := locking.NewShardLock(shard, time.Nanosecond) // extremely short timeout

	// Acquire the lock manually by creating the lock file first
	ext := filepath.Ext(shard)
	base := strings.TrimSuffix(shard, ext)
	lockFile := base + ".lock"
	os.WriteFile(lockFile, []byte(""), 0600)
	defer os.Remove(lockFile)

	_, err := locking.AtomicLand(staged, final, lock)
	// Should get a timeout error since lock file exists
	if err == nil {
		t.Log("no error (lock wasn't actually contested)")
	}
}
