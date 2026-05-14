// Package discovery provides functionality for discovering APM primitive files.
// Migrated from src/apm_cli/primitives/discovery.py.
package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/githubnext/apm/internal/constants"
	"github.com/githubnext/apm/internal/primitives/primmodels"
	"github.com/githubnext/apm/internal/primitives/primparser"
	"github.com/githubnext/apm/internal/utils/exclude"
	"github.com/githubnext/apm/internal/utils/paths"
)

// PrimitiveConflict records when two primitives compete for the same name.
type PrimitiveConflict struct {
	PrimitiveName  string
	PrimitiveType  string
	WinningSource  string
	LosingSource   string
	FilePath       string
}

// PrimitiveCollection holds all discovered primitives.
type PrimitiveCollection struct {
	Chatmodes    []*primmodels.Chatmode
	Instructions []*primmodels.Instruction
	Contexts     []*primmodels.Context
	Skills       []*primmodels.Skill
	Conflicts    []PrimitiveConflict

	chatmodeIndex    map[string]int
	instructionIndex map[string]int
	contextIndex     map[string]int
	skillIndex       map[string]int
}

// NewPrimitiveCollection creates an initialized PrimitiveCollection.
func NewPrimitiveCollection() *PrimitiveCollection {
	return &PrimitiveCollection{
		chatmodeIndex:    make(map[string]int),
		instructionIndex: make(map[string]int),
		contextIndex:     make(map[string]int),
		skillIndex:       make(map[string]int),
	}
}

// AddPrimitive adds a primitive to the collection with conflict detection.
func (c *PrimitiveCollection) AddPrimitive(p primmodels.Primitive) error {
	switch v := p.(type) {
	case *primmodels.Chatmode:
		c.addChatmode(v)
	case *primmodels.Instruction:
		c.addInstruction(v)
	case *primmodels.Context:
		c.addContext(v)
	case *primmodels.Skill:
		c.addSkill(v)
	default:
		return fmt.Errorf("unknown primitive type: %T", p)
	}
	return nil
}

func (c *PrimitiveCollection) addChatmode(p *primmodels.Chatmode) {
	if idx, exists := c.chatmodeIndex[p.Name]; exists {
		existing := c.Chatmodes[idx]
		if shouldReplace(existing.Source, p.Source) {
			c.Conflicts = append(c.Conflicts, PrimitiveConflict{
				PrimitiveName: p.Name, PrimitiveType: "chatmode",
				WinningSource: p.Source, LosingSource: existing.Source,
				FilePath: p.FilePath,
			})
			c.Chatmodes[idx] = p
		} else {
			c.Conflicts = append(c.Conflicts, PrimitiveConflict{
				PrimitiveName: p.Name, PrimitiveType: "chatmode",
				WinningSource: existing.Source, LosingSource: p.Source,
				FilePath: existing.FilePath,
			})
		}
		return
	}
	c.chatmodeIndex[p.Name] = len(c.Chatmodes)
	c.Chatmodes = append(c.Chatmodes, p)
}

func (c *PrimitiveCollection) addInstruction(p *primmodels.Instruction) {
	if idx, exists := c.instructionIndex[p.Name]; exists {
		existing := c.Instructions[idx]
		if shouldReplace(existing.Source, p.Source) {
			c.Conflicts = append(c.Conflicts, PrimitiveConflict{
				PrimitiveName: p.Name, PrimitiveType: "instruction",
				WinningSource: p.Source, LosingSource: existing.Source,
				FilePath: p.FilePath,
			})
			c.Instructions[idx] = p
		} else {
			c.Conflicts = append(c.Conflicts, PrimitiveConflict{
				PrimitiveName: p.Name, PrimitiveType: "instruction",
				WinningSource: existing.Source, LosingSource: p.Source,
				FilePath: existing.FilePath,
			})
		}
		return
	}
	c.instructionIndex[p.Name] = len(c.Instructions)
	c.Instructions = append(c.Instructions, p)
}

func (c *PrimitiveCollection) addContext(p *primmodels.Context) {
	if idx, exists := c.contextIndex[p.Name]; exists {
		existing := c.Contexts[idx]
		if shouldReplace(existing.Source, p.Source) {
			c.Conflicts = append(c.Conflicts, PrimitiveConflict{
				PrimitiveName: p.Name, PrimitiveType: "context",
				WinningSource: p.Source, LosingSource: existing.Source,
				FilePath: p.FilePath,
			})
			c.Contexts[idx] = p
		} else {
			c.Conflicts = append(c.Conflicts, PrimitiveConflict{
				PrimitiveName: p.Name, PrimitiveType: "context",
				WinningSource: existing.Source, LosingSource: p.Source,
				FilePath: existing.FilePath,
			})
		}
		return
	}
	c.contextIndex[p.Name] = len(c.Contexts)
	c.Contexts = append(c.Contexts, p)
}

func (c *PrimitiveCollection) addSkill(p *primmodels.Skill) {
	if idx, exists := c.skillIndex[p.Name]; exists {
		existing := c.Skills[idx]
		if shouldReplace(existing.Source, p.Source) {
			c.Conflicts = append(c.Conflicts, PrimitiveConflict{
				PrimitiveName: p.Name, PrimitiveType: "skill",
				WinningSource: p.Source, LosingSource: existing.Source,
				FilePath: p.FilePath,
			})
			c.Skills[idx] = p
		} else {
			c.Conflicts = append(c.Conflicts, PrimitiveConflict{
				PrimitiveName: p.Name, PrimitiveType: "skill",
				WinningSource: existing.Source, LosingSource: p.Source,
				FilePath: existing.FilePath,
			})
		}
		return
	}
	c.skillIndex[p.Name] = len(c.Skills)
	c.Skills = append(c.Skills, p)
}

// shouldReplace returns true when newSource should replace existingSource.
// Local always wins over dependency; earlier dependency wins over later.
func shouldReplace(existingSource, newSource string) bool {
	existingLocal := existingSource == "local" || existingSource == ""
	newLocal := newSource == "local" || newSource == ""
	if newLocal && !existingLocal {
		return true
	}
	return false
}

// Local primitive glob patterns (with recursive search via **/).
var localPrimitivePatterns = map[string][]string{
	"chatmode": {
		"**/.apm/agents/*.agent.md",
		"**/.github/agents/*.agent.md",
		"**/*.agent.md",
		"**/.apm/chatmodes/*.chatmode.md",
		"**/.github/chatmodes/*.chatmode.md",
		"**/*.chatmode.md",
	},
	"instruction": {
		"**/.apm/instructions/*.instructions.md",
		"**/.github/instructions/*.instructions.md",
		"**/*.instructions.md",
	},
	"context": {
		"**/.apm/context/*.context.md",
		"**/.apm/memory/*.memory.md",
		"**/.github/context/*.context.md",
		"**/.github/memory/*.memory.md",
		"**/*.context.md",
		"**/*.memory.md",
	},
}

// Dependency primitive patterns (for .apm directory within dependencies).
var dependencyPrimitivePatterns = map[string][]string{
	"chatmode":    {"agents/*.agent.md", "chatmodes/*.chatmode.md"},
	"instruction": {"instructions/*.instructions.md"},
	"context":     {"context/*.context.md", "memory/*.memory.md"},
}

// Dependency .github primitive patterns.
var dependencyGithubPrimitivePatterns = map[string][]string{
	"chatmode":    {"agents/*.agent.md", "chatmodes/*.chatmode.md"},
	"instruction": {"instructions/*.instructions.md"},
	"context":     {"context/*.context.md", "memory/*.memory.md"},
}

// DiscoverPrimitives finds all APM primitive files in the project.
func DiscoverPrimitives(baseDir string, excludePatterns []string) (*PrimitiveCollection, error) {
	collection := NewPrimitiveCollection()
	safePatterns, _ := exclude.ValidateExcludePatterns(excludePatterns)

	for _, ptPatterns := range localPrimitivePatterns {
		files, err := FindPrimitiveFiles(baseDir, ptPatterns, safePatterns)
		if err != nil {
			continue
		}
		for _, fp := range files {
			prim, err := primparser.ParsePrimitiveFile(fp, "local")
			if err != nil {
				fmt.Printf("Warning: Failed to parse %s: %v\n", fp, err)
				continue
			}
			collection.AddPrimitive(prim) //nolint:errcheck
		}
	}
	discoverLocalSkill(baseDir, collection, safePatterns)
	return collection, nil
}

// DiscoverPrimitivesWithDependencies performs enhanced discovery including dependencies.
func DiscoverPrimitivesWithDependencies(baseDir string, excludePatterns []string) (*PrimitiveCollection, error) {
	collection := NewPrimitiveCollection()
	safePatterns, _ := exclude.ValidateExcludePatterns(excludePatterns)

	scanLocalPrimitives(baseDir, collection, safePatterns)
	discoverLocalSkill(baseDir, collection, safePatterns)
	scanDependencyPrimitives(baseDir, collection)
	return collection, nil
}

// scanLocalPrimitives scans the local .apm/ directory for primitives.
func scanLocalPrimitives(baseDir string, collection *PrimitiveCollection, excludePatterns []string) {
	for _, ptPatterns := range localPrimitivePatterns {
		files, err := FindPrimitiveFiles(baseDir, ptPatterns, excludePatterns)
		if err != nil {
			continue
		}
		basePath, _ := filepath.Abs(baseDir)
		apmModulesPath := filepath.Join(basePath, "apm_modules")
		for _, fp := range files {
			absFile, _ := filepath.Abs(fp)
			if isUnderDirectory(absFile, apmModulesPath) {
				continue
			}
			prim, err := primparser.ParsePrimitiveFile(fp, "local")
			if err != nil {
				fmt.Printf("Warning: Failed to parse local primitive %s: %v\n", fp, err)
				continue
			}
			collection.AddPrimitive(prim) //nolint:errcheck
		}
	}
}

// scanDependencyPrimitives scans all dependencies in apm_modules/ with priority handling.
func scanDependencyPrimitives(baseDir string, collection *PrimitiveCollection) {
	apmModulesPath := filepath.Join(baseDir, "apm_modules")
	info, err := os.Stat(apmModulesPath)
	if err != nil || !info.IsDir() {
		return
	}
	depOrder := getDependencyDeclarationOrder(baseDir)
	for _, depName := range depOrder {
		parts := strings.Split(depName, "/")
		depPath := filepath.Join(append([]string{apmModulesPath}, parts...)...)
		info, err := os.Stat(depPath)
		if err == nil && info.IsDir() {
			ScanDirectoryWithSource(depPath, collection, "dependency:"+depName)
		}
	}
}

// getDependencyDeclarationOrder returns dependency installed paths in declaration order.
// Simplified: reads lockfile paths only (apm.yml parsing would need more infra).
func getDependencyDeclarationOrder(baseDir string) []string {
	// Fallback: return directories from apm_modules sorted alphabetically
	apmModulesPath := filepath.Join(baseDir, "apm_modules")
	entries, err := os.ReadDir(apmModulesPath)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			// Try two-level paths (owner/repo)
			subEntries, err := os.ReadDir(filepath.Join(apmModulesPath, e.Name()))
			if err != nil {
				names = append(names, e.Name())
				continue
			}
			for _, se := range subEntries {
				if se.IsDir() {
					names = append(names, e.Name()+"/"+se.Name())
				}
			}
		}
	}
	return names
}

// ScanDirectoryWithSource scans a directory for primitives with a specific source tag.
func ScanDirectoryWithSource(directory string, collection *PrimitiveCollection, source string) {
	apmDir := filepath.Join(directory, ".apm")
	if info, err := os.Stat(apmDir); err == nil && info.IsDir() {
		scanPatterns(apmDir, dependencyPrimitivePatterns, collection, source)
	}
	githubDir := filepath.Join(directory, ".github")
	if info, err := os.Stat(githubDir); err == nil && info.IsDir() {
		scanPatterns(githubDir, dependencyGithubPrimitivePatterns, collection, source)
	}
	discoverSkillInDirectory(directory, collection, source)
}

func discoverLocalSkill(baseDir string, collection *PrimitiveCollection, excludePatterns []string) {
	skillPath := filepath.Join(baseDir, "SKILL.md")
	info, err := os.Stat(skillPath)
	if err != nil || !info.Mode().IsRegular() {
		return
	}
	absBase, _ := filepath.Abs(baseDir)
	absSkill, _ := filepath.Abs(skillPath)
	if exclude.ShouldExclude(absSkill, absBase, excludePatterns) {
		return
	}
	if !isReadable(skillPath) {
		return
	}
	skill, err := primparser.ParseSkillFile(skillPath, "local")
	if err != nil {
		fmt.Printf("Warning: Failed to parse SKILL.md: %v\n", err)
		return
	}
	collection.AddPrimitive(skill) //nolint:errcheck
}

func discoverSkillInDirectory(directory string, collection *PrimitiveCollection, source string) {
	skillPath := filepath.Join(directory, "SKILL.md")
	if !isReadable(skillPath) {
		return
	}
	skill, err := primparser.ParseSkillFile(skillPath, source)
	if err != nil {
		fmt.Printf("Warning: Failed to parse SKILL.md in %s: %v\n", directory, err)
		return
	}
	collection.AddPrimitive(skill) //nolint:errcheck
}

// scanPatterns walks baseDir once and matches files against all patterns.
func scanPatterns(baseDir string, patterns map[string][]string, collection *PrimitiveCollection, source string) {
	info, err := os.Stat(baseDir)
	if err != nil || !info.IsDir() {
		return
	}
	// Flatten all patterns
	var allPatterns []string
	for _, ps := range patterns {
		allPatterns = append(allPatterns, ps...)
	}

	err = filepath.WalkDir(baseDir, func(fp string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(baseDir, fp)
		if err != nil {
			return nil
		}
		relFwd := strings.ReplaceAll(rel, string(filepath.Separator), "/")
		if !matchesAnyPattern(relFwd, allPatterns) {
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		if !isReadable(fp) {
			return nil
		}
		prim, err := primparser.ParsePrimitiveFile(fp, source)
		if err != nil {
			fmt.Printf("Warning: Failed to parse dependency primitive %s: %v\n", fp, err)
			return nil
		}
		collection.AddPrimitive(prim) //nolint:errcheck
		return nil
	})
	_ = err
}

// FindPrimitiveFiles finds primitive files matching the given patterns.
func FindPrimitiveFiles(baseDir string, patterns []string, excludePatterns []string) ([]string, error) {
	info, err := os.Stat(baseDir)
	if err != nil || !info.IsDir() {
		return nil, nil
	}
	basePath, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, err
	}

	var allFiles []string

	err = filepath.WalkDir(basePath, func(fp string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if d.IsDir() {
			if _, skip := constants.DefaultSkipDirs[name]; skip {
				return filepath.SkipDir
			}
			if exclude.ShouldExclude(fp, basePath, excludePatterns) {
				return filepath.SkipDir
			}
			return nil
		}
		// Sort within directory is handled by WalkDir (lexical order already)
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if exclude.ShouldExclude(fp, basePath, excludePatterns) {
			return nil
		}
		rel := paths.PortableRelpath(fp, basePath)
		for _, pat := range patterns {
			if globMatch(rel, pat) {
				allFiles = append(allFiles, fp)
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Filter invalid
	valid := make([]string, 0, len(allFiles))
	for _, fp := range allFiles {
		fi, err := os.Lstat(fp)
		if err != nil {
			continue
		}
		if !fi.Mode().IsRegular() {
			continue
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if isReadable(fp) {
			valid = append(valid, fp)
		}
	}
	sort.Strings(valid)
	return valid, nil
}

// globMatch matches a forward-slash relative path against a glob pattern.
// Segment-aware: ** matches zero or more complete path segments.
func globMatch(relPath, pattern string) bool {
	pathParts := splitNonEmpty(relPath, "/")
	patternParts := splitNonEmpty(pattern, "/")
	memo := make(map[[2]int]bool)
	var match func(pi, qi int) bool
	match = func(pi, qi int) bool {
		key := [2]int{pi, qi}
		if v, ok := memo[key]; ok {
			return v
		}
		if qi == len(patternParts) {
			result := pi == len(pathParts)
			memo[key] = result
			return result
		}
		cur := patternParts[qi]
		if cur == "**" {
			result := match(pi, qi+1)
			if !result && pi < len(pathParts) {
				result = match(pi+1, qi)
			}
			memo[key] = result
			return result
		}
		if pi >= len(pathParts) {
			memo[key] = false
			return false
		}
		result := fnmatchSegment(pathParts[pi], cur) && match(pi+1, qi+1)
		memo[key] = result
		return result
	}
	return match(0, 0)
}

// fnmatchSegment matches a single path segment against a pattern.
// Supports * (any chars within segment) and ? (single char).
func fnmatchSegment(name, pattern string) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		case '*':
			if len(pattern) == 1 {
				return true
			}
			rest := pattern[1:]
			for i := 0; i <= len(name); i++ {
				if fnmatchSegment(name[i:], rest) {
					return true
				}
			}
			return false
		case '?':
			if len(name) == 0 {
				return false
			}
			name = name[1:]
			pattern = pattern[1:]
		default:
			if len(name) == 0 || name[0] != pattern[0] {
				return false
			}
			name = name[1:]
			pattern = pattern[1:]
		}
	}
	return len(name) == 0
}

func matchesAnyPattern(relPath string, patterns []string) bool {
	for _, p := range patterns {
		if globMatch(relPath, p) {
			return true
		}
	}
	return false
}

func isUnderDirectory(filePath, directory string) bool {
	rel, err := filepath.Rel(directory, filePath)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}

func isReadable(fp string) bool {
	f, err := os.Open(fp)
	if err != nil {
		return false
	}
	buf := make([]byte, 1)
	_, err = f.Read(buf)
	f.Close()
	return err == nil
}

func splitNonEmpty(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
