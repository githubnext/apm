package locking_test

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/githubnext/apm/internal/cache/locking"
)

func TestNewShardLock_NotNil_Extra4(t *testing.T) {
	dir := t.TempDir()
	l := locking.NewShardLock(dir, time.Second)
	if l == nil {
		t.Error("expected non-nil ShardLock")
	}
}

func TestNewShardLock_ZeroTimeout_Extra4(t *testing.T) {
	dir := t.TempDir()
	l := locking.NewShardLock(dir, 0)
	if l == nil {
		t.Error("expected non-nil ShardLock even with zero timeout")
	}
}

func TestStagePath_HasIncomplete_Extra4(t *testing.T) {
	p := locking.StagePath("/some/path/shard")
	if !strings.Contains(p, "incomplete") {
		t.Errorf("expected 'incomplete' in stage path, got %q", p)
	}
}

func TestStagePath_SameDir_Extra4(t *testing.T) {
	p := locking.StagePath("/dir/shard")
	if filepath.Dir(p) != "/dir" {
		t.Errorf("expected stage path in same dir /dir, got %q", filepath.Dir(p))
	}
}

func TestStagePath_NotEqualsInput_Extra4(t *testing.T) {
	input := "/dir/shard"
	p := locking.StagePath(input)
	if p == input {
		t.Error("expected stage path to differ from input")
	}
}

func TestLockUnlock_NoError_Extra4(t *testing.T) {
	dir := t.TempDir()
	shardDir := filepath.Join(dir, "shard.d")
	l := locking.NewShardLock(shardDir, 5*time.Second)
	if err := l.Lock(); err != nil {
		t.Fatalf("Lock() error: %v", err)
	}
	l.Unlock()
}

func TestSafeRemoveAll_NonexistentNoError_Extra4(t *testing.T) {
	err := locking.SafeRemoveAll("/nonexistent/path/that/does/not/exist")
	if err != nil {
		t.Errorf("expected no error for nonexistent path, got %v", err)
	}
}

func TestCleanupIncomplete_EmptyDir_Extra4(t *testing.T) {
	dir := t.TempDir()
	n := locking.CleanupIncomplete(dir)
	if n != 0 {
		t.Errorf("expected 0 cleaned in empty dir, got %d", n)
	}
}
