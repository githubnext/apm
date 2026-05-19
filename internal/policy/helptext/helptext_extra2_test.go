package helptext_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/policy/helptext"
)

func TestPolicySourceFormsHelp_MentionsLocalFile(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "local") && !strings.Contains(h, "file") && !strings.Contains(h, "path") {
		t.Error("PolicySourceFormsHelp should mention local file path option")
	}
}

func TestPolicySourceFormsHelp_IsNotBlank(t *testing.T) {
	h := strings.TrimSpace(helptext.PolicySourceFormsHelp)
	if h == "" {
		t.Error("PolicySourceFormsHelp should not be blank")
	}
}

func TestPolicySourceFormsHelp_NoControlChars(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	for i, r := range h {
		if r < 32 && r != '\n' && r != '\t' {
			t.Errorf("unexpected control character at position %d: %q", i, r)
		}
	}
}

func TestPolicySourceFormsHelp_MentionsHttps(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "https://") {
		t.Error("PolicySourceFormsHelp should include an https:// example")
	}
}

func TestPolicySourceFormsHelp_MentionsOwnerRepoSlash(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "/") {
		t.Error("PolicySourceFormsHelp should reference owner/repo format with slash")
	}
}

func TestPolicySourceFormsHelp_LongerThan20Chars(t *testing.T) {
	if len(helptext.PolicySourceFormsHelp) <= 20 {
		t.Error("PolicySourceFormsHelp should be a meaningful help string, not a stub")
	}
}

func TestPolicySourceFormsHelp_ContainsOrgKeyword(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "org") && !strings.Contains(h, "Org") && !strings.Contains(h, "ORG") {
		t.Error("PolicySourceFormsHelp should mention 'org' auto-discovery")
	}
}
