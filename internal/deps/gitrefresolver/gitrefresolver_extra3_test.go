package gitrefresolver

import "testing"

func TestRemoteRef_IsTagTrue_Extra3(t *testing.T) {
	r := RemoteRef{Name: "refs/tags/v2.0.0", SHA: "abc123", IsTag: true, IsBranch: false}
	if !r.IsTag {
		t.Error("IsTag should be true")
	}
	if r.IsBranch {
		t.Error("IsBranch should be false for tag")
	}
}

func TestRemoteRef_IsBranchTrue_Extra3(t *testing.T) {
	r := RemoteRef{Name: "refs/heads/main", SHA: "def456", IsTag: false, IsBranch: true}
	if r.IsTag {
		t.Error("IsTag should be false for branch")
	}
	if !r.IsBranch {
		t.Error("IsBranch should be true")
	}
}

func TestRemoteRef_BothFalse_Extra3(t *testing.T) {
	r := RemoteRef{Name: "HEAD"}
	if r.IsTag || r.IsBranch {
		t.Error("HEAD ref should not be tag or branch")
	}
}

func TestResolvedReference_AllFields_Extra3(t *testing.T) {
	rr := ResolvedReference{
		SHA:     "abc1234567890abcdef1234567890abcdef123456",
		RefType: ReferenceTypeBranch,
		Ref:     "main",
	}
	if rr.SHA == "" {
		t.Error("SHA should not be empty")
	}
	if rr.Ref != "main" {
		t.Errorf("Ref = %q, want main", rr.Ref)
	}
}

func TestResolvedReference_Zero_Extra3(t *testing.T) {
	var rr ResolvedReference
	if rr.SHA != "" || rr.Ref != "" {
		t.Error("zero ResolvedReference should have empty fields")
	}
}

func TestGitHubAPIResult_AllFields_Extra3(t *testing.T) {
	r := GitHubAPIResult{SHA: "abc1234"}
	if r.SHA != "abc1234" {
		t.Errorf("SHA = %q, want abc1234", r.SHA)
	}
}

func TestIsFullSHA_ValidLength_Extra3(t *testing.T) {
	sha := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
	if !IsFullSHA(sha) {
		t.Errorf("IsFullSHA(%q) = false, want true", sha)
	}
}

func TestIsFullSHA_TooShort_Extra3(t *testing.T) {
	if IsFullSHA("abc123") {
		t.Error("IsFullSHA(abc123) = true, want false for short SHA")
	}
}

func TestIsShortSHA_SevenChars_Extra3(t *testing.T) {
	if !IsShortSHA("abc1234") {
		t.Error("IsShortSHA(abc1234) should be true")
	}
}

func TestIsShortSHA_TooShort_Extra3(t *testing.T) {
	if IsShortSHA("abc") {
		t.Error("IsShortSHA(abc) should be false (< 7 chars)")
	}
}

func TestFindRef_CaseSensitive_Extra3(t *testing.T) {
	refs := []RemoteRef{
		{Name: "refs/heads/Main"},
		{Name: "refs/heads/main"},
	}
	r, ok := FindRef(refs, "refs/heads/main")
	if !ok {
		t.Fatal("FindRef should find refs/heads/main")
	}
	if r.Name != "refs/heads/main" {
		t.Errorf("found %q, want refs/heads/main", r.Name)
	}
}

func TestParseLsRemoteOutput_TagAndBranch_Extra3(t *testing.T) {
	out := "abc1234567890abcdef1234567890abcdef123456\trefs/tags/v1.0\ndef1234567890abcdef1234567890abcdef123456\trefs/heads/main\n"
	refs := ParseLsRemoteOutput(out)
	if len(refs) != 2 {
		t.Fatalf("parsed %d refs, want 2", len(refs))
	}
	found := map[string]bool{}
	for _, r := range refs {
		found[r.Name] = true
	}
	// ParseLsRemoteOutput strips the refs/tags/ and refs/heads/ prefixes
	if !found["v1.0"] {
		t.Error("missing v1.0 tag")
	}
	if !found["main"] {
		t.Error("missing main branch")
	}
}
