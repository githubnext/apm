package deptypes

import "testing"

func TestParseGitReference_Empty(t *testing.T) {
	refType, name := ParseGitReference("")
	if refType != GitRefBranch {
		t.Errorf("empty ref: got type %d, want GitRefBranch", refType)
	}
	if name != "main" {
		t.Errorf("empty ref: got name %q, want %q", name, "main")
	}
}

func TestParseGitReference_Commit(t *testing.T) {
	commits := []string{"abc1234", "deadbeef1234567", "a1b2c3d4e5f6a7b8"}
	for _, c := range commits {
		refType, name := ParseGitReference(c)
		if refType != GitRefCommit {
			t.Errorf("ParseGitReference(%q): got type %d, want GitRefCommit", c, refType)
		}
		if name != c {
			t.Errorf("ParseGitReference(%q): got name %q, want %q", c, name, c)
		}
	}
}

func TestParseGitReference_Tag(t *testing.T) {
	tags := []string{"v1.2.3", "1.0.0", "v2.0.0-beta"}
	for _, tag := range tags {
		refType, name := ParseGitReference(tag)
		if refType != GitRefTag {
			t.Errorf("ParseGitReference(%q): got type %d, want GitRefTag", tag, refType)
		}
		if name != tag {
			t.Errorf("ParseGitReference(%q): got name %q, want %q", tag, name, tag)
		}
	}
}

func TestParseGitReference_Branch(t *testing.T) {
	branches := []string{"main", "feature/my-branch", "develop"}
	for _, b := range branches {
		refType, name := ParseGitReference(b)
		if refType != GitRefBranch {
			t.Errorf("ParseGitReference(%q): got type %d, want GitRefBranch", b, refType)
		}
		if name != b {
			t.Errorf("ParseGitReference(%q): got name %q, want %q", b, name, b)
		}
	}
}

func TestRemoteRefStruct(t *testing.T) {
	r := RemoteRef{Name: "main", RefType: GitRefBranch, CommitSHA: "abc1234"}
	if r.Name != "main" || r.RefType != GitRefBranch || r.CommitSHA != "abc1234" {
		t.Error("RemoteRef fields not set correctly")
	}
}

func TestResolvedReferenceStruct(t *testing.T) {
	rr := ResolvedReference{
		OriginalRef:    "v1.0.0",
		RefType:        GitRefTag,
		ResolvedCommit: "abc1234",
		RefName:        "v1.0.0",
	}
	if rr.OriginalRef != "v1.0.0" || rr.RefType != GitRefTag {
		t.Error("ResolvedReference fields not set correctly")
	}
}

func TestGitRefType_constants(t *testing.T) {
if GitRefBranch == GitRefTag {
t.Error("GitRefBranch must differ from GitRefTag")
}
if GitRefBranch == GitRefCommit {
t.Error("GitRefBranch must differ from GitRefCommit")
}
if GitRefTag == GitRefCommit {
t.Error("GitRefTag must differ from GitRefCommit")
}
}

func TestParseGitReference_shortHex_isBranch(t *testing.T) {
// 5-char hex is too short to be a commit; should be treated as a branch name
refType, name := ParseGitReference("abcde")
if refType != GitRefBranch {
t.Errorf("5-char hex: expected GitRefBranch, got %d", refType)
}
if name != "abcde" {
t.Errorf("expected name=abcde, got %q", name)
}
}

func TestParseGitReference_40charHex(t *testing.T) {
sha := "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0"
refType, _ := ParseGitReference(sha)
if refType != GitRefCommit {
t.Errorf("40-char hex: expected GitRefCommit, got %d", refType)
}
}

func TestParseGitReference_semverVariants(t *testing.T) {
cases := []string{"v1.0.0", "2.3.4", "v10.20.30-alpha.1", "1.0.0-rc.1"}
for _, c := range cases {
refType, name := ParseGitReference(c)
if refType != GitRefTag {
t.Errorf("ParseGitReference(%q): expected GitRefTag, got %d", c, refType)
}
if name != c {
t.Errorf("ParseGitReference(%q): name mismatch %q", c, name)
}
}
}

func TestRemoteRef_zeroValue(t *testing.T) {
var r RemoteRef
if r.Name != "" || r.CommitSHA != "" {
t.Error("zero-value RemoteRef should have empty fields")
}
}

func TestResolvedReference_zeroValue(t *testing.T) {
var rr ResolvedReference
if rr.OriginalRef != "" || rr.ResolvedCommit != "" {
t.Error("zero-value fields should be empty")
}
}

func TestVirtualPackageType_constants(t *testing.T) {
if VirtualPackageFile == VirtualPackageSubdirectory {
t.Error("VirtualPackageFile must differ from VirtualPackageSubdirectory")
}
}
