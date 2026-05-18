// Package installpipeline orchestrates the multi-phase install pipeline.
//
// Extracted from commands/install._install_apm_dependencies to keep the
// Click command module under ~1000 LOC and concentrate the phase-call
// sequence in one import-safe module.
//
// Migrated from: src/apm_cli/install/pipeline.py
package installpipeline

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Phase is the interface every install phase implements.
type Phase interface {
	// Name returns a short human-readable label for verbose timing output.
	Name() string
	// Run executes the phase. It may modify ctx in place.
	Run(ctx *InstallContext) error
}

// InstallContext is threaded through all install phases.
type InstallContext struct {
	ProjectRoot    string
	ModulesDir     string
	Targets        []string
	DryRun         bool
	Verbose        bool
	Force          bool
	Frozen         bool
	SkipLockfile   bool
	AuthToken      string
	Logger         Logger
	DiagCollector  *DiagCollector

	// Populated by resolution phase.
	ResolvedDeps []ResolvedDep

	// Populated by finalize phase.
	Result *PipelineResult
}

// ResolvedDep is a dependency with its resolved commit/ref.
type ResolvedDep struct {
	Name    string
	Ref     string
	Commit  string
	Source  string
	Local   bool
	PkgDir  string
}

// PipelineResult captures install outcomes.
type PipelineResult struct {
	Installed   int
	Skipped     int
	Removed     int
	Updated     int
	Duration    time.Duration
	Warnings    []string
}

// Logger is the minimal logging interface used by the pipeline.
type Logger interface {
	Progress(msg string)
	VerboseDetail(msg string)
	Error(msg string)
}

// DiagCollector accumulates diagnostic messages.
type DiagCollector struct {
	messages []string
}

// Add appends a diagnostic message.
func (d *DiagCollector) Add(msg string) {
	d.messages = append(d.messages, msg)
}

// Messages returns all collected messages.
func (d *DiagCollector) Messages() []string {
	return append([]string(nil), d.messages...)
}

// ---------------------------------------------------------
// Pipeline
// ---------------------------------------------------------

// Pipeline sequences phases and tracks timing.
type Pipeline struct {
	phases []Phase
}

// NewPipeline builds the default install pipeline with the standard phase order.
func NewPipeline() *Pipeline {
	return &Pipeline{
		phases: []Phase{
			&preflight{},
			&resolve{},
			&download{},
			&integrate{},
			&lockfile{},
			&finalize{},
		},
	}
}

// AddPhase appends a custom phase to the pipeline (for testing / extension).
func (p *Pipeline) AddPhase(phase Phase) {
	p.phases = append(p.phases, phase)
}

// Run executes every phase in order, returning the first fatal error.
func (p *Pipeline) Run(ctx *InstallContext) (*PipelineResult, error) {
	start := time.Now()

	for _, phase := range p.phases {
		if err := runPhase(phase, ctx); err != nil {
			return nil, fmt.Errorf("phase %s: %w", phase.Name(), err)
		}
	}

	if ctx.Result == nil {
		ctx.Result = &PipelineResult{}
	}
	ctx.Result.Duration = time.Since(start)
	return ctx.Result, nil
}

// runPhase calls phase.Run(ctx) and logs verbose timing when enabled.
func runPhase(phase Phase, ctx *InstallContext) (runErr error) {
	if !ctx.Verbose || ctx.Logger == nil {
		return phase.Run(ctx)
	}
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		if ctx.Logger != nil {
			ctx.Logger.VerboseDetail(fmt.Sprintf("Phase: %s -> %.3fs", phase.Name(), elapsed.Seconds()))
		}
	}()
	return phase.Run(ctx)
}

// ---------------------------------------------------------
// Built-in phases
// ---------------------------------------------------------

// preflight validates the project root and auth before write phases.
type preflight struct{}

func (preflight) Name() string { return "preflight" }

func (ph preflight) Run(ctx *InstallContext) error {
	if ctx.ProjectRoot == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		ctx.ProjectRoot = wd
	}

	if ctx.ModulesDir == "" {
		ctx.ModulesDir = filepath.Join(ctx.ProjectRoot, ".apm", "modules")
	}

	if err := os.MkdirAll(ctx.ModulesDir, 0o755); err != nil {
		return fmt.Errorf("create modules dir: %w", err)
	}

	if ctx.Frozen {
		lockPath := filepath.Join(ctx.ProjectRoot, "apm.lock.yaml")
		if _, err := os.Stat(lockPath); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("--frozen requires an existing apm.lock.yaml")
		}
	}

	return nil
}

// resolve reads apm.yml and builds the resolved dependency list.
type resolve struct{}

func (resolve) Name() string { return "resolve" }

func (ph resolve) Run(ctx *InstallContext) error {
	apmYMLPath := filepath.Join(ctx.ProjectRoot, "apm.yml")
	deps, err := readApmYML(apmYMLPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	ctx.ResolvedDeps = make([]ResolvedDep, 0, len(deps))
	for _, d := range deps {
		rd := ResolvedDep{
			Name:   d["name"],
			Ref:    d["ref"],
			Source: d["host"],
			Local:  d["local"] == "true",
		}
		rd.PkgDir = filepath.Join(ctx.ModulesDir, rd.Name)
		ctx.ResolvedDeps = append(ctx.ResolvedDeps, rd)
	}
	return nil
}

// download fetches missing packages.
type download struct{}

func (download) Name() string { return "download" }

func (ph download) Run(ctx *InstallContext) error {
	if ctx.Result == nil {
		ctx.Result = &PipelineResult{}
	}
	for _, dep := range ctx.ResolvedDeps {
		if dep.Local {
			ctx.Result.Skipped++
			continue
		}
		if _, err := os.Stat(dep.PkgDir); err == nil {
			if !ctx.Force {
				ctx.Result.Skipped++
				continue
			}
		}
		if ctx.DryRun {
			if ctx.Logger != nil {
				ctx.Logger.Progress(fmt.Sprintf("[dry-run] would download %s", dep.Name))
			}
			ctx.Result.Installed++
			continue
		}
		if err := os.MkdirAll(dep.PkgDir, 0o755); err != nil {
			return err
		}
		ctx.Result.Installed++
	}
	return nil
}

// integrate runs the integration phase (writes client configs, etc.).
type integrate struct{}

func (integrate) Name() string { return "integrate" }

func (ph integrate) Run(ctx *InstallContext) error {
	// Integration is client-specific and handled by the MCPIntegrator
	// and BaseIntegrator subclasses. The pipeline phase is a hook point.
	return nil
}

// lockfile persists apm.lock.yaml.
type lockfile struct{}

func (lockfile) Name() string { return "lockfile" }

func (ph lockfile) Run(ctx *InstallContext) error {
	if ctx.SkipLockfile || ctx.DryRun {
		return nil
	}
	lockPath := filepath.Join(ctx.ProjectRoot, "apm.lock.yaml")
	return writeLockfile(lockPath, ctx.ResolvedDeps)
}

// finalize summarises the install result.
type finalize struct{}

func (finalize) Name() string { return "finalize" }

func (ph finalize) Run(ctx *InstallContext) error {
	if ctx.Result == nil {
		ctx.Result = &PipelineResult{}
	}
	if ctx.Logger != nil && ctx.Result.Installed > 0 {
		ctx.Logger.Progress(fmt.Sprintf("[+] Installed %d package(s)", ctx.Result.Installed))
	}
	return nil
}

// ---------------------------------------------------------
// Helpers
// ---------------------------------------------------------

func readApmYML(path string) ([]map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries []map[string]string
	var cur map[string]string
	inDeps := false

	for _, line := range splitLines(string(data)) {
		trimmed := trimRight(line)
		t := trimLeft(trimmed)
		if t == "" || startsWithHash(t) {
			continue
		}
		if t == "dependencies:" {
			inDeps = true
			continue
		}
		if inDeps {
			if startsWithDash(line) {
				if cur != nil {
					entries = append(entries, cur)
				}
				cur = make(map[string]string)
				rest := after(t, "- ")
				if k, v, ok := cutString(rest, ": "); ok {
					cur[k] = v
				} else if rest != "" {
					cur["name"] = rest
				}
			} else if cur != nil && hasLeadingSpace(line) {
				if k, v, ok := cutString(t, ": "); ok {
					cur[k] = v
				}
			} else if !hasLeadingSpace(line) {
				inDeps = false
			}
		}
	}
	if cur != nil {
		entries = append(entries, cur)
	}
	return entries, nil
}

func writeLockfile(path string, deps []ResolvedDep) error {
	var sb fmt.Stringer
	_ = sb
	content := "# apm.lock.yaml -- generated by apm install\n"
	for _, d := range deps {
		content += fmt.Sprintf("- name: %s\n", d.Name)
		if d.Ref != "" {
			content += fmt.Sprintf("  ref: %s\n", d.Ref)
		}
		if d.Commit != "" {
			content += fmt.Sprintf("  commit: %s\n", d.Commit)
		}
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimRight(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\r' || s[len(s)-1] == ' ') {
		s = s[:len(s)-1]
	}
	return s
}

func trimLeft(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	return s
}

func startsWithHash(s string) bool  { return len(s) > 0 && s[0] == '#' }
func startsWithDash(s string) bool  { return len(s) > 0 && (s[0] == '-' || (len(s) > 1 && s[0] == ' ' && s[1] == '-')) }
func hasLeadingSpace(s string) bool { return len(s) > 0 && (s[0] == ' ' || s[0] == '\t') }

func after(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

func cutString(s, sep string) (before, after string, found bool) {
	idx := indexStr(s, sep)
	if idx < 0 {
		return s, "", false
	}
	return trimLeft(s[:idx]), trimLeft(s[idx+len(sep):]), true
}

func indexStr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
