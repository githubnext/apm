package deptypes

import (
	"testing"
)

func TestParseGitReference_BranchWithDot(t *testing.T) {
	ref := "feature/add.stuff"
	refType, name := ParseGitReference(ref)
	if refType != GitRefBranch {
		t.Errorf("expected GitRefBranch for branch-like ref, got %v", refType)
	}
	if name != ref {
		t.Errorf("expected name=%q, got %q", ref, name)
	}
}

func TestParseGitReference_41CharHex_NotCommit(t *testing.T) {
	// 41 chars -- too long for commit, should be branch
	ref := "a1b2c3d4e5a1b2c3d4e5a1b2c3d4e5a1b2c3d4e5a"
	refType, _ := ParseGitReference(ref)
	if refType == GitRefCommit {
		t.Errorf("41-char hex should not be GitRefCommit")
	}
}

func TestParseGitReference_SemverNoPrefix(t *testing.T) {
	ref := "1.2.3"
	refType, name := ParseGitReference(ref)
	if refType != GitRefTag {
		t.Errorf("expected GitRefTag for semver without v-prefix, got %v", refType)
	}
	if name != ref {
		t.Errorf("expected name=%q, got %q", ref, name)
	}
}

func TestParseGitReference_7CharHex(t *testing.T) {
	ref := "abc1234"
	refType, _ := ParseGitReference(ref)
	if refType != GitRefCommit {
		t.Errorf("expected GitRefCommit for 7-char hex, got %v", refType)
	}
}

func TestRemoteRef_ZeroValue(t *testing.T) {
	var r RemoteRef
	if r.Name != "" || r.CommitSHA != "" {
		t.Error("zero-value RemoteRef should have empty Name and CommitSHA")
	}
}

func TestResolvedReference_ZeroValue(t *testing.T) {
	var r ResolvedReference
	if r.OriginalRef != "" || r.ResolvedCommit != "" {
		t.Error("zero-value ResolvedReference should have empty fields")
	}
}

func TestVirtualPackageType_Distinct(t *testing.T) {
	if VirtualPackageFile == VirtualPackageSubdirectory {
		t.Error("VirtualPackageFile and VirtualPackageSubdirectory should be distinct")
	}
}

func TestGitReferenceType_ThreeDistinct(t *testing.T) {
	types := []GitReferenceType{GitRefBranch, GitRefTag, GitRefCommit}
	seen := map[GitReferenceType]bool{}
	for _, t2 := range types {
		if seen[t2] {
			t.Errorf("duplicate GitReferenceType value: %v", t2)
		}
		seen[t2] = true
	}
}

func TestParseGitReference_TagWithMinorPatch(t *testing.T) {
	for _, ref := range []string{"v0.0.1", "v10.20.30", "1.0.0"} {
		refType, _ := ParseGitReference(ref)
		if refType != GitRefTag {
			t.Errorf("expected GitRefTag for %q, got %v", ref, refType)
		}
	}
}

func TestRemoteRef_Slice(t *testing.T) {
	refs := []RemoteRef{
		{Name: "main", RefType: GitRefBranch, CommitSHA: "abc"},
		{Name: "v1.0.0", RefType: GitRefTag, CommitSHA: "def"},
	}
	if len(refs) != 2 {
		t.Errorf("expected 2 refs, got %d", len(refs))
	}
}
