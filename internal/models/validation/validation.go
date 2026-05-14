// Package validation provides validation logic and type enums for APM packages.
//
// Mirrors src/apm_cli/models/validation.py.
package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// PackageType classifies packages based on their content.
type PackageType int

const (
	PackageTypeAPMPackage        PackageType = iota // Has apm.yml
	PackageTypeClaudeSkill                          // Has SKILL.md, no apm.yml
	PackageTypeHookPackage                          // Has hooks/hooks.json, no apm.yml or SKILL.md
	PackageTypeHybrid                               // Has both apm.yml and SKILL.md (root)
	PackageTypeMarketplacePlugin                    // Has plugin.json or .claude-plugin/
	PackageTypeSkillBundle                          // Has skills/<name>/SKILL.md (nested)
	PackageTypeInvalid                              // None of the above
)

// String returns a human-readable name for the package type.
func (t PackageType) String() string {
	switch t {
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
	default:
		return "invalid"
	}
}

// PackageContentType is the user-facing type field in apm.yml.
type PackageContentType int

const (
	PackageContentTypeInstructions PackageContentType = iota // Compile to AGENTS.md only
	PackageContentTypeSkill                                  // Install as native skill only
	PackageContentTypeHybrid                                 // Both (default)
	PackageContentTypePrompts                                // Commands/prompts only
)

// String returns the string value of the content type.
func (t PackageContentType) String() string {
	switch t {
	case PackageContentTypeInstructions:
		return "instructions"
	case PackageContentTypeSkill:
		return "skill"
	case PackageContentTypeHybrid:
		return "hybrid"
	case PackageContentTypePrompts:
		return "prompts"
	default:
		return "hybrid"
	}
}

// PackageContentTypeFromString parses a string into a PackageContentType.
func PackageContentTypeFromString(value string) (PackageContentType, error) {
	if value == "" {
		return 0, fmt.Errorf("package type cannot be empty")
	}
	v := strings.ToLower(strings.TrimSpace(value))
	switch v {
	case "instructions":
		return PackageContentTypeInstructions, nil
	case "skill":
		return PackageContentTypeSkill, nil
	case "hybrid":
		return PackageContentTypeHybrid, nil
	case "prompts":
		return PackageContentTypePrompts, nil
	default:
		return 0, fmt.Errorf("invalid package type '%s'. Valid types are: 'instructions', 'skill', 'hybrid', 'prompts'", value)
	}
}

// ValidationError enumerates types of validation errors for APM packages.
type ValidationError int

const (
	ValidationErrorMissingAPMYml           ValidationError = iota
	ValidationErrorMissingAPMDir
	ValidationErrorInvalidYmlFormat
	ValidationErrorMissingRequiredField
	ValidationErrorInvalidVersionFormat
	ValidationErrorInvalidDependencyFormat
	ValidationErrorEmptyAPMDir
	ValidationErrorInvalidPrimitiveStructure
)

// ValidationResult holds the result of APM package validation.
type ValidationResult struct {
	IsValid     bool
	Errors      []string
	Warnings    []string
	PackageType PackageType
}

// NewValidationResult creates an empty (valid) ValidationResult.
func NewValidationResult() *ValidationResult {
	return &ValidationResult{IsValid: true}
}

// AddError adds a validation error and marks the result as invalid.
func (r *ValidationResult) AddError(err string) {
	r.Errors = append(r.Errors, err)
	r.IsValid = false
}

// AddWarning adds a validation warning.
func (r *ValidationResult) AddWarning(warning string) {
	r.Warnings = append(r.Warnings, warning)
}

// HasIssues returns true if there are any errors or warnings.
func (r *ValidationResult) HasIssues() bool {
	return len(r.Errors) > 0 || len(r.Warnings) > 0
}

// Summary returns a human-readable summary of validation results.
func (r *ValidationResult) Summary() string {
	if r.IsValid && len(r.Warnings) == 0 {
		return "[+] Package is valid"
	} else if r.IsValid && len(r.Warnings) > 0 {
		return fmt.Sprintf("[!] Package is valid with %d warning(s)", len(r.Warnings))
	}
	return fmt.Sprintf("[x] Package is invalid with %d error(s)", len(r.Errors))
}

// pluginDirs defines the canonical order of plugin content directories.
var pluginDirs = []string{"agents", "skills", "commands"}

// DetectionEvidence is a snapshot of file-system signals for package classification.
type DetectionEvidence struct {
	HasAPMYml          bool
	HasSkillMD         bool
	HasHookJSON        bool
	PluginJSONPath     string // empty if not found
	PluginDirsPresent  []string
	HasClaudePluginDir bool
	NestedSkillDirs    []string
	HasPluginManifest  bool
}

// HasPluginEvidence returns true if a real plugin manifest is present.
func (e *DetectionEvidence) HasPluginEvidence() bool {
	return e.HasPluginManifest
}

// hasHookJSON checks if the package has hook JSON files in hooks/ or .apm/hooks/.
func hasHookJSON(packagePath string) bool {
	for _, dir := range []string{filepath.Join(packagePath, "hooks"), filepath.Join(packagePath, ".apm", "hooks")} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
				return true
			}
		}
	}
	return false
}

// findPluginJSON searches for plugin.json in the package root.
func findPluginJSON(packagePath string) string {
	p := filepath.Join(packagePath, "plugin.json")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return ""
}

// GatherDetectionEvidence collects all package-type signals from a directory.
func GatherDetectionEvidence(packagePath string) *DetectionEvidence {
	ev := &DetectionEvidence{}

	// Check apm.yml
	if _, err := os.Stat(filepath.Join(packagePath, "apm.yml")); err == nil {
		ev.HasAPMYml = true
	}

	// Check SKILL.md
	if _, err := os.Stat(filepath.Join(packagePath, "SKILL.md")); err == nil {
		ev.HasSkillMD = true
	}

	// Check hook JSON
	ev.HasHookJSON = hasHookJSON(packagePath)

	// Check plugin dirs
	for _, dir := range pluginDirs {
		if info, err := os.Stat(filepath.Join(packagePath, dir)); err == nil && info.IsDir() {
			ev.PluginDirsPresent = append(ev.PluginDirsPresent, dir)
		}
	}

	// Check plugin.json
	ev.PluginJSONPath = findPluginJSON(packagePath)

	// Check .claude-plugin/
	if info, err := os.Stat(filepath.Join(packagePath, ".claude-plugin")); err == nil && info.IsDir() {
		ev.HasClaudePluginDir = true
	}

	// Plugin manifest = plugin.json OR .claude-plugin/
	ev.HasPluginManifest = ev.PluginJSONPath != "" || ev.HasClaudePluginDir

	// Nested skill dirs: directories under skills/ that contain a SKILL.md
	skillsDir := filepath.Join(packagePath, "skills")
	if entries, err := os.ReadDir(skillsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			skillMD := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
			if _, err := os.Stat(skillMD); err == nil {
				ev.NestedSkillDirs = append(ev.NestedSkillDirs, entry.Name())
			}
		}
	}

	return ev
}

// DetectPackageType classifies a package directory into a PackageType.
// Returns (packageType, pluginJSONPath). pluginJSONPath is non-empty only
// when MARKETPLACE_PLUGIN was matched via an actual plugin.json file.
func DetectPackageType(packagePath string) (PackageType, string) {
	ev := GatherDetectionEvidence(packagePath)

	// 1. Plugin manifest present -> MARKETPLACE_PLUGIN
	if ev.HasPluginManifest {
		return PackageTypeMarketplacePlugin, ev.PluginJSONPath
	}

	// 2. Root SKILL.md + apm.yml -> HYBRID
	if ev.HasAPMYml && ev.HasSkillMD {
		return PackageTypeHybrid, ""
	}

	// 3. Root SKILL.md only -> CLAUDE_SKILL
	if ev.HasSkillMD {
		return PackageTypeClaudeSkill, ""
	}

	// 4. Nested skills/<x>/SKILL.md -> SKILL_BUNDLE
	if len(ev.NestedSkillDirs) > 0 {
		return PackageTypeSkillBundle, ""
	}

	// 5. apm.yml present -> APM_PACKAGE or INVALID
	if ev.HasAPMYml {
		apmDir := filepath.Join(packagePath, ".apm")
		if info, err := os.Stat(apmDir); err == nil && info.IsDir() {
			return PackageTypeAPMPackage, ""
		}
		if apmYMLDeclaresDependencies(filepath.Join(packagePath, "apm.yml")) {
			return PackageTypeAPMPackage, ""
		}
		return PackageTypeInvalid, ""
	}

	// 6. hooks/*.json only -> HOOK_PACKAGE
	if ev.HasHookJSON {
		return PackageTypeHookPackage, ""
	}

	// 7. Nothing recognisable -> INVALID
	return PackageTypeInvalid, ""
}

// apmYMLDeclaresDependencies returns true iff apm.yml declares at least one dependency.
func apmYMLDeclaresDependencies(apmYMLPath string) bool {
	data, err := os.ReadFile(apmYMLPath)
	if err != nil {
		return false
	}
	// Simple heuristic: look for "apm:" or "mcp:" under dependencies/devDependencies
	// with at least one list item. A full YAML parse is not available without external libs.
	content := string(data)
	// Look for a non-empty apm: or mcp: list under dependencies or devDependencies
	depSection := extractYAMLSection(content, "dependencies")
	devSection := extractYAMLSection(content, "devDependencies")
	return hasListedDeps(depSection) || hasListedDeps(devSection)
}

// extractYAMLSection extracts a named top-level section from simple YAML.
func extractYAMLSection(content, key string) string {
	lines := strings.Split(content, "\n")
	inSection := false
	var result []string
	prefix := key + ":"
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == prefix || strings.HasPrefix(trimmed, prefix+" ") {
			inSection = true
			result = append(result, line)
			continue
		}
		if inSection {
			// Stop when we hit another top-level key (no leading space)
			if len(line) > 0 && line[0] != ' ' && line[0] != '\t' && line[0] != '#' && trimmed != "" {
				break
			}
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// hasListedDeps checks if the section has apm: or mcp: lists with entries.
func hasListedDeps(section string) bool {
	lines := strings.Split(section, "\n")
	inAPMorMCP := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "apm:" || trimmed == "mcp:" {
			inAPMorMCP = true
			continue
		}
		if inAPMorMCP {
			if strings.HasPrefix(trimmed, "- ") {
				return true
			}
			if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
				inAPMorMCP = false
			}
		}
	}
	return false
}

// semverRe matches a semantic version string (x.y.z).
var semverRe = regexp.MustCompile(`^\d+\.\d+\.\d+`)

// ValidateAPMPackage validates that a directory contains a valid APM package.
func ValidateAPMPackage(packagePath string) *ValidationResult {
	result := NewValidationResult()

	// Check if directory exists
	info, err := os.Stat(packagePath)
	if err != nil {
		result.AddError(fmt.Sprintf("Package directory does not exist: %s", packagePath))
		return result
	}
	if !info.IsDir() {
		result.AddError(fmt.Sprintf("Package path is not a directory: %s", packagePath))
		return result
	}

	// Detect package type
	pkgType, pluginJSONPath := DetectPackageType(packagePath)
	result.PackageType = pkgType

	if pkgType == PackageTypeInvalid {
		apmYMLPath := filepath.Join(packagePath, "apm.yml")
		if _, err := os.Stat(apmYMLPath); err == nil {
			apmPath := filepath.Join(packagePath, ".apm")
			if apmInfo, err := os.Stat(apmPath); err == nil && !apmInfo.IsDir() {
				result.AddError(".apm must be a directory")
			} else {
				dirName := filepath.Base(packagePath)
				result.AddError(fmt.Sprintf(
					"Not a valid APM package: %s has apm.yml but is missing the required .apm/ directory. "+
						"Add .apm/ with primitives (instructions, skills, etc.), "+
						"declare dependencies in apm.yml (curated aggregator), "+
						"or add skills/<name>/SKILL.md for a skill bundle.", dirName))
			}
		} else {
			dirName := filepath.Base(packagePath)
			result.AddError(fmt.Sprintf(
				"Not a valid APM package: no apm.yml, SKILL.md, hooks, or plugin structure found in %s. "+
					"Ensure the package has SKILL.md (skill bundle), "+
					"apm.yml + .apm/ (APM package), or plugin.json (Claude plugin) at its root.", dirName))
		}
		return result
	}

	switch pkgType {
	case PackageTypeHookPackage:
		return validateHookPackage(packagePath, result)
	case PackageTypeClaudeSkill:
		return validateClaudeSkill(packagePath, result)
	case PackageTypeMarketplacePlugin:
		return validateMarketplacePlugin(packagePath, pluginJSONPath, result)
	case PackageTypeSkillBundle:
		return validateSkillBundle(packagePath, result)
	case PackageTypeHybrid:
		return validateHybridPackage(packagePath, result)
	default:
		return validateAPMPackageWithYML(packagePath, result)
	}
}

func validateHookPackage(packagePath string, result *ValidationResult) *ValidationResult {
	// Hook package is valid as-is -- just has hooks/*.json
	return result
}

func validateClaudeSkill(packagePath string, result *ValidationResult) *ValidationResult {
	// Check SKILL.md is readable
	skillMD := filepath.Join(packagePath, "SKILL.md")
	if _, err := os.ReadFile(skillMD); err != nil {
		result.AddError(fmt.Sprintf("Failed to read SKILL.md: %v", err))
	}
	return result
}

func validateMarketplacePlugin(packagePath, pluginJSONPath string, result *ValidationResult) *ValidationResult {
	// Check plugin.json or .claude-plugin/ is present and readable
	if pluginJSONPath != "" {
		if _, err := os.ReadFile(pluginJSONPath); err != nil {
			result.AddError(fmt.Sprintf("Failed to read plugin.json: %v", err))
		}
	}
	return result
}

func validateSkillBundle(packagePath string, result *ValidationResult) *ValidationResult {
	skillsDir := filepath.Join(packagePath, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		result.AddError(fmt.Sprintf("SKILL_BUNDLE detected but could not read skills/ directory: %v", err))
		return result
	}

	var skillNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		skillMD := filepath.Join(skillsDir, name, "SKILL.md")
		if _, err := os.Stat(skillMD); err != nil {
			continue
		}

		// Path safety: reject traversal
		if strings.Contains(name, "..") || strings.Contains(name, "/") {
			result.AddError(fmt.Sprintf("Invalid skill directory name: %s", name))
			continue
		}

		skillNames = append(skillNames, name)
	}

	if len(skillNames) == 0 {
		result.AddError(fmt.Sprintf("SKILL_BUNDLE detected but no valid skills/<name>/SKILL.md found in %s/skills/", filepath.Base(packagePath)))
		return result
	}

	return result
}

func validateHybridPackage(packagePath string, result *ValidationResult) *ValidationResult {
	apmDir := filepath.Join(packagePath, ".apm")
	if info, err := os.Stat(apmDir); err == nil && info.IsDir() {
		return validateAPMPackageWithYML(packagePath, result)
	}

	// Skill-bundle path (no .apm/)
	apmYMLPath := filepath.Join(packagePath, "apm.yml")
	if _, err := os.Stat(apmYMLPath); err != nil {
		result.AddError("HYBRID package missing apm.yml")
		return result
	}

	// Check SKILL.md is present
	skillMD := filepath.Join(packagePath, "SKILL.md")
	if _, err := os.Stat(skillMD); err != nil {
		result.AddError("HYBRID package missing SKILL.md")
		return result
	}

	return result
}

func validateAPMPackageWithYML(packagePath string, result *ValidationResult) *ValidationResult {
	apmYMLPath := filepath.Join(packagePath, "apm.yml")

	// Parse apm.yml basic fields
	data, err := os.ReadFile(apmYMLPath)
	if err != nil {
		result.AddError(fmt.Sprintf("Invalid apm.yml: %v", err))
		return result
	}

	// Check for .apm directory
	apmDir := filepath.Join(packagePath, ".apm")
	apmDirInfo, apmDirErr := os.Stat(apmDir)
	if apmDirErr != nil {
		// No .apm/ -- check if dep-only (curated aggregator)
		if apmYMLDeclaresDependencies(apmYMLPath) {
			return result
		}
		result.AddError(fmt.Sprintf("Missing required directory: .apm/ -- "+
			"an APM package with apm.yml needs either a .apm/ directory "+
			"containing primitives, or dependencies declared in apm.yml. "+
			"Alternatively, add a SKILL.md to make this a skill bundle."))
		return result
	}

	if !apmDirInfo.IsDir() {
		result.AddError(".apm must be a directory")
		return result
	}

	// Check for primitives in .apm/
	primitiveTypes := []string{"instructions", "chatmodes", "contexts", "prompts"}
	hasPrimitives := false
	for _, pt := range primitiveTypes {
		ptDir := filepath.Join(apmDir, pt)
		entries, err := os.ReadDir(ptDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				hasPrimitives = true
				// Check for empty files
				content, err := os.ReadFile(filepath.Join(ptDir, e.Name()))
				if err == nil && strings.TrimSpace(string(content)) == "" {
					result.AddWarning(fmt.Sprintf("Empty primitive file: .apm/%s/%s", pt, e.Name()))
				}
			}
		}
	}

	if !hasPrimitives {
		hasPrimitives = hasHookJSON(packagePath)
	}

	if !hasPrimitives {
		result.AddWarning("No primitive files found in .apm/ directory")
	}

	// Version format validation (basic semver check)
	// Extract version from apm.yml content
	version := extractYAMLField(string(data), "version")
	if version != "" && !semverRe.MatchString(version) {
		result.AddWarning(fmt.Sprintf("Version '%s' doesn't follow semantic versioning (x.y.z)", version))
	}

	return result
}

// extractYAMLField extracts a simple scalar field value from YAML content.
func extractYAMLField(content, key string) string {
	prefix := key + ":"
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			val := strings.TrimSpace(trimmed[len(prefix):])
			// Strip quotes
			if len(val) >= 2 && (val[0] == '"' || val[0] == '\'') {
				val = val[1 : len(val)-1]
			}
			return val
		}
	}
	return ""
}
