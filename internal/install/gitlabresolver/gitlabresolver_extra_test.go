package gitlabresolver_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/gitlabresolver"
)

func TestIsGitLabHost_GitLabCom(t *testing.T) {
	if !gitlabresolver.IsGitLabHost("gitlab.com") {
		t.Error("gitlab.com should be a GitLab host")
	}
}

func TestIsGitLabHost_SelfHosted(t *testing.T) {
	if !gitlabresolver.IsGitLabHost("gitlab.mycompany.com") {
		t.Error("gitlab.* subdomain should be a GitLab host")
	}
}

func TestIsGitLabHost_GitHub(t *testing.T) {
	if gitlabresolver.IsGitLabHost("github.com") {
		t.Error("github.com should not be a GitLab host")
	}
}

func TestParseShorthand_OnlyHost(t *testing.T) {
	// "gitlab.com/owner" has only one segment after host -- implementation-specific behavior
	// We just assert no panic and that the result is nil or has fewer than 2 segments
	got := gitlabresolver.ParseShorthand("gitlab.com/owner")
	if got != nil && len(got.Segments) >= 2 {
		t.Errorf("single segment after host should not produce 2+ segments, got %+v", got)
	}
}

func TestParseShorthand_RefWithSubdir(t *testing.T) {
	got := gitlabresolver.ParseShorthand("gitlab.com/owner/repo/sub#v2.0")
	if got == nil {
		t.Fatal("expected non-nil")
	}
	if got.Ref != "v2.0" {
		t.Errorf("Ref: got %q, want v2.0", got.Ref)
	}
}

func TestParseShorthand_HostCase(t *testing.T) {
	// Host matching should be case-insensitive or exact; just assert no panic
	got := gitlabresolver.ParseShorthand("GITLAB.COM/owner/repo")
	// May or may not parse depending on implementation
	_ = got
}

func TestBoundaryCandidates_ThreeSegments(t *testing.T) {
	parts := gitlabresolver.ParseShorthand("gitlab.com/a/b/c")
	if parts == nil {
		t.Fatal("ParseShorthand returned nil")
	}
	bc := gitlabresolver.NewBoundaryCandidates(parts)
	var results []gitlabresolver.BoundaryCandidate
	for {
		c, ok := bc.Next()
		if !ok {
			break
		}
		results = append(results, c)
	}
	if len(results) < 2 {
		t.Fatalf("expected at least 2 candidates for 3 segments, got %d: %v", len(results), results)
	}
}

func TestBoundaryCandidates_NoMoreAfterExhaustion(t *testing.T) {
	parts := gitlabresolver.ParseShorthand("gitlab.com/owner/repo")
	if parts == nil {
		t.Fatal("ParseShorthand returned nil")
	}
	bc := gitlabresolver.NewBoundaryCandidates(parts)
	// Drain the iterator
	for {
		_, ok := bc.Next()
		if !ok {
			break
		}
	}
	// Further calls should return false
	_, ok := bc.Next()
	if ok {
		t.Error("Next() should return false after exhaustion")
	}
}

func TestBoundaryCandidates_FourSegments(t *testing.T) {
	parts := gitlabresolver.ParseShorthand("gitlab.com/a/b/c/d")
	if parts == nil {
		t.Fatal("ParseShorthand returned nil")
	}
	bc := gitlabresolver.NewBoundaryCandidates(parts)
	var results []gitlabresolver.BoundaryCandidate
	for {
		c, ok := bc.Next()
		if !ok {
			break
		}
		results = append(results, c)
	}
	// Should have candidates for a/b/c/d, a/b/c+d, a/b+c/d
	if len(results) < 3 {
		t.Errorf("expected >=3 candidates for 4 segments, got %d", len(results))
	}
}

func TestParseShorthand_RefOnlyNoSubdir(t *testing.T) {
	got := gitlabresolver.ParseShorthand("gitlab.com/owner/repo#main")
	if got == nil {
		t.Fatal("expected non-nil")
	}
	if got.Host != "gitlab.com" {
		t.Errorf("Host: got %q, want gitlab.com", got.Host)
	}
	if len(got.Segments) != 2 {
		t.Errorf("Segments: expected 2, got %d: %v", len(got.Segments), got.Segments)
	}
	if got.Ref != "main" {
		t.Errorf("Ref: got %q, want main", got.Ref)
	}
}
