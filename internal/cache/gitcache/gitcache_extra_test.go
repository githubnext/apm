package gitcache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew_CreatesDirectories(t *testing.T) {
	tmp := t.TempDir()
	_, err := New(tmp, false)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Verify subdirectories were created
	entries, err := os.ReadDir(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Error("expected at least one subdirectory created by New")
	}
}

func TestGetCacheStats_AfterCleanAll(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	c.CleanAll()
	stats := c.GetCacheStats()
	if stats.DBCount != 0 {
		t.Errorf("DBCount after CleanAll: got %d, want 0", stats.DBCount)
	}
	if stats.CheckoutCount != 0 {
		t.Errorf("CheckoutCount after CleanAll: got %d, want 0", stats.CheckoutCount)
	}
}

func TestCleanAll_Idempotent(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	c.CleanAll()
	c.CleanAll() // should not panic or error
}

func TestPrune_OldEntryRemoved(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	// Create a fake checkout dir with an old modification time
	checkoutsRoot := filepath.Join(tmp, "git", "checkouts_v1")
	oldDir := filepath.Join(checkoutsRoot, "old-entry")
	if err := os.MkdirAll(oldDir, 0o700); err != nil {
		t.Fatal(err)
	}
	// Set mtime to 60 days ago
	pastTime := time.Now().AddDate(0, 0, -60)
	if err := os.Chtimes(oldDir, pastTime, pastTime); err != nil {
		t.Fatal(err)
	}
	removed := c.Prune(30)
	if removed < 1 {
		t.Errorf("expected at least 1 pruned, got %d", removed)
	}
}

func TestPrune_RecentEntryKept(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	checkoutsRoot := filepath.Join(tmp, "git", "checkouts_v1")
	newDir := filepath.Join(checkoutsRoot, "recent-entry")
	if err := os.MkdirAll(newDir, 0o700); err != nil {
		t.Fatal(err)
	}
	removed := c.Prune(30)
	if removed != 0 {
		t.Errorf("expected 0 pruned for recent entry, got %d", removed)
	}
}

func TestSanitizeURL_NoCredentials(t *testing.T) {
	got := sanitizeURL("https://github.com/org/repo")
	if got != "https://github.com/org/repo" {
		t.Errorf("sanitizeURL without credentials changed URL: %q", got)
	}
}

func TestSanitizeURL_EmptyString(t *testing.T) {
	got := sanitizeURL("")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestSanitizeURL_SSHNoCredentials(t *testing.T) {
	url := "git@github.com:org/repo.git"
	got := sanitizeURL(url)
	if got != url {
		t.Errorf("SSH URL should not be modified: %q", got)
	}
}

func TestMergeEnv_NilExtra(t *testing.T) {
	base := []string{"A=1"}
	result := mergeEnv(base, nil)
	if len(result) != 1 || result[0] != "A=1" {
		t.Errorf("nil extra: got %v, want [A=1]", result)
	}
}

func TestMergeEnv_BothPresent(t *testing.T) {
	base := []string{"A=1"}
	extra := []string{"B=2"}
	result := mergeEnv(base, extra)
	if len(result) != 2 {
		t.Errorf("expected 2 elements, got %v", result)
	}
}
