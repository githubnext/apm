package gitremoteops

import (
	"testing"
)

func TestParseLsRemoteOutput_Empty(t *testing.T) {
	refs := ParseLsRemoteOutput("")
	if len(refs) != 0 {
		t.Errorf("expected 0 refs, got %d", len(refs))
	}
}

func TestParseLsRemoteOutput_Branches(t *testing.T) {
	input := "abc123\trefs/heads/main\ndef456\trefs/heads/feature/foo\n"
	refs := ParseLsRemoteOutput(input)
	branches := make(map[string]string)
	for _, r := range refs {
		if r.RefType == GitRefBranch {
			branches[r.Name] = r.CommitSHA
		}
	}
	if branches["main"] != "abc123" {
		t.Errorf("expected main=abc123, got %s", branches["main"])
	}
	if branches["feature/foo"] != "def456" {
		t.Errorf("expected feature/foo=def456, got %s", branches["feature/foo"])
	}
}

func TestParseLsRemoteOutput_Tags(t *testing.T) {
	input := "aaa111\trefs/tags/v1.0.0\nbbb222\trefs/tags/v1.0.0^{}\n"
	refs := ParseLsRemoteOutput(input)
	tags := make(map[string]string)
	for _, r := range refs {
		if r.RefType == GitRefTag {
			tags[r.Name] = r.CommitSHA
		}
	}
	// ^{} dereferenced commit should take precedence
	if tags["v1.0.0"] != "bbb222" {
		t.Errorf("expected v1.0.0=bbb222 (dereferenced), got %s", tags["v1.0.0"])
	}
}

func TestParseLsRemoteOutput_MalformedLines(t *testing.T) {
	input := "noop\n\tabc\tabc\trefs/heads/bad\n"
	refs := ParseLsRemoteOutput(input)
	// malformed lines should be skipped gracefully
	_ = refs
}

func TestSortRefsBySemver_SemverFirst(t *testing.T) {
	refs := []RemoteRef{
		{Name: "latest", RefType: GitRefTag},
		{Name: "v2.0.0", RefType: GitRefTag},
		{Name: "v1.0.0", RefType: GitRefTag},
		{Name: "v1.2.3", RefType: GitRefTag},
	}
	sorted := SortRefsBySemver(refs)
	// semver tags should come before non-semver
	for i, r := range sorted {
		if r.Name == "latest" {
			if i < 3 {
				t.Errorf("non-semver 'latest' should sort last, got index %d", i)
			}
		}
	}
	// first element should be highest semver
	if sorted[0].Name != "v2.0.0" {
		t.Errorf("expected v2.0.0 first, got %s", sorted[0].Name)
	}
}

func TestSortRefsBySemver_PreservesOriginal(t *testing.T) {
	original := []RemoteRef{
		{Name: "v1.0.0", RefType: GitRefTag},
		{Name: "v2.0.0", RefType: GitRefTag},
	}
	sorted := SortRefsBySemver(original)
	// original should not be mutated
	if original[0].Name != "v1.0.0" {
		t.Error("original slice should not be modified")
	}
	if sorted[0].Name != "v2.0.0" {
		t.Errorf("expected v2.0.0 first, got %s", sorted[0].Name)
	}
}

func TestParseLsRemoteOutput_Mixed(t *testing.T) {
	input := `abc000	refs/heads/main
bcd111	refs/tags/v1.0.0
cde222	refs/tags/v1.0.0^{}
eff333	refs/heads/develop
`
	refs := ParseLsRemoteOutput(input)
	branchCount, tagCount := 0, 0
	for _, r := range refs {
		if r.RefType == GitRefBranch {
			branchCount++
		} else {
			tagCount++
		}
	}
	if branchCount != 2 {
		t.Errorf("expected 2 branches, got %d", branchCount)
	}
	if tagCount != 1 {
		t.Errorf("expected 1 tag, got %d", tagCount)
	}
}
