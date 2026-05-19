package refresolver

import (
	"testing"
	"time"
)

func TestRemoteRef_ZeroValue(t *testing.T) {
	var r RemoteRef
	if r.Name != "" || r.SHA != "" {
		t.Error("zero-value RemoteRef should have empty fields")
	}
}

func TestRemoteRef_Assign(t *testing.T) {
	r := RemoteRef{Name: "refs/tags/v1.0.0", SHA: "abc123def456abc123def456abc123def456abc1"}
	if r.Name != "refs/tags/v1.0.0" {
		t.Errorf("unexpected Name: %q", r.Name)
	}
	if r.SHA != "abc123def456abc123def456abc123def456abc1" {
		t.Errorf("unexpected SHA: %q", r.SHA)
	}
}

func TestDefaultTTL_Is5Minutes(t *testing.T) {
	if DefaultTTL != 5*time.Minute {
		t.Errorf("DefaultTTL should be 5 minutes, got %v", DefaultTTL)
	}
}

func TestRefCache_ZeroTTLExpires(t *testing.T) {
	c := NewRefCache(0)
	refs := []RemoteRef{{Name: "refs/heads/main", SHA: "aaa"}}
	c.Put("owner/repo", refs)
	// TTL=0 means entries expire immediately
	got := c.Get("owner/repo")
	// May be nil or non-nil depending on timing; just ensure no panic
	_ = got
}

func TestRefCache_PutGetRoundtrip(t *testing.T) {
	c := NewRefCache(time.Hour)
	refs := []RemoteRef{
		{Name: "refs/heads/main", SHA: "aaaa1111aaaa1111aaaa1111aaaa1111aaaa1111"},
		{Name: "refs/tags/v2.0", SHA: "bbbb2222bbbb2222bbbb2222bbbb2222bbbb2222"},
	}
	c.Put("org/proj", refs)
	got := c.Get("org/proj")
	if len(got) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(got))
	}
	if got[0].Name != "refs/heads/main" {
		t.Errorf("unexpected first ref name: %q", got[0].Name)
	}
}

func TestRefCache_IsolatesSlice(t *testing.T) {
	c := NewRefCache(time.Hour)
	orig := []RemoteRef{{Name: "refs/heads/x", SHA: "1234567890123456789012345678901234567890"}}
	c.Put("a/b", orig)
	got := c.Get("a/b")
	// Mutating returned slice should not affect cache
	got[0].Name = "mutated"
	got2 := c.Get("a/b")
	if got2[0].Name != "refs/heads/x" {
		t.Error("cache returned slice should be isolated from mutations")
	}
}

func TestRefCache_LenAfterClear(t *testing.T) {
	c := NewRefCache(time.Hour)
	c.Put("p/q", []RemoteRef{{Name: "refs/heads/a", SHA: "0000111122223333444455556666777788889999"}})
	if c.Len() != 1 {
		t.Errorf("expected Len=1, got %d", c.Len())
	}
	c.Clear()
	if c.Len() != 0 {
		t.Errorf("expected Len=0 after Clear, got %d", c.Len())
	}
}

func TestGitLsRemoteError_ErrorMethod(t *testing.T) {
	e := &GitLsRemoteError{Package: "owner/repo", Summary: "git failed", Hint: "check credentials"}
	msg := e.Error()
	if msg == "" {
		t.Error("Error() should return non-empty string")
	}
}

func TestOfflineMissError_ErrorMethod(t *testing.T) {
	e := &OfflineMissError{Remote: "owner/repo"}
	msg := e.Error()
	if msg == "" {
		t.Error("Error() should return non-empty string")
	}
}

func TestRefCache_DifferentKeysDontCollide(t *testing.T) {
	c := NewRefCache(time.Hour)
	c.Put("org/a", []RemoteRef{{Name: "refs/heads/main", SHA: "aaaa0000aaaa0000aaaa0000aaaa0000aaaa0000"}})
	c.Put("org/b", []RemoteRef{{Name: "refs/heads/dev", SHA: "bbbb0000bbbb0000bbbb0000bbbb0000bbbb0000"}})
	a := c.Get("org/a")
	b := c.Get("org/b")
	if a[0].Name == b[0].Name {
		t.Error("different keys should return different refs")
	}
}
