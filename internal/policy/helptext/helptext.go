// Package helptext contains shared help text for policy-related CLI commands.
// Migrated from src/apm_cli/policy/_help_text.py.
package helptext

// PolicySourceFormsHelp is the canonical user-facing description of the
// --policy / --policy-source argument formats accepted by discover_policy.
const PolicySourceFormsHelp = "Accepts: 'org' (auto-discover from your project's git remote), " +
	"'owner/repo' (defaults to github.com), an https:// URL, or a " +
	"local file path."
