// Package models defines core data structures for APM packages.
// Mirrors src/apm_cli/models/results.py and src/apm_cli/models/validation.py.
package models

// InstallResult holds the result of an APM install operation.
// Mirrors src/apm_cli/models/results.py:InstallResult.
type InstallResult struct {
	InstalledCount     int
	PromptsIntegrated  int
	AgentsIntegrated   int
	Diagnostics        interface{}
	PackageTypes       map[string]string
}

// PrimitiveCounts holds counts of primitives in a package.
// Mirrors src/apm_cli/models/results.py:PrimitiveCounts.
type PrimitiveCounts struct {
	Prompts      int
	Agents       int
	Instructions int
	Skills       int
	Hooks        int
	Commands     int
}

// PackageType classifies a package by its content.
// Mirrors src/apm_cli/models/validation.py:PackageType.
type PackageType int

const (
	PackageTypeAPMPackage         PackageType = iota
	PackageTypeClaudeSkill
	PackageTypeHookPackage
	PackageTypeHybrid
	PackageTypeMarketplacePlugin
	PackageTypeSkillBundle
	PackageTypeInvalid
)

func (p PackageType) String() string {
	switch p {
	case PackageTypeAPMPackage:
		return "apm_package"
	case PackageTypeClaudeSkill:
		return "claude_skill"
	case PackageTypeHookPackage:
		return "hook_package"
	case PackageTypeHybrid:
		return "hybrid"
	case PackageTypeMarketplacePlugin:
		return "marketplace_plugin"
	case PackageTypeSkillBundle:
		return "skill_bundle"
	case PackageTypeInvalid:
		return "invalid"
	default:
		return "unknown"
	}
}

// PackageContentType is the explicit package content type declared in apm.yml.
// Mirrors src/apm_cli/models/validation.py:PackageContentType.
type PackageContentType int

const (
	PackageContentTypeInstructions PackageContentType = iota
	PackageContentTypeSkill
	PackageContentTypeHybrid
	PackageContentTypePrompts
)

func (p PackageContentType) String() string {
	switch p {
	case PackageContentTypeInstructions:
		return "instructions"
	case PackageContentTypeSkill:
		return "skill"
	case PackageContentTypeHybrid:
		return "hybrid"
	case PackageContentTypePrompts:
		return "prompts"
	default:
		return "unknown"
	}
}

// ParsePackageContentType parses a string into a PackageContentType.
// Mirrors src/apm_cli/models/validation.py:PackageContentType.from_string.
func ParsePackageContentType(value string) (PackageContentType, error) {
	if value == "" {
		return 0, errorf("Package type cannot be empty")
	}
	switch value {
	case "instructions":
		return PackageContentTypeInstructions, nil
	case "skill":
		return PackageContentTypeSkill, nil
	case "hybrid":
		return PackageContentTypeHybrid, nil
	case "prompts":
		return PackageContentTypePrompts, nil
	default:
		return 0, errorf("Invalid package type '%s'. Valid types are: 'instructions', 'skill', 'hybrid', 'prompts'", value)
	}
}

// ValidationError enumerates types of validation errors for APM packages.
// Mirrors src/apm_cli/models/validation.py:ValidationError.
type ValidationErrorCode string

const (
	ValidationErrMissingAPMYml     ValidationErrorCode = "missing_apm_yml"
	ValidationErrMissingAPMDir     ValidationErrorCode = "missing_apm_dir"
	ValidationErrInvalidYMLFormat  ValidationErrorCode = "invalid_yml_format"
	ValidationErrMissingRequired   ValidationErrorCode = "missing_required_field"
)

// ValidationResult holds the result of a package validation.
type ValidationResult struct {
	Valid       bool
	Errors      []ValidationErrorCode
	PackageType PackageType
}

// PluginMetadata holds metadata for a plugin.
// Mirrors src/apm_cli/models/plugin.py:PluginMetadata.
type PluginMetadata struct {
	ID           string
	Name         string
	Version      string
	Description  string
	Author       string
	Repository   string
	Homepage     string
	License      string
	Tags         []string
	Dependencies []string
}
