package install

// InstallRequest holds typed user intent for one install invocation.
// It is immutable: built once by the CLI handler and consumed by the
// install service.
type InstallRequest struct {
	// Required
	ApmPackage interface{} // *models.APMPackage

	UpdateRefs            bool
	Verbose               bool
	OnlyPackages          []string
	Force                 bool
	ParallelDownloads     int
	Scope                 string // InstallScope
	AuthResolver          interface{}
	Target                string
	AllowInsecure         bool
	AllowInsecureHosts    []string
	MarketplaceProvenance map[string]interface{}
	ProtocolPref          interface{} // ProtocolPreference
	AllowProtocolFallback *bool       // nil => read env
	NoPolicy              bool
	SkillSubset           []string
	SkillSubsetFromCLI    bool
	LegacySkillPaths      bool
	Frozen                bool

	// PlanCallback is invoked after resolve completes and before downloads
	// begin.  Return true to proceed or false to abort cleanly.
	PlanCallback func(plan *UpdatePlan) bool
}

// NewInstallRequest creates an InstallRequest with sensible defaults.
func NewInstallRequest(apmPackage interface{}) *InstallRequest {
	return &InstallRequest{
		ApmPackage:        apmPackage,
		ParallelDownloads: 4,
	}
}
