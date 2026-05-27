package install

import "path/filepath"

// InstallContext holds mutable state passed between install pipeline phases.
// Each phase reads inputs populated by earlier phases and writes its own
// outputs here.  This makes implicit shared state explicit and auditable.
type InstallContext struct {
	// Required on construction
	ProjectRoot string
	ApmDir      string

	// Inputs: from CLI args / APMPackage
	UpdateRefs             bool
	Scope                  string // InstallScope string value
	ParallelDownloads      int
	TargetOverride         string
	AllowInsecure          bool
	AllowInsecureHosts     []string
	DryRun                 bool
	Force                  bool
	Verbose                bool
	Dev                    bool
	OnlyPackages           []string
	NoPolicy               bool
	LegacySkillPaths       bool
	SkillSubset            []string
	SkillSubsetFromCLI     bool
	AllowProtocolFallback  *bool // nil means read APM_ALLOW_PROTOCOL_FALLBACK env

	// Resolve phase outputs
	AllApmDeps         []interface{}
	RootHasLocalPrims  bool
	DepsToInstall      []interface{}
	DependencyGraph    interface{}
	ExistingLockfile   interface{}
	LockfilePath       string
	ApmModulesDir      string
	CallbackDownloaded map[string]interface{}
	CallbackFailures   map[string]bool
	TransitiveFailures [][]string // pairs of [dep, reason]

	// Targets phase outputs
	Targets     []interface{}
	Integrators map[string]interface{}

	// Download phase outputs
	PreDownloadResults  map[string]interface{}
	PreDownloadedKeys   map[string]bool

	// Pre-integrate inputs
	Diagnostics    interface{}
	RegistryConfig interface{}
	ManagedFiles   map[string]bool

	// Integrate phase outputs
	IntendedDepKeys          map[string]bool
	PackageDeployedFiles     map[string][]string
	PackageTypes             map[string]string
	PackageHashes            map[string]string
	ExpectedHashChangeDeps   map[string]bool
	InstalledCount           int
	UnpinnedCount            int
	InstalledPackages        []interface{}
	TotalPromptsIntegrated   int
	TotalAgentsIntegrated    int
	TotalSkillsIntegrated    int
	TotalSubSkillsPromoted   int
	TotalInstructionsInteg   int
	TotalCommandsIntegrated  int
	TotalHooksIntegrated     int
	TotalLinksResolved       int
	DirectDepFailed          bool

	// Policy gate
	PolicyFetch             interface{}
	PolicyEnforcementActive bool
	EarlyLockfile           interface{}
	DirectMCPDeps           []interface{}

	// Post-deps local content tracking
	OldLocalDeployed       []string
	LocalDeployedFiles     []string
	LocalContentErrsBefore int

	// Cowork guard
	CoworkNonsupportedWarned bool
}

// NewInstallContext creates a minimal InstallContext with required fields.
func NewInstallContext(projectRoot, apmDir string) *InstallContext {
	return &InstallContext{
		ProjectRoot:            filepath.Clean(projectRoot),
		ApmDir:                 filepath.Clean(apmDir),
		ParallelDownloads:      4,
		CallbackDownloaded:     make(map[string]interface{}),
		CallbackFailures:       make(map[string]bool),
		Integrators:            make(map[string]interface{}),
		PreDownloadResults:     make(map[string]interface{}),
		PreDownloadedKeys:      make(map[string]bool),
		ManagedFiles:           make(map[string]bool),
		IntendedDepKeys:        make(map[string]bool),
		PackageDeployedFiles:   make(map[string][]string),
		PackageTypes:           make(map[string]string),
		PackageHashes:          make(map[string]string),
		ExpectedHashChangeDeps: make(map[string]bool),
	}
}
