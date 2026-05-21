package gitstderr

import (
	"strings"
	"testing"
)

func TestTranslate_ForbiddenURL_IsAuth(t *testing.T) {
	got := Translate("fatal: The requested URL returned error: 403 Forbidden", Options{})
	if got.Kind != KindAuth {
		t.Errorf("expected KindAuth for 403, got %v", got.Kind)
	}
}

func TestTranslate_UnauthorizedURL_IsAuth(t *testing.T) {
	got := Translate("The requested URL returned error: 401 Unauthorized", Options{})
	if got.Kind != KindAuth {
		t.Errorf("expected KindAuth for 401, got %v", got.Kind)
	}
}

func TestTranslate_RepositoryNotFound_IsNotFound(t *testing.T) {
	got := Translate("fatal: repository not found", Options{})
	if got.Kind != KindNotFound {
		t.Errorf("expected KindNotFound, got %v", got.Kind)
	}
}

func TestTranslate_CouldNotResolveRef_IsNotFound(t *testing.T) {
	got := Translate("error: couldn't find remote ref main", Options{})
	if got.Kind != KindNotFound {
		t.Errorf("expected KindNotFound for missing ref, got %v", got.Kind)
	}
}

func TestTranslate_OperationInSummary(t *testing.T) {
	got := Translate("authentication failed", Options{Operation: "clone"})
	if !strings.Contains(got.Summary, "clone") {
		t.Errorf("expected 'clone' in summary: %q", got.Summary)
	}
}

func TestTranslate_RemoteInHint(t *testing.T) {
	got := Translate("repository not found", Options{Remote: "origin"})
	if !strings.Contains(got.Hint, "origin") {
		t.Errorf("expected remote 'origin' in hint: %q", got.Hint)
	}
}

func TestTranslate_UnknownKind_HintHasOperation(t *testing.T) {
	got := Translate("some unknown error", Options{Operation: "fetch"})
	if !strings.Contains(got.Hint, "fetch") {
		t.Errorf("expected 'fetch' in hint for unknown error: %q", got.Hint)
	}
}

func TestTranslate_SummaryBelowMaxLen(t *testing.T) {
	got := Translate("authentication failed", Options{Operation: strings.Repeat("x", 200)})
	if len(got.Summary) > 90 {
		t.Errorf("summary too long: %d chars", len(got.Summary))
	}
}

func TestGitErrorKind_UnknownString(t *testing.T) {
	var k GitErrorKind = 999
	if k.String() != "unknown" {
		t.Errorf("expected 'unknown' for unrecognized kind, got %q", k.String())
	}
}

func TestTranslate_RawPreservedExact(t *testing.T) {
	input := "fatal: authentication failed for some host"
	got := Translate(input, Options{})
	if got.Raw != input {
		t.Errorf("raw should match input when under limit: got %q", got.Raw)
	}
}

func TestTranslate_SSLError_IsTimeout(t *testing.T) {
	got := Translate("OpenSSL SSL_read: Connection reset by peer", Options{})
	if got.Kind != KindTimeout {
		t.Errorf("expected KindTimeout for SSL_read reset, got %v", got.Kind)
	}
}
