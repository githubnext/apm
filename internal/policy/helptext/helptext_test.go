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
