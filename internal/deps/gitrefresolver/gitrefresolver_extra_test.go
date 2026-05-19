package gitrefresolver

import (
	"testing"
)

func TestIsFullSHA_EdgeCases(t *testing.T) {
	// All zeros is valid full SHA
	if !IsFullSHA("0000000000000000000000000000000000000000") {
		t.Error("all-zeros 40-char SHA should be valid")
	}
	// All f's
	if !IsFullSHA("ffffffffffffffffffffffffffffffffffffffff") {
		t.Error("all-f 40-char SHA should be valid")
	}
}

func TestIsShortSHA_BoundaryLengths(t *testing.T) {
	// exactly 7 chars
	if !IsShortSHA("abcdef1") {
		t.Error("7-char hex should be valid short SHA")
	}
	// 6 chars is too short
	if IsShortSHA("abcde1") {
		t.Error("6-char hex should NOT be valid short SHA")
	}
}

func TestIsShortSHA_MixedCase(t *testing.T) {
	if IsShortSHA("ABCDEF1") {
		t.Error("uppercase hex should not be valid short SHA")
	}
}

func TestNew_FieldsSet(t *testing.T) {
	r := New("github.com", "mytoken123")
	if r.Host != "github.com" {
		t.Errorf("Host = %q, want github.com", r.Host)
	}
	if r.AuthToken != "mytoken123" {
		t.Errorf("AuthToken = %q, want mytoken123", r.AuthToken)
	}
	if r.Timeout == 0 {
		t.Error("Timeout should be non-zero")
	}
}

func TestNew_EmptyFields(t *testing.T) {
	r := New("", "")
	if r == nil {
		t.Fatal("New should return non-nil resolver")
	}
	if r.Host != "" {
		t.Errorf("expected empty host, got %q", r.Host)
	}
}

func TestReferenceTypeConstants(t *testing.T) {
	// iota order: Branch, Tag, Commit, Unknown
	if ReferenceTypeBranch == ReferenceTypeTag {
		t.Error("Branch and Tag should be distinct")
	}
	if ReferenceTypeTag == ReferenceTypeCommit {
		t.Error("Tag and Commit should be distinct")
	}
	if ReferenceTypeCommit == ReferenceTypeUnknown {
		t.Error("Commit and Unknown should be distinct")
	}
}

func TestIsFullSHA_NotHex(t *testing.T) {
	cases := []string{
		"zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", // non-hex
		"abcdef123456789012345678901234567890xyz1", // contains xyz
	}
	for _, c := range cases {
		if IsFullSHA(c) {
			t.Errorf("IsFullSHA(%q) should be false", c)
		}
	}
}

func TestIsShortSHA_ContainsNonHex(t *testing.T) {
	if IsShortSHA("abcdefg") {
		t.Error("'g' is not hex, should return false")
	}
}
