package helptext_test

import (
"strings"
"testing"

"github.com/githubnext/apm/internal/policy/helptext"
)

func TestPolicySourceFormsHelp_StartsWithAccepts(t *testing.T) {
h := helptext.PolicySourceFormsHelp
if !strings.HasPrefix(h, "Accepts") {
t.Errorf("PolicySourceFormsHelp should start with 'Accepts', got: %q", h[:min(20, len(h))])
}
}

func TestPolicySourceFormsHelp_MentionsGitHub(t *testing.T) {
h := helptext.PolicySourceFormsHelp
if !strings.Contains(h, "github.com") && !strings.Contains(h, "GitHub") && !strings.Contains(h, "git") {
t.Error("PolicySourceFormsHelp should reference git/github hosting")
}
}

func TestPolicySourceFormsHelp_HasCommaList(t *testing.T) {
h := helptext.PolicySourceFormsHelp
// The help string should list multiple options (at least one comma)
if !strings.Contains(h, ",") {
t.Error("PolicySourceFormsHelp should list multiple options separated by commas")
}
}

func TestPolicySourceFormsHelp_NoLeadingSpace(t *testing.T) {
h := helptext.PolicySourceFormsHelp
if strings.HasPrefix(h, " ") || strings.HasPrefix(h, "\t") {
t.Error("PolicySourceFormsHelp should not have leading whitespace")
}
}

func TestPolicySourceFormsHelp_NoTrailingNewline(t *testing.T) {
h := helptext.PolicySourceFormsHelp
if strings.HasSuffix(h, "\n") {
t.Error("PolicySourceFormsHelp should not end with a newline")
}
}

func TestPolicySourceFormsHelp_SingleLine(t *testing.T) {
h := helptext.PolicySourceFormsHelp
if strings.Contains(h, "\n") {
t.Error("PolicySourceFormsHelp should fit on a single line (no embedded newlines)")
}
}

func TestPolicySourceFormsHelp_OrgMentionedFirst(t *testing.T) {
h := helptext.PolicySourceFormsHelp
orgIdx := strings.Index(h, "org")
ownerIdx := strings.Index(h, "owner/repo")
if orgIdx < 0 {
t.Skip("'org' not found in help text")
}
if ownerIdx < 0 {
t.Skip("'owner/repo' not found in help text")
}
if orgIdx > ownerIdx {
t.Errorf("'org' form should appear before 'owner/repo' form in help text")
}
}

func TestPolicySourceFormsHelp_MentionsDefaultHost(t *testing.T) {
h := helptext.PolicySourceFormsHelp
if !strings.Contains(h, "github.com") {
t.Error("PolicySourceFormsHelp should mention github.com as the default host")
}
}

func TestPolicySourceFormsHelp_MentionsFilePath(t *testing.T) {
h := helptext.PolicySourceFormsHelp
if !strings.Contains(h, "file") && !strings.Contains(h, "path") && !strings.Contains(h, "local") {
t.Error("PolicySourceFormsHelp should mention local file path option")
}
}

func TestPolicySourceFormsHelp_MinWordCount(t *testing.T) {
h := helptext.PolicySourceFormsHelp
words := strings.Fields(h)
if len(words) < 10 {
t.Errorf("PolicySourceFormsHelp too short (%d words); expected at least 10", len(words))
}
}

func min(a, b int) int {
if a < b {
return a
}
return b
}
