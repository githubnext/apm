// Package template implements the shared post-acquire integration flow for all DependencySources.
// This is the Template Method companion to the Strategy pattern in install/sources.
package template

// Deltas holds counter-deltas from integration of one package.
type Deltas map[string]int

// PackageInfo is a minimal representation of a resolved package.
type PackageInfo struct {
Name string
Path string
}

// Materialization represents the result of a DependencySource.acquire() call.
type Materialization struct {
InstallPath string
DepKey      string
PackageInfo *PackageInfo
Deltas      Deltas
}

// IntegrationResult holds integration counts for one package.
type IntegrationResult struct {
Prompts        int
Agents         int
Skills         int
SubSkills      int
Instructions   int
Commands       int
Hooks          int
LinksResolved  int
DeployedFiles  []string
}

// SecurityGateFunc is the signature of the pre-deploy security gate.
type SecurityGateFunc func(installPath, packageName string, force bool) bool

// IntegrateFunc is the signature of the primitive integrator.
type IntegrateFunc func(info *PackageInfo, projectRoot string) (*IntegrationResult, error)

// DiagnosticsCounter supports per-package diagnostic counts.
type DiagnosticsCounter interface {
CountForPackage(depKey, kind string) int
AddError(msg, pkg string)
}

// Logger supports verbose package-inline warnings.
type Logger interface {
Verbose() bool
PackageInlineWarning(msg string)
}

// Config holds all dependencies for RunIntegrationTemplate.
type Config struct {
SecurityGate SecurityGateFunc
Integrate    IntegrateFunc
Diagnostics  DiagnosticsCounter
Logger       Logger
ProjectRoot  string
HasTargets   bool
Force        bool
// IntegrateErrorPrefix is the per-source error prefix (Strategy pattern).
IntegrateErrorPrefix string
// IsLocal indicates whether the dep ref is local (for error key selection).
IsLocal  bool
LocalPath string
// PackageDeployedFiles is updated in place.
PackageDeployedFiles map[string][]string
}

// RunIntegrationTemplate runs the shared post-acquire integration flow.
// Returns a counter-delta map, or nil if the materialization is nil (source declined).
func RunIntegrationTemplate(m *Materialization, cfg *Config) Deltas {
if m == nil {
return nil
}
return integrateMaterilaization(m, cfg)
}

func integrateMaterilaization(m *Materialization, cfg *Config) Deltas {
deltas := m.Deltas
if deltas == nil {
deltas = Deltas{}
}

// No-op when targets are empty or acquire decided to skip integration.
if m.PackageInfo == nil || !cfg.HasTargets {
cfg.PackageDeployedFiles[m.DepKey] = []string{}
return deltas
}

defer func() {
// Verbose: inline skip / error count for this package.
if cfg.Logger != nil && cfg.Logger.Verbose() {
skipCount := cfg.Diagnostics.CountForPackage(m.DepKey, "collision")
errCount := cfg.Diagnostics.CountForPackage(m.DepKey, "error")
if skipCount > 0 {
noun := "file"
if skipCount != 1 {
noun = "files"
}
cfg.Logger.PackageInlineWarning(
"    [!] " + itoa(skipCount) + " " + noun + " skipped (local files exist)",
)
}
if errCount > 0 {
noun := "error"
if errCount != 1 {
noun = "errors"
}
cfg.Logger.PackageInlineWarning(
"    [!] " + itoa(errCount) + " integration " + noun,
)
}
}
}()

// Pre-deploy security gate.
if cfg.SecurityGate != nil {
if !cfg.SecurityGate(m.InstallPath, m.DepKey, cfg.Force) {
cfg.PackageDeployedFiles[m.DepKey] = []string{}
return deltas
}
}

// Primitive integration.
if cfg.Integrate != nil {
result, err := cfg.Integrate(m.PackageInfo, cfg.ProjectRoot)
if err != nil {
packageKey := m.DepKey
if cfg.IsLocal && cfg.LocalPath != "" {
packageKey = cfg.LocalPath
}
cfg.Diagnostics.AddError(cfg.IntegrateErrorPrefix+": "+err.Error(), packageKey)
} else if result != nil {
deltas["prompts"] = result.Prompts
deltas["agents"] = result.Agents
deltas["skills"] = result.Skills
deltas["sub_skills"] = result.SubSkills
deltas["instructions"] = result.Instructions
deltas["commands"] = result.Commands
deltas["hooks"] = result.Hooks
deltas["links_resolved"] = result.LinksResolved
cfg.PackageDeployedFiles[m.DepKey] = result.DeployedFiles
}
}

return deltas
}

// itoa converts an int to a string without importing strconv at call sites.
func itoa(n int) string {
if n == 0 {
return "0"
}
neg := n < 0
if neg {
n = -n
}
buf := make([]byte, 20)
i := len(buf)
for n >= 10 {
i--
buf[i] = byte('0' + n%10)
n /= 10
}
i--
buf[i] = byte('0' + n)
if neg {
i--
buf[i] = '-'
}
return string(buf[i:])
}
