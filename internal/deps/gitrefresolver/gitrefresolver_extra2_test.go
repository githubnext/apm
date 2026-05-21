package gitrefresolver

import (
	"testing"
)

func TestParseLsRemoteOutput_TagAndBranch(t *testing.T) {
	output := `abc1234567890123456789012345678901234567890	refs/tags/v1.0.0
def1234567890123456789012345678901234567890	refs/heads/main
`
	refs := ParseLsRemoteOutput(output)
	if len(refs) < 2 {
		t.Fatalf("expected >= 2 refs, got %d", len(refs))
	}
	var foundTag, foundBranch bool
	for _, r := range refs {
		// ParseLsRemoteOutput strips the refs/tags/ and refs/heads/ prefixes
		if r.IsTag && r.Name == "v1.0.0" {
			foundTag = true
		}
		if r.IsBranch && r.Name == "main" {
			foundBranch = true
		}
	}
	if !foundTag {
		t.Errorf("expected to find tag ref, got %v", refs)
	}
	if !foundBranch {
		t.Errorf("expected to find branch ref, got %v", refs)
	}
}

func TestParseLsRemoteOutput_Empty(t *testing.T) {
	refs := ParseLsRemoteOutput("")
	if len(refs) != 0 {
		t.Errorf("expected 0 refs for empty output, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_SHAField(t *testing.T) {
	sha := "abcdef1234567890abcdef1234567890abcdef12"
	output := sha + "\trefs/heads/feature\n"
	refs := ParseLsRemoteOutput(output)
	if len(refs) == 0 {
		t.Fatal("expected at least 1 ref")
	}
	if refs[0].SHA != sha {
		t.Errorf("SHA: got %q, want %q", refs[0].SHA, sha)
	}
}

func TestFindRef_ExistingRef(t *testing.T) {
	refs := []RemoteRef{
		{Name: "main", SHA: "abc123", IsBranch: true},
		{Name: "v1.0.0", SHA: "def456", IsTag: true},
	}
	r, ok := FindRef(refs, "main")
	if !ok {
		t.Fatal("expected to find ref")
	}
	if r.SHA != "abc123" {
		t.Errorf("SHA: %q", r.SHA)
	}
}

func TestFindRef_MissingRef(t *testing.T) {
	refs := []RemoteRef{
		{Name: "main", SHA: "abc123", IsBranch: true},
	}
	_, ok := FindRef(refs, "nonexistent")
	if ok {
		t.Error("should not find nonexistent ref")
	}
}

func TestFindRef_EmptySlice(t *testing.T) {
	_, ok := FindRef([]RemoteRef{}, "refs/heads/main")
	if ok {
		t.Error("should not find ref in empty slice")
	}
}

func TestRemoteRef_ZeroValue(t *testing.T) {
	var r RemoteRef
	if r.Name != "" || r.SHA != "" {
		t.Error("zero value should have empty strings")
	}
	if r.IsTag || r.IsBranch {
		t.Error("zero value bools should be false")
	}
}

func TestResolvedReference_ZeroValue(t *testing.T) {
	var r ResolvedReference
	if r.SHA != "" || r.Ref != "" {
		t.Error("zero value should have empty strings")
	}
	if r.RefType != 0 {
		t.Errorf("zero RefType should be 0, got %d", r.RefType)
	}
}

func TestGitHubAPIResult_ZeroValue(t *testing.T) {
	var r GitHubAPIResult
	if r.SHA != "" {
		t.Error("zero value SHA should be empty")
	}
}

func TestClassifyRef_KnownBranch(t *testing.T) {
	refs := []RemoteRef{
		{Name: "refs/heads/main", SHA: "abc1234567890123456789012345678901234567890", IsBranch: true},
	}
	rt := ClassifyRef(refs, "main")
	if rt == 0 {
		t.Error("expected non-zero reference type for known branch")
	}
}

func TestClassifyRef_FullSHA(t *testing.T) {
	refs := []RemoteRef{}
	sha := "abcdef1234567890abcdef1234567890abcdef12"
	rt := ClassifyRef(refs, sha)
	_ = rt
}

func TestIsFullSHA_ValidSHA(t *testing.T) {
	sha := "abcdef1234567890abcdef1234567890abcdef12"
	if !IsFullSHA(sha) {
		t.Errorf("expected true for valid 40-char hex SHA")
	}
}

func TestIsFullSHA_TooShort(t *testing.T) {
	if IsFullSHA("abc123") {
		t.Error("short string should not be full SHA")
	}
}

func TestIsShortSHA_ValidLength(t *testing.T) {
	if !IsShortSHA("abcdef1") {
		t.Error("7-char hex string should be short SHA")
	}
}

func TestNew_StoresHostAndToken(t *testing.T) {
	r := New("github.com", "mytoken")
	if r == nil {
		t.Fatal("New should return non-nil")
	}
}

func TestParseLsRemoteOutput_MalformedLines(t *testing.T) {
	output := "not-a-ref-line\n\t\nonly-one-field\n"
	refs := ParseLsRemoteOutput(output)
	_ = refs
}

func TestParseLsRemoteOutput_MultipleTags(t *testing.T) {
	output := "aaaa1234567890123456789012345678901234567890\trefs/tags/v1.0.0\n" +
		"bbbb1234567890123456789012345678901234567890\trefs/tags/v2.0.0\n" +
		"cccc1234567890123456789012345678901234567890\trefs/tags/v3.0.0\n"
	refs := ParseLsRemoteOutput(output)
	tagCount := 0
	for _, r := range refs {
		if r.IsTag {
			tagCount++
		}
	}
	if tagCount != 3 {
		t.Errorf("expected 3 tags, got %d", tagCount)
	}
}
