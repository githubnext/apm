package refresolver_test

import (
	"testing"
	"time"

	"github.com/githubnext/apm/internal/marketplace/refresolver"
)

// --- RefCache tests ---

func TestRefCache_MissOnEmpty(t *testing.T) {
	c := refresolver.NewRefCache(5 * time.Minute)
	if got := c.Get("owner/repo"); got != nil {
		t.Fatalf("expected nil on miss, got %v", got)
	}
}

func TestRefCache_PutAndGet(t *testing.T) {
	c := refresolver.NewRefCache(5 * time.Minute)
	refs := []refresolver.RemoteRef{
		{Name: "refs/tags/v1.0.0", SHA: "aabbccdd" + "aabbccddaabbccddaabbccddaabbccdd"},
	}
	c.Put("owner/repo", refs)
	got := c.Get("owner/repo")
	if len(got) != 1 || got[0].Name != "refs/tags/v1.0.0" {
		t.Fatalf("unexpected refs: %v", got)
	}
}

func TestRefCache_Expiry(t *testing.T) {
	c := refresolver.NewRefCache(1 * time.Millisecond)
	refs := []refresolver.RemoteRef{{Name: "refs/heads/main", SHA: "aabbccddaabbccddaabbccddaabbccddaabbccdd"}}
	c.Put("owner/repo", refs)
	time.Sleep(10 * time.Millisecond)
	if got := c.Get("owner/repo"); got != nil {
		t.Fatalf("expected nil after expiry, got %v", got)
	}
}

func TestRefCache_Len(t *testing.T) {
	c := refresolver.NewRefCache(time.Minute)
	if c.Len() != 0 {
		t.Fatal("expected 0 len")
	}
	c.Put("a/b", nil)
	c.Put("c/d", nil)
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestRefCache_Clear(t *testing.T) {
	c := refresolver.NewRefCache(time.Minute)
	c.Put("a/b", nil)
	c.Clear()
	if c.Len() != 0 {
		t.Fatal("expected 0 after clear")
	}
}

func TestRefCache_CopyIsolation(t *testing.T) {
	c := refresolver.NewRefCache(time.Minute)
	orig := []refresolver.RemoteRef{{Name: "refs/heads/main", SHA: "aabbccddaabbccddaabbccddaabbccddaabbccdd"}}
	c.Put("owner/repo", orig)
	orig[0].Name = "mutated"
	got := c.Get("owner/repo")
	if got[0].Name == "mutated" {
		t.Fatal("cache did not copy on Put -- mutation leaked")
	}
}

// --- Offline mode ---

func TestNew_OfflineMode_ReturnsOfflineMissError(t *testing.T) {
	r := refresolver.New(5, true, "github.com", "")
	defer r.Close()
	_, err := r.ListRemoteRefs("owner/repo")
	if err == nil {
		t.Fatal("expected error in offline mode")
	}
	ome, ok := err.(*refresolver.OfflineMissError)
	if !ok {
		t.Fatalf("expected *OfflineMissError, got %T: %v", err, err)
	}
	if ome.Remote != "owner/repo" {
		t.Errorf("unexpected remote: %q", ome.Remote)
	}
}

func TestOfflineMissError_Message(t *testing.T) {
	e := &refresolver.OfflineMissError{Remote: "foo/bar"}
	if e.Error() == "" {
		t.Fatal("error message must not be empty")
	}
}

func TestGitLsRemoteError_WithHint(t *testing.T) {
	e := &refresolver.GitLsRemoteError{Summary: "failed", Hint: "try again"}
	if e.Error() != "failed try again" {
		t.Errorf("unexpected: %q", e.Error())
	}
}

func TestGitLsRemoteError_NoHint(t *testing.T) {
	e := &refresolver.GitLsRemoteError{Summary: "failed"}
	if e.Error() != "failed" {
		t.Errorf("unexpected: %q", e.Error())
	}
}
