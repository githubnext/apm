// Package models provides data models for compilation output and results.
package models

// PlacementStrategy represents how instructions are placed across the project.
type PlacementStrategy string

const (
	PlacementStrategySinglePoint   PlacementStrategy = "Single Point"
	PlacementStrategySelectiveMulti PlacementStrategy = "Selective Multi"
	PlacementStrategyDistributed   PlacementStrategy = "Distributed"
)

// ProjectAnalysis holds analysis of the project structure and file distribution.
type ProjectAnalysis struct {
	DirectoriesScanned          int
	FilesAnalyzed               int
	FileTypesDetected           []string
	InstructionPatternsDetected int
	MaxDepth                    int
	ConstitutionDetected        bool
	ConstitutionPath            string
}

// GetFileTypesSummary returns a concise summary of detected file types.
func (p *ProjectAnalysis) GetFileTypesSummary() string {
	if len(p.FileTypesDetected) == 0 {
		return "none"
	}
	types := make([]string, 0, len(p.FileTypesDetected))
	for _, t := range p.FileTypesDetected {
		stripped := t
		for len(stripped) > 0 && stripped[0] == '.' {
			stripped = stripped[1:]
		}
		if stripped != "" {
			types = append(types, stripped)
		}
	}
	// Simple sort
	for i := 0; i < len(types); i++ {
		for j := i + 1; j < len(types); j++ {
			if types[j] < types[i] {
				types[i], types[j] = types[j], types[i]
			}
		}
	}
	if len(types) <= 3 {
		result := ""
		for i, t := range types {
			if i > 0 {
				result += ", "
			}
			result += t
		}
		return result
	}
	result := types[0] + ", " + types[1] + ", " + types[2]
	return result + " and " + itoa(len(types)-3) + " more"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}

// OptimizationDecision holds details about a specific optimization decision for an instruction.
type OptimizationDecision struct {
	InstructionName      string
	Pattern              string
	MatchingDirectories  int
	TotalDirectories     int
	DistributionScore    float64
	Strategy             PlacementStrategy
	PlacementDirectories []string
	Reasoning            string
	RelevanceScore       float64
}

// DistributionRatio returns matching/total directories ratio.
func (o *OptimizationDecision) DistributionRatio() float64 {
	if o.TotalDirectories == 0 {
		return 0.0
	}
	return float64(o.MatchingDirectories) / float64(o.TotalDirectories)
}

// PlacementSummary summarizes a single AGENTS.md file placement.
type PlacementSummary struct {
	Path             string
	InstructionCount int
	SourceCount      int
	Sources          []string
}

// OptimizationStats holds performance and efficiency statistics from optimization.
type OptimizationStats struct {
	AverageContextEfficiency float64
	PollutionImprovement     *float64
	BaselineEfficiency       *float64
	PlacementAccuracy        *float64
	GenerationTimeMs         *int
	TotalAgentsFiles         int
	DirectoriesAnalyzed      int
}

// EfficiencyImprovement calculates efficiency improvement percentage.
func (o *OptimizationStats) EfficiencyImprovement() *float64 {
	if o.BaselineEfficiency != nil && *o.BaselineEfficiency != 0 {
		v := (o.AverageContextEfficiency - *o.BaselineEfficiency) / *o.BaselineEfficiency * 100
		return &v
	}
	return nil
}

// EfficiencyPercentage returns efficiency as percentage.
func (o *OptimizationStats) EfficiencyPercentage() float64 {
	return o.AverageContextEfficiency * 100
}

// CompilationResults holds complete results from the compilation process.
type CompilationResults struct {
	ProjectAnalysis      *ProjectAnalysis
	OptimizationDecisions []OptimizationDecision
	PlacementSummaries   []PlacementSummary
	OptimizationStats    *OptimizationStats
	Warnings             []string
	Errors               []string
	IsDryRun             bool
	TargetName           string
}

// TotalInstructions returns the total number of instructions processed.
func (c *CompilationResults) TotalInstructions() int {
	total := 0
	for _, s := range c.PlacementSummaries {
		total += s.InstructionCount
	}
	return total
}

// HasIssues returns true if there are any warnings or errors.
func (c *CompilationResults) HasIssues() bool {
	return len(c.Warnings) > 0 || len(c.Errors) > 0
}

// NewCompilationResults creates a new CompilationResults with defaults.
func NewCompilationResults() *CompilationResults {
	return &CompilationResults{
		TargetName: "AGENTS.md",
	}
}
