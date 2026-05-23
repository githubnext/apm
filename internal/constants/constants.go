// Package constants provides shared constants for the APM CLI.
package constants

// File and directory names
const (
	APMYMLFilename            = "apm.yml"
	APMLockFilename           = "apm.lock"
	APMModulesDir             = "apm_modules"
	APMDir                    = ".apm"
	SkillMDFilename           = "SKILL.md"
	AgentsMDFilename          = "AGENTS.md"
	ClaudeMDFilename          = "CLAUDE.md"
	GithubDir                 = ".github"
	ClaudeDir                 = ".claude"
	GitignoreFilename         = ".gitignore"
	APMModulesGitignorePattern = "apm_modules/"
)

// InstallMode controls which dependency types are installed.
type InstallMode string

const (
	InstallModeAll InstallMode = "all"
	InstallModeAPM InstallMode = "apm"
	InstallModeMCP InstallMode = "mcp"
)

// DefaultSkipDirs lists directories unconditionally skipped during
// primitive-file discovery. These never contain APM primitives or user
// source files and can be very large (e.g. node_modules, .git objects).
// NOTE: .apm is intentionally absent -- it is where primitives live.
var DefaultSkipDirs = map[string]bool{
	".git":          true,
	"node_modules":  true,
	"__pycache__":   true,
	".pytest_cache": true,
	".venv":         true,
	"venv":          true,
	".tox":          true,
	"build":         true,
	"dist":          true,
	".mypy_cache":   true,
	"apm_modules":   true,
}
