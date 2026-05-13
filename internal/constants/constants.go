// Package constants defines shared constants for the APM CLI.
// Migrated from src/apm_cli/constants.py
package constants

// InstallMode controls which dependency types are installed.
type InstallMode string

const (
	InstallModeAll InstallMode = "all"
	InstallModeAPM InstallMode = "apm"
	InstallModeMCP InstallMode = "mcp"
)

// File and directory names.
const (
	APMYMLFilename            = "apm.yml"
	APMLockFilename           = "apm.lock"
	APMModulesDir             = "apm_modules"
	APMDir                    = ".apm"
	SkillMDFilename           = "SKILL.md"
	AgentsMDFilename          = "AGENTS.md"
	ClaudeMDFilename          = "CLAUDE.md"
	GitHubDir                 = ".github"
	ClaudeDir                 = ".claude"
	GitignoreFilename         = ".gitignore"
	APMModulesGitignorePattern = "apm_modules/"
)

// DefaultSkipDirs lists directory names unconditionally skipped during
// primitive-file discovery. These never contain APM primitives and can
// be very large (e.g. node_modules, .git objects).
var DefaultSkipDirs = map[string]struct{}{
	".git":          {},
	"node_modules":  {},
	"__pycache__":   {},
	".pytest_cache": {},
	".venv":         {},
	"venv":          {},
	".tox":          {},
	"build":         {},
	"dist":          {},
	".mypy_cache":   {},
	"apm_modules":   {},
}
