package gitremoteops

import (
	"testing"
)

func TestRemoteRef_BranchNameField(t *testing.T) {
	r := RemoteRef{Name: "main", RefType: GitRefBranch, CommitSHA: "abc"}
	if r.Name != "main" {
		t.Errorf("expected name=main")
	}
}

func TestRemoteRef_TagNameField(t *testing.T) {
	r := RemoteRef{Name: "v1.0.0", RefType: GitRefTag, CommitSHA: "def"}
	if r.CommitSHA != "def" {
		t.Errorf("expected CommitSHA=def")
	}
}

func TestParseLsRemoteOutput_SingleBranchName(t *testing.T) {
	input := "abc123\trefs/heads/main\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].Name != "main" {
		t.Errorf("expected name=main, got %s", refs[0].Name)
	}
}

func TestParseLsRemoteOutput_SingleTagType(t *testing.T) {
	input := "def456\trefs/tags/v1.0.0\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
	if refs[0].RefType != GitRefTag {
		t.Errorf("expected tag ref type")
	}
}

func TestSortRefsBySemver_EmptySlice(t *testing.T) {
	refs := SortRefsBySemver([]RemoteRef{})
	if len(refs) != 0 {
		t.Errorf("expected empty, got %d", len(refs))
	}
}

func TestSortRefsBySemver_OnlyOne(t *testing.T) {
	refs := SortRefsBySemver([]RemoteRef{{Name: "v1.0.0", RefType: GitRefTag}})
	if len(refs) != 1 {
		t.Errorf("expected 1, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_TwoBranches(t *testing.T) {
	input := "aaa\trefs/heads/main\nbbb\trefs/heads/dev\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 2 {
		t.Errorf("expected 2 refs, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_FeatureBranchSHA(t *testing.T) {
	input := "deadbeef\trefs/heads/feature\n"
	refs := ParseLsRemoteOutput(input)
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref")
	}
	if refs[0].CommitSHA != "deadbeef" {
		t.Errorf("expected deadbeef, got %s", refs[0].CommitSHA)
	}
}

func TestParseLsRemoteOutput_OnlySpaces(t *testing.T) {
	refs := ParseLsRemoteOutput("   \n  ")
	if len(refs) != 0 {
		t.Errorf("expected 0 refs, got %d", len(refs))
	}
}
