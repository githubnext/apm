// Package apmpackage provides the APMPackage and PackageInfo data models.
// Migrated from src/apm_cli/models/apm_package.py.
package apmpackage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PackageContentType represents the content type of a package.
type PackageContentType int

const (
	ContentTypeInstructions PackageContentType = iota
	ContentTypeSkill
	ContentTypeHybrid
	ContentTypePrompts
)

// String returns the string representation of a PackageContentType.
func (t PackageContentType) String() string {
	switch t {
	case ContentTypeInstructions:
		return "instructions"
	case ContentTypeSkill:
		return "skill"
	case ContentTypeHybrid:
		return "hybrid"
	case ContentTypePrompts:
		return "prompts"
	default:
		return "unknown"
	}
}

// ParseContentType parses a string content type.
func ParseContentType(s string) (PackageContentType, error) {
	switch strings.ToLower(s) {
	case "instructions":
		return ContentTypeInstructions, nil
	case "skill":
		return ContentTypeSkill, nil
	case "hybrid":
		return ContentTypeHybrid, nil
	case "prompts":
		return ContentTypePrompts, nil
	default:
		return 0, fmt.Errorf("unknown content type: %s", s)
	}
}

// APMPackage represents an APM package with metadata.
type APMPackage struct {
	Name             string
	Version          string
	Description      string
	Author           string
	License          string
	Source           string
	ResolvedCommit   string
	Dependencies     map[string][]interface{}
	DevDependencies  map[string][]interface{}
	Scripts          map[string]string
	PackagePath      string
	SourcePath       string
	Target           interface{} // string or []string
	Type             *PackageContentType
	Includes         interface{} // string "auto" or []string
}

// PackageInfo contains information about a downloaded/installed package.
type PackageInfo struct {
	Package           *APMPackage
	InstallPath       string
	InstalledAt       string
	PackageType       string // "APM_PACKAGE", "CLAUDE_SKILL", or "HYBRID"
}

// GetPrimitivesPath returns the path to the .apm directory for this package.
func (p *PackageInfo) GetPrimitivesPath() string {
	return filepath.Join(p.InstallPath, ".apm")
}

// HasPrimitives checks if the package has any primitives.
func (p *PackageInfo) HasPrimitives() bool {
	apmDir := p.GetPrimitivesPath()
	for _, pt := range []string{"instructions", "chatmodes", "contexts", "prompts", "hooks"} {
		dir := filepath.Join(apmDir, pt)
		if entries, err := os.ReadDir(dir); err == nil && len(entries) > 0 {
			return true
		}
	}
	hooksDir := filepath.Join(p.InstallPath, "hooks")
	if entries, err := os.ReadDir(hooksDir); err == nil {
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".json") {
				return true
			}
		}
	}
	return false
}

// LoadFromApmYml loads basic package metadata from an apm.yml file.
// This is a lightweight loader that extracts name/version/description/target
// without full dependency parsing.
func LoadFromApmYml(apmYmlPath string) (*APMPackage, error) {
	f, err := os.Open(apmYmlPath)
	if err != nil {
		return nil, fmt.Errorf("apm.yml not found: %s", apmYmlPath)
	}
	defer f.Close()

	data := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			// Strip inline YAML quotes
			val = strings.Trim(val, "\"'")
			if val != "" && !strings.HasPrefix(val, "{") && !strings.HasPrefix(val, "[") {
				data[key] = val
			}
		}
	}

	name := data["name"]
	version := data["version"]
	if name == "" {
		return nil, fmt.Errorf("missing required field 'name' in apm.yml")
	}
	if version == "" {
		return nil, fmt.Errorf("missing required field 'version' in apm.yml")
	}

	pkg := &APMPackage{
		Name:        name,
		Version:     version,
		Description: data["description"],
		Author:      data["author"],
		License:     data["license"],
		PackagePath: filepath.Dir(apmYmlPath),
		SourcePath:  filepath.Dir(apmYmlPath),
	}

	if t := data["target"]; t != "" {
		pkg.Target = t
	}

	if typeStr := data["type"]; typeStr != "" {
		ct, err := ParseContentType(typeStr)
		if err == nil {
			pkg.Type = &ct
		}
	}

	return pkg, nil
}
