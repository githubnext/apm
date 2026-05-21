package helptext_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/policy/helptext"
)

func TestPolicySourceFormsHelp_NoLeadingWhitespace_v4(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if len(h) > 0 && (h[0] == ' ' || h[0] == '\t' || h[0] == '\n') {
		t.Errorf("help text starts with whitespace: %q", h[:10])
	}
}

func TestPolicySourceFormsHelp_ContainsOrgKeyword_v4(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "org") {
		t.Error("expected help text to mention 'org'")
	}
}

func TestPolicySourceFormsHelp_ContainsURLFormat_v4(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "https://") {
		t.Error("expected help text to mention 'https://'")
	}
}

func TestPolicySourceFormsHelp_ContainsOwnerRepo_v4(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(h, "owner/repo") {
		t.Error("expected help text to mention 'owner/repo'")
	}
}

func TestPolicySourceFormsHelp_ContainsLocalPath_v4(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if !strings.Contains(strings.ToLower(h), "local") && !strings.Contains(h, "path") {
		t.Error("expected help text to mention local file path")
	}
}

func TestPolicySourceFormsHelp_LengthReasonable_v4(t *testing.T) {
	h := helptext.PolicySourceFormsHelp
	if len(h) < 50 {
		t.Errorf("help text too short: %d chars", len(h))
	}
	if len(h) > 2000 {
		t.Errorf("help text too long: %d chars", len(h))
	}
}
