// Package output provides data models and formatters for APM CLI compilation
// output. It mirrors the Python apm_cli.output module.
package output

// PlacementStrategy represents the placement strategy type for optimization decisions.
type PlacementStrategy string

const (
	// PlacementSinglePoint mirrors Python's PlacementStrategy.SINGLE_POINT.
	PlacementSinglePoint PlacementStrategy = "Single Point"
	// PlacementSelectiveMulti mirrors Python's PlacementStrategy.SELECTIVE_MULTI.
	PlacementSelectiveMulti PlacementStrategy = "Selective Multi"
	// PlacementDistributed mirrors Python's PlacementStrategy.DISTRIBUTED.
	PlacementDistributed PlacementStrategy = "Distributed"
)

// ProjectAnalysis holds the analysis of the project structure and file distribution.
type ProjectAnalysis struct {
	// DirectoriesScanned is the number of directories examined.
	DirectoriesScanned int
	// FilesAnalyzed is the number of files examined.
	FilesAnalyzed int
	// FileTypesDetected is the set of file extensions detected.
	FileTypesDetected map[string]bool
	// InstructionPatternsDetected is the count of distinct instruction patterns.
	InstructionPatternsDetected int
	// MaxDepth is the maximum directory depth encountered.
	MaxDepth int
	// ConstitutionDetected reports whether a constitution file was found.
	ConstitutionDetected bool
	// ConstitutionPath is the path to the detected constitution, or empty string.
	ConstitutionPath string
}

// GetFileTypesSummary returns a concise summary of detected file types.
func (p *ProjectAnalysis) GetFileTypesSummary() string {
	if len(p.FileTypesDetected) == 0 {
		return "none"
	}
	types := make([]string, 0, len(p.FileTypesDetected))
	for t := range p.FileTypesDetected {
		cleaned := t
		if len(cleaned) > 0 && cleaned[0] == '.' {
			cleaned = cleaned[1:]
		}
		if cleaned != "" {
			types = append(types, cleaned)
		}
	}
	sortStrings(types)
	if len(types) <= 3 {
		return joinStrings(types)
	}
	base := joinStrings(types[:3])
	return base + " and " + itoa(len(types)-3) + " more"
}

// OptimizationDecision holds details about a specific placement decision.
type OptimizationDecision struct {
	// InstructionID uniquely identifies the instruction being placed.
	InstructionID string
	// Pattern is the glob/path pattern associated with the instruction.
	Pattern string
	// MatchingDirectories is the count of directories that match this instruction.
	MatchingDirectories int
	// TotalDirectories is the total directory count in the project.
	TotalDirectories int
	// DistributionScore is the fraction of directories the instruction covers.
	DistributionScore float64
	// Strategy is the placement strategy chosen for this instruction.
	Strategy PlacementStrategy
	// PlacementDirectories is the list of directories selected for placement.
	PlacementDirectories []string
	// Reasoning is a human-readable explanation for the placement choice.
	Reasoning string
	// RelevanceScore is the coverage efficiency for the primary placement directory.
	RelevanceScore float64
}

// DistributionRatio returns MatchingDirectories / TotalDirectories, or 0 when TotalDirectories == 0.
func (o *OptimizationDecision) DistributionRatio() float64 {
	if o.TotalDirectories == 0 {
		return 0.0
	}
	return float64(o.MatchingDirectories) / float64(o.TotalDirectories)
}

// PlacementSummary summarizes a single target-file placement.
type PlacementSummary struct {
	// Path is the absolute or relative filesystem path to the placed file.
	Path string
	// InstructionCount is the number of instructions placed in this file.
	InstructionCount int
	// SourceCount is the number of source instruction files contributing.
	SourceCount int
	// Sources is the list of source file paths contributing to this placement.
	Sources []string
}

// OptimizationStats holds performance and efficiency statistics from an optimization run.
type OptimizationStats struct {
	// AverageContextEfficiency is the mean context coverage ratio across all placements.
	AverageContextEfficiency float64
	// PollutionImprovement is the reduction in context pollution, or nil.
	PollutionImprovement *float64
	// BaselineEfficiency is the pre-optimization efficiency, or nil.
	BaselineEfficiency *float64
	// PlacementAccuracy is the fraction of correctly placed instructions, or nil.
	PlacementAccuracy *float64
	// GenerationTimeMs is the wall-clock time in milliseconds, or nil.
	GenerationTimeMs *int
	// TotalAgentsFiles is the count of target AGENTS.md files written.
	TotalAgentsFiles int
	// DirectoriesAnalyzed is the count of directories considered.
	DirectoriesAnalyzed int
}

// EfficiencyImprovement returns (AverageContextEfficiency-BaselineEfficiency)/BaselineEfficiency*100,
// or nil when BaselineEfficiency is nil.
func (o *OptimizationStats) EfficiencyImprovement() *float64 {
	if o.BaselineEfficiency == nil {
		return nil
	}
	v := (o.AverageContextEfficiency - *o.BaselineEfficiency) / *o.BaselineEfficiency * 100
	return &v
}

// EfficiencyPercentage returns AverageContextEfficiency * 100.
func (o *OptimizationStats) EfficiencyPercentage() float64 {
	return o.AverageContextEfficiency * 100
}

// CompilationResults holds the complete results from a compilation process.
type CompilationResults struct {
	// ProjectAnalysis describes the scanned project.
	ProjectAnalysis ProjectAnalysis
	// OptimizationDecisions lists per-instruction placement decisions.
	OptimizationDecisions []OptimizationDecision
	// PlacementSummaries lists per-output-file summaries.
	PlacementSummaries []PlacementSummary
	// OptimizationStats holds performance statistics.
	OptimizationStats OptimizationStats
	// Warnings lists non-fatal issues encountered.
	Warnings []string
	// Errors lists fatal issues encountered.
	Errors []string
	// IsDryRun reports whether this was a dry-run (no files written).
	IsDryRun bool
	// TargetName is the output filename (default "AGENTS.md").
	TargetName string
}

// TotalInstructions returns the sum of InstructionCount across all PlacementSummaries.
func (c *CompilationResults) TotalInstructions() int {
	total := 0
	for _, s := range c.PlacementSummaries {
		total += s.InstructionCount
	}
	return total
}

// HasIssues returns true when there are any warnings or errors.
func (c *CompilationResults) HasIssues() bool {
	return len(c.Warnings) > 0 || len(c.Errors) > 0
}

// --- helpers (no external deps) ---

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func joinStrings(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}
