package locking_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/githubnext/apm/internal/cache/locking"
)

func TestNewShardLock_NotNil_Extra3(t *testing.T) {
	sl := locking.NewShardLock("/tmp/testdir/shard", 5*time.Second)
	if sl == nil {
		t.Fatal("NewShardLock must not return nil")
	}
}

func TestNewShardLock_ZeroTimeout_UsesDefault_Extra3(t *testing.T) {
	// Zero timeout should not panic; it falls back to the default.
	sl := locking.NewShardLock("/tmp/testdir/shard", 0)
	if sl == nil {
		t.Fatal("NewShardLock with zero timeout must not return nil")
	}
}

func TestStagePath_ContainsFinalName_Extra3(t *testing.T) {
	final := "/tmp/testcache/myentry"
	stage := locking.StagePath(final)
	if !strings.Contains(filepath.Base(stage), "myentry") {
		t.Errorf("stage path should contain base name of final path, got %q", stage)
	}
}

func TestStagePath_ContainsIncomplete_Extra3(t *testing.T) {
	stage := locking.StagePath("/tmp/testcache/someentry")
	if !strings.Contains(stage, "incomplete") {
		t.Errorf("stage path should contain 'incomplete', got %q", stage)
	}
}

func TestStagePath_SameDir_Extra3(t *testing.T) {
	final := "/tmp/testcache/entry"
	stage := locking.StagePath(final)
	if filepath.Dir(stage) != filepath.Dir(final) {
		t.Errorf("stage path should be in same dir as final: stage=%q final=%q", stage, final)
	}
}

func TestStagePath_Unique_Extra3(t *testing.T) {
	a := locking.StagePath("/tmp/testcache/entry")
	b := locking.StagePath("/tmp/testcache/entry")
	if a == b {
		t.Error("two calls to StagePath should produce different paths")
	}
}

func TestLockUnlock_FileCreatedAndRemoved_Extra3(t *testing.T) {
	tmp := t.TempDir()
	shardDir := filepath.Join(tmp, "shard.bin")
	sl := locking.NewShardLock(shardDir, 2*time.Second)

	if err := sl.Lock(); err != nil {
		t.Fatalf("Lock failed: %v", err)
	}
	// Lock file should exist while locked.
	lockFile := filepath.Join(tmp, "shard.lock")
	if _, err := os.Stat(lockFile); err != nil {
		t.Errorf("lock file should exist while locked: %v", err)
	}
	sl.Unlock()
	// After unlock the lock file should be gone.
	if _, err := os.Stat(lockFile); err == nil {
		t.Error("lock file should be removed after Unlock")
	}
}
