// Package request defines InstallRequest, the typed input for the install pipeline.
package request

// InstallRequest bundles user intent for one install invocation.
type InstallRequest struct {
ApmPackagePath        string
UpdateRefs            bool
Verbose               bool
OnlyPackages          []string
Force                 bool
ParallelDownloads     int
Target                string
AllowInsecure         bool
AllowInsecureHosts    []string
NoPolicy              bool
SkillSubset           []string
SkillSubsetFromCLI    bool
LegacySkillPaths      bool
Frozen                bool
ProtocolPref          string
AllowProtocolFallback *bool
}

// DefaultInstallRequest returns an InstallRequest with sensible defaults.
func DefaultInstallRequest() InstallRequest {
return InstallRequest{
ParallelDownloads: 4,
}
}
