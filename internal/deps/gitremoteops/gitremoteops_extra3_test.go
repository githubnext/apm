package gitremoteops

import (
	"testing"
)

func TestParseLsRemoteOutput_SingleBranch(t *testing.T) {
	input := "abc123\trefs/heads/main\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].Name != "main" {
		t.Errorf("expected main, got %q", refs[0].Name)
	}
	if refs[0].RefType != GitRefBranch {
		t.Error("expected branch ref type")
	}
}

func TestParseLsRemoteOutput_SingleTag(t *testing.T) {
	input := "def456\trefs/tags/v1.0.0\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].Name != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %q", refs[0].Name)
	}
	if refs[0].RefType != GitRefTag {
		t.Error("expected tag ref type")
	}
}

func TestParseLsRemoteOutput_CommitSHAStored(t *testing.T) {
	input := "aabbccdd\trefs/heads/feature\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].CommitSHA != "aabbccdd" {
		t.Errorf("expected aabbccdd, got %q", refs[0].CommitSHA)
	}
}

func TestSortRefsBySemver_ThreeSemverVersions(t *testing.T) {
	refs := []RemoteRef{
		{Name: "v1.0.0", RefType: GitRefTag},
		{Name: "v3.0.0", RefType: GitRefTag},
		{Name: "v2.0.0", RefType: GitRefTag},
	}
	sorted := SortRefsBySemver(refs)
	if sorted[0].Name != "v3.0.0" {
		t.Errorf("expected v3.0.0 first, got %q", sorted[0].Name)
	}
}

func TestSortRefsBySemver_WithMinorAndPatch(t *testing.T) {
	refs := []RemoteRef{
		{Name: "v1.0.1", RefType: GitRefTag},
		{Name: "v1.1.0", RefType: GitRefTag},
		{Name: "v1.0.0", RefType: GitRefTag},
	}
	sorted := SortRefsBySemver(refs)
	if sorted[0].Name != "v1.1.0" {
		t.Errorf("expected v1.1.0 first, got %q", sorted[0].Name)
	}
}

func TestRemoteRef_ZeroValue(t *testing.T) {
	var r RemoteRef
	if r.Name != "" {
		t.Error("expected empty Name")
	}
	if r.CommitSHA != "" {
		t.Error("expected empty CommitSHA")
	}
}

func TestGitRefBranchConstantZero(t *testing.T) {
	if GitRefBranch != 0 {
		t.Errorf("expected GitRefBranch=0, got %d", GitRefBranch)
	}
}

func TestGitRefTagConstantOne(t *testing.T) {
	if GitRefTag != 1 {
		t.Errorf("expected GitRefTag=1, got %d", GitRefTag)
	}
}

func TestParseLsRemoteOutput_UnknownRefIgnored(t *testing.T) {
	input := "abc123\trefs/notes/commits\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 0 {
		t.Errorf("expected 0 refs for unknown ref type, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_TablesWithSpaces(t *testing.T) {
	input := "  abc123  \t  refs/heads/spacious  \n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].Name != "spacious" {
		t.Errorf("expected 'spacious', got %q", refs[0].Name)
	}
}

func TestSortRefsBySemver_ReturnsNewSlice(t *testing.T) {
	refs := []RemoteRef{
		{Name: "v1.0.0"},
		{Name: "v2.0.0"},
	}
	sorted := SortRefsBySemver(refs)
	if &sorted[0] == &refs[0] {
		t.Error("expected new slice, not same backing array")
	}
}
