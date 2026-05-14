// Package skillintegrator provides skill integration for APM packages.
// Deploys SKILL.md-based packages to .github/skills/, .claude/skills/, etc.
// Ported from src/apm_cli/integration/skill_integrator.py
package skillintegrator

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/githubnext/apm/internal/integration/targets"
)

// SkillIntegrationResult holds results of a skill integration operation.
type SkillIntegrationResult struct {
	SkillCreated      bool
	SkillUpdated      bool
	SkillSkipped      bool
	SkillPath         string // path to deployed SKILL.md, empty if not deployed
	ReferencesCopied  int    // total files copied to skill directory
	LinksResolved     int    // always 0 (kept for backward compat)
	SubSkillsPromoted int    // number of sub-skills promoted to top-level
	TargetPaths       []string
}

// nameRe matches valid agentskills.io skill names.
var nameRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
var camelRe = regexp.MustCompile(`([a-z])([A-Z])`)
var badCharsRe = regexp.MustCompile(`[^a-z0-9-]`)
var multiHyphenRe = regexp.MustCompile(`-+`)

// ToHyphenCase converts a package name to hyphen-case (max 64 chars).
func ToHyphenCase(name string) string {
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	name = strings.NewReplacer("_", "-", " ", "-").Replace(name)
	name = camelRe.ReplaceAllString(name, "${1}-${2}")
	name = strings.ToLower(name)
	name = badCharsRe.ReplaceAllString(name, "")
	name = multiHyphenRe.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	if len(name) > 64 {
		name = name[:64]
	}
	return name
}

// ValidateSkillName validates a skill name per agentskills.io spec.
// Returns (valid, errorMessage).
func ValidateSkillName(name string) (bool, string) {
	if len(name) == 0 {
		return false, "Skill name cannot be empty"
	}
	if len(name) > 64 {
		return false, "Skill name must be 1-64 characters"
	}
	if strings.Contains(name, "--") {
		return false, "Skill name cannot contain consecutive hyphens (--)"
	}
	if strings.HasPrefix(name, "-") {
		return false, "Skill name cannot start with a hyphen"
	}
	if strings.HasSuffix(name, "-") {
		return false, "Skill name cannot end with a hyphen"
	}
	if !nameRe.MatchString(name) {
		return false, "Skill name must be lowercase alphanumeric with hyphens only"
	}
	return true, ""
}

// NormalizeSkillName converts any package name to a valid skill name.
func NormalizeSkillName(name string) string {
	return ToHyphenCase(name)
}

// ignoreNonContent returns true for paths that should not be copied
// (hidden files/dirs except SKILL.md, .git, __pycache__, *.pyc).
func ignoreNonContent(name string) bool {
	if name == ".git" || name == "__pycache__" || strings.HasSuffix(name, ".pyc") {
		return true
	}
	return false
}

// copyDirSkill copies src directory to dst, skipping non-content files.
func copyDirSkill(src, dst string) (int, error) {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return 0, err
	}
	count := 0
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, werr error) error {
		if werr != nil {
			return nil
		}
		rel, _ := filepath.Rel(src, path)
		if rel == "." {
			return nil
		}
		parts := strings.SplitN(rel, string(filepath.Separator), 2)
		if len(parts) > 0 && ignoreNonContent(parts[0]) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		if err := os.WriteFile(target, data, 0o644); err != nil {
			return err
		}
		count++
		return nil
	})
	return count, err
}

// dirsEqual returns true if two directory trees have identical file contents.
func dirsEqual(a, b string) bool {
	aFiles := map[string][]byte{}
	bFiles := map[string][]byte{}
	collectFiles := func(root string, m map[string][]byte) {
		_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, _ error) error {
			if d == nil || d.IsDir() {
				return nil
			}
			rel, _ := filepath.Rel(root, path)
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			m[rel] = data
			return nil
		})
	}
	collectFiles(a, aFiles)
	collectFiles(b, bFiles)
	if len(aFiles) != len(bFiles) {
		return false
	}
	for k, va := range aFiles {
		vb, ok := bFiles[k]
		if !ok || string(va) != string(vb) {
			return false
		}
	}
	return true
}

// SkillIntegrator handles integration of SKILL.md-based packages.
type SkillIntegrator struct {
	mu                       sync.Mutex
	nativeSkillSessionOwners map[string]string
}

// New returns a new SkillIntegrator.
func New() *SkillIntegrator {
	return &SkillIntegrator{
		nativeSkillSessionOwners: map[string]string{},
	}
}

// allKnownTargets returns a slice of all known target profiles.
func allKnownTargets() []*targets.TargetProfile {
	out := make([]*targets.TargetProfile, 0, len(targets.KnownTargets))
	for _, t := range targets.KnownTargets {
		out = append(out, t)
	}
	return out
}

// FindInstructionFiles returns all .instructions.md files from .apm/instructions/.
func FindInstructionFiles(packagePath string) []string {
	dir := filepath.Join(packagePath, ".apm", "instructions")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".instructions.md") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	return out
}

// FindAgentFiles returns all .agent.md files from .apm/agents/.
func FindAgentFiles(packagePath string) []string {
	dir := filepath.Join(packagePath, ".apm", "agents")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".agent.md") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	return out
}

// FindPromptFiles returns all .prompt.md files from package root and .apm/prompts/.
func FindPromptFiles(packagePath string) []string {
	var out []string
	entries, err := os.ReadDir(packagePath)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".prompt.md") {
				out = append(out, filepath.Join(packagePath, e.Name()))
			}
		}
	}
	dir := filepath.Join(packagePath, ".apm", "prompts")
	if entries2, err := os.ReadDir(dir); err == nil {
		for _, e := range entries2 {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".prompt.md") {
				out = append(out, filepath.Join(dir, e.Name()))
			}
		}
	}
	return out
}

// FindContextFiles returns all context and memory files.
func FindContextFiles(packagePath string) []string {
	var out []string
	for _, sub := range []string{".apm/context", ".apm/memory"} {
		dir := filepath.Join(packagePath, sub)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		suffix := ".context.md"
		if strings.HasSuffix(sub, "memory") {
			suffix = ".memory.md"
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), suffix) {
				out = append(out, filepath.Join(dir, e.Name()))
			}
		}
	}
	return out
}

// PackageInfo is a minimal interface for package metadata used by skill integration.
type PackageInfo struct {
	InstallPath string
	PackageType string // "CLAUDE_SKILL", "HYBRID", "SKILL_BUNDLE", "MARKETPLACE_PLUGIN", "INSTRUCTIONS", "PROMPTS"
	IsVirtual   bool
	IsSubdir    bool
	UniqueKey   string
}

// shouldInstallSkill returns true for packages that should be installed as skills.
func shouldInstallSkill(pkg *PackageInfo) bool {
	switch pkg.PackageType {
	case "CLAUDE_SKILL", "HYBRID", "SKILL_BUNDLE", "MARKETPLACE_PLUGIN":
		return true
	}
	return false
}

// promoteSubSkills promotes sub-skills from .apm/skills/ to a target skills root.
func promoteSubSkills(
	subSkillsDir string,
	targetSkillsRoot string,
	parentName string,
	ownedBy map[string]string,
	managedFiles map[string]struct{},
	force bool,
	nameFilter map[string]struct{},
) (int, []string) {
	entries, err := os.ReadDir(subSkillsDir)
	if err != nil {
		return 0, nil
	}
	promoted := 0
	var deployed []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		subPath := filepath.Join(subSkillsDir, e.Name())
		if _, err := os.Stat(filepath.Join(subPath, "SKILL.md")); err != nil {
			continue
		}
		rawName := e.Name()
		if nameFilter != nil {
			if _, ok := nameFilter[rawName]; !ok {
				continue
			}
		}
		valid, _ := ValidateSkillName(rawName)
		subName := rawName
		if !valid {
			subName = NormalizeSkillName(rawName)
		}
		target := filepath.Join(targetSkillsRoot, subName)
		if _, err := os.Stat(target); err == nil {
			if dirsEqual(subPath, target) {
				promoted++
				deployed = append(deployed, target)
				continue
			}
			relPath := filepath.Join(filepath.Base(targetSkillsRoot), subName)
			isManaged := false
			if managedFiles != nil {
				norm := strings.ReplaceAll(relPath, "\\", "/")
				_, isManaged = managedFiles[norm]
			}
			prevOwner := ownedBy[subName]
			isSelfOverwrite := prevOwner != "" && prevOwner == parentName
			if managedFiles != nil && !isManaged && !isSelfOverwrite && !force {
				continue
			}
			_ = os.RemoveAll(target)
		}
		if err := os.MkdirAll(target, 0o755); err != nil {
			continue
		}
		if _, err := copyDirSkill(subPath, target); err != nil {
			continue
		}
		promoted++
		deployed = append(deployed, target)
	}
	return promoted, deployed
}

// IntegrateNativeSkill deploys a package with a root SKILL.md to all active targets.
func (si *SkillIntegrator) IntegrateNativeSkill(
	pkg *PackageInfo,
	projectRoot string,
	force bool,
	managedFiles map[string]struct{},
	allTargets []*targets.TargetProfile,
) *SkillIntegrationResult {
	packagePath := pkg.InstallPath
	rawSkillName := filepath.Base(packagePath)
	valid, _ := ValidateSkillName(rawSkillName)
	skillName := rawSkillName
	if !valid {
		skillName = NormalizeSkillName(rawSkillName)
	}

	if allTargets == nil {
		allTargets = targets.ActiveTargets(projectRoot, nil)
	}
	skillCreated := false
	skillUpdated := false
	filesCopied := 0
	var allTargetPaths []string
	var primarySkillMD string

	seen := map[string]bool{}

	for idx, tgt := range allTargets {
		if !tgt.Supports("skills") {
			continue
		}
		sm := tgt.Primitives["skills"]
		effectiveRoot := sm.DeployRoot
		if effectiveRoot == "" {
			effectiveRoot = tgt.RootDir
		}
		targetSkillDir := filepath.Join(projectRoot, effectiveRoot, "skills", skillName)
		// path security: no traversal
		if strings.Contains(skillName, "..") {
			continue
		}
		resolved, _ := filepath.EvalSymlinks(targetSkillDir)
		if resolved == "" {
			resolved = targetSkillDir
		}
		if seen[resolved] {
			continue
		}
		seen[resolved] = true

		isPrimary := idx == 0
		if isPrimary {
			if _, err := os.Stat(targetSkillDir); os.IsNotExist(err) {
				skillCreated = true
			} else {
				skillUpdated = true
			}
			primarySkillMD = filepath.Join(targetSkillDir, "SKILL.md")
		}

		_ = os.RemoveAll(targetSkillDir)
		_ = os.MkdirAll(filepath.Dir(targetSkillDir), 0o755)
		n, _ := copyDirSkill(packagePath, targetSkillDir)
		allTargetPaths = append(allTargetPaths, targetSkillDir)
		if isPrimary {
			filesCopied = n
		}

		// Promote sub-skills
		subSkillsDir := filepath.Join(packagePath, ".apm", "skills")
		targetSkillsRoot := filepath.Join(projectRoot, effectiveRoot, "skills")
		_, subDeployed := promoteSubSkills(subSkillsDir, targetSkillsRoot, skillName, nil, managedFiles, force, nil)
		allTargetPaths = append(allTargetPaths, subDeployed...)
		_ = subDeployed
	}

	si.mu.Lock()
	if pkg.UniqueKey != "" {
		si.nativeSkillSessionOwners[skillName] = pkg.UniqueKey
	}
	si.mu.Unlock()

	primaryRoot := filepath.Join(projectRoot, ".github", "skills")
	subSkillsCount := 0
	for _, p := range allTargetPaths {
		if filepath.Dir(p) == primaryRoot && filepath.Base(p) != skillName {
			subSkillsCount++
		}
	}

	return &SkillIntegrationResult{
		SkillCreated:      skillCreated,
		SkillUpdated:      skillUpdated,
		SkillSkipped:      false,
		SkillPath:         primarySkillMD,
		ReferencesCopied:  filesCopied,
		SubSkillsPromoted: subSkillsCount,
		TargetPaths:       allTargetPaths,
	}
}

// IntegrateSkillBundle promotes every skill in a root-level skills/ directory.
func (si *SkillIntegrator) IntegrateSkillBundle(
	pkg *PackageInfo,
	projectRoot string,
	skillsDir string,
	force bool,
	managedFiles map[string]struct{},
	allTargets []*targets.TargetProfile,
	nameFilter map[string]struct{},
) *SkillIntegrationResult {
	if allTargets == nil {
		allTargets = targets.ActiveTargets(projectRoot, nil)
	}
	parentName := filepath.Base(pkg.InstallPath)
	totalPromoted := 0
	var allDeployed []string
	anyCreated := false
	seen := map[string]bool{}

	for idx, tgt := range allTargets {
		if !tgt.Supports("skills") {
			continue
		}
		sm := tgt.Primitives["skills"]
		effectiveRoot := sm.DeployRoot
		if effectiveRoot == "" {
			effectiveRoot = tgt.RootDir
		}
		targetSkillsRoot := filepath.Join(projectRoot, effectiveRoot, "skills")
		resolved, _ := filepath.EvalSymlinks(targetSkillsRoot)
		if resolved == "" {
			resolved = targetSkillsRoot
		}
		if seen[resolved] {
			continue
		}
		seen[resolved] = true
		_ = os.MkdirAll(targetSkillsRoot, 0o755)

		isPrimary := idx == 0
		n, deployed := promoteSubSkills(skillsDir, targetSkillsRoot, parentName, nil, managedFiles, force, nameFilter)
		if isPrimary {
			totalPromoted = n
			if n > 0 {
				anyCreated = true
			}
		}
		allDeployed = append(allDeployed, deployed...)
	}

	return &SkillIntegrationResult{
		SkillCreated:      anyCreated,
		SkillSkipped:      false,
		SubSkillsPromoted: totalPromoted,
		TargetPaths:       allDeployed,
	}
}

// PromoteSubSkillsStandalone promotes sub-skills for non-skill packages.
func (si *SkillIntegrator) PromoteSubSkillsStandalone(
	pkg *PackageInfo,
	projectRoot string,
	force bool,
	managedFiles map[string]struct{},
	allTargets []*targets.TargetProfile,
) (int, []string) {
	subSkillsDir := filepath.Join(pkg.InstallPath, ".apm", "skills")
	if _, err := os.Stat(subSkillsDir); err != nil {
		return 0, nil
	}
	if allTargets == nil {
		allTargets = targets.ActiveTargets(projectRoot, nil)
	}
	parentName := filepath.Base(pkg.InstallPath)
	count := 0
	var allDeployed []string
	seen := map[string]bool{}

	for idx, tgt := range allTargets {
		if !tgt.Supports("skills") {
			continue
		}
		sm := tgt.Primitives["skills"]
		effectiveRoot := sm.DeployRoot
		if effectiveRoot == "" {
			effectiveRoot = tgt.RootDir
		}
		targetSkillsRoot := filepath.Join(projectRoot, effectiveRoot, "skills")
		resolved, _ := filepath.EvalSymlinks(targetSkillsRoot)
		if resolved == "" {
			resolved = targetSkillsRoot
		}
		if seen[resolved] {
			continue
		}
		seen[resolved] = true
		_ = os.MkdirAll(targetSkillsRoot, 0o755)

		isPrimary := idx == 0
		n, deployed := promoteSubSkills(subSkillsDir, targetSkillsRoot, parentName, nil, managedFiles, force, nil)
		if isPrimary {
			count = n
		}
		allDeployed = append(allDeployed, deployed...)
	}
	return count, allDeployed
}

// IntegratePackageSkill is the main entry point for skill integration.
func (si *SkillIntegrator) IntegratePackageSkill(
	pkg *PackageInfo,
	projectRoot string,
	force bool,
	managedFiles map[string]struct{},
	allTargets []*targets.TargetProfile,
	skillSubset []string,
) *SkillIntegrationResult {
	if !shouldInstallSkill(pkg) {
		subCount, subDeployed := si.PromoteSubSkillsStandalone(pkg, projectRoot, force, managedFiles, allTargets)
		return &SkillIntegrationResult{
			SkillSkipped:      true,
			SubSkillsPromoted: subCount,
			TargetPaths:       subDeployed,
		}
	}

	if pkg.IsVirtual && !pkg.IsSubdir {
		return &SkillIntegrationResult{SkillSkipped: true}
	}

	sourceSkillMD := filepath.Join(pkg.InstallPath, "SKILL.md")
	if _, err := os.Stat(sourceSkillMD); err == nil {
		return si.IntegrateNativeSkill(pkg, projectRoot, force, managedFiles, allTargets)
	}

	// Check for SKILL_BUNDLE
	rootSkillsDir := filepath.Join(pkg.InstallPath, "skills")
	if info, err := os.Stat(rootSkillsDir); err == nil && info.IsDir() {
		var nameFilter map[string]struct{}
		if len(skillSubset) > 0 {
			nameFilter = make(map[string]struct{}, len(skillSubset))
			for _, s := range skillSubset {
				nameFilter[s] = struct{}{}
			}
		}
		hasSkill := false
		entries, _ := os.ReadDir(rootSkillsDir)
		for _, e := range entries {
			if e.IsDir() {
				if _, err := os.Stat(filepath.Join(rootSkillsDir, e.Name(), "SKILL.md")); err == nil {
					hasSkill = true
					break
				}
			}
		}
		if hasSkill {
			return si.IntegrateSkillBundle(pkg, projectRoot, rootSkillsDir, force, managedFiles, allTargets, nameFilter)
		}
	}

	subCount, subDeployed := si.PromoteSubSkillsStandalone(pkg, projectRoot, force, managedFiles, allTargets)
	return &SkillIntegrationResult{
		SkillSkipped:      true,
		SubSkillsPromoted: subCount,
		TargetPaths:       subDeployed,
	}
}

// SyncStats holds cleanup statistics.
type SyncStats struct {
	FilesRemoved int
	Errors       int
}

// SyncIntegration removes orphaned skill directories.
func (si *SkillIntegrator) SyncIntegration(
	installedSkillNames map[string]struct{},
	projectRoot string,
	managedFiles map[string]struct{},
	allTargets []*targets.TargetProfile,
) SyncStats {
	if allTargets == nil {
		allTargets = allKnownTargets()
	}
	var stats SyncStats

	if managedFiles != nil {
		skillPrefixes := skillPrefixList(allTargets)
		projectResolved, _ := filepath.EvalSymlinks(projectRoot)
		if projectResolved == "" {
			projectResolved = projectRoot
		}
		for relPath := range managedFiles {
			norm := strings.ReplaceAll(relPath, "\\", "/")
			if strings.Contains(norm, "..") {
				continue
			}
			if !hasAnyPrefix(norm, skillPrefixes) {
				continue
			}
			target := filepath.Join(projectRoot, relPath)
			if _, err := os.Stat(target); err != nil {
				continue
			}
			info, err := os.Lstat(target)
			if err != nil {
				continue
			}
			if info.IsDir() {
				if err := os.RemoveAll(target); err != nil {
					stats.Errors++
				} else {
					stats.FilesRemoved++
				}
			} else {
				if err := os.Remove(target); err != nil {
					stats.Errors++
				} else {
					stats.FilesRemoved++
				}
			}
		}
		return stats
	}

	// Legacy: npm-style orphan detection
	seen := map[string]bool{}
	for _, tgt := range allTargets {
		if !tgt.Supports("skills") {
			continue
		}
		sm := tgt.Primitives["skills"]
		effectiveRoot := sm.DeployRoot
		if effectiveRoot == "" {
			effectiveRoot = tgt.RootDir
		}
		skillsDir := filepath.Join(projectRoot, effectiveRoot, "skills")
		resolved, _ := filepath.EvalSymlinks(skillsDir)
		if resolved == "" {
			resolved = skillsDir
		}
		if seen[resolved] {
			continue
		}
		seen[resolved] = true
		entries, err := os.ReadDir(skillsDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			if _, ok := installedSkillNames[e.Name()]; ok {
				continue
			}
			target := filepath.Join(skillsDir, e.Name())
			if err := os.RemoveAll(target); err != nil {
				stats.Errors++
			} else {
				stats.FilesRemoved++
			}
		}
	}
	return stats
}

func skillPrefixList(allTargets []*targets.TargetProfile) []string {
	var out []string
	for _, tgt := range allTargets {
		if !tgt.Supports("skills") {
			continue
		}
		sm := tgt.Primitives["skills"]
		effectiveRoot := sm.DeployRoot
		if effectiveRoot == "" {
			effectiveRoot = tgt.RootDir
		}
		out = append(out, effectiveRoot+"/skills/")
	}
	return out
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
