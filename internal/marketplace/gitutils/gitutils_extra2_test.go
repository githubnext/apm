package gitutils

import (
	"strings"
	"testing"
)

func TestRedactToken_NoToken_Unchanged(t *testing.T) {
	input := "https://github.com/org/repo"
	got := RedactToken(input)
	if got != input {
		t.Errorf("no-token URL should be unchanged; got %q", got)
	}
}

func TestRedactToken_HttpScheme(t *testing.T) {
	input := "http://tok@github.com/repo"
	got := RedactToken(input)
	if strings.Contains(got, "tok") {
		t.Errorf("token still visible in http URL: %q", got)
	}
}

func TestRedactToken_QueryParamToken(t *testing.T) {
	input := "https://github.com/repo?token=secretval"
	got := RedactToken(input)
	if strings.Contains(got, "secretval") {
		t.Errorf("query token still visible: %q", got)
	}
	if !strings.Contains(got, "token=***") {
		t.Errorf("expected token=*** in output: %q", got)
	}
}

func TestRedactToken_AmpersandTokenParam(t *testing.T) {
	input := "https://github.com/repo?foo=bar&token=abc123"
	got := RedactToken(input)
	if strings.Contains(got, "abc123") {
		t.Errorf("token still visible: %q", got)
	}
}

func TestRedactToken_EmptyString(t *testing.T) {
	got := RedactToken("")
	if got != "" {
		t.Errorf("empty input should yield empty output, got %q", got)
	}
}

func TestRedactToken_PlainText(t *testing.T) {
	input := "some plain text without a URL"
	got := RedactToken(input)
	if got != input {
		t.Errorf("plain text should be unchanged, got %q", got)
	}
}

func TestRedactToken_RedactedPlaceholder(t *testing.T) {
	input := "https://ghp_token123@github.com/org/repo"
	got := RedactToken(input)
	if !strings.Contains(got, "***") {
		t.Errorf("expected *** placeholder, got %q", got)
	}
}

func TestRedactToken_HostPreserved(t *testing.T) {
	input := "https://tok@github.com/org/repo"
	got := RedactToken(input)
	if !strings.Contains(got, "github.com") {
		t.Errorf("host should be preserved: %q", got)
	}
}
