package gitutils

import (
	"strings"
	"testing"
)

func TestRedactToken_multipleTokensInLine(t *testing.T) {
	input := "https://tok1@github.com clone && https://tok2@gitlab.com"
	got := RedactToken(input)
	if strings.Contains(got, "tok1") || strings.Contains(got, "tok2") {
		t.Errorf("tokens still visible: %q", got)
	}
}

func TestRedactToken_plainText(t *testing.T) {
	input := "no tokens here, just plain text"
	got := RedactToken(input)
	if got != input {
		t.Errorf("plain text modified unexpectedly: %q", got)
	}
}

func TestRedactToken_httpsAt(t *testing.T) {
	input := "https://mytoken@github.com/owner/repo.git"
	got := RedactToken(input)
	if got != "https://***@github.com/owner/repo.git" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestRedactToken_httpAt(t *testing.T) {
	input := "http://secret@example.com/repo"
	got := RedactToken(input)
	if got != "https://***@example.com/repo" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestRedactToken_queryParam(t *testing.T) {
	input := "https://api.github.com/repos/a/b?token=abc123&other=val"
	got := RedactToken(input)
	if got != "https://api.github.com/repos/a/b?token=***&other=val" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestRedactToken_ampersandParam(t *testing.T) {
	input := "https://example.com/path?foo=1&token=secret"
	got := RedactToken(input)
	if got != "https://example.com/path?foo=1&token=***" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestRedactToken_noToken(t *testing.T) {
	input := "https://github.com/owner/repo"
	got := RedactToken(input)
	if got != input {
		t.Errorf("unexpected modification: %q", got)
	}
}

func TestRedactToken_empty(t *testing.T) {
	got := RedactToken("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}
