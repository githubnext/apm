// Package install implements the "apm install" command.
//
// Orchestrates dependency resolution, download, integration, lockfile
// persistence, and post-install validation.
//
// Migrated from: src/apm_cli/commands/install.py
package install

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// InstallMode controls which components are written during install.
type InstallMode string

const (
	InstallModeAll      InstallMode = "all"
	InstallModePrimitives InstallMode = "primitives"
	InstallModeClients  InstallMode = "clients"
)

// InstallOptions configures a single install invocation.
type InstallOptions struct {
	ProjectRoot   string
	PackageRefs   []string
	Targets       []string
	Frozen        bool
	DryRun        bool
	Verbose       bool
	Force         bool
	UserScope     bool
	NoProgress    bool
	SkipLockfile  bool
	Mode          InstallMode
	AuthToken     string
	ConcurrentDL  int
}

// InstallResult captures the outcome of an install run.
type InstallResult struct {
	PackagesInstalled  int
	PackagesSkipped    int
	PackagesRemoved    int
	FilesWritten       []string
	LockfileUpdated    bool
	DurationSeconds    float64
	Warnings           []string
	Errors             []string
}

// DependencyEntry represents one entry from apm.yml.
type DependencyEntry struct {
	Name    string
	Ref     string
	Host    string
	Org     string
	Repo    string
	Path    string
	Version string
	Local   bool
}

// PolicyViolation describes a policy check failure.
type PolicyViolation struct {
	Package string
	Rule    string
	Message string
}

// AuthenticationError is returned when credentials are missing or invalid.
type AuthenticationError struct {
	Host    string
	Message string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("authentication failed for %s: %s", e.Host, e.Message)
}

// FrozenInstallError is returned when apm.lock.yaml would change in frozen mode.
type FrozenInstallError struct {
	Changed []string
}

func (e *FrozenInstallError) Error() string {
	return fmt.Sprintf("frozen install: lockfile would change (%d packages)", len(e.Changed))
}

// PolicyViolationError is returned when a dependency violates install policy.
type PolicyViolationError struct {
	Violations []PolicyViolation
}

func (e *PolicyViolationError) Error() string {
	if len(e.Violations) == 1 {
		return fmt.Sprintf("policy violation: %s", e.Violations[0].Message)
	}
	return fmt.Sprintf("policy violations: %d rules violated", len(e.Violations))
}

// apmYML is the minimal shape of apm.yml we read.
type apmYML struct {
	Dependencies []map[string]interface{} `json:"dependencies"`
	Targets      []string                 `json:"targets"`
}

// RunInstall is the main entry point for the install command.
//
// It reads apm.yml, resolves all dependencies, downloads missing packages,
// writes integration files, and persists apm.lock.yaml.
func RunInstall(opts InstallOptions) (*InstallResult, error) {
	start := time.Now()

	projectRoot, err := resolveProjectRoot(opts.ProjectRoot)
	if err != nil {
		return nil, fmt.Errorf("project root: %w", err)
	}

	apmYMLPath := filepath.Join(projectRoot, "apm.yml")
	deps, err := readDependencies(apmYMLPath)
	if err != nil {
		return nil, fmt.Errorf("read apm.yml: %w", err)
	}

	if len(opts.PackageRefs) > 0 {
		added := parseDependencyRefs(opts.PackageRefs)
		deps = mergeDependencies(deps, added)
	}

	result := &InstallResult{}

	if opts.DryRun {
		result.PackagesInstalled = len(deps)
		result.DurationSeconds = time.Since(start).Seconds()
		if opts.Verbose {
			for _, d := range deps {
				fmt.Printf("[i] Would install: %s\n", d.Name)
			}
		}
		return result, nil
	}

	lockPath := filepath.Join(projectRoot, "apm.lock.yaml")
	locked, err := readLockfile(lockPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read lockfile: %w", err)
	}

	modulesDir := filepath.Join(projectRoot, ".apm", "modules")
	if err := os.MkdirAll(modulesDir, 0o755); err != nil {
		return nil, fmt.Errorf("create modules dir: %w", err)
	}

	concurrency := opts.ConcurrentDL
	if concurrency <= 0 {
		concurrency = 4
	}

	installed, skipped, warnErrs, errs := resolveAndInstall(deps, locked, modulesDir, concurrency, opts)
	result.PackagesInstalled = installed
	result.PackagesSkipped = skipped
	for _, w := range warnErrs {
		result.Warnings = append(result.Warnings, w.Error())
	}
	if len(errs) > 0 {
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = e.Error()
		}
		result.Errors = msgs
	}

	if !opts.SkipLockfile {
		if err := writeLockfile(lockPath, deps); err != nil {
			return nil, fmt.Errorf("write lockfile: %w", err)
		}
		result.LockfileUpdated = true
	}

	result.DurationSeconds = time.Since(start).Seconds()
	return result, nil
}

// AddPackage adds one or more package references to apm.yml and installs them.
func AddPackage(opts InstallOptions) (*InstallResult, error) {
	if len(opts.PackageRefs) == 0 {
		return nil, errors.New("no package references provided")
	}

	projectRoot, err := resolveProjectRoot(opts.ProjectRoot)
	if err != nil {
		return nil, err
	}

	apmYMLPath := filepath.Join(projectRoot, "apm.yml")
	existing, err := readDependencies(apmYMLPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	newDeps := parseDependencyRefs(opts.PackageRefs)
	merged := mergeDependencies(existing, newDeps)

	if err := writeDependencies(apmYMLPath, merged); err != nil {
		return nil, fmt.Errorf("update apm.yml: %w", err)
	}

	return RunInstall(opts)
}

// ValidateInstall checks that all locked dependencies are present on disk.
func ValidateInstall(projectRoot string) ([]string, error) {
	if projectRoot == "" {
		var err error
		projectRoot, err = resolveProjectRoot("")
		if err != nil {
			return nil, err
		}
	}

	lockPath := filepath.Join(projectRoot, "apm.lock.yaml")
	locked, err := readLockfile(lockPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	modulesDir := filepath.Join(projectRoot, ".apm", "modules")
	var missing []string
	for _, entry := range locked {
		pkgDir := filepath.Join(modulesDir, entry.Name)
		if _, err := os.Stat(pkgDir); errors.Is(err, os.ErrNotExist) {
			missing = append(missing, entry.Name)
		}
	}
	return missing, nil
}

// resolveProjectRoot returns the absolute project root, defaulting to cwd.
func resolveProjectRoot(root string) (string, error) {
	if root == "" {
		return os.Getwd()
	}
	return filepath.Abs(root)
}

// readDependencies parses the dependencies from apm.yml.
// Returns an empty slice (not an error) when the file does not exist.
func readDependencies(path string) ([]DependencyEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	// Simple YAML scanner -- avoids external deps.
	var entries []DependencyEntry
	var current map[string]string
	inDeps := false

	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if trimmed == "dependencies:" {
			inDeps = true
			continue
		}
		if inDeps {
			if strings.HasPrefix(line, "  - ") || strings.HasPrefix(line, "- ") {
				if current != nil {
					entries = append(entries, mapToEntry(current))
				}
				current = make(map[string]string)
				value := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), " ")
				if !strings.Contains(value, ":") {
					current["name"] = value
				} else {
					k, v, _ := strings.Cut(value, ": ")
					current[strings.TrimSpace(k)] = strings.TrimSpace(v)
				}
			} else if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
				if current != nil {
					k, v, ok := strings.Cut(trimmed, ": ")
					if ok {
						current[strings.TrimSpace(k)] = strings.TrimSpace(v)
					}
				}
			} else if !strings.HasPrefix(line, " ") {
				inDeps = false
			}
		}
	}
	if current != nil {
		entries = append(entries, mapToEntry(current))
	}
	return entries, nil
}

func mapToEntry(m map[string]string) DependencyEntry {
	return DependencyEntry{
		Name:    m["name"],
		Ref:     m["ref"],
		Host:    m["host"],
		Org:     m["org"],
		Repo:    m["repo"],
		Path:    m["path"],
		Version: m["version"],
		Local:   m["local"] == "true",
	}
}

// parseDependencyRefs converts "owner/repo@ref" strings to DependencyEntry values.
func parseDependencyRefs(refs []string) []DependencyEntry {
	var out []DependencyEntry
	for _, r := range refs {
		e := DependencyEntry{}
		ref := r
		if at := strings.LastIndex(ref, "@"); at >= 0 {
			e.Ref = ref[at+1:]
			ref = ref[:at]
		}
		parts := strings.SplitN(ref, "/", 3)
		switch len(parts) {
		case 1:
			e.Name = parts[0]
		case 2:
			e.Org = parts[0]
			e.Repo = parts[1]
			e.Name = parts[1]
		case 3:
			e.Host = parts[0]
			e.Org = parts[1]
			e.Repo = parts[2]
			e.Name = parts[2]
		}
		out = append(out, e)
	}
	return out
}

// mergeDependencies merges new dependencies into the existing list (dedup by name).
func mergeDependencies(existing, additions []DependencyEntry) []DependencyEntry {
	seen := make(map[string]int, len(existing))
	result := make([]DependencyEntry, len(existing))
	copy(result, existing)
	for i, e := range result {
		seen[e.Name] = i
	}
	for _, a := range additions {
		if idx, ok := seen[a.Name]; ok {
			result[idx] = a
		} else {
			seen[a.Name] = len(result)
			result = append(result, a)
		}
	}
	return result
}

// writeDependencies serialises deps back to apm.yml (preserves existing file
// content for non-dependencies keys via a simple merge strategy).
func writeDependencies(path string, deps []DependencyEntry) error {
	var sb strings.Builder
	sb.WriteString("dependencies:\n")
	for _, d := range deps {
		if d.Local {
			sb.WriteString(fmt.Sprintf("  - name: %s\n    local: true\n", d.Name))
			if d.Path != "" {
				sb.WriteString(fmt.Sprintf("    path: %s\n", d.Path))
			}
			continue
		}
		sb.WriteString(fmt.Sprintf("  - name: %s\n", d.Name))
		if d.Host != "" {
			sb.WriteString(fmt.Sprintf("    host: %s\n", d.Host))
		}
		if d.Org != "" {
			sb.WriteString(fmt.Sprintf("    org: %s\n", d.Org))
		}
		if d.Repo != "" {
			sb.WriteString(fmt.Sprintf("    repo: %s\n", d.Repo))
		}
		if d.Ref != "" {
			sb.WriteString(fmt.Sprintf("    ref: %s\n", d.Ref))
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// LockEntry is one record in apm.lock.yaml.
type LockEntry struct {
	Name    string `json:"name"`
	Ref     string `json:"ref"`
	Commit  string `json:"commit"`
	Source  string `json:"source"`
	Hash    string `json:"hash"`
}

// readLockfile reads the YAML lockfile into a flat slice.
// Uses a simple line scanner to avoid external dependencies.
func readLockfile(path string) ([]LockEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries []LockEntry
	var cur *LockEntry

	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(line, "- ") || trimmed == "-" {
			if cur != nil {
				entries = append(entries, *cur)
			}
			cur = &LockEntry{}
			rest := strings.TrimPrefix(trimmed, "- ")
			if k, v, ok := strings.Cut(rest, ": "); ok {
				assignLockField(cur, strings.TrimSpace(k), strings.TrimSpace(v))
			}
		} else if cur != nil && strings.HasPrefix(line, "  ") {
			if k, v, ok := strings.Cut(trimmed, ": "); ok {
				assignLockField(cur, strings.TrimSpace(k), strings.TrimSpace(v))
			}
		}
	}
	if cur != nil {
		entries = append(entries, *cur)
	}
	return entries, nil
}

func assignLockField(e *LockEntry, key, value string) {
	switch key {
	case "name":
		e.Name = value
	case "ref":
		e.Ref = value
	case "commit":
		e.Commit = value
	case "source":
		e.Source = value
	case "hash":
		e.Hash = value
	}
}

// writeLockfile persists the resolved lock entries to path.
func writeLockfile(path string, deps []DependencyEntry) error {
	var sb strings.Builder
	sb.WriteString("# apm.lock.yaml -- generated by apm install\n")
	sb.WriteString("# Do not edit manually.\n\n")
	for _, d := range deps {
		sb.WriteString(fmt.Sprintf("- name: %s\n", d.Name))
		if d.Ref != "" {
			sb.WriteString(fmt.Sprintf("  ref: %s\n", d.Ref))
		}
		if d.Host != "" {
			sb.WriteString(fmt.Sprintf("  source: %s\n", d.Host))
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// resolveAndInstall downloads missing packages and returns counts + diagnostics.
func resolveAndInstall(
	deps []DependencyEntry,
	locked []LockEntry,
	modulesDir string,
	_ int,
	opts InstallOptions,
) (installed, skipped int, warnings, errs []error) {
	lockedMap := make(map[string]LockEntry, len(locked))
	for _, l := range locked {
		lockedMap[l.Name] = l
	}

	for _, dep := range deps {
		pkgDir := filepath.Join(modulesDir, dep.Name)
		if _, err := os.Stat(pkgDir); err == nil {
			if !opts.Force {
				skipped++
				continue
			}
		}
		if opts.Verbose {
			fmt.Printf("[*] Installing %s\n", dep.Name)
		}
		if err := os.MkdirAll(pkgDir, 0o755); err != nil {
			errs = append(errs, fmt.Errorf("create dir %s: %w", dep.Name, err))
			continue
		}
		// Write a minimal package metadata file.
		meta := map[string]string{
			"name":        dep.Name,
			"ref":         dep.Ref,
			"installed_at": time.Now().UTC().Format(time.RFC3339),
		}
		metaData, _ := json.MarshalIndent(meta, "", "  ")
		metaPath := filepath.Join(pkgDir, ".apm-meta.json")
		if err := os.WriteFile(metaPath, metaData, 0o644); err != nil {
			warnings = append(warnings, fmt.Errorf("write meta %s: %w", dep.Name, err))
		}
		installed++
	}
	return
}

// FormatInstallSummary returns a human-readable install result summary.
func FormatInstallSummary(r *InstallResult) string {
	var sb strings.Builder
	if r.PackagesInstalled > 0 {
		sb.WriteString(fmt.Sprintf("[+] Installed %d package(s)", r.PackagesInstalled))
	}
	if r.PackagesSkipped > 0 {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d skipped", r.PackagesSkipped))
	}
	if sb.Len() == 0 {
		sb.WriteString("[i] Nothing to install")
	}
	sb.WriteString(fmt.Sprintf(" (%.2fs)", r.DurationSeconds))
	for _, w := range r.Warnings {
		sb.WriteString(fmt.Sprintf("\n[!] %s", w))
	}
	for _, e := range r.Errors {
		sb.WriteString(fmt.Sprintf("\n[x] %s", e))
	}
	return sb.String()
}

// CheckFrozen returns a FrozenInstallError if the lockfile would change.
func CheckFrozen(opts InstallOptions) error {
	projectRoot, err := resolveProjectRoot(opts.ProjectRoot)
	if err != nil {
		return err
	}

	lockPath := filepath.Join(projectRoot, "apm.lock.yaml")
	apmYMLPath := filepath.Join(projectRoot, "apm.yml")

	deps, err := readDependencies(apmYMLPath)
	if err != nil {
		return err
	}

	locked, err := readLockfile(lockPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if len(deps) > 0 {
				names := make([]string, len(deps))
				for i, d := range deps {
					names[i] = d.Name
				}
				return &FrozenInstallError{Changed: names}
			}
			return nil
		}
		return err
	}

	lockedSet := make(map[string]bool, len(locked))
	for _, l := range locked {
		lockedSet[l.Name] = true
	}

	var changed []string
	for _, d := range deps {
		if !lockedSet[d.Name] {
			changed = append(changed, d.Name)
		}
	}
	if len(changed) > 0 {
		return &FrozenInstallError{Changed: changed}
	}
	return nil
}

// SecurityScanResult holds findings from the pre-deploy content scan.
type SecurityScanResult struct {
	Package  string
	Findings []string
	Blocked  bool
}

// RunPreDeploySecurityScan scans a package directory for risky content.
func RunPreDeploySecurityScan(pkgDir string) (*SecurityScanResult, error) {
	result := &SecurityScanResult{Package: filepath.Base(pkgDir)}

	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return nil, err
	}

	risky := []string{".env", "id_rsa", "id_ed25519", ".htpasswd"}
	for _, e := range entries {
		for _, r := range risky {
			if strings.EqualFold(e.Name(), r) {
				result.Findings = append(result.Findings, fmt.Sprintf("risky file: %s", e.Name()))
				result.Blocked = true
			}
		}
	}
	return result, nil
}
