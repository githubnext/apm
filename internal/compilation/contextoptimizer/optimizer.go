// Package contextoptimizer implements the Context Optimization Engine.
//
// Minimizes irrelevant context loaded by agents working in specific
// directories, following the Minimal Context Principle: place each
// instruction at the shallowest directory that covers all files
// matching its pattern, without bleeding into unrelated subtrees.
//
// Migrated from: src/apm_cli/compilation/context_optimizer.py
package contextoptimizer

import (
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// DefaultExcludedDirnames is the set of directory names always skipped during
// project-structure analysis.
var DefaultExcludedDirnames = map[string]bool{
	"node_modules": true,
	"__pycache__":  true,
	".git":         true,
	"dist":         true,
	"build":        true,
	"apm_modules":  true,
}

// -------------------------------------------------------------------
// Data types
// -------------------------------------------------------------------

// DirectoryAnalysis summarises a directory's file distribution.
type DirectoryAnalysis struct {
	Directory     string
	Depth         int
	TotalFiles    int
	PatternCounts map[string]int // pattern -> matching-file count
	FileTypes     map[string]bool
}

// RelevanceScore returns the fraction of files in this directory that match pattern.
func (d *DirectoryAnalysis) RelevanceScore(pattern string) float64 {
	if d.TotalFiles == 0 {
		return 0
	}
	return float64(d.PatternCounts[pattern]) / float64(d.TotalFiles)
}

// InheritanceAnalysis captures the inheritance chain for a working directory.
type InheritanceAnalysis struct {
	WorkingDirectory     string
	InheritanceChain     []string // most-specific first
	TotalContextLoad     int
	RelevantContextLoad  int
	PollutionScore       float64
}

// EfficiencyRatio returns the fraction of loaded context that is relevant.
func (a *InheritanceAnalysis) EfficiencyRatio() float64 {
	if a.TotalContextLoad == 0 {
		return 1
	}
	return float64(a.RelevantContextLoad) / float64(a.TotalContextLoad)
}

// PlacementCandidate is a candidate directory for placing an instruction.
type PlacementCandidate struct {
	Directory         string
	Score             float64
	CoverageRatio     float64
	PollutionScore    float64
	MaintenanceScore  float64
	Depth             int
	IsLeaf            bool
}

// PlacementDecision is the final placement recommendation for one instruction.
type PlacementDecision struct {
	InstructionPath string
	TargetDirectory string
	Strategy        string  // "single_point" | "distributed" | "selective" | "unchanged"
	Score           float64
	Candidates      []PlacementCandidate
	Reason          string
}

// OptimizationResult holds the full output of an optimization pass.
type OptimizationResult struct {
	Decisions        []PlacementDecision
	Stats            OptimizationStats
	ElapsedSeconds   float64
}

// OptimizationStats holds summary metrics from an optimization run.
type OptimizationStats struct {
	TotalInstructions    int
	Optimized            int
	Unchanged            int
	PollutionReduction   float64
	CoverageGain         float64
	PhaseTimings         map[string]float64
}

// ProjectStructure holds the cached analysis of the project file tree.
type ProjectStructure struct {
	Dirs        map[string]*DirectoryAnalysis
	AllFiles    []string
	MaxDepth    int
}

// -------------------------------------------------------------------
// ContextOptimizer
// -------------------------------------------------------------------

// ContextOptimizer is the main engine.
type ContextOptimizer struct {
	BaseDir         string
	ExcludePatterns []string

	mu          sync.Mutex
	structure   *ProjectStructure
	globCache   map[string][]string
	timingData  map[string]float64
	timingEnabled bool
}

// New constructs a ContextOptimizer.
func New(baseDir string, excludePatterns []string) *ContextOptimizer {
	if baseDir == "" {
		baseDir = "."
	}
	abs, err := filepath.Abs(baseDir)
	if err != nil {
		abs = baseDir
	}
	return &ContextOptimizer{
		BaseDir:         abs,
		ExcludePatterns: excludePatterns,
		globCache:       make(map[string][]string),
		timingData:      make(map[string]float64),
	}
}

// EnableTiming turns on per-phase timing collection.
func (c *ContextOptimizer) EnableTiming(verbose bool) {
	c.timingEnabled = verbose
}

// -------------------------------------------------------------------
// File enumeration
// -------------------------------------------------------------------

func (c *ContextOptimizer) getAllFiles() []string {
	var files []string
	_ = filepath.WalkDir(c.BaseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if DefaultExcludedDirnames[d.Name()] {
				return filepath.SkipDir
			}
			if c.shouldExcludePath(path) {
				return filepath.SkipDir
			}
			return nil
		}
		rel, _ := filepath.Rel(c.BaseDir, path)
		if !c.shouldExcludePath(rel) {
			files = append(files, rel)
		}
		return nil
	})
	return files
}

func (c *ContextOptimizer) shouldExcludePath(path string) bool {
	for _, pat := range c.ExcludePatterns {
		if matched, _ := filepath.Match(pat, filepath.Base(path)); matched {
			return true
		}
		if strings.Contains(path, pat) {
			return true
		}
	}
	return false
}

// -------------------------------------------------------------------
// Project structure analysis
// -------------------------------------------------------------------

func (c *ContextOptimizer) analyzeProjectStructure() *ProjectStructure {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.structure != nil {
		return c.structure
	}
	files := c.getAllFiles()
	dirs := make(map[string]*DirectoryAnalysis)
	maxDepth := 0

	for _, f := range files {
		dir := filepath.Dir(f)
		if dir == "." {
			dir = ""
		}
		// Add this dir and all ancestor dirs
		parts := strings.Split(dir, string(filepath.Separator))
		for depth := 0; depth <= len(parts); depth++ {
			var d string
			if depth == 0 {
				d = ""
			} else {
				d = filepath.Join(parts[:depth]...)
			}
			if _, ok := dirs[d]; !ok {
				dirs[d] = &DirectoryAnalysis{
					Directory:     d,
					Depth:         depth,
					PatternCounts: make(map[string]int),
					FileTypes:     make(map[string]bool),
				}
			}
			dirs[d].TotalFiles++
			ext := filepath.Ext(f)
			dirs[d].FileTypes[ext] = true
			if depth > maxDepth {
				maxDepth = depth
			}
		}
	}
	c.structure = &ProjectStructure{Dirs: dirs, AllFiles: files, MaxDepth: maxDepth}
	return c.structure
}

// -------------------------------------------------------------------
// Pattern matching
// -------------------------------------------------------------------

func (c *ContextOptimizer) fileMatchesPattern(filePath, pattern string) bool {
	base := filepath.Base(filePath)
	if matched, _ := filepath.Match(pattern, base); matched {
		return true
	}
	if matched, _ := filepath.Match(pattern, filePath); matched {
		return true
	}
	// Handle glob-style with **
	if strings.Contains(pattern, "**") {
		parts := strings.SplitN(pattern, "**", 2)
		if strings.HasPrefix(filePath, strings.TrimPrefix(parts[0], "/")) {
			suffix := strings.TrimPrefix(parts[1], "/")
			if suffix == "" || strings.HasSuffix(filePath, suffix) {
				return true
			}
		}
	}
	return false
}

func (c *ContextOptimizer) findMatchingDirectories(pattern string) map[string]bool {
	struct_ := c.analyzeProjectStructure()
	result := make(map[string]bool)
	for _, f := range struct_.AllFiles {
		if c.fileMatchesPattern(f, pattern) {
			dir := filepath.Dir(f)
			if dir == "." {
				dir = ""
			}
			result[dir] = true
		}
	}
	return result
}

// -------------------------------------------------------------------
// Scoring helpers
// -------------------------------------------------------------------

func (c *ContextOptimizer) calculateInheritancePollution(dir, pattern string) float64 {
	struct_ := c.analyzeProjectStructure()
	analysis, ok := struct_.Dirs[dir]
	if !ok || analysis.TotalFiles == 0 {
		return 0
	}
	matching := c.findMatchingDirectories(pattern)
	// Count files in subtree that DON'T match the pattern
	var unrelated int
	for _, f := range struct_.AllFiles {
		fDir := filepath.Dir(f)
		if fDir == "." {
			fDir = ""
		}
		if fDir == dir || strings.HasPrefix(fDir, dir+string(filepath.Separator)) {
			if !matching[fDir] {
				unrelated++
			}
		}
	}
	return float64(unrelated) / float64(analysis.TotalFiles)
}

func (c *ContextOptimizer) calculateDistributionScore(matchingDirs map[string]bool) float64 {
	if len(matchingDirs) == 0 {
		return 0
	}
	// Higher score = more evenly distributed
	return math.Min(1.0, float64(len(matchingDirs))/10.0)
}

func (c *ContextOptimizer) calculateCoverageEfficiency(dir, pattern string) float64 {
	matching := c.findMatchingDirectories(pattern)
	if len(matching) == 0 {
		return 0
	}
	// How many of the matching dirs are covered by placing at dir?
	covered := 0
	for d := range matching {
		if d == dir || strings.HasPrefix(d, dir+string(filepath.Separator)) {
			covered++
		}
	}
	return float64(covered) / float64(len(matching))
}

// -------------------------------------------------------------------
// Placement optimization
// -------------------------------------------------------------------

// OptimizeInstructionPlacement returns the best directory for each instruction pattern.
func (c *ContextOptimizer) OptimizeInstructionPlacement(patterns []string) *OptimizationResult {
	t0 := time.Now()
	result := &OptimizationResult{
		Stats: OptimizationStats{
			TotalInstructions: len(patterns),
			PhaseTimings:      make(map[string]float64),
		},
	}

	for _, pat := range patterns {
		decision := c.optimizeSinglePattern(pat)
		result.Decisions = append(result.Decisions, decision)
		if decision.Strategy != "unchanged" {
			result.Stats.Optimized++
		} else {
			result.Stats.Unchanged++
		}
	}

	result.ElapsedSeconds = time.Since(t0).Seconds()
	return result
}

func (c *ContextOptimizer) optimizeSinglePattern(pattern string) PlacementDecision {
	matchingDirs := c.findMatchingDirectories(pattern)
	if len(matchingDirs) == 0 {
		return PlacementDecision{
			InstructionPath: pattern,
			TargetDirectory: "",
			Strategy:        "unchanged",
			Reason:          "no matching files found",
		}
	}

	candidates := c.generateCandidates(pattern, matchingDirs)
	if len(candidates) == 0 {
		return PlacementDecision{
			InstructionPath: pattern,
			TargetDirectory: "",
			Strategy:        "unchanged",
			Reason:          "no viable placement candidates",
		}
	}

	best := candidates[0]
	strategy := "single_point"
	if len(matchingDirs) > 5 {
		strategy = "distributed"
	}

	return PlacementDecision{
		InstructionPath: pattern,
		TargetDirectory: best.Directory,
		Strategy:        strategy,
		Score:           best.Score,
		Candidates:      candidates,
		Reason:          "optimized placement",
	}
}

func (c *ContextOptimizer) generateCandidates(pattern string, matchingDirs map[string]bool) []PlacementCandidate {
	struct_ := c.analyzeProjectStructure()
	var candidates []PlacementCandidate

	for dir := range struct_.Dirs {
		coverage := c.calculateCoverageEfficiency(dir, pattern)
		if coverage == 0 {
			continue
		}
		pollution := c.calculateInheritancePollution(dir, pattern)
		depth := struct_.Dirs[dir].Depth
		maintenanceScore := 1.0 / float64(depth+1)

		score := coverage*0.5 + (1-pollution)*0.3 + maintenanceScore*0.2

		candidates = append(candidates, PlacementCandidate{
			Directory:        dir,
			Score:            score,
			CoverageRatio:    coverage,
			PollutionScore:   pollution,
			MaintenanceScore: maintenanceScore,
			Depth:            depth,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	if len(candidates) > 10 {
		candidates = candidates[:10]
	}
	return candidates
}

// -------------------------------------------------------------------
// Inheritance analysis
// -------------------------------------------------------------------

// AnalyzeContextInheritance computes the inheritance chain for workingDir.
func (c *ContextOptimizer) AnalyzeContextInheritance(workingDir string) *InheritanceAnalysis {
	chain := c.getInheritanceChain(workingDir)
	struct_ := c.analyzeProjectStructure()

	total := 0
	for _, d := range chain {
		if a, ok := struct_.Dirs[d]; ok {
			total += a.TotalFiles
		}
	}

	return &InheritanceAnalysis{
		WorkingDirectory:    workingDir,
		InheritanceChain:    chain,
		TotalContextLoad:    total,
		RelevantContextLoad: total, // conservative: assume all relevant
	}
}

func (c *ContextOptimizer) getInheritanceChain(dir string) []string {
	var chain []string
	parts := strings.Split(dir, string(filepath.Separator))
	for i := len(parts); i >= 0; i-- {
		var d string
		if i == 0 {
			d = ""
		} else {
			d = filepath.Join(parts[:i]...)
		}
		chain = append(chain, d)
	}
	return chain
}

// -------------------------------------------------------------------
// Stats
// -------------------------------------------------------------------

// GetOptimizationStats returns summary stats from a completed optimization.
func (c *ContextOptimizer) GetOptimizationStats(result *OptimizationResult) OptimizationStats {
	return result.Stats
}
