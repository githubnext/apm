// Package compilationformatter formats compilation output for APM.
package compilationformatter

import (
	"fmt"
	"path/filepath"
	"strings"
)

// PlacementStrategy describes the optimization strategy used.
type PlacementStrategy string

const (
	StrategySinglePoint   PlacementStrategy = "Single Point"
	StrategySelectiveMulti PlacementStrategy = "Selective Multi"
	StrategyDistributed   PlacementStrategy = "Distributed"
)

// ProjectAnalysis holds analysis of the project structure.
type ProjectAnalysis struct {
	DirectoriesScanned         int
	FilesAnalyzed              int
	FileTypesDetected          []string
	InstructionPatternsDetected int
	MaxDepth                   int
	ConstitutionDetected       bool
	ConstitutionPath           string
}

// FileTypesSummary returns a concise summary of detected file types.
func (p *ProjectAnalysis) FileTypesSummary() string {
	if len(p.FileTypesDetected) == 0 {
		return "none"
	}
	types := make([]string, 0, len(p.FileTypesDetected))
	for _, t := range p.FileTypesDetected {
		types = append(types, strings.TrimPrefix(t, "."))
	}
	if len(types) <= 3 {
		return strings.Join(types, ", ")
	}
	return fmt.Sprintf("%s and %d more", strings.Join(types[:3], ", "), len(types)-3)
}

// OptimizationDecision holds details about a placement decision for one instruction.
type OptimizationDecision struct {
	Pattern               string
	InstructionFilePath   string // file_path.name equivalent
	MatchingDirectories   int
	TotalDirectories      int
	DistributionScore     float64
	Strategy              PlacementStrategy
	PlacementDirectories  []string
	Reasoning             string
	RelevanceScore        float64
}

// PlacementSummary summarises a single AGENTS.md file placement.
type PlacementSummary struct {
	Path             string
	InstructionCount int
	SourceCount      int
	Sources          []string
}

// RelativePath returns path relative to base, prefixed with "./" when at root.
func (s *PlacementSummary) RelativePath(base string) string {
	rel, err := filepath.Rel(base, s.Path)
	if err != nil {
		return s.Path
	}
	if rel == "." {
		return "."
	}
	return rel
}

// OptimizationStats holds efficiency statistics.
type OptimizationStats struct {
	AverageContextEfficiency float64
	PollutionImprovement     *float64
	BaselineEfficiency       *float64
	PlacementAccuracy        *float64
	GenerationTimeMs         *int
	TotalAgentsFiles         int
	DirectoriesAnalyzed      int
}

// EfficiencyPercentage returns efficiency as a percentage.
func (s *OptimizationStats) EfficiencyPercentage() float64 {
	return s.AverageContextEfficiency * 100
}

// EfficiencyImprovement returns efficiency improvement over baseline, if available.
func (s *OptimizationStats) EfficiencyImprovement() *float64 {
	if s.BaselineEfficiency == nil {
		return nil
	}
	v := (s.AverageContextEfficiency - *s.BaselineEfficiency) * 100
	return &v
}

// CompilationResults holds all results from a compilation run.
type CompilationResults struct {
	TargetName          string
	PlacementSummaries  []PlacementSummary
	OptimizationDecisions []OptimizationDecision
	ProjectAnalysis     *ProjectAnalysis
	OptimizationStats   OptimizationStats
	Warnings            []string
	Errors              []string
	IsDryRun            bool
}

// HasIssues returns true if there are any warnings or errors.
func (r *CompilationResults) HasIssues() bool {
	return len(r.Warnings) > 0 || len(r.Errors) > 0
}

// CompilationFormatter formats compilation output for the CLI.
type CompilationFormatter struct {
	UseColor   bool
	targetName string
}

// New creates a new CompilationFormatter.
func New(useColor bool) *CompilationFormatter {
	return &CompilationFormatter{UseColor: useColor, targetName: "AGENTS.md"}
}

// FormatDefault formats standard compilation output.
func (f *CompilationFormatter) FormatDefault(results *CompilationResults) string {
	f.targetName = results.TargetName
	var lines []string

	lines = append(lines, f.formatProjectDiscovery(results.ProjectAnalysis)...)
	lines = append(lines, "")
	lines = append(lines, f.formatOptimizationProgress(results.OptimizationDecisions, results.ProjectAnalysis)...)
	lines = append(lines, "")
	lines = append(lines, f.formatResultsSummary(results)...)

	if results.HasIssues() {
		lines = append(lines, "")
		lines = append(lines, f.formatIssues(results.Warnings, results.Errors)...)
	}

	return strings.Join(lines, "\n")
}

// FormatVerbose formats verbose compilation output with mathematical details.
func (f *CompilationFormatter) FormatVerbose(results *CompilationResults) string {
	f.targetName = results.TargetName
	var lines []string

	lines = append(lines, f.formatProjectDiscovery(results.ProjectAnalysis)...)
	lines = append(lines, "")
	lines = append(lines, f.formatOptimizationProgress(results.OptimizationDecisions, results.ProjectAnalysis)...)
	lines = append(lines, "")
	lines = append(lines, f.formatMathematicalAnalysis(results.OptimizationDecisions)...)
	lines = append(lines, "")
	lines = append(lines, f.formatCoverageExplanation(results.OptimizationStats)...)
	lines = append(lines, "")
	lines = append(lines, f.formatDetailedMetrics(results.OptimizationStats)...)
	lines = append(lines, "")
	lines = append(lines, f.formatFinalSummary(results)...)

	if results.HasIssues() {
		lines = append(lines, "")
		lines = append(lines, f.formatIssues(results.Warnings, results.Errors)...)
	}

	return strings.Join(lines, "\n")
}

// FormatDryRun formats dry-run output.
func (f *CompilationFormatter) FormatDryRun(results *CompilationResults) string {
	f.targetName = results.TargetName
	var lines []string

	lines = append(lines, f.formatProjectDiscovery(results.ProjectAnalysis)...)
	lines = append(lines, "")
	lines = append(lines, f.formatOptimizationProgress(results.OptimizationDecisions, results.ProjectAnalysis)...)
	lines = append(lines, "")
	lines = append(lines, f.formatDryRunSummary(results)...)

	if results.HasIssues() {
		lines = append(lines, "")
		lines = append(lines, f.formatIssues(results.Warnings, results.Errors)...)
	}

	return strings.Join(lines, "\n")
}

func (f *CompilationFormatter) formatProjectDiscovery(analysis *ProjectAnalysis) []string {
	lines := []string{"Analyzing project structure..."}

	if analysis == nil {
		return lines
	}

	if analysis.ConstitutionDetected {
		lines = append(lines, fmt.Sprintf("|- Constitution detected: %s", analysis.ConstitutionPath))
	}

	fileTypesSummary := analysis.FileTypesSummary()
	lines = append(lines,
		fmt.Sprintf("|- %d directories scanned (max depth: %d)", analysis.DirectoriesScanned, analysis.MaxDepth),
		fmt.Sprintf("|- %d files analyzed across %d file types (%s)", analysis.FilesAnalyzed, len(analysis.FileTypesDetected), fileTypesSummary),
		fmt.Sprintf("+- %d instruction patterns detected", analysis.InstructionPatternsDetected),
	)
	return lines
}

func (f *CompilationFormatter) formatOptimizationProgress(decisions []OptimizationDecision, analysis *ProjectAnalysis) []string {
	lines := []string{"Optimizing placements..."}

	if analysis != nil && analysis.ConstitutionDetected {
		lines = append(lines,
			fmt.Sprintf("%-25s %-15s %-10s -> %-25s (rel: 100%%)", "**", "constitution.md", "ALL", "./AGENTS.md"),
		)
	}

	for _, d := range decisions {
		pattern := d.Pattern
		if pattern == "" {
			pattern = "(global)"
		}

		source := "unknown"
		if d.InstructionFilePath != "" {
			source = d.InstructionFilePath
		}

		ratio := fmt.Sprintf("%d/%d dirs", d.MatchingDirectories, d.TotalDirectories)

		if len(d.PlacementDirectories) == 1 {
			placement := f.getRelativeDisplayPath(d.PlacementDirectories[0])
			relevance := d.RelevanceScore
			if relevance == 0 {
				relevance = 1.0
			}
			line := fmt.Sprintf("%-25s %-15s %-10s -> %-25s (rel: %.0f%%)",
				pattern, source, ratio, placement, relevance*100)
			lines = append(lines, line)
		} else {
			line := fmt.Sprintf("%-25s %-15s %-10s -> %d locations",
				pattern, source, ratio, len(d.PlacementDirectories))
			lines = append(lines, line)
		}
	}
	return lines
}

func (f *CompilationFormatter) formatResultsSummary(results *CompilationResults) []string {
	var lines []string

	fileCount := len(results.PlacementSummaries)
	plural := "s"
	if fileCount == 1 {
		plural = ""
	}
	summaryLine := fmt.Sprintf("Generated %d %s file%s", fileCount, results.TargetName, plural)
	if results.IsDryRun {
		summaryLine = fmt.Sprintf("[DRY RUN] Would generate %d %s file%s", fileCount, results.TargetName, plural)
	}
	lines = append(lines, summaryLine)

	stats := results.OptimizationStats
	effPct := stats.EfficiencyPercentage()
	metricLines := []string{fmt.Sprintf("+- Context efficiency:    %.1f%%", effPct)}

	if imp := stats.EfficiencyImprovement(); imp != nil {
		if *imp > 0 {
			metricLines[0] += fmt.Sprintf(" (baseline: %.1f%%, improvement: +%.0f%%)", *stats.BaselineEfficiency*100, *imp)
		} else {
			metricLines[0] += fmt.Sprintf(" (baseline: %.1f%%, change: %.0f%%)", *stats.BaselineEfficiency*100, *imp)
		}
	}

	if stats.PollutionImprovement != nil {
		pollutionPct := (1.0 - *stats.PollutionImprovement) * 100
		var improvementPct string
		if *stats.PollutionImprovement > 0 {
			improvementPct = fmt.Sprintf("-%.0f%%", *stats.PollutionImprovement*100)
		} else {
			improvementPct = fmt.Sprintf("+%.0f%%", -(*stats.PollutionImprovement)*100)
		}
		metricLines = append(metricLines, fmt.Sprintf("|- Average pollution:     %.1f%% (improvement: %s)", pollutionPct, improvementPct))
	}

	if stats.PlacementAccuracy != nil {
		metricLines = append(metricLines, fmt.Sprintf("|- Placement accuracy:    %.1f%% (mathematical optimum)", *stats.PlacementAccuracy*100))
	}

	if stats.GenerationTimeMs != nil {
		metricLines = append(metricLines, fmt.Sprintf("+- Generation time:       %dms", *stats.GenerationTimeMs))
	} else if len(metricLines) > 1 {
		metricLines[len(metricLines)-1] = strings.Replace(metricLines[len(metricLines)-1], "|-", "+-", 1)
	}

	lines = append(lines, metricLines...)
	lines = append(lines, "", "Placement Distribution")

	for i, summary := range results.PlacementSummaries {
		relPath := summary.RelativePath(".")
		contentText := f.getPlacementDescription(&summary)
		sourceText := fmt.Sprintf("%d source", summary.SourceCount)
		if summary.SourceCount != 1 {
			sourceText += "s"
		}
		prefix := "|-"
		if i == len(results.PlacementSummaries)-1 {
			prefix = "+-"
		}
		line := fmt.Sprintf("%s %-30s %s from %s", prefix, relPath, contentText, sourceText)
		lines = append(lines, line)
	}
	return lines
}

func (f *CompilationFormatter) formatFinalSummary(results *CompilationResults) []string {
	// In verbose mode use same structure as results summary with placement distribution.
	return f.formatResultsSummary(results)
}

func (f *CompilationFormatter) formatDryRunSummary(results *CompilationResults) []string {
	lines := []string{"[DRY RUN] File generation preview:"}

	for i, summary := range results.PlacementSummaries {
		relPath := summary.RelativePath(".")
		instrText := fmt.Sprintf("%d instruction", summary.InstructionCount)
		if summary.InstructionCount != 1 {
			instrText += "s"
		}
		srcText := fmt.Sprintf("%d source", summary.SourceCount)
		if summary.SourceCount != 1 {
			srcText += "s"
		}
		prefix := "|-"
		if i == len(results.PlacementSummaries)-1 {
			prefix = "+-"
		}
		lines = append(lines, fmt.Sprintf("%s %-30s %s, %s", prefix, relPath, instrText, srcText))
	}

	lines = append(lines, "", "[DRY RUN] No files written. Run 'apm compile' to apply changes.")
	return lines
}

func (f *CompilationFormatter) formatMathematicalAnalysis(decisions []OptimizationDecision) []string {
	lines := []string{"Mathematical Optimization Analysis", ""}
	lines = append(lines, "Coverage-First Strategy Analysis:")

	for _, d := range decisions {
		pattern := d.Pattern
		if pattern == "" {
			pattern = "(global)"
		}
		score := fmt.Sprintf("%.3f", d.DistributionScore)
		strategy := string(d.Strategy)
		var coverage string
		if d.DistributionScore < 0.7 {
			coverage = "[+] Verified"
		} else {
			coverage = "[!] Root Fallback"
		}
		lines = append(lines, fmt.Sprintf("  %-30s %-8s %-15s %s", pattern, score, strategy, coverage))
	}

	lines = append(lines, "",
		"Mathematical Foundation:",
		"  Objective: minimize sum(context_pollution x directory_weight)",
		"  Constraints: for_allfile_matching_pattern -> can_inherit_instruction",
		"  Algorithm: Three-tier strategy with coverage verification",
		"  Principle: Coverage guarantee takes priority over efficiency",
	)
	return lines
}

func (f *CompilationFormatter) formatCoverageExplanation(stats OptimizationStats) []string {
	lines := []string{"Coverage vs. Efficiency Analysis", ""}

	efficiency := stats.EfficiencyPercentage()

	if efficiency < 30 {
		lines = append(lines,
			"[!] Low Efficiency Detected:",
			"   * Coverage guarantee requires some instructions at root level",
			"   * This creates pollution for specialized directories",
			"   * Trade-off: Guaranteed coverage vs. optimal efficiency",
			"   * Alternative: Higher efficiency with coverage violations (data loss)",
			"",
			"This may be mathematically optimal given coverage constraints",
		)
	} else if efficiency < 60 {
		lines = append(lines,
			"[+] Moderate Efficiency:",
			"   * Good balance between coverage and efficiency",
			"   * Some coverage-driven pollution is acceptable",
			"   * Most patterns are well-localized",
		)
	} else {
		lines = append(lines,
			"High Efficiency:",
			"   * Excellent pattern locality achieved",
			"   * Minimal coverage conflicts",
			"   * Instructions are optimally placed",
		)
	}

	lines = append(lines, "",
		"Why Coverage Takes Priority:",
		"   * Every file must access applicable instructions",
		"   * Hierarchical inheritance prevents data loss",
		"   * Better low efficiency than missing instructions",
	)
	return lines
}

func (f *CompilationFormatter) formatDetailedMetrics(stats OptimizationStats) []string {
	lines := []string{"Performance Metrics"}

	efficiency := stats.EfficiencyPercentage()
	pollution := 100 - efficiency

	effAssessment := assessEfficiency(efficiency)
	pollAssessment := assessPollution(pollution)

	lines = append(lines,
		fmt.Sprintf("Context Efficiency: %.1f%% (%s)", efficiency, effAssessment),
		fmt.Sprintf("Pollution Level:    %.1f%% (%s)", pollution, pollAssessment),
		"Guide: 80-100% Excellent | 60-80% Good | 40-60% Fair | 20-40% Poor | <20% Very Poor",
	)
	return lines
}

func assessEfficiency(v float64) string {
	switch {
	case v >= 80:
		return "Excellent"
	case v >= 60:
		return "Good"
	case v >= 40:
		return "Fair"
	case v >= 20:
		return "Poor"
	default:
		return "Very Poor"
	}
}

func assessPollution(v float64) string {
	switch {
	case v <= 10:
		return "Excellent"
	case v <= 25:
		return "Good"
	case v <= 50:
		return "Fair"
	default:
		return "Poor"
	}
}

func (f *CompilationFormatter) formatIssues(warnings, errors []string) []string {
	var lines []string
	for _, e := range errors {
		lines = append(lines, "x Error: "+e)
	}
	for _, w := range warnings {
		if strings.Contains(w, "\n") {
			wLines := strings.Split(w, "\n")
			lines = append(lines, "[!] Warning: "+wLines[0])
			for _, wl := range wLines[1:] {
				if strings.TrimSpace(wl) != "" {
					lines = append(lines, "           "+wl)
				}
			}
		} else {
			lines = append(lines, "[!] Warning: "+w)
		}
	}
	return lines
}

func (f *CompilationFormatter) getRelativeDisplayPath(path string) string {
	rel, err := filepath.Rel(".", path)
	if err != nil {
		return filepath.Join(path, f.targetName)
	}
	if rel == "." {
		return "./" + f.targetName
	}
	return filepath.ToSlash(filepath.Join(rel, f.targetName))
}

func (f *CompilationFormatter) getPlacementDescription(summary *PlacementSummary) string {
	hasConstitution := false
	for _, src := range summary.Sources {
		if strings.Contains(src, "constitution.md") {
			hasConstitution = true
			break
		}
	}

	var parts []string
	if hasConstitution {
		parts = append(parts, "Constitution")
	}
	if summary.InstructionCount > 0 {
		plural := "s"
		if summary.InstructionCount == 1 {
			plural = ""
		}
		parts = append(parts, fmt.Sprintf("%d instruction%s", summary.InstructionCount, plural))
	}
	if len(parts) > 0 {
		return strings.Join(parts, " and ")
	}
	return "content"
}
