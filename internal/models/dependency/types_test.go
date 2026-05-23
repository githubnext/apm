package dependency_test

import (
	"testing"

	"github.com/githubnext/apm/internal/models/dependency"
)

// TestParityGitRefBranch mirrors test: parse_git_reference("main") -> (BRANCH, "main")
func TestParityGitRefBranch(t *testing.T) {
	refType, ref := dependency.ParseGitReference("main")
	if refType != dependency.GitRefBranch {
		t.Errorf("expected BRANCH, got %s", refType)
	}
	if ref != "main" {
		t.Errorf("expected 'main', got %s", ref)
	}
}

// TestParityGitRefEmpty mirrors test: parse_git_reference("") -> (BRANCH, "main")
func TestParityGitRefEmpty(t *testing.T) {
	refType, ref := dependency.ParseGitReference("")
	if refType != dependency.GitRefBranch {
		t.Errorf("expected BRANCH, got %s", refType)
	}
	if ref != "main" {
		t.Errorf("expected 'main', got %s", ref)
	}
}

// TestParityGitRefCommitSHA mirrors: parse_git_reference("abc1234") -> (COMMIT, "abc1234")
func TestParityGitRefCommitSHA(t *testing.T) {
	refType, ref := dependency.ParseGitReference("abc1234")
	if refType != dependency.GitRefCommit {
		t.Errorf("expected COMMIT, got %s", refType)
	}
	if ref != "abc1234" {
		t.Errorf("expected 'abc1234', got %s", ref)
	}
}

// TestParityGitRefFullSHA mirrors: full 40-char SHA -> COMMIT
func TestParityGitRefFullSHA(t *testing.T) {
	sha := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
	refType, ref := dependency.ParseGitReference(sha)
	if refType != dependency.GitRefCommit {
		t.Errorf("expected COMMIT, got %s", refType)
	}
	if ref != sha {
		t.Errorf("expected full SHA, got %s", ref)
	}
}

// TestParityGitRefSemver mirrors: "v1.2.3" -> (TAG, "v1.2.3")
func TestParityGitRefSemver(t *testing.T) {
	refType, ref := dependency.ParseGitReference("v1.2.3")
	if refType != dependency.GitRefTag {
		t.Errorf("expected TAG, got %s", refType)
	}
	if ref != "v1.2.3" {
		t.Errorf("expected 'v1.2.3', got %s", ref)
	}
}

// TestParityGitRefSemverNoV mirrors: "1.2.3" -> (TAG, "1.2.3")
func TestParityGitRefSemverNoV(t *testing.T) {
	refType, ref := dependency.ParseGitReference("1.2.3")
	if refType != dependency.GitRefTag {
		t.Errorf("expected TAG, got %s", refType)
	}
	if ref != "1.2.3" {
		t.Errorf("expected '1.2.3', got %s", ref)
	}
}

// TestParityGitRefTypeString validates string representations
func TestParityGitRefTypeString(t *testing.T) {
	cases := []struct {
		refType dependency.GitReferenceType
		want    string
	}{
		{dependency.GitRefBranch, "branch"},
		{dependency.GitRefTag, "tag"},
		{dependency.GitRefCommit, "commit"},
	}
	for _, c := range cases {
		if got := c.refType.String(); got != c.want {
			t.Errorf("GitReferenceType.String() = %s, want %s", got, c.want)
		}
	}
}

// TestParityVirtualPackageTypeString mirrors VirtualPackageType string values
func TestParityVirtualPackageTypeString(t *testing.T) {
	if dependency.VirtualPackageFile.String() != "file" {
		t.Errorf("expected 'file', got %s", dependency.VirtualPackageFile.String())
	}
	if dependency.VirtualPackageSubdirectory.String() != "subdirectory" {
		t.Errorf("expected 'subdirectory', got %s", dependency.VirtualPackageSubdirectory.String())
	}
}

// TestParityResolvedReferenceString mirrors ResolvedReference.__str__
func TestParityResolvedReferenceString(t *testing.T) {
	// No resolved commit: just refname
	r := dependency.ResolvedReference{RefName: "main", RefType: dependency.GitRefBranch}
	if r.String() != "main" {
		t.Errorf("expected 'main', got %s", r.String())
	}

	// Commit type: short SHA
	r2 := dependency.ResolvedReference{
		RefType:        dependency.GitRefCommit,
		ResolvedCommit: "abc1234def567890",
		RefName:        "abc1234",
	}
	if r2.String() != "abc1234d" {
		t.Errorf("expected 'abc1234d', got %s", r2.String())
	}

	// Branch with commit: "main (abc1234de)"
	r3 := dependency.ResolvedReference{
		RefType:        dependency.GitRefBranch,
		RefName:        "main",
		ResolvedCommit: "abc1234def567890",
	}
	if r3.String() != "main (abc1234d)" {
		t.Errorf("expected 'main (abc1234d)', got %s", r3.String())
	}
}
