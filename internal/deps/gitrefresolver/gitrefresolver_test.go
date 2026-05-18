package gitrefresolver

import (
	"testing"
)

func TestIsFullSHA(t *testing.T) {
	tests := []struct {
		ref  string
		want bool
	}{
		{"abcdef1234567890abcdef1234567890abcdef12", true},
		{"0000000000000000000000000000000000000000", true},
		{"abc123", false},
		{"", false},
		{"abcdef1234567890abcdef1234567890abcdef1", false},  // 39 chars
		{"abcdef1234567890abcdef1234567890abcdef123", false}, // 41 chars
		{"ABCDEF1234567890abcdef1234567890abcdef12", false},  // uppercase
		{"main", false},
		{"v1.2.3", false},
	}

	for _, tc := range tests {
		t.Run(tc.ref, func(t *testing.T) {
			got := IsFullSHA(tc.ref)
			if got != tc.want {
				t.Errorf("IsFullSHA(%q) = %v, want %v", tc.ref, got, tc.want)
			}
		})
	}
}

func TestIsShortSHA(t *testing.T) {
	tests := []struct {
		ref  string
		want bool
	}{
		{"abcdef1", true},
		{"abcdef1234567", true},
		{"abcdef1234567890abcdef1234567890abcdef12", true}, // full sha also matches
		{"abc", false},                                     // too short (< 7)
		{"", false},
		{"main", false},
		{"ABCDEF1", false}, // uppercase not hex
		{"1234567", true},
		{"123456g", false}, // non-hex char
	}

	for _, tc := range tests {
		t.Run(tc.ref, func(t *testing.T) {
			got := IsShortSHA(tc.ref)
			if got != tc.want {
				t.Errorf("IsShortSHA(%q) = %v, want %v", tc.ref, got, tc.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	r := New("github.com", "mytoken")
	if r.Host != "github.com" {
		t.Errorf("expected Host=github.com, got %q", r.Host)
	}
	if r.AuthToken != "mytoken" {
		t.Errorf("expected AuthToken=mytoken, got %q", r.AuthToken)
	}
	if r.Timeout == 0 {
		t.Error("expected non-zero timeout")
	}
}

func TestGitReferenceTypeConstants(t *testing.T) {
if ReferenceTypeBranch != 0 {
t.Errorf("ReferenceTypeBranch should be 0, got %d", ReferenceTypeBranch)
}
if ReferenceTypeTag == ReferenceTypeBranch {
t.Error("ReferenceTypeTag and ReferenceTypeBranch should differ")
}
if ReferenceTypeCommit == ReferenceTypeTag {
t.Error("ReferenceTypeCommit and ReferenceTypeTag should differ")
}
if ReferenceTypeUnknown == ReferenceTypeCommit {
t.Error("ReferenceTypeUnknown and ReferenceTypeCommit should differ")
}
}

func TestRemoteRef_Fields(t *testing.T) {
r := RemoteRef{
Name:     "refs/heads/main",
SHA:      "abcdef1234567890abcdef1234567890abcdef12",
IsTag:    false,
IsBranch: true,
}
if r.Name != "refs/heads/main" {
t.Errorf("unexpected Name: %q", r.Name)
}
if !r.IsBranch {
t.Error("IsBranch should be true")
}
if r.IsTag {
t.Error("IsTag should be false")
}
if !IsFullSHA(r.SHA) {
t.Error("SHA should be a valid full SHA")
}
}

func TestResolvedReference_Fields(t *testing.T) {
rr := ResolvedReference{
SHA:     "abcdef1234567890abcdef1234567890abcdef12",
RefType: ReferenceTypeBranch,
Ref:     "main",
}
if rr.Ref != "main" {
t.Errorf("unexpected Ref: %q", rr.Ref)
}
if rr.RefType != ReferenceTypeBranch {
t.Errorf("unexpected RefType: %d", rr.RefType)
}
}

func TestGitHubAPIResult_Fields(t *testing.T) {
r := GitHubAPIResult{SHA: "abcdef1234567890abcdef1234567890abcdef12"}
if r.SHA == "" {
t.Error("SHA should not be empty")
}
if !IsFullSHA(r.SHA) {
t.Error("SHA should be a valid full SHA")
}
}

func TestNew_DefaultTimeout(t *testing.T) {
r := New("ghe.example.com", "token")
if r.Timeout <= 0 {
t.Error("expected positive default timeout")
}
if r.Host != "ghe.example.com" {
t.Errorf("expected Host=ghe.example.com, got %q", r.Host)
}
}

func TestIsFullSHA_AllHexChars(t *testing.T) {
// All valid hex chars
sha := "0123456789abcdef01234567890123456789abcd"
if !IsFullSHA(sha) {
t.Errorf("expected true for valid hex SHA, got false")
}
}

func TestIsShortSHA_ExactlySevenChars(t *testing.T) {
if !IsShortSHA("abcdef1") {
t.Error("7-char hex string should be short SHA")
}
if IsShortSHA("abcde1") {
t.Error("6-char string should not be short SHA")
}
}
