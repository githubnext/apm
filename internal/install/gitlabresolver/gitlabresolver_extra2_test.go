package gitlabresolver

import (
	"strings"
	"testing"
)

func TestGitLabDirectShorthandUnresolved_NotEmpty(t *testing.T) {
	if GitLabDirectShorthandUnresolved == "" {
		t.Error("GitLabDirectShorthandUnresolved should not be empty")
	}
}

func TestGitLabDirectShorthandUnresolved_ContainsDirect(t *testing.T) {
	if !strings.Contains(GitLabDirectShorthandUnresolved, "Direct GitLab") {
		t.Error("expected 'Direct GitLab' in error constant")
	}
}

func TestShorthandParts_ZeroValue(t *testing.T) {
	var sp ShorthandParts
	if sp.Host != "" || sp.Ref != "" || sp.Segments != nil {
		t.Error("zero ShorthandParts should have empty Host, Ref, nil Segments")
	}
}

func TestBoundaryCandidate_ZeroValue(t *testing.T) {
	var bc BoundaryCandidate
	if bc.RepoPath != "" || bc.VirtualPath != "" {
		t.Error("zero BoundaryCandidate fields should be empty")
	}
}

func TestParseShorthand_NoDot(t *testing.T) {
	// host without dot should return nil
	result := ParseShorthand("noDotHost/owner/repo")
	if result != nil {
		t.Error("expected nil for host with no dot")
	}
}

func TestParseShorthand_WithRef(t *testing.T) {
	sp := ParseShorthand("gitlab.example.com/owner/repo#v1.2")
	if sp == nil {
		t.Fatal("expected non-nil")
	}
	if sp.Ref != "v1.2" {
		t.Errorf("expected Ref=v1.2, got %q", sp.Ref)
	}
	if sp.Host != "gitlab.example.com" {
		t.Errorf("expected Host=gitlab.example.com, got %q", sp.Host)
	}
}

func TestParseShorthand_NoRef(t *testing.T) {
	sp := ParseShorthand("gitlab.com/org/repo")
	if sp == nil {
		t.Fatal("expected non-nil")
	}
	if sp.Ref != "" {
		t.Errorf("expected empty Ref, got %q", sp.Ref)
	}
}

func TestParseShorthand_SingleSegment(t *testing.T) {
	// "gitlab.com" alone has no slash-separated path
	result := ParseShorthand("gitlab.com")
	if result != nil {
		t.Error("expected nil when no path segments")
	}
}

func TestIsGitLabHost_ADOVisualStudio(t *testing.T) {
	if IsGitLabHost("myorg.visualstudio.com") {
		t.Error("visualstudio.com should not be considered GitLab")
	}
}

func TestIsGitLabHost_DevAzureCom(t *testing.T) {
	if IsGitLabHost("dev.azure.com") {
		t.Error("dev.azure.com should not be considered GitLab")
	}
}

func TestIsGitLabHost_GHECom(t *testing.T) {
	if IsGitLabHost("myorg.ghe.com") {
		t.Error("*.ghe.com should not be considered GitLab")
	}
}

func TestBoundaryCandidates_VirtualPath(t *testing.T) {
	sp := &ShorthandParts{Host: "gl.io", Segments: []string{"org", "repo", "sub"}, Ref: ""}
	bc := NewBoundaryCandidates(sp)
	first, ok := bc.Next()
	if !ok {
		t.Fatal("expected at least one candidate")
	}
	// First candidate: all 3 segments as repo path, virtual=""
	if first.RepoPath != "org/repo/sub" {
		t.Errorf("expected RepoPath=org/repo/sub, got %q", first.RepoPath)
	}
	if first.VirtualPath != "" {
		t.Errorf("expected VirtualPath empty, got %q", first.VirtualPath)
	}
}

func TestBoundaryCandidates_SecondCandidate(t *testing.T) {
	sp := &ShorthandParts{Host: "gl.io", Segments: []string{"org", "repo", "sub"}, Ref: ""}
	bc := NewBoundaryCandidates(sp)
	_, _ = bc.Next() // skip first
	second, ok := bc.Next()
	if !ok {
		t.Fatal("expected second candidate")
	}
	if second.RepoPath != "org/repo" {
		t.Errorf("expected org/repo, got %q", second.RepoPath)
	}
	if second.VirtualPath != "sub" {
		t.Errorf("expected VirtualPath=sub, got %q", second.VirtualPath)
	}
}
