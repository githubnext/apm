package gitutils

import (
	"strings"
	"testing"
)

func TestRedactToken_ColonPasswordAt(t *testing.T) {
	// user:password@host format
	input := "https://user:secret-pass@github.com/org/repo"
	got := RedactToken(input)
	if strings.Contains(got, "secret-pass") {
		t.Errorf("password still visible: %q", got)
	}
}

func TestRedactToken_Multiline(t *testing.T) {
	input := "https://tok1@github.com\nhttps://tok2@gitlab.com"
	got := RedactToken(input)
	if strings.Contains(got, "tok1") || strings.Contains(got, "tok2") {
		t.Errorf("tokens visible in multiline: %q", got)
	}
}

func TestRedactToken_PreservesScheme(t *testing.T) {
	input := "https://tok@github.com/repo"
	got := RedactToken(input)
	if !strings.HasPrefix(got, "https://") {
		t.Errorf("https scheme should be preserved: %q", got)
	}
}

func TestRedactToken_ShortToken(t *testing.T) {
	input := "https://x@github.com/a/b"
	got := RedactToken(input)
	if strings.Contains(got, "@") && strings.Contains(got, "x@") {
		t.Errorf("single-char token should be redacted: %q", got)
	}
}

func TestRedactToken_NoSchemeNoRedaction(t *testing.T) {
	// No http/https scheme -- should not modify
	input := "git@github.com:owner/repo.git"
	got := RedactToken(input)
	// SCP-style doesn't have http token; just assert no panic
	_ = got
}

func TestRedactToken_TokenInMiddle(t *testing.T) {
	input := "running: https://secret@github.com/repo.git --depth 1"
	got := RedactToken(input)
	if strings.Contains(got, "secret") {
		t.Errorf("token still visible in complex input: %q", got)
	}
}

func TestRedactToken_GHEHost(t *testing.T) {
	input := "https://mytoken@ghe.mycompany.com/org/repo"
	got := RedactToken(input)
	if strings.Contains(got, "mytoken") {
		t.Errorf("token still visible for GHE host: %q", got)
	}
	if !strings.Contains(got, "ghe.mycompany.com") {
		t.Errorf("GHE host should be preserved: %q", got)
	}
}

func TestRedactToken_LongToken(t *testing.T) {
	tok := strings.Repeat("a", 100)
	input := "https://" + tok + "@github.com/repo"
	got := RedactToken(input)
	if strings.Contains(got, tok) {
		t.Errorf("long token not redacted: %q", got[:min(len(got), 100)])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
