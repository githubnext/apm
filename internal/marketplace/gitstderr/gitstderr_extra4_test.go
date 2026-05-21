package gitstderr_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/gitstderr"
)

func TestGitErrorKindString_Auth(t *testing.T) {
	if gitstderr.KindAuth.String() != "auth" {
		t.Errorf("got %q", gitstderr.KindAuth.String())
	}
}

func TestGitErrorKindString_NotFound(t *testing.T) {
	if gitstderr.KindNotFound.String() != "not_found" {
		t.Errorf("got %q", gitstderr.KindNotFound.String())
	}
}

func TestGitErrorKindString_Timeout(t *testing.T) {
	if gitstderr.KindTimeout.String() != "timeout" {
		t.Errorf("got %q", gitstderr.KindTimeout.String())
	}
}

func TestGitErrorKindString_Unknown(t *testing.T) {
	if gitstderr.KindUnknown.String() != "unknown" {
		t.Errorf("got %q", gitstderr.KindUnknown.String())
	}
}

func TestTranslate_EmptyStderr_IsUnknown(t *testing.T) {
	r := gitstderr.Translate("", gitstderr.Options{})
	if r.Kind != gitstderr.KindUnknown {
		t.Errorf("empty stderr: got kind %v, want unknown", r.Kind)
	}
}

func TestTranslate_SummaryNonEmptyForAuth(t *testing.T) {
	r := gitstderr.Translate("authentication failed for https://github.com/x/y", gitstderr.Options{
		Operation: "clone",
		Remote:    "https://github.com/x/y",
	})
	if r.Summary == "" {
		t.Error("expected non-empty summary for auth error")
	}
}

func TestTranslate_RawTruncated_v4(t *testing.T) {
	long := strings.Repeat("x", 2000)
	r := gitstderr.Translate(long, gitstderr.Options{})
	if len(r.Raw) > 600 {
		t.Errorf("expected Raw to be truncated, got len %d", len(r.Raw))
	}
}

func TestTranslate_NetworkTimeout_IsTimeout(t *testing.T) {
	r := gitstderr.Translate("fatal: unable to connect to github.com: connection timed out", gitstderr.Options{})
	if r.Kind != gitstderr.KindTimeout {
		t.Errorf("expected KindTimeout, got %v", r.Kind)
	}
}
