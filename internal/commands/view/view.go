// Package view implements the "apm view" / "apm info" command.
//
// Shows detailed metadata for an installed APM package: version history,
// source repository, installed files, and optional field filters.
//
// Migrated from: src/apm_cli/commands/view.py
package view

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ViewOptions configures what apm view displays.
type ViewOptions struct {
	ProjectRoot string
	Package     string
	Field       string // optional: "versions" | ""
	Format      string // "text" | "json"
	Verbose     bool
}

// PackageInfo holds metadata for an installed package.
type PackageInfo struct {
	Name          string            `json:"name"`
	InstalledPath string            `json:"installed_path"`
	Ref           string            `json:"ref,omitempty"`
	Commit        string            `json:"commit,omitempty"`
	Source        string            `json:"source,omitempty"`
	ApmYML        map[string]interface{} `json:"apm_yml,omitempty"`
	Files         []string          `json:"files,omitempty"`
	Versions      []string          `json:"versions,omitempty"`
}

const (
	apmModulesDir = ".apm_modules"
	apmYMLFile    = "apm.yml"
	skillMDFile   = "SKILL.md"
)

// Run executes the view command and prints output.
func Run(opts ViewOptions) error {
	apmModules := filepath.Join(opts.ProjectRoot, apmModulesDir)

	pkgPath, err := resolvePackagePath(opts.Package, apmModules)
	if err != nil {
		return err
	}

	info, err := buildPackageInfo(opts.Package, pkgPath)
	if err != nil {
		return fmt.Errorf("read package info: %w", err)
	}

	if opts.Field != "" {
		return printField(info, opts.Field, opts.Format)
	}

	if opts.Format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}

	printText(info, opts.Verbose)
	return nil
}

// resolvePackagePath locates the package directory inside apmModulesDir.
func resolvePackagePath(pkg, apmModules string) (string, error) {
	if pkg == "" {
		return "", fmt.Errorf("package name is required")
	}
	// Guard against traversal
	if strings.Contains(pkg, "..") {
		return "", fmt.Errorf("invalid package name: %q", pkg)
	}

	// Direct path match (handles org/repo)
	direct := filepath.Join(apmModules, filepath.FromSlash(pkg))
	if fi, err := os.Stat(direct); err == nil && fi.IsDir() {
		return direct, nil
	}

	// Fallback: two-level scan for short (repo-only) names
	entries, err := os.ReadDir(apmModules)
	if err != nil {
		return "", fmt.Errorf("cannot read %s: %w", apmModules, err)
	}
	for _, org := range entries {
		if !org.IsDir() {
			continue
		}
		candidate := filepath.Join(apmModules, org.Name(), pkg)
		if fi, err := os.Stat(candidate); err == nil && fi.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("package %q not found in %s", pkg, apmModules)
}

// buildPackageInfo collects metadata from the package directory.
func buildPackageInfo(name, pkgPath string) (*PackageInfo, error) {
	info := &PackageInfo{
		Name:          name,
		InstalledPath: pkgPath,
	}

	// Read apm.yml if present
	ymlPath := filepath.Join(pkgPath, apmYMLFile)
	if data, err := os.ReadFile(ymlPath); err == nil {
		var yml map[string]interface{}
		if err := parseSimpleYAML(data, &yml); err == nil {
			info.ApmYML = yml
			if src, ok := yml["source"].(string); ok {
				info.Source = src
			}
			if ref, ok := yml["ref"].(string); ok {
				info.Ref = ref
			}
		}
	}

	// List installed files
	var files []string
	_ = filepath.WalkDir(pkgPath, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(pkgPath, path)
		files = append(files, rel)
		return nil
	})
	info.Files = files

	return info, nil
}

// printField prints only the requested field.
func printField(info *PackageInfo, field, format string) error {
	switch field {
	case "versions":
		if format == "json" {
			enc := json.NewEncoder(os.Stdout)
			return enc.Encode(info.Versions)
		}
		if len(info.Versions) == 0 {
			fmt.Println("(no version history available)")
		}
		for _, v := range info.Versions {
			fmt.Println(v)
		}
		return nil
	default:
		return fmt.Errorf("unknown field %q; valid fields: versions", field)
	}
}

// printText renders a human-readable summary.
func printText(info *PackageInfo, verbose bool) {
	fmt.Printf("Package: %s\n", info.Name)
	fmt.Printf("  Path:   %s\n", info.InstalledPath)
	if info.Ref != "" {
		fmt.Printf("  Ref:    %s\n", info.Ref)
	}
	if info.Commit != "" {
		fmt.Printf("  Commit: %s\n", info.Commit)
	}
	if info.Source != "" {
		fmt.Printf("  Source: %s\n", info.Source)
	}
	if verbose && len(info.Files) > 0 {
		fmt.Printf("  Files (%d):\n", len(info.Files))
		for _, f := range info.Files {
			fmt.Printf("    %s\n", f)
		}
	}
}

// parseSimpleYAML does minimal key:value YAML parsing into a map.
func parseSimpleYAML(data []byte, out *map[string]interface{}) error {
	m := make(map[string]interface{})
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	*out = m
	return nil
}
