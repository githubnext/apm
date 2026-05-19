package gitremoteops

import (
	"testing"
)

func TestParseLsRemoteOutput_WhitespaceOnly(t *testing.T) {
	refs := ParseLsRemoteOutput("   \n\t\n")
	if len(refs) != 0 {
		t.Errorf("expected 0 refs for whitespace-only input, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_TagWithoutDeref(t *testing.T) {
	input := "abc123\trefs/tags/v2.0.0\n"
	refs := ParseLsRemoteOutput(input)
	found := false
	for _, r := range refs {
		if r.RefType == GitRefTag && r.Name == "v2.0.0" && r.CommitSHA == "abc123" {
			found = true
		}
	}
	if !found {
		t.Error("expected v2.0.0 tag with sha abc123")
	}
}

func TestParseLsRemoteOutput_DerefOverridesAnnotated(t *testing.T) {
	// annotated sha first, then ^{} sha
	input := "tag111\trefs/tags/v3.0.0\ncommit222\trefs/tags/v3.0.0^{}\n"
	refs := ParseLsRemoteOutput(input)
	for _, r := range refs {
		if r.RefType == GitRefTag && r.Name == "v3.0.0" {
			if r.CommitSHA != "commit222" {
				t.Errorf("expected dereferenced sha commit222, got %s", r.CommitSHA)
			}
			return
		}
	}
	t.Error("v3.0.0 tag not found")
}

func TestParseLsRemoteOutput_BranchWithSlash(t *testing.T) {
	input := "sha111\trefs/heads/feature/my-feature\n"
	refs := ParseLsRemoteOutput(input)
	found := false
	for _, r := range refs {
		if r.RefType == GitRefBranch && r.Name == "feature/my-feature" {
			found = true
		}
	}
	if !found {
		t.Error("expected branch feature/my-feature")
	}
}

func TestParseLsRemoteOutput_OnlyTags(t *testing.T) {
	input := "sha1\trefs/tags/v0.1.0\nsha2\trefs/tags/v0.2.0\n"
	refs := ParseLsRemoteOutput(input)
	for _, r := range refs {
		if r.RefType == GitRefBranch {
			t.Errorf("unexpected branch in tag-only input: %s", r.Name)
		}
	}
	if len(refs) != 2 {
		t.Errorf("expected 2 tags, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_OnlyBranches(t *testing.T) {
	input := "sha1\trefs/heads/main\nsha2\trefs/heads/dev\n"
	refs := ParseLsRemoteOutput(input)
	for _, r := range refs {
		if r.RefType == GitRefTag {
			t.Errorf("unexpected tag in branch-only input: %s", r.Name)
		}
	}
	if len(refs) != 2 {
		t.Errorf("expected 2 branches, got %d", len(refs))
	}
}

func TestSortRefsBySemver_Empty(t *testing.T) {
	sorted := SortRefsBySemver(nil)
	if len(sorted) != 0 {
		t.Errorf("expected empty, got %d", len(sorted))
	}
}

func TestSortRefsBySemver_Single(t *testing.T) {
	refs := []RemoteRef{{Name: "v1.0.0", RefType: GitRefTag}}
	sorted := SortRefsBySemver(refs)
	if len(sorted) != 1 || sorted[0].Name != "v1.0.0" {
		t.Errorf("single ref sort failed: %v", sorted)
	}
}

func TestSortRefsBySemver_NonSemverLast(t *testing.T) {
	refs := []RemoteRef{
		{Name: "stable", RefType: GitRefTag},
		{Name: "v1.0.0", RefType: GitRefTag},
		{Name: "nightly", RefType: GitRefTag},
	}
	sorted := SortRefsBySemver(refs)
	if sorted[0].Name != "v1.0.0" {
		t.Errorf("expected semver tag first, got %s", sorted[0].Name)
	}
}

func TestSortRefsBySemver_MultipleNonSemver(t *testing.T) {
	refs := []RemoteRef{
		{Name: "alpha", RefType: GitRefTag},
		{Name: "beta", RefType: GitRefTag},
		{Name: "v1.2.3", RefType: GitRefTag},
	}
	sorted := SortRefsBySemver(refs)
	if sorted[0].Name != "v1.2.3" {
		t.Errorf("semver tag should be first: %s", sorted[0].Name)
	}
}

func TestSortRefsBySemver_AllNonSemver(t *testing.T) {
	refs := []RemoteRef{
		{Name: "beta", RefType: GitRefTag},
		{Name: "alpha", RefType: GitRefTag},
	}
	sorted := SortRefsBySemver(refs)
	if len(sorted) != 2 {
		t.Errorf("expected 2 refs, got %d", len(sorted))
	}
}

func TestSortRefsBySemver_SemverDescending(t *testing.T) {
	refs := []RemoteRef{
		{Name: "v1.0.0", RefType: GitRefTag},
		{Name: "v3.0.0", RefType: GitRefTag},
		{Name: "v2.0.0", RefType: GitRefTag},
	}
	sorted := SortRefsBySemver(refs)
	if sorted[0].Name != "v3.0.0" {
		t.Errorf("expected v3.0.0 first (descending), got %s", sorted[0].Name)
	}
	if sorted[len(sorted)-1].Name != "v1.0.0" {
		t.Errorf("expected v1.0.0 last, got %s", sorted[len(sorted)-1].Name)
	}
}

func TestRemoteRef_Fields(t *testing.T) {
	r := RemoteRef{Name: "main", RefType: GitRefBranch, CommitSHA: "abc123"}
	if r.Name != "main" {
		t.Errorf("Name = %q, want main", r.Name)
	}
	if r.RefType != GitRefBranch {
		t.Errorf("RefType = %v, want GitRefBranch", r.RefType)
	}
	if r.CommitSHA != "abc123" {
		t.Errorf("CommitSHA = %q, want abc123", r.CommitSHA)
	}
}

func TestGitRefType_Constants(t *testing.T) {
	if GitRefBranch == GitRefTag {
		t.Error("GitRefBranch and GitRefTag must be distinct")
	}
}
