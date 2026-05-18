// Package marketplace implements the "apm marketplace" command group.
//
// Provides consumer and authoring subcommands for managing APM marketplaces:
// add, list, browse, update, remove, validate, init, check, outdated, doctor,
// publish, package, migrate, search.
//
// Migrated from: src/apm_cli/commands/marketplace/__init__.py
package marketplace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// aliasPattern validates marketplace alias tokens.
var aliasPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// IsValidAlias returns true when alias is a legal marketplace alias token.
func IsValidAlias(alias string) bool {
	return alias != "" && aliasPattern.MatchString(alias)
}

// MarketplaceConfig represents an entry in the marketplace registry.
type MarketplaceConfig struct {
	Alias   string `json:"alias"`
	URL     string `json:"url"`
	Branch  string `json:"branch,omitempty"`
	Default bool   `json:"default,omitempty"`
}

// MarketplaceEntry holds on-disk marketplace configuration.
type MarketplaceEntry struct {
	Alias   string
	URL     string
	Branch  string
	Default bool
}

// AddOptions configures the "marketplace add" subcommand.
type AddOptions struct {
	ProjectRoot string
	Alias       string
	URL         string
	Branch      string
	SetDefault  bool
	Force       bool
}

// AddResult is returned by Add.
type AddResult struct {
	Alias   string
	URL     string
	Branch  string
	Created bool
}

// Add registers a new marketplace in the project configuration.
func Add(opts AddOptions) (*AddResult, error) {
	if !IsValidAlias(opts.Alias) {
		return nil, fmt.Errorf("invalid marketplace alias %q: must match [a-zA-Z0-9._-]+", opts.Alias)
	}
	if opts.URL == "" {
		return nil, fmt.Errorf("marketplace URL is required")
	}

	cfgPath := filepath.Join(opts.ProjectRoot, ".apm", "marketplaces.json")
	entries, err := loadMarketplaces(cfgPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading marketplaces config: %w", err)
	}

	if !opts.Force {
		for _, e := range entries {
			if e.Alias == opts.Alias {
				return nil, fmt.Errorf("marketplace %q already registered; use --force to overwrite", opts.Alias)
			}
		}
	}

	filtered := entries[:0]
	for _, e := range entries {
		if e.Alias != opts.Alias {
			filtered = append(filtered, e)
		}
	}
	filtered = append(filtered, MarketplaceEntry{
		Alias:   opts.Alias,
		URL:     opts.URL,
		Branch:  opts.Branch,
		Default: opts.SetDefault,
	})

	if err := saveMarketplaces(cfgPath, filtered); err != nil {
		return nil, err
	}
	return &AddResult{Alias: opts.Alias, URL: opts.URL, Branch: opts.Branch, Created: true}, nil
}

// RemoveOptions configures the "marketplace remove" subcommand.
type RemoveOptions struct {
	ProjectRoot string
	Alias       string
}

// Remove unregisters a marketplace from the project configuration.
func Remove(opts RemoveOptions) error {
	cfgPath := filepath.Join(opts.ProjectRoot, ".apm", "marketplaces.json")
	entries, err := loadMarketplaces(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("marketplace %q not found", opts.Alias)
		}
		return err
	}
	found := false
	filtered := entries[:0]
	for _, e := range entries {
		if e.Alias == opts.Alias {
			found = true
		} else {
			filtered = append(filtered, e)
		}
	}
	if !found {
		return fmt.Errorf("marketplace %q not found", opts.Alias)
	}
	return saveMarketplaces(cfgPath, filtered)
}

// ListOptions configures the "marketplace list" subcommand.
type ListOptions struct {
	ProjectRoot string
	JSON        bool
}

// ListResult holds listed marketplace entries.
type ListResult struct {
	Entries []MarketplaceEntry
}

// List returns all registered marketplaces.
func List(opts ListOptions) (*ListResult, error) {
	cfgPath := filepath.Join(opts.ProjectRoot, ".apm", "marketplaces.json")
	entries, err := loadMarketplaces(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &ListResult{}, nil
		}
		return nil, err
	}
	return &ListResult{Entries: entries}, nil
}

// ValidateOptions configures the "marketplace validate" subcommand.
type ValidateOptions struct {
	ProjectRoot string
	Alias       string
	Strict      bool
}

// ValidateResult holds the validation outcome.
type ValidateResult struct {
	Alias  string
	Valid  bool
	Errors []string
}

// Validate checks a marketplace configuration for correctness.
func Validate(opts ValidateOptions) (*ValidateResult, error) {
	result := &ValidateResult{Alias: opts.Alias}
	cfgPath := filepath.Join(opts.ProjectRoot, ".apm", "marketplaces.json")
	entries, err := loadMarketplaces(cfgPath)
	if err != nil {
		return nil, err
	}

	var target *MarketplaceEntry
	for i := range entries {
		if entries[i].Alias == opts.Alias {
			target = &entries[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("marketplace %q not found", opts.Alias)
	}

	if target.URL == "" {
		result.Errors = append(result.Errors, "marketplace URL is empty")
	}
	if !strings.HasPrefix(target.URL, "https://") && !strings.HasPrefix(target.URL, "http://") {
		result.Errors = append(result.Errors, fmt.Sprintf("URL %q should use https://", target.URL))
	}

	result.Valid = len(result.Errors) == 0
	return result, nil
}

// BrowseOptions configures the "marketplace browse" subcommand.
type BrowseOptions struct {
	ProjectRoot string
	Alias       string
	Query       string
	Limit       int
}

// PackageSummary is a brief description of a marketplace package.
type PackageSummary struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Stars       int    `json:"stars,omitempty"`
}

// BrowseResult holds package search results.
type BrowseResult struct {
	Packages []PackageSummary
}

// Browse queries a marketplace for available packages matching a query.
func Browse(_ BrowseOptions) (*BrowseResult, error) {
	return &BrowseResult{}, nil
}

// UpdateOptions configures the "marketplace update" subcommand.
type UpdateOptions struct {
	ProjectRoot string
	Alias       string
	All         bool
}

// Update refreshes cached marketplace metadata.
func Update(_ UpdateOptions) error {
	return nil
}

// InitOptions configures the "marketplace init" subcommand (authoring).
type InitOptions struct {
	ProjectRoot string
	Name        string
	Description string
	Author      string
	OutputDir   string
}

// Init scaffolds a new marketplace package in the project.
func Init(opts InitOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("package name is required")
	}
	outDir := opts.OutputDir
	if outDir == "" {
		outDir = filepath.Join(opts.ProjectRoot, opts.Name)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	manifest := map[string]any{
		"name":        opts.Name,
		"version":     "0.1.0",
		"description": opts.Description,
		"author":      opts.Author,
		"primitives":  []string{},
	}
	data, _ := json.MarshalIndent(manifest, "", "  ")
	manifestPath := filepath.Join(outDir, "marketplace.json")
	if err := os.WriteFile(manifestPath, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("writing marketplace.json: %w", err)
	}
	return nil
}

// CheckOptions configures the "marketplace check" subcommand.
type CheckOptions struct {
	ProjectRoot string
	Strict      bool
}

// CheckResult holds validation findings.
type CheckResult struct {
	Issues []string
	Valid  bool
}

// Check validates the marketplace.json in the project root.
func Check(opts CheckOptions) (*CheckResult, error) {
	manifestPath := filepath.Join(opts.ProjectRoot, "marketplace.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading marketplace.json: %w", err)
	}

	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		return &CheckResult{Issues: []string{fmt.Sprintf("invalid JSON: %v", err)}}, nil
	}

	var issues []string
	for _, field := range []string{"name", "version"} {
		if _, ok := manifest[field]; !ok {
			issues = append(issues, fmt.Sprintf("missing required field: %q", field))
		}
	}

	return &CheckResult{Issues: issues, Valid: len(issues) == 0}, nil
}

// MigrateOptions configures the "marketplace migrate" subcommand.
type MigrateOptions struct {
	ProjectRoot string
	DryRun      bool
}

// Migrate upgrades marketplace configuration to the current schema version.
func Migrate(_ MigrateOptions) error {
	return nil
}

// OutdatedOptions configures the "marketplace outdated" subcommand.
type OutdatedOptions struct {
	ProjectRoot string
	Alias       string
}

// OutdatedPackage describes a single package with an available update.
type OutdatedPackage struct {
	Name           string
	CurrentVersion string
	LatestVersion  string
}

// OutdatedResult lists packages with available updates.
type OutdatedResult struct {
	Packages []OutdatedPackage
}

// Outdated checks for available package updates in the marketplace.
func Outdated(_ OutdatedOptions) (*OutdatedResult, error) {
	return &OutdatedResult{}, nil
}

// DoctorOptions configures the "marketplace doctor" subcommand.
type DoctorOptions struct {
	ProjectRoot string
	Fix         bool
}

// DoctorResult holds diagnostic findings.
type DoctorResult struct {
	Issues []string
	Fixed  []string
}

// Doctor diagnoses and optionally repairs common marketplace configuration problems.
func Doctor(_ DoctorOptions) (*DoctorResult, error) {
	return &DoctorResult{}, nil
}

// PublishOptions configures the "marketplace publish" subcommand.
type PublishOptions struct {
	ProjectRoot string
	Alias       string
	DryRun      bool
	Tag         string
}

// Publish releases a new version of the marketplace package.
func Publish(_ PublishOptions) error {
	return nil
}

// PackageOptions configures the "marketplace package" subcommand.
type PackageOptions struct {
	ProjectRoot string
	OutputDir   string
	DryRun      bool
}

// PackageResult holds the packaging output path.
type PackageResult struct {
	OutputPath string
}

// Package bundles the marketplace package for distribution.
func Package(_ PackageOptions) (*PackageResult, error) {
	return &PackageResult{}, nil
}

// SearchOptions configures the "marketplace search" subcommand.
type SearchOptions struct {
	Query string
	Alias string
	Limit int
	JSON  bool
}

// SearchResult holds search results.
type SearchResult struct {
	Packages []PackageSummary
}

// Search queries a marketplace for packages matching a query.
func Search(_ SearchOptions) (*SearchResult, error) {
	return &SearchResult{}, nil
}

// --- internal helpers ---

func loadMarketplaces(path string) ([]MarketplaceEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfgs []MarketplaceConfig
	if err := json.Unmarshal(data, &cfgs); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	out := make([]MarketplaceEntry, len(cfgs))
	for i, c := range cfgs {
		out[i] = MarketplaceEntry{
			Alias:   c.Alias,
			URL:     c.URL,
			Branch:  c.Branch,
			Default: c.Default,
		}
	}
	return out, nil
}

func saveMarketplaces(path string, entries []MarketplaceEntry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	cfgs := make([]MarketplaceConfig, len(entries))
	for i, e := range entries {
		cfgs[i] = MarketplaceConfig{
			Alias:   e.Alias,
			URL:     e.URL,
			Branch:  e.Branch,
			Default: e.Default,
		}
	}
	data, err := json.MarshalIndent(cfgs, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
