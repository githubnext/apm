package gitcache

import (
	"testing"
)

func TestNew(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp, false)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c == nil {
		t.Fatal("New returned nil")
	}
	if c.cacheRoot != tmp {
		t.Errorf("cacheRoot: got %q, want %q", c.cacheRoot, tmp)
	}
	if c.refresh {
		t.Error("refresh should be false")
	}

	// With refresh=true
	c2, err := New(tmp, true)
	if err != nil {
		t.Fatalf("New(refresh=true): %v", err)
	}
	if !c2.refresh {
		t.Error("refresh should be true")
	}
}

func TestGetCacheStats(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	stats := c.GetCacheStats()
	if stats.DBCount < 0 || stats.CheckoutCount < 0 {
		t.Error("stats counts should be non-negative")
	}
}

func TestCleanAll(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	// Should not panic or error on an empty cache
	c.CleanAll()
}

func TestPruneEmptyCache(t *testing.T) {
	tmp := t.TempDir()
	c, err := New(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	removed := c.Prune(30)
	if removed != 0 {
		t.Errorf("Prune on empty cache: got %d, want 0", removed)
	}
}

func TestSanitizeURL(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"https://user:pass@github.com/org/repo", "https://***@github.com/org/repo"},
		{"https://github.com/org/repo", "https://github.com/org/repo"},
		{"git@github.com:org/repo.git", "git@github.com:org/repo.git"},
		{"", ""},
	}
	for _, tc := range cases {
		got := sanitizeURL(tc.in)
		if got != tc.want {
			t.Errorf("sanitizeURL(%q): got %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestMergeEnv(t *testing.T) {
	base := []string{"A=1", "B=2"}
	extra := []string{"C=3"}
	merged := mergeEnv(base, extra)
	if len(merged) != 3 {
		t.Errorf("mergeEnv: got %d elements, want 3", len(merged))
	}
	// Empty extra
	merged2 := mergeEnv(base, nil)
	if len(merged2) != 2 {
		t.Errorf("mergeEnv(nil extra): got %d, want 2", len(merged2))
	}
}
