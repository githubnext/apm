package gitrefresolver

import "testing"

func TestIsFullSHA_ExactLength_Extra4(t *testing.T) {
sha := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
if !IsFullSHA(sha) {
t.Errorf("expected full SHA recognized: %s", sha)
}
}

func TestIsFullSHA_TooLong_Extra4(t *testing.T) {
sha := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2ab"
if IsFullSHA(sha) {
t.Error("expected too-long SHA not recognized")
}
}

func TestIsShortSHA_EightChars_Extra4(t *testing.T) {
sha := "a1b2c3d4"
if !IsShortSHA(sha) {
t.Errorf("expected short SHA recognized: %s", sha)
}
}

func TestIsShortSHA_FifteenChars_Extra4(t *testing.T) {
sha := "a1b2c3d4e5f6a1b"
if !IsShortSHA(sha) {
t.Errorf("expected short SHA recognized: %s", sha)
}
}

func TestNew_NonNil_Extra4(t *testing.T) {
r := New("github.com", "mytoken")
if r == nil {
t.Fatal("expected non-nil resolver")
}
}

func TestNew_EmptyToken_Extra4(t *testing.T) {
r := New("github.com", "")
if r == nil {
t.Fatal("expected non-nil resolver even with empty token")
}
}

func TestRemoteRef_IsBranch_Extra4(t *testing.T) {
ref := RemoteRef{IsBranch: true}
if !ref.IsBranch {
t.Error("expected IsBranch true")
}
}

func TestRemoteRef_IsTag_Extra4(t *testing.T) {
ref := RemoteRef{IsTag: true}
if !ref.IsTag {
t.Error("expected IsTag true")
}
}

func TestRemoteRef_SHAField_Extra4(t *testing.T) {
ref := RemoteRef{SHA: "abc123def456abc123def456abc123def456abc1"}
if ref.SHA != "abc123def456abc123def456abc123def456abc1" {
t.Errorf("unexpected SHA: %s", ref.SHA)
}
}

func TestGitHubAPIResult_SHAField_Extra4(t *testing.T) {
r := GitHubAPIResult{SHA: "deadbeef1234abc0deadbeef1234abc0deadbeef"}
if r.SHA != "deadbeef1234abc0deadbeef1234abc0deadbeef" {
t.Errorf("unexpected SHA: %s", r.SHA)
}
}

func TestResolvedReference_RefType_Extra4(t *testing.T) {
r := ResolvedReference{RefType: ReferenceTypeBranch}
if r.RefType != ReferenceTypeBranch {
t.Errorf("unexpected ref type: %v", r.RefType)
}
}

func TestReferenceTypeDistinct_Extra4(t *testing.T) {
if ReferenceTypeBranch == ReferenceTypeTag {
t.Error("branch and tag constants should be distinct")
}
if ReferenceTypeCommit == ReferenceTypeUnknown {
t.Error("commit and unknown constants should be distinct")
}
}

func TestParseLsRemoteOutput_WithBranch_Extra4(t *testing.T) {
raw := "abc1234567890abcdef1234567890abcdef123456\trefs/heads/main\n"
refs := ParseLsRemoteOutput(raw)
if len(refs) != 1 {
t.Fatalf("expected 1 ref, got %d", len(refs))
}
if !refs[0].IsBranch {
t.Error("expected IsBranch true for refs/heads/main")
}
}
