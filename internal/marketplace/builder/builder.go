// Package builder provides the MarketplaceBuilder: load, resolve, compose, and write marketplace.json.
// Migrated from src/apm_cli/marketplace/builder.py.
package builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/githubnext/apm/internal/marketplace/mkio"
	"github.com/githubnext/apm/internal/marketplace/refresolver"
	"github.com/githubnext/apm/internal/marketplace/semver"
	"github.com/githubnext/apm/internal/marketplace/tagpattern"
	"github.com/githubnext/apm/internal/marketplace/ymlschema"

	"os"
	"path"
)

// BuildDiagnostic is a structured diagnostic emitted during marketplace.json composition.
type BuildDiagnostic struct {
	Level   string // "warning" | "verbose"
	Message string
}

// ResolvedPackage is a package entry after ref resolution.
type ResolvedPackage struct {
	Name             string
	SourceRepo       string // "owner/repo" only
	Subdir           string // APM-only (for git-subdir source object)
	Ref              string // resolved tag name, e.g. "v1.2.0"
	SHA              string // 40-char git SHA
	RequestedVersion string // original APM-only range (for diagnostics)
	Tags             []string
	IsPrerelease     bool // True if the resolved ref was a prerelease semver
}

// ResolveResult is the result of resolving package refs in a marketplace build.
type ResolveResult struct {
	Entries []ResolvedPackage
	Errors  [][2]string // (package name, error message) pairs
}

// OK returns true when every package resolved without error.
func (r ResolveResult) OK() bool { return len(r.Errors) == 0 }

// BuildReport summarizes a build run.
type BuildReport struct {
	Resolved       []ResolvedPackage
	Errors         [][2]string
	Warnings       []string
	Diagnostics    []BuildDiagnostic
	UnchangedCount int
	AddedCount     int
	UpdatedCount   int
	RemovedCount   int
	OutputPath     string
	DryRun         bool
}

// BuildOptions holds configuration knobs for MarketplaceBuilder.
type BuildOptions struct {
	Concurrency      int
	TimeoutSeconds   float64
	IncludePrerelease bool
	AllowHead        bool
	ContinueOnError  bool
	Offline          bool
	OutputOverride   string
	DryRun           bool
}

// DefaultBuildOptions returns sensible defaults.
func DefaultBuildOptions() BuildOptions {
	return BuildOptions{
		Concurrency:    8,
		TimeoutSeconds: 10.0,
	}
}

// sha40RE matches a 40-char hex SHA.
var sha40RE = regexp.MustCompile(`^[0-9a-f]{40}$`)

// versionRangeChars are chars that indicate a range constraint rather than a display version.
var versionRangeChars = []byte{'^', '~', '>', '<', '='}

func isDisplayVersion(version string) bool {
	if version == "" {
		return false
	}
	v := strings.TrimSpace(version)
	for _, c := range versionRangeChars {
		if v[0] == c {
			return false
		}
	}
	if strings.ContainsAny(v, " *") {
		return false
	}
	parts := strings.Split(v, ".")
	if len(parts) == 0 {
		return false
	}
	last := strings.ToLower(parts[len(parts)-1])
	if last == "x" {
		return false
	}
	return true
}

// subtractPluginRoot removes pluginRoot prefix from a local source path.
func subtractPluginRoot(src, pluginRoot string) (string, error) {
	normSrc := strings.TrimRight(strings.TrimLeft(src, "./"), "/")
	normRoot := strings.TrimRight(strings.TrimLeft(pluginRoot, "./"), "/")
	if !strings.HasPrefix(normSrc, normRoot) {
		return "", fmt.Errorf("source '%s' does not start with pluginRoot '%s'", src, pluginRoot)
	}
	rel := strings.TrimPrefix(normSrc, normRoot)
	rel = strings.TrimLeft(rel, "/")
	if rel == "" || rel == "." {
		return "", fmt.Errorf("subtracting pluginRoot '%s' from source '%s' yields empty path", pluginRoot, src)
	}
	if strings.HasPrefix(rel, "/") {
		return "", fmt.Errorf("pluginRoot subtraction produced absolute path: '%s'", rel)
	}
	for _, seg := range strings.Split(rel, "/") {
		if seg == ".." {
			return "", fmt.Errorf("pluginRoot subtraction produced path with traversal: '%s'", rel)
		}
	}
	return "./" + rel, nil
}

// BuildError is raised on build failures.
type BuildError struct {
	Msg     string
	Package string
}

func (e *BuildError) Error() string { return e.Msg }

// HeadNotAllowedError is raised when a branch ref is resolved without allow_head.
type HeadNotAllowedError struct {
	Package string
	Ref     string
}

func (e *HeadNotAllowedError) Error() string {
	return fmt.Sprintf("package '%s': ref '%s' is a branch head; use allow_head to allow it", e.Package, e.Ref)
}

// RefNotFoundError is raised when a ref cannot be found on the remote.
type RefNotFoundError struct {
	Package    string
	Ref        string
	OwnerRepo  string
}

func (e *RefNotFoundError) Error() string {
	return fmt.Sprintf("package '%s': ref '%s' not found on remote '%s'", e.Package, e.Ref, e.OwnerRepo)
}

// NoMatchingVersionError is raised when no tag satisfies the semver range.
type NoMatchingVersionError struct {
	Package      string
	VersionRange string
	Detail       string
}

func (e *NoMatchingVersionError) Error() string {
	return fmt.Sprintf("package '%s': no tag satisfies '%s' (%s)", e.Package, e.VersionRange, e.Detail)
}

// MarketplaceBuilder loads, resolves, composes, and writes marketplace.json.
type MarketplaceBuilder struct {
	ymlPath     string
	projectRoot string
	options     BuildOptions
	yml         *ymlschema.MarketplaceConfig
	resolver    *refresolver.RefResolver
	githubToken string
	host        string
	authResolved bool

	composeWarnings    []string
	composeDiagnostics []BuildDiagnostic
}

// New constructs a MarketplaceBuilder for the given marketplace.yml path.
func New(marketplaceYMLPath string, options BuildOptions) *MarketplaceBuilder {
	return &MarketplaceBuilder{
		ymlPath:     marketplaceYMLPath,
		projectRoot: filepath.Dir(marketplaceYMLPath),
		options:     options,
		host:        "github.com",
	}
}

// FromConfig constructs a builder from an already-loaded MarketplaceConfig.
func FromConfig(config *ymlschema.MarketplaceConfig, projectRoot string, options BuildOptions) *MarketplaceBuilder {
	b := &MarketplaceBuilder{
		ymlPath:     filepath.Join(projectRoot, "apm.yml"),
		projectRoot: projectRoot,
		options:     options,
		yml:         config,
		host:        "github.com",
	}
	return b
}

func (b *MarketplaceBuilder) loadYML() (*ymlschema.MarketplaceConfig, error) {
	if b.yml != nil {
		return b.yml, nil
	}
	isLegacy := path.Base(b.ymlPath) != "apm.yml"
	cfg, err := ymlschema.LoadFromFile(b.ymlPath, isLegacy)
	if err != nil {
		return nil, err
	}
	b.yml = cfg
	return b.yml, nil
}

func (b *MarketplaceBuilder) ensureAuth() {
	if b.authResolved {
		return
	}
	if b.options.Offline {
		b.authResolved = true
		return
	}
	// Resolve GitHub token from env
	for _, envVar := range []string{"GITHUB_APM_PAT", "GITHUB_TOKEN", "GH_TOKEN"} {
		if t := os.Getenv(envVar); t != "" {
			b.githubToken = t
			break
		}
	}
	b.authResolved = true
}

func (b *MarketplaceBuilder) getResolver() *refresolver.RefResolver {
	if b.resolver == nil {
		b.ensureAuth()
		b.resolver = refresolver.New(b.options.TimeoutSeconds, b.options.Offline, b.host, b.githubToken)
	}
	return b.resolver
}

func (b *MarketplaceBuilder) outputPath(yml *ymlschema.MarketplaceConfig) (string, error) {
	if b.options.OutputOverride != "" {
		return b.options.OutputOverride, nil
	}
	outputPath := filepath.Join(b.projectRoot, yml.Output)
	// containment guard
	rel, err := filepath.Rel(b.projectRoot, outputPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", &BuildError{Msg: fmt.Sprintf("output path '%s' escapes project root", outputPath)}
	}
	return outputPath, nil
}

// stripRefPrefix removes refs/tags/ or refs/heads/ prefix.
func stripRefPrefix(refname string) string {
	if strings.HasPrefix(refname, "refs/tags/") {
		return refname[len("refs/tags/"):]
	}
	if strings.HasPrefix(refname, "refs/heads/") {
		return refname[len("refs/heads/"):]
	}
	return refname
}

// resolveExplicitRef resolves an entry with an explicit ref: field.
func (b *MarketplaceBuilder) resolveExplicitRef(entry ymlschema.PackageEntry, resolver *refresolver.RefResolver) (ResolvedPackage, error) {
	refText := entry.Ref
	ownerRepo := entry.Source

	if sha40RE.MatchString(refText) {
		sv, _ := semver.Parse(strings.TrimLeft(refText, "vV"))
		isPrerelease := sv.Prerelease != ""
		return ResolvedPackage{
			Name:             entry.Name,
			SourceRepo:       ownerRepo,
			Subdir:           entry.Subdir,
			Ref:              refText,
			SHA:              refText,
			RequestedVersion: entry.Version,
			Tags:             entry.Tags,
			IsPrerelease:     isPrerelease,
		}, nil
	}

	refs, err := resolver.ListRemoteRefs(ownerRepo)
	if err != nil {
		return ResolvedPackage{}, &BuildError{Msg: err.Error(), Package: entry.Name}
	}

	// Try as tag first
	for _, rr := range refs {
		if !strings.HasPrefix(rr.Name, "refs/tags/") {
			continue
		}
		tagName := stripRefPrefix(rr.Name)
		if tagName == refText {
			sv, _ := semver.Parse(strings.TrimLeft(tagName, "vV"))
			return ResolvedPackage{
				Name:             entry.Name,
				SourceRepo:       ownerRepo,
				Subdir:           entry.Subdir,
				Ref:              tagName,
				SHA:              rr.SHA,
				RequestedVersion: entry.Version,
				Tags:             entry.Tags,
				IsPrerelease:     sv.Prerelease != "",
			}, nil
		}
	}

	// Try as full refname
	for _, rr := range refs {
		if rr.Name == refText {
			short := stripRefPrefix(rr.Name)
			isBranch := strings.HasPrefix(rr.Name, "refs/heads/")
			if isBranch && !b.options.AllowHead {
				return ResolvedPackage{}, &HeadNotAllowedError{Package: entry.Name, Ref: short}
			}
			sv, _ := semver.Parse(strings.TrimLeft(short, "vV"))
			return ResolvedPackage{
				Name:             entry.Name,
				SourceRepo:       ownerRepo,
				Subdir:           entry.Subdir,
				Ref:              short,
				SHA:              rr.SHA,
				RequestedVersion: entry.Version,
				Tags:             entry.Tags,
				IsPrerelease:     sv.Prerelease != "",
			}, nil
		}
	}

	// Try as branch name
	for _, rr := range refs {
		if rr.Name == "refs/heads/"+refText {
			if !b.options.AllowHead {
				return ResolvedPackage{}, &HeadNotAllowedError{Package: entry.Name, Ref: refText}
			}
			return ResolvedPackage{
				Name:             entry.Name,
				SourceRepo:       ownerRepo,
				Subdir:           entry.Subdir,
				Ref:              refText,
				SHA:              rr.SHA,
				RequestedVersion: entry.Version,
				Tags:             entry.Tags,
				IsPrerelease:     false,
			}, nil
		}
	}

	if strings.ToUpper(refText) == "HEAD" && !b.options.AllowHead {
		return ResolvedPackage{}, &HeadNotAllowedError{Package: entry.Name, Ref: "HEAD"}
	}
	return ResolvedPackage{}, &RefNotFoundError{Package: entry.Name, Ref: refText, OwnerRepo: ownerRepo}
}

// resolveVersionRange resolves an entry using its version: semver range.
func (b *MarketplaceBuilder) resolveVersionRange(entry ymlschema.PackageEntry, resolver *refresolver.RefResolver, yml *ymlschema.MarketplaceConfig) (ResolvedPackage, error) {
	versionRange := entry.Version
	ownerRepo := entry.Source

	pattern := entry.TagPattern
	if pattern == "" {
		pattern = yml.Build.TagPattern
	}
	if pattern == "" {
		pattern = "v{version}"
	}

	tagRx, err := tagpattern.BuildTagRegex(pattern)
	if err != nil {
		return ResolvedPackage{}, &BuildError{Msg: fmt.Sprintf("invalid tag pattern '%s': %v", pattern, err), Package: entry.Name}
	}

	refs, err := resolver.ListRemoteRefs(ownerRepo)
	if err != nil {
		return ResolvedPackage{}, &BuildError{Msg: err.Error(), Package: entry.Name}
	}

	type candidate struct {
		sv      semver.SemVer
		tagName string
		sha     string
	}
	var candidates []candidate

	for _, rr := range refs {
		if !strings.HasPrefix(rr.Name, "refs/tags/") {
			continue
		}
		tagName := rr.Name[len("refs/tags/"):]
		versionStr, ok := tagpattern.ExtractVersion(tagRx, tagName)
		if !ok {
			continue
		}
		sv, err := semver.Parse(versionStr)
		if err != nil {
			continue
		}
		includePrerelease := entry.IncludePrerelease || b.options.IncludePrerelease
		if sv.Prerelease != "" && !includePrerelease {
			continue
		}
		if semver.SatisfiesRange(sv, versionRange) {
			candidates = append(candidates, candidate{sv: sv, tagName: tagName, sha: rr.SHA})
		}
	}

	if len(candidates) == 0 {
		return ResolvedPackage{}, &NoMatchingVersionError{
			Package:      entry.Name,
			VersionRange: versionRange,
			Detail:       fmt.Sprintf("pattern='%s', remote='%s'", pattern, ownerRepo),
		}
	}

	// Pick highest
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.sv.Compare(best.sv) > 0 {
			best = c
		}
	}

	return ResolvedPackage{
		Name:             entry.Name,
		SourceRepo:       ownerRepo,
		Subdir:           entry.Subdir,
		Ref:              best.tagName,
		SHA:              best.sha,
		RequestedVersion: versionRange,
		Tags:             entry.Tags,
		IsPrerelease:     best.sv.Prerelease != "",
	}, nil
}

// resolveEntry resolves a single package entry to a concrete tag + SHA.
func (b *MarketplaceBuilder) resolveEntry(entry ymlschema.PackageEntry, yml *ymlschema.MarketplaceConfig) (ResolvedPackage, error) {
	if entry.IsLocal {
		return ResolvedPackage{
			Name:             entry.Name,
			SourceRepo:       "",
			Subdir:           entry.Source,
			Ref:              "",
			SHA:              "",
			RequestedVersion: entry.Version,
			Tags:             entry.Tags,
			IsPrerelease:     false,
		}, nil
	}
	resolver := b.getResolver()
	if entry.Ref != "" {
		return b.resolveExplicitRef(entry, resolver)
	}
	return b.resolveVersionRange(entry, resolver, yml)
}

// Resolve resolves every entry concurrently.
func (b *MarketplaceBuilder) Resolve() (ResolveResult, error) {
	yml, err := b.loadYML()
	if err != nil {
		return ResolveResult{}, err
	}
	entries := yml.Packages
	if len(entries) == 0 {
		return ResolveResult{}, nil
	}

	// Eagerly create the resolver before spawning goroutines
	b.getResolver()

	type indexedResult struct {
		idx     int
		pkg     ResolvedPackage
		errPair [2]string
		hasErr  bool
	}

	sem := make(chan struct{}, b.options.Concurrency)
	if b.options.Concurrency <= 0 {
		sem = make(chan struct{}, 8)
	}

	resultCh := make(chan indexedResult, len(entries))
	var wg sync.WaitGroup

	for i, entry := range entries {
		wg.Add(1)
		go func(idx int, e ymlschema.PackageEntry) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			pkg, resolveErr := b.resolveEntry(e, yml)
			if resolveErr != nil {
				var buildErr *BuildError
				var headErr *HeadNotAllowedError
				var refErr *RefNotFoundError
				var noMatchErr *NoMatchingVersionError
				if errors.As(resolveErr, &buildErr) || errors.As(resolveErr, &headErr) ||
					errors.As(resolveErr, &refErr) || errors.As(resolveErr, &noMatchErr) {
					resultCh <- indexedResult{idx: idx, errPair: [2]string{e.Name, resolveErr.Error()}, hasErr: true}
					return
				}
				resultCh <- indexedResult{idx: idx, errPair: [2]string{e.Name, resolveErr.Error()}, hasErr: true}
				return
			}
			resultCh <- indexedResult{idx: idx, pkg: pkg}
		}(i, entry)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	results := make(map[int]ResolvedPackage)
	var errs [][2]string
	var firstErr error

	for r := range resultCh {
		if r.hasErr {
			errs = append(errs, r.errPair)
			if !b.options.ContinueOnError && firstErr == nil {
				firstErr = fmt.Errorf("error resolving '%s': %s", r.errPair[0], r.errPair[1])
			}
		} else {
			results[r.idx] = r.pkg
		}
	}

	if firstErr != nil {
		return ResolveResult{}, firstErr
	}

	ordered := make([]ResolvedPackage, 0, len(results))
	for idx := range entries {
		if pkg, ok := results[idx]; ok {
			ordered = append(ordered, pkg)
		}
	}

	return ResolveResult{Entries: ordered, Errors: errs}, nil
}

// ComposeMarketplaceJSON produces an Anthropic-compliant marketplace.json dict.
func (b *MarketplaceBuilder) ComposeMarketplaceJSON(resolved []ResolvedPackage) (map[string]interface{}, error) {
	yml, err := b.loadYML()
	if err != nil {
		return nil, err
	}

	entryByName := make(map[string]*ymlschema.PackageEntry)
	for i := range yml.Packages {
		entryByName[yml.Packages[i].Name] = &yml.Packages[i]
	}

	doc := make(map[string]interface{})
	doc["name"] = yml.Name
	if yml.DescriptionOverridden && yml.Description != "" {
		doc["description"] = yml.Description
	}
	if yml.VersionOverridden && yml.Version != "" {
		doc["version"] = yml.Version
	}

	ownerDict := make(map[string]interface{})
	ownerDict["name"] = yml.Owner.Name
	if yml.Owner.Email != "" {
		ownerDict["email"] = yml.Owner.Email
	}
	if yml.Owner.URL != "" {
		ownerDict["url"] = yml.Owner.URL
	}
	doc["owner"] = ownerDict

	if len(yml.Metadata) > 0 {
		doc["metadata"] = yml.Metadata
	}

	var plugins []interface{}
	var diagnostics []BuildDiagnostic
	pluginRoot := ""
	if m, ok := yml.Metadata["pluginRoot"]; ok {
		if s, ok := m.(string); ok {
			pluginRoot = s
		}
	}
	stripCount := 0
	overrideCount := 0

	for _, pkg := range resolved {
		plugin := make(map[string]interface{})
		plugin["name"] = pkg.Name

		entry := entryByName[pkg.Name]
		isLocal := entry != nil && entry.IsLocal

		if isLocal {
			if entry.Description != "" {
				plugin["description"] = entry.Description
			}
			if entry.Version != "" {
				plugin["version"] = entry.Version
			}
		} else {
			if entry != nil && entry.Description != "" {
				plugin["description"] = entry.Description
			}
			if entry != nil && isDisplayVersion(entry.Version) {
				plugin["version"] = entry.Version
			} else if pkg.Ref != "" && isDisplayVersion(pkg.Ref) {
				// Fallback: use resolved ref as display version if applicable
			}
		}

		if entry != nil && len(entry.Author) > 0 {
			plugin["author"] = entry.Author
		}
		if entry != nil && entry.License != "" {
			plugin["license"] = entry.License
		}
		if entry != nil && entry.Repository != "" {
			plugin["repository"] = entry.Repository
		}
		if len(pkg.Tags) > 0 {
			plugin["tags"] = pkg.Tags
		}
		if isLocal && entry != nil && entry.Homepage != "" {
			plugin["homepage"] = entry.Homepage
		}

		// source
		if isLocal {
			sourceValue := entry.Source
			if pluginRoot != "" {
				stripped, err := subtractPluginRoot(entry.Source, pluginRoot)
				if err != nil {
					// W1: source outside pluginRoot -- emit as-is
					diagnostics = append(diagnostics, BuildDiagnostic{
						Level:   "warning",
						Message: fmt.Sprintf("[!] Package '%s': source '%s' is outside pluginRoot '%s' -- emitted as-is", pkg.Name, entry.Source, pluginRoot),
					})
				} else {
					sourceValue = stripped
					stripCount++
					diagnostics = append(diagnostics, BuildDiagnostic{
						Level:   "verbose",
						Message: fmt.Sprintf("[i] Package '%s': stripped pluginRoot -- '%s' -> '%s'", pkg.Name, entry.Source, sourceValue),
					})
				}
			}
			plugin["source"] = sourceValue
		} else {
			srcObj := make(map[string]interface{})
			if pkg.Subdir != "" {
				srcObj["source"] = "git-subdir"
				srcObj["url"] = pkg.SourceRepo
				srcObj["path"] = pkg.Subdir
			} else {
				srcObj["source"] = "github"
				srcObj["repo"] = pkg.SourceRepo
			}
			if pkg.Ref != "" {
				srcObj["ref"] = pkg.Ref
			}
			if pkg.SHA != "" {
				srcObj["sha"] = pkg.SHA
			}
			plugin["source"] = srcObj
		}

		plugins = append(plugins, plugin)
	}

	_ = overrideCount
	_ = stripCount

	// Build verbose summary
	if pluginRoot != "" && stripCount > 0 {
		diagnostics = append(diagnostics, BuildDiagnostic{
			Level:   "verbose",
			Message: fmt.Sprintf("pluginRoot: stripped from %d local source(s)", stripCount),
		})
	}

	// Duplicate name check
	var buildWarnings []string
	seenNames := make(map[string]string)
	for _, p := range plugins {
		pm := p.(map[string]interface{})
		pname := pm["name"].(string)
		srcLabel := "?"
		if src, ok := pm["source"]; ok {
			switch s := src.(type) {
			case string:
				srcLabel = s
			case map[string]interface{}:
				if v, ok := s["path"]; ok {
					srcLabel = fmt.Sprintf("%v", v)
				} else if v, ok := s["repo"]; ok {
					srcLabel = fmt.Sprintf("%v", v)
				}
			}
		}
		if prev, exists := seenNames[pname]; exists {
			buildWarnings = append(buildWarnings, fmt.Sprintf("Duplicate package name '%s': '%s' and '%s'. Consumers will see duplicate entries in browse.", pname, prev, srcLabel))
		} else {
			seenNames[pname] = srcLabel
		}
	}

	b.composeWarnings = buildWarnings
	b.composeDiagnostics = diagnostics
	doc["plugins"] = plugins
	return doc, nil
}

type pluginSHAs map[string]string

func extractPluginSHAs(data map[string]interface{}) pluginSHAs {
	out := make(pluginSHAs)
	rawPlugins, _ := data["plugins"].([]interface{})
	for _, p := range rawPlugins {
		pm, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := pm["name"].(string)
		sha := ""
		switch s := pm["source"].(type) {
		case string:
			sha = s
		case map[string]interface{}:
			if v, ok := s["sha"].(string); ok {
				sha = v
			} else if v, ok := s["commit"].(string); ok {
				sha = v
			}
		}
		out[name] = sha
	}
	return out
}

func computeDiff(oldJSON, newJSON map[string]interface{}) (unchanged, added, updated, removed int) {
	if oldJSON == nil {
		return 0, len(extractPluginSHAs(newJSON)), 0, 0
	}
	oldPlugins := extractPluginSHAs(oldJSON)
	newPlugins := extractPluginSHAs(newJSON)

	for name, sha := range newPlugins {
		if _, exists := oldPlugins[name]; !exists {
			added++
		} else if oldPlugins[name] == sha {
			unchanged++
		} else {
			updated++
		}
	}
	for name := range oldPlugins {
		if _, exists := newPlugins[name]; !exists {
			removed++
		}
	}
	return
}

func serializeJSON(data map[string]interface{}) ([]byte, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

func loadExistingJSON(p string) map[string]interface{} {
	data, err := os.ReadFile(p)
	if err != nil {
		return nil
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil
	}
	return doc
}

// Build runs the full pipeline: load -> resolve -> compose -> write.
func (b *MarketplaceBuilder) Build() (BuildReport, error) {
	result, err := b.Resolve()
	if err != nil {
		return BuildReport{}, err
	}

	newJSON, err := b.ComposeMarketplaceJSON(result.Entries)
	if err != nil {
		return BuildReport{}, err
	}

	buildWarnings := b.composeWarnings
	buildDiagnostics := b.composeDiagnostics

	yml, err := b.loadYML()
	if err != nil {
		return BuildReport{}, err
	}
	outPath, err := b.outputPath(yml)
	if err != nil {
		return BuildReport{}, err
	}

	oldJSON := loadExistingJSON(outPath)
	unchanged, added, updated, removed := computeDiff(oldJSON, newJSON)

	if !b.options.DryRun {
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return BuildReport{}, err
		}
		content, err := serializeJSON(newJSON)
		if err != nil {
			return BuildReport{}, err
		}
		if err := mkio.AtomicWrite(outPath, content); err != nil {
			return BuildReport{}, err
		}
	}

	if b.resolver != nil {
		b.resolver.Close()
	}

	return BuildReport{
		Resolved:       result.Entries,
		Errors:         result.Errors,
		Warnings:       buildWarnings,
		Diagnostics:    buildDiagnostics,
		UnchangedCount: unchanged,
		AddedCount:     added,
		UpdatedCount:   updated,
		RemovedCount:   removed,
		OutputPath:     outPath,
		DryRun:         b.options.DryRun,
	}, nil
}
