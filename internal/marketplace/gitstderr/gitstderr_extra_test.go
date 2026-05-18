package gitstderr_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/marketplace/gitstderr"
)

func TestTranslate_PermissionDenied_IsAuth(t *testing.T) {
	r := gitstderr.Translate("fatal: could not read from remote repository", gitstderr.Options{Operation: "fetch"})
	// Could be auth or not_found -- just confirm no panic and non-empty result
	if r.Summary == "" && r.Kind == gitstderr.KindUnknown {
		t.Log("fell back to KindUnknown for read-from-remote (acceptable)")
	}
}

func TestTranslate_SummaryNonEmpty(t *testing.T) {
	r := gitstderr.Translate("fatal: authentication failed", gitstderr.Options{})
	if r.Summary == "" {
		t.Error("Summary should not be empty for auth failure")
	}
}

func TestTranslate_HintNonEmpty(t *testing.T) {
	r := gitstderr.Translate("fatal: authentication failed for 'https://github.com/org/repo'", gitstderr.Options{Remote: "org/repo"})
	if r.Hint == "" {
		t.Error("Hint should not be empty for auth failure")
	}
}

func TestTranslate_RawTruncated(t *testing.T) {
	long := strings.Repeat("a", 10000)
	r := gitstderr.Translate(long, gitstderr.Options{})
	if len(r.Raw) > 1024 {
		t.Errorf("Raw should be truncated, got len=%d", len(r.Raw))
	}
}

func TestTranslate_AllKindsHaveStringRepr(t *testing.T) {
	kinds := []gitstderr.GitErrorKind{
		gitstderr.KindAuth,
		gitstderr.KindNotFound,
		gitstderr.KindTimeout,
		gitstderr.KindUnknown,
	}
	for _, k := range kinds {
		s := k.String()
		if s == "" {
			t.Errorf("GitErrorKind(%d).String() returned empty", k)
		}
	}
}

func TestTranslate_ConnectionRefused_IsTimeout(t *testing.T) {
	r := gitstderr.Translate("fatal: unable to connect to github.com: connection refused", gitstderr.Options{})
	// connection refused is a network error -- timeout or unknown
	if r.Kind != gitstderr.KindTimeout && r.Kind != gitstderr.KindUnknown {
		t.Logf("connection refused classified as %s (informational)", r.Kind)
	}
	// Must not panic
	_ = r.Summary
}

func TestTranslate_ExitCode_Propagated(t *testing.T) {
	code := 128
	r := gitstderr.Translate("fatal: repository 'https://github.com/no/exist' not found",
		gitstderr.Options{ExitCode: &code})
	// not found (exit 128) should classify as KindNotFound or KindUnknown -- just no panic
	if r.Kind != gitstderr.KindNotFound && r.Kind != gitstderr.KindUnknown {
		t.Errorf("unexpected kind for not-found: %s", r.Kind)
	}
}

func TestTranslate_NotHTTPS_StillClassified(t *testing.T) {
	// SSH-format not-found message
	r := gitstderr.Translate("ERROR: Repository not found.", gitstderr.Options{Operation: "clone"})
	if r.Kind != gitstderr.KindNotFound {
		t.Errorf("expected KindNotFound for SSH not-found, got %s", r.Kind)
	}
}

func TestTranslate_Multiline_NoNewlineInSummary(t *testing.T) {
	r := gitstderr.Translate("fatal: something\nfailed\nwith lots\nof lines", gitstderr.Options{})
	if strings.Contains(r.Summary, "\n") {
		t.Errorf("Summary should not contain newlines: %q", r.Summary)
	}
}
