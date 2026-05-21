package deptypes

import "testing"

func TestParseGitReference_UppercaseHex_IsBranch(t *testing.T) {
	// Uppercase hex is not matched by the lowercase commit regex; treated as branch.
	refType, name := ParseGitReference("ABCDEF1234567")
	if refType != GitRefBranch {
		t.Errorf("uppercase hex: expected GitRefBranch, got %d", refType)
	}
	if name != "ABCDEF1234567" {
		t.Errorf("expected name=ABCDEF1234567, got %q", name)
	}
}

func TestParseGitReference_SixCharHex_IsBranch(t *testing.T) {
	refType, _ := ParseGitReference("abc123")
	if refType != GitRefBranch {
		t.Errorf("6-char hex: expected GitRefBranch, got %d", refType)
	}
}

func TestParseGitReference_SevenCharHex_IsCommit(t *testing.T) {
	refType, _ := ParseGitReference("abc1234")
	if refType != GitRefCommit {
		t.Errorf("7-char hex: expected GitRefCommit, got %d", refType)
	}
}

func TestParseGitReference_PrefixedSemver(t *testing.T) {
	cases := []struct {
		input   string
		refType GitReferenceType
	}{
		{"v0.0.1", GitRefTag},
		{"v100.200.300", GitRefTag},
		{"0.0.0", GitRefTag},
	}
	for _, tc := range cases {
		got, _ := ParseGitReference(tc.input)
		if got != tc.refType {
			t.Errorf("ParseGitReference(%q): got %d, want %d", tc.input, got, tc.refType)
		}
	}
}

func TestRemoteRef_FieldAssignment(t *testing.T) {
	cases := []struct {
		name    string
		refType GitReferenceType
		sha     string
	}{
		{"main", GitRefBranch, "abc1234"},
		{"v1.2.3", GitRefTag, "def5678"},
		{"abc1234abc1234a", GitRefCommit, "abc1234abc1234a"},
	}
	for _, tc := range cases {
		r := RemoteRef{Name: tc.name, RefType: tc.refType, CommitSHA: tc.sha}
		if r.Name != tc.name || r.RefType != tc.refType || r.CommitSHA != tc.sha {
			t.Errorf("RemoteRef field mismatch for %q", tc.name)
		}
	}
}

func TestGitReferenceType_Iota(t *testing.T) {
	if int(GitRefBranch) != 0 {
		t.Error("GitRefBranch should be 0")
	}
	if int(GitRefTag) != 1 {
		t.Error("GitRefTag should be 1")
	}
	if int(GitRefCommit) != 2 {
		t.Error("GitRefCommit should be 2")
	}
}

func TestVirtualPackageType_Iota(t *testing.T) {
	if int(VirtualPackageFile) != 0 {
		t.Error("VirtualPackageFile should be 0")
	}
	if int(VirtualPackageSubdirectory) != 1 {
		t.Error("VirtualPackageSubdirectory should be 1")
	}
}

func TestResolvedReference_AllFields(t *testing.T) {
	rr := ResolvedReference{
		OriginalRef:    "feature/x",
		RefType:        GitRefBranch,
		ResolvedCommit: "abc1234567890",
		RefName:        "feature/x",
	}
	if rr.OriginalRef != "feature/x" {
		t.Errorf("OriginalRef mismatch: %q", rr.OriginalRef)
	}
	if rr.RefType != GitRefBranch {
		t.Errorf("RefType mismatch: %d", rr.RefType)
	}
	if rr.ResolvedCommit != "abc1234567890" {
		t.Errorf("ResolvedCommit mismatch: %q", rr.ResolvedCommit)
	}
	if rr.RefName != "feature/x" {
		t.Errorf("RefName mismatch: %q", rr.RefName)
	}
}

func TestParseGitReference_BranchWithSlash(t *testing.T) {
	refType, name := ParseGitReference("feature/my-branch")
	if refType != GitRefBranch {
		t.Errorf("expected GitRefBranch, got %d", refType)
	}
	if name != "feature/my-branch" {
		t.Errorf("expected 'feature/my-branch', got %q", name)
	}
}

func TestParseGitReference_Distinctness(t *testing.T) {
	commitRef, _ := ParseGitReference("a1b2c3d4e5f6a7b")
	tagRef, _ := ParseGitReference("v1.2.3")
	branchRef, _ := ParseGitReference("main")

	if commitRef == tagRef || commitRef == branchRef || tagRef == branchRef {
		t.Error("all three reference types should be distinct")
	}
}
