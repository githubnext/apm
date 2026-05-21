package gitutils_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/gitutils"
)

func TestRedactToken_EmptyString(t *testing.T) {
	got := gitutils.RedactToken("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestRedactToken_NoURL(t *testing.T) {
	input := "just some plain text without any URL"
	got := gitutils.RedactToken(input)
	if got != input {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestRedactToken_HTTPScheme(t *testing.T) {
	input := "http://mytoken@github.com/org/repo"
	got := gitutils.RedactToken(input)
	if strings.Contains(got, "mytoken") {
		t.Errorf("token not redacted: %q", got)
	}
}

func TestRedactToken_QueryParamAmpersand(t *testing.T) {
	input := "https://api.github.com/repos?foo=bar&token=secretval"
	got := gitutils.RedactToken(input)
	if strings.Contains(got, "secretval") {
		t.Errorf("query token not redacted: %q", got)
	}
	if !strings.Contains(got, "&token=***") {
		t.Errorf("expected &token=***, got %q", got)
	}
}

func TestRedactToken_MultipleURLs(t *testing.T) {
	input := "https://tok1@host1.com/a and https://tok2@host2.com/b"
	got := gitutils.RedactToken(input)
	if strings.Contains(got, "tok1") || strings.Contains(got, "tok2") {
		t.Errorf("not all tokens redacted: %q", got)
	}
}

func TestRedactToken_NoAtSign(t *testing.T) {
	input := "https://github.com/owner/repo"
	got := gitutils.RedactToken(input)
	if got != input {
		t.Errorf("expected unchanged (no token), got %q", got)
	}
}
