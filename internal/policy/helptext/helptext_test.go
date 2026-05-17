package helptext_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/policy/helptext"
)

func TestPolicySourceFormsHelp_NotEmpty(t *testing.T) {
	if helptext.PolicySourceFormsHelp == "" {
		t.Fatal("PolicySourceFormsHelp must not be empty")
	}
}

func TestPolicySourceFormsHelp_ContainsKeywords(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	for _, kw := range []string{"org", "owner/repo", "https://", "local"} {
		if !strings.Contains(h, kw) {
			t.Errorf("PolicySourceFormsHelp missing expected keyword %q", kw)
		}
	}
}

func TestPolicySourceFormsHelp_ContainsAcceptsWord(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "Accepts") && !strings.Contains(h, "accepts") {
		t.Error("PolicySourceFormsHelp should mention accepted formats")
	}
}

func TestPolicySourceFormsHelp_ContainsHttpsScheme(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "https://") {
		t.Error("PolicySourceFormsHelp should mention https:// URL format")
	}
}

func TestPolicySourceFormsHelp_ContainsOrgAutoDiscovery(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "org") {
		t.Error("PolicySourceFormsHelp should mention org auto-discovery")
	}
}

func TestPolicySourceFormsHelp_ContainsOwnerRepoForm(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "owner/repo") {
		t.Error("PolicySourceFormsHelp should mention owner/repo form")
	}
}

func TestPolicySourceFormsHelp_ContainsLocalPath(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "local") && !strings.Contains(h, "file") {
		t.Error("PolicySourceFormsHelp should mention local file path option")
	}
}

func TestPolicySourceFormsHelp_IsASCII(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	for i, ch := range h {
		if ch > 127 {
			t.Errorf("PolicySourceFormsHelp contains non-ASCII character %q at position %d", ch, i)
		}
	}
}

func TestPolicySourceFormsHelp_ReasonableLength(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if len(h) < 50 {
		t.Errorf("PolicySourceFormsHelp too short (%d chars); expected at least 50", len(h))
	}
	if len(h) > 1000 {
		t.Errorf("PolicySourceFormsHelp too long (%d chars); expected at most 1000", len(h))
	}
}

func TestPolicySourceFormsHelp_DoesNotContainHTML(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if strings.Contains(h, "<html") || strings.Contains(h, "<br>") {
		t.Error("PolicySourceFormsHelp should not contain HTML markup")
	}
}

func TestPolicySourceFormsHelp_Stable(t *testing.T) {
	// Calling the constant twice returns the same value (it is constant).
	h1 := helptext.PolicySourceFormsHelp
	h2 := helptext.PolicySourceFormsHelp
	if h1 != h2 {
		t.Error("PolicySourceFormsHelp is not stable across accesses")
	}
}
