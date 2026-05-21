package gitremoteops

import (
	"testing"
)

func TestParseLsRemoteOutput_MixedTagsAndBranches(t *testing.T) {
	input := "abc123\trefs/heads/main\ndef456\trefs/tags/v1.0.0\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_EmptyInput(t *testing.T) {
	refs := ParseLsRemoteOutput("")
	if len(refs) != 0 {
		t.Errorf("expected 0 refs for empty input, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_DerefTakesPrecedence(t *testing.T) {
	// deref sha should override annotated tag sha
	input := "aaa111\trefs/tags/v2.0.0\nbbb222\trefs/tags/v2.0.0^{}\n"
	refs := ParseLsRemoteOutput(input)
	// should have one tag entry
	var tags []RemoteRef
	for _, r := range refs {
		if r.RefType == GitRefTag && r.Name == "v2.0.0" {
			tags = append(tags, r)
		}
	}
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag ref for v2.0.0, got %d", len(tags))
	}
	if tags[0].CommitSHA != "bbb222" {
		t.Errorf("expected deref sha bbb222, got %q", tags[0].CommitSHA)
	}
}

func TestParseLsRemoteOutput_MalformedLineSkipped(t *testing.T) {
	input := "notatabseparatedline\nabc123\trefs/heads/main\n"
	refs := ParseLsRemoteOutput(input)
	// only the valid line should be parsed
	if len(refs) != 1 {
		t.Errorf("expected 1 ref, got %d", len(refs))
	}
}

func TestSortRefsBySemver_SemverBeforeNonSemver(t *testing.T) {
	refs := []RemoteRef{
		{Name: "nightly", RefType: GitRefTag},
		{Name: "v1.0.0", RefType: GitRefTag},
	}
	sorted := SortRefsBySemver(refs)
	if sorted[0].Name != "v1.0.0" {
		t.Errorf("expected semver tag first, got %q", sorted[0].Name)
	}
}

func TestSortRefsBySemver_DoesNotMutateOriginal(t *testing.T) {
	refs := []RemoteRef{
		{Name: "nightly", RefType: GitRefTag},
		{Name: "v2.0.0", RefType: GitRefTag},
	}
	origFirst := refs[0].Name
	_ = SortRefsBySemver(refs)
	if refs[0].Name != origFirst {
		t.Error("SortRefsBySemver should not mutate original slice")
	}
}

func TestParseLsRemoteOutput_BranchRefType(t *testing.T) {
	input := "abc123\trefs/heads/feature/x\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].RefType != GitRefBranch {
		t.Errorf("expected GitRefBranch, got %d", refs[0].RefType)
	}
	if refs[0].Name != "feature/x" {
		t.Errorf("expected 'feature/x', got %q", refs[0].Name)
	}
}

func TestParseLsRemoteOutput_TagRefType(t *testing.T) {
	input := "sha999\trefs/tags/v0.1.2\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].RefType != GitRefTag {
		t.Errorf("expected GitRefTag, got %d", refs[0].RefType)
	}
}
