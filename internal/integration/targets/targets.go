// Package targets defines the registry of known integration target profiles
// (Copilot, Claude, Cursor, etc.) and helpers for target resolution.
//
// Migrated from src/apm_cli/integration/targets.py
package targets

import (
"os"
"path/filepath"
"strings"
)

// PrimitiveMapping describes where a single primitive type is deployed.
type PrimitiveMapping struct {
Subdir     string // subdirectory under target root
Extension  string // file extension or suffix
FormatID   string // opaque transformer tag
DeployRoot string // optional root override (empty = use target root)
}

// TargetProfile describes capabilities and layout of a single target tool.
type TargetProfile struct {
Name      string
RootDir   string
Primitives map[string]PrimitiveMapping

AutoCreate  bool
DetectByDir bool

UserSupported            interface{} // bool or "partial"
UserRootDir              string
UnsupportedUserPrimitives []string
RequiresFlag             string
GeneratedFiles           []string
PackPrefixes             []string
CompileFamily            string
HooksConfigDisplay       string

// Set by ForScope for dynamic-root targets.
ResolvedDeployRoot string
}

// Prefix returns the path prefix for this target (e.g. ".github/").
func (t *TargetProfile) Prefix() string {
return t.RootDir + "/"
}

// EffectivePackPrefixes returns the path prefixes used by pack-time filtering.
func (t *TargetProfile) EffectivePackPrefixes() []string {
if len(t.PackPrefixes) > 0 {
return t.PackPrefixes
}
return []string{t.Prefix()}
}

// Supports returns true if this target accepts the primitive.
func (t *TargetProfile) Supports(primitive string) bool {
_, ok := t.Primitives[primitive]
return ok
}

// EffectiveRoot returns the root directory for the given scope.
func (t *TargetProfile) EffectiveRoot(userScope bool) string {
if userScope && t.UserRootDir != "" {
return t.UserRootDir
}
return t.RootDir
}

// SupportsAtUserScope returns true if the primitive can be deployed at user scope.
func (t *TargetProfile) SupportsAtUserScope(primitive string) bool {
if t.UserSupported == false || t.UserSupported == nil {
return false
}
for _, u := range t.UnsupportedUserPrimitives {
if u == primitive {
return false
}
}
return t.Supports(primitive)
}

// DeployPath returns the filesystem path for deployment.
func (t *TargetProfile) DeployPath(projectRoot string, parts ...string) string {
if t.ResolvedDeployRoot != "" {
base := t.ResolvedDeployRoot
if len(parts) > 0 {
return filepath.Join(append([]string{base}, parts...)...)
}
return base
}
base := filepath.Join(projectRoot, t.RootDir)
if len(parts) > 0 {
return filepath.Join(append([]string{base}, parts...)...)
}
return base
}

// ForScope returns a scope-resolved copy of this profile.
// Returns nil if the target does not support user scope.
func (t *TargetProfile) ForScope(userScope bool) *TargetProfile {
if !userScope {
cp := *t
return &cp
}

// Check user_supported
switch v := t.UserSupported.(type) {
case bool:
if !v {
return nil
}
case string:
if v != "partial" {
return nil
}
case nil:
return nil
}

cp := *t
newRoot := t.UserRootDir
if newRoot == "" {
newRoot = t.RootDir
}

// Claude Code honors CLAUDE_CONFIG_DIR
if t.Name == "claude" {
if env := strings.TrimSpace(os.Getenv("CLAUDE_CONFIG_DIR")); env != "" {
home, _ := os.UserHomeDir()
abs := filepath.Clean(env)
if rel, err := filepath.Rel(home, abs); err == nil && !strings.HasPrefix(rel, "..") {
newRoot = filepath.ToSlash(rel)
} else {
newRoot = abs
}
}
}

cp.RootDir = newRoot

// Filter unsupported user primitives
if len(t.UnsupportedUserPrimitives) > 0 {
filtered := make(map[string]PrimitiveMapping)
unsup := make(map[string]bool, len(t.UnsupportedUserPrimitives))
for _, u := range t.UnsupportedUserPrimitives {
unsup[u] = true
}
for k, v := range t.Primitives {
if !unsup[k] {
filtered[k] = v
}
}
cp.Primitives = filtered
}

return &cp
}

// ShouldUseLegacySkillPaths returns true when APM_LEGACY_SKILL_PATHS is set.
func ShouldUseLegacySkillPaths() bool {
val := strings.ToLower(strings.TrimSpace(os.Getenv("APM_LEGACY_SKILL_PATHS")))
return val == "1" || val == "true" || val == "yes"
}

// ApplyLegacySkillPaths resets deploy_root on every skills primitive.
func ApplyLegacySkillPaths(profiles []*TargetProfile) []*TargetProfile {
result := make([]*TargetProfile, len(profiles))
for i, p := range profiles {
if pm, ok := p.Primitives["skills"]; ok && pm.DeployRoot != "" {
cp := *p
prims := make(map[string]PrimitiveMapping, len(p.Primitives))
for k, v := range p.Primitives {
prims[k] = v
}
pm.DeployRoot = ""
prims["skills"] = pm
cp.Primitives = prims
result[i] = &cp
} else {
result[i] = p
}
}
return result
}

// KnownTargets is the registry of all known integration targets.
var KnownTargets = map[string]*TargetProfile{
"copilot": {
Name:    "copilot",
RootDir: ".github",
Primitives: map[string]PrimitiveMapping{
"instructions": {Subdir: "instructions", Extension: ".instructions.md", FormatID: "github_instructions"},
"prompts":      {Subdir: "prompts", Extension: ".prompt.md", FormatID: "github_prompt"},
"agents":       {Subdir: "agents", Extension: ".agent.md", FormatID: "github_agent"},
"skills":       {Subdir: "skills", Extension: "/SKILL.md", FormatID: "skill_standard", DeployRoot: ".agents"},
"hooks":        {Subdir: "hooks", Extension: ".json", FormatID: "github_hooks"},
},
AutoCreate:               true,
DetectByDir:              true,
UserSupported:            "partial",
UserRootDir:              ".copilot",
UnsupportedUserPrimitives: []string{"prompts", "instructions"},
GeneratedFiles:           []string{"copilot-instructions.md"},
CompileFamily:            "vscode",
},
"claude": {
Name:    "claude",
RootDir: ".claude",
Primitives: map[string]PrimitiveMapping{
"instructions": {Subdir: "rules", Extension: ".md", FormatID: "claude_rules"},
"agents":       {Subdir: "agents", Extension: ".md", FormatID: "claude_agent"},
"commands":     {Subdir: "commands", Extension: ".md", FormatID: "claude_command"},
"skills":       {Subdir: "skills", Extension: "/SKILL.md", FormatID: "skill_standard"},
"hooks":        {Subdir: "hooks", Extension: ".json", FormatID: "claude_hooks"},
},
AutoCreate:         false,
DetectByDir:        true,
UserSupported:      true,
CompileFamily:      "claude",
HooksConfigDisplay: ".claude/settings.json",
},
"cursor": {
Name:    "cursor",
RootDir: ".cursor",
Primitives: map[string]PrimitiveMapping{
"instructions": {Subdir: "rules", Extension: ".mdc", FormatID: "cursor_rules"},
"agents":       {Subdir: "agents", Extension: ".md", FormatID: "cursor_agent"},
"commands":     {Subdir: "commands", Extension: ".md", FormatID: "claude_command"},
"skills":       {Subdir: "skills", Extension: "/SKILL.md", FormatID: "skill_standard", DeployRoot: ".agents"},
"hooks":        {Subdir: "hooks", Extension: ".json", FormatID: "cursor_hooks"},
},
AutoCreate:               false,
DetectByDir:              true,
UserSupported:            "partial",
UserRootDir:              ".cursor",
UnsupportedUserPrimitives: []string{"instructions"},
CompileFamily:            "agents",
HooksConfigDisplay:       ".cursor/hooks.json",
},
"opencode": {
Name:    "opencode",
RootDir: ".opencode",
Primitives: map[string]PrimitiveMapping{
"agents":   {Subdir: "agents", Extension: ".md", FormatID: "opencode_agent"},
"commands": {Subdir: "commands", Extension: ".md", FormatID: "opencode_command"},
"skills":   {Subdir: "skills", Extension: "/SKILL.md", FormatID: "skill_standard", DeployRoot: ".agents"},
},
AutoCreate:               false,
DetectByDir:              true,
UserSupported:            "partial",
UserRootDir:              ".config/opencode",
UnsupportedUserPrimitives: []string{"hooks"},
CompileFamily:            "agents",
},
"gemini": {
Name:    "gemini",
RootDir: ".gemini",
Primitives: map[string]PrimitiveMapping{
"commands": {Subdir: "commands", Extension: ".toml", FormatID: "gemini_command"},
"skills":   {Subdir: "skills", Extension: "/SKILL.md", FormatID: "skill_standard", DeployRoot: ".agents"},
"hooks":    {Subdir: "hooks", Extension: ".json", FormatID: "gemini_hooks"},
},
AutoCreate:         false,
DetectByDir:        true,
UserSupported:      true,
UserRootDir:        ".gemini",
CompileFamily:      "gemini",
HooksConfigDisplay: ".gemini/settings.json",
},
"codex": {
Name:    "codex",
RootDir: ".codex",
Primitives: map[string]PrimitiveMapping{
"agents": {Subdir: "agents", Extension: ".toml", FormatID: "codex_agent"},
"skills": {Subdir: "skills", Extension: "/SKILL.md", FormatID: "skill_standard", DeployRoot: ".agents"},
"hooks":  {Subdir: "", Extension: "hooks.json", FormatID: "codex_hooks"},
},
AutoCreate:         false,
DetectByDir:        true,
UserSupported:      "partial",
PackPrefixes:       []string{".codex/", ".agents/"},
CompileFamily:      "agents",
HooksConfigDisplay: ".codex/hooks.json",
},
"windsurf": {
Name:    "windsurf",
RootDir: ".windsurf",
Primitives: map[string]PrimitiveMapping{
"instructions": {Subdir: "rules", Extension: ".md", FormatID: "windsurf_rules"},
"agents":       {Subdir: "skills", Extension: "/SKILL.md", FormatID: "windsurf_agent_skill"},
"skills":       {Subdir: "skills", Extension: "/SKILL.md", FormatID: "skill_standard"},
"commands":     {Subdir: "workflows", Extension: ".md", FormatID: "windsurf_workflow"},
"hooks":        {Subdir: "", Extension: "hooks.json", FormatID: "windsurf_hooks"},
},
AutoCreate:               false,
DetectByDir:              true,
UserSupported:            "partial",
UserRootDir:              ".codeium/windsurf",
UnsupportedUserPrimitives: []string{"instructions"},
CompileFamily:            "agents",
HooksConfigDisplay:       ".windsurf/hooks.json",
},
"agent-skills": {
Name:    "agent-skills",
RootDir: ".agents",
Primitives: map[string]PrimitiveMapping{
"skills": {Subdir: "skills", Extension: "/SKILL.md", FormatID: "skill_standard"},
},
AutoCreate:    true,
DetectByDir:   false,
UserSupported: true,
UserRootDir:   ".agents",
},
"copilot-cowork": {
Name:    "copilot-cowork",
RootDir: "copilot-cowork",
Primitives: map[string]PrimitiveMapping{
"skills": {Subdir: "skills", Extension: "/SKILL.md", FormatID: "skill_standard"},
},
AutoCreate:    false,
DetectByDir:   false,
UserSupported: true,
RequiresFlag:  "copilot_cowork",
},
}

// GetIntegrationPrefixes returns all known target root prefixes.
func GetIntegrationPrefixes(profiles []*TargetProfile) []string {
source := profiles
if source == nil {
for _, p := range KnownTargets {
source = append(source, p)
}
}
seen := make(map[string]bool)
var prefixes []string
for _, t := range source {
// Dynamic-root targets (cowork) use cowork:// prefix
if t.RequiresFlag == "copilot_cowork" {
const coworkPrefix = "cowork://"
if !seen[coworkPrefix] {
seen[coworkPrefix] = true
prefixes = append(prefixes, coworkPrefix)
}
continue
}
if !seen[t.Prefix()] {
seen[t.Prefix()] = true
prefixes = append(prefixes, t.Prefix())
}
for _, m := range t.Primitives {
if m.DeployRoot != "" {
dp := m.DeployRoot + "/"
if !seen[dp] {
seen[dp] = true
prefixes = append(prefixes, dp)
}
}
}
}
return prefixes
}

// ActiveTargets returns the target profiles that should be deployed into projectRoot.
// Resolution order: explicit target -> directory detection -> fallback (copilot).
func ActiveTargets(projectRoot string, explicitTargets []string) []*TargetProfile {
if len(explicitTargets) > 0 {
profiles := make([]*TargetProfile, 0)
seen := make(map[string]bool)
for _, t := range explicitTargets {
canonical := t
if t == "vscode" || t == "agents" {
canonical = "copilot"
}
if canonical == "all" {
var all []*TargetProfile
for _, p := range KnownTargets {
if p.Name != "agent-skills" && p.Name != "copilot-cowork" {
all = append(all, p)
}
}
return all
}
if p, ok := KnownTargets[canonical]; ok && !seen[canonical] {
seen[canonical] = true
profiles = append(profiles, p)
}
}
return profiles
}

// Auto-detect by directory presence
var detected []*TargetProfile
for _, p := range KnownTargets {
if p.DetectByDir {
if fi, err := os.Stat(filepath.Join(projectRoot, p.RootDir)); err == nil && fi.IsDir() {
detected = append(detected, p)
}
}
}
if len(detected) > 0 {
return detected
}
return []*TargetProfile{KnownTargets["copilot"]}
}

// ResolveTargets returns scope-resolved target profiles.
func ResolveTargets(projectRoot string, userScope bool, explicitTargets []string) []*TargetProfile {
var raw []*TargetProfile
if userScope {
raw = activeTargetsUserScope(explicitTargets)
} else {
raw = ActiveTargets(projectRoot, explicitTargets)
}
resolved := make([]*TargetProfile, 0, len(raw))
for _, t := range raw {
scoped := t.ForScope(userScope)
if scoped != nil {
resolved = append(resolved, scoped)
}
}
return resolved
}

func activeTargetsUserScope(explicitTargets []string) []*TargetProfile {
home, _ := os.UserHomeDir()

if len(explicitTargets) > 0 {
profiles := make([]*TargetProfile, 0)
seen := make(map[string]bool)
for _, t := range explicitTargets {
canonical := t
if t == "vscode" || t == "agents" {
canonical = "copilot"
}
if canonical == "all" {
var all []*TargetProfile
for _, p := range KnownTargets {
if p.UserSupported != nil && p.UserSupported != false && p.Name != "copilot-cowork" {
all = append(all, p)
}
}
return all
}
if p, ok := KnownTargets[canonical]; ok {
us := p.UserSupported
if (us == true || us == "partial") && !seen[canonical] {
seen[canonical] = true
profiles = append(profiles, p)
}
}
}
return profiles
}

var detected []*TargetProfile
for _, p := range KnownTargets {
us := p.UserSupported
if (us == true || us == "partial") && p.DetectByDir {
root := p.EffectiveRoot(true)
if fi, err := os.Stat(filepath.Join(home, root)); err == nil && fi.IsDir() {
detected = append(detected, p)
}
}
}
if len(detected) > 0 {
return detected
}
return []*TargetProfile{KnownTargets["copilot"]}
}
