// Package installctx provides the mutable state passed between install pipeline phases.
//
// Each phase is a function run(ctx *InstallContext) that reads inputs already
// populated by earlier phases and writes its own outputs to the context.
package installctx

import (
	"path/filepath"
	"sync"
)

// InstallContext holds state shared across install pipeline phases.
// Fields are grouped by the phase that first populates them.
type InstallContext struct {
	mu sync.RWMutex

	// Required on construction
	ProjectRoot string
	ApmDir      string

	// Inputs: populated by the caller from CLI args / APMPackage
	UpdateRefs             bool
	ParallelDownloads      int
	TargetOverride         string
	AllowInsecure          bool
	AllowInsecureHosts     []string
	DryRun                 bool
	Force                  bool
	Verbose                bool
	Dev                    bool
	OnlyPackages           []string
	AllowProtocolFallback  *bool // nil => read env

	// Resolve phase outputs
	RootHasLocalPrimitives bool
	LockfilePath           string
	ApmModulesDir          string
	InstalledCount         int
	UnpinnedCount          int

	// Integrate phase outputs
	IntendedDepKeys        map[string]bool
	PackageDeployedFiles   map[string][]string
	PackageTypes           map[string]string
	PackageHashes          map[string]string
	ExpectedHashChangeDeps map[string]bool
	TotalPromptsIntegrated int
	TotalAgentsIntegrated  int
	TotalSkillsIntegrated  int
	TotalSubSkillsPromoted int
	TotalInstructionsIntegrated int
	TotalCommandsIntegrated int
	TotalHooksIntegrated   int
	TotalLinksResolved     int
	DirectDepFailed        bool

	// Policy gate
	PolicyEnforcementActive bool
	NoPolicy                bool
	SkillSubset             []string
	SkillSubsetFromCLI      bool

	// Local content tracking
	OldLocalDeployed    []string
	LocalDeployedFiles  []string
	LocalContentErrorsBefore int

	// Cowork integration
	CoworkNonsupportedWarned bool

	// Legacy opt-out
	LegacySkillPaths bool
}

// New creates an InstallContext with all maps and slices initialised.
func New(projectRoot, apmDir string) *InstallContext {
	return &InstallContext{
		ProjectRoot:             projectRoot,
		ApmDir:                  apmDir,
		ParallelDownloads:       4,
		AllowInsecureHosts:      make([]string, 0),
		OnlyPackages:            make([]string, 0),
		IntendedDepKeys:         make(map[string]bool),
		PackageDeployedFiles:    make(map[string][]string),
		PackageTypes:            make(map[string]string),
		PackageHashes:           make(map[string]string),
		ExpectedHashChangeDeps:  make(map[string]bool),
		OldLocalDeployed:        make([]string, 0),
		LocalDeployedFiles:      make([]string, 0),
	}
}

// ApmModulesDirOrDefault returns ApmModulesDir or the default path.
func (ctx *InstallContext) ApmModulesDirOrDefault() string {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	if ctx.ApmModulesDir != "" {
		return ctx.ApmModulesDir
	}
	return filepath.Join(ctx.ProjectRoot, "apm_modules")
}

// LockfilePathOrDefault returns LockfilePath or the default path.
func (ctx *InstallContext) LockfilePathOrDefault() string {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	if ctx.LockfilePath != "" {
		return ctx.LockfilePath
	}
	return filepath.Join(ctx.ProjectRoot, "apm.lock.yaml")
}
