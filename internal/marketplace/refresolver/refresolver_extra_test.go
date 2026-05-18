package refresolver_test

import (
	"strings"
	"testing"
	"time"

	"github.com/githubnext/apm/internal/marketplace/refresolver"
)

func TestRefCache_GetMiss(t *testing.T) {
	c := refresolver.NewRefCache(5 * time.Minute)
	refs := c.Get("owner/repo")
	if refs != nil {
		t.Errorf("expected nil on cache miss, got %v", refs)
	}
}

func TestRefCache_PutAndGetTwo(t *testing.T) {
	c := refresolver.NewRefCache(5 * time.Minute)
	refs := []refresolver.RemoteRef{
		{Name: "refs/heads/main", SHA: strings.Repeat("a", 40)},
		{Name: "refs/tags/v1.0.0", SHA: strings.Repeat("b", 40)},
	}
	c.Put("owner/repo", refs)
	got := c.Get("owner/repo")
	if got == nil {
		t.Fatal("expected cache hit")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 refs, got %d", len(got))
	}
}

func TestRefCache_ExpiryShort(t *testing.T) {
	c := refresolver.NewRefCache(1 * time.Millisecond)
	c.Put("owner/repo", []refresolver.RemoteRef{{Name: "refs/heads/main", SHA: strings.Repeat("b", 40)}})
	time.Sleep(5 * time.Millisecond)
	got := c.Get("owner/repo")
	if got != nil {
		t.Errorf("expected expired entry to be nil, got %v", got)
	}
}

func TestRefCache_ClearAll(t *testing.T) {
	c := refresolver.NewRefCache(5 * time.Minute)
	c.Put("owner/repo", []refresolver.RemoteRef{{Name: "refs/heads/main", SHA: strings.Repeat("c", 40)}})
	c.Clear()
	if c.Len() != 0 {
		t.Errorf("expected empty cache after Clear, got %d", c.Len())
	}
}

func TestRefCache_LenGrows(t *testing.T) {
	c := refresolver.NewRefCache(5 * time.Minute)
	if c.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", c.Len())
	}
	c.Put("a/b", []refresolver.RemoteRef{{Name: "refs/heads/main", SHA: strings.Repeat("d", 40)}})
	c.Put("c/d", []refresolver.RemoteRef{{Name: "refs/heads/main", SHA: strings.Repeat("e", 40)}})
	if c.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", c.Len())
	}
}

func TestRemoteRef_fields(t *testing.T) {
	r := refresolver.RemoteRef{
		Name: "refs/tags/v1.2.3",
		SHA:  strings.Repeat("f", 40),
	}
	if r.Name != "refs/tags/v1.2.3" {
		t.Errorf("unexpected Name: %s", r.Name)
	}
	if len(r.SHA) != 40 {
		t.Errorf("SHA should be 40 chars, got %d", len(r.SHA))
	}
}

func TestGitLsRemoteError(t *testing.T) {
	err := &refresolver.GitLsRemoteError{Package: "owner/repo", Summary: "fatal: not a git repo"}
	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}
}

func TestOfflineMissError(t *testing.T) {
	err := &refresolver.OfflineMissError{Package: "org/project", Remote: "https://github.com/org/project"}
	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}
}

func TestRefCache_PutOverwrites(t *testing.T) {
	c := refresolver.NewRefCache(5 * time.Minute)
	sha1 := strings.Repeat("1", 40)
	sha2 := strings.Repeat("2", 40)
	c.Put("owner/repo", []refresolver.RemoteRef{{Name: "refs/heads/main", SHA: sha1}})
	c.Put("owner/repo", []refresolver.RemoteRef{{Name: "refs/heads/main", SHA: sha2}})
	got := c.Get("owner/repo")
	if got == nil || got[0].SHA != sha2 {
		t.Errorf("expected overwritten SHA %s, got %v", sha2, got)
	}
}

func TestRefCache_MultipleKeys(t *testing.T) {
	c := refresolver.NewRefCache(5 * time.Minute)
	c.Put("org1/repo1", []refresolver.RemoteRef{{Name: "refs/heads/main", SHA: strings.Repeat("a", 40)}})
	c.Put("org2/repo2", []refresolver.RemoteRef{{Name: "refs/heads/main", SHA: strings.Repeat("b", 40)}})
	if c.Get("org1/repo1") == nil {
		t.Error("expected hit for org1/repo1")
	}
	if c.Get("org2/repo2") == nil {
		t.Error("expected hit for org2/repo2")
	}
	if c.Get("org3/repo3") != nil {
		t.Error("expected miss for org3/repo3")
	}
}
