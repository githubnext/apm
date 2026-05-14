package targets

import (
"testing"
)

func TestKnownTargetsRegistered(t *testing.T) {
expected := []string{"copilot", "claude", "cursor", "opencode", "gemini", "codex", "windsurf", "agent-skills", "copilot-cowork"}
for _, name := range expected {
if _, ok := KnownTargets[name]; !ok {
t.Errorf("missing target %q", name)
}
}
}

func TestTargetPrefix(t *testing.T) {
tgt := KnownTargets["copilot"]
if got := tgt.Prefix(); got != ".github/" {
t.Errorf("expected .github/, got %s", got)
}
}

func TestTargetSupports(t *testing.T) {
tgt := KnownTargets["copilot"]
if !tgt.Supports("skills") {
t.Error("copilot should support skills")
}
if tgt.Supports("nonexistent") {
t.Error("copilot should not support nonexistent")
}
}

func TestForScopeProjectScope(t *testing.T) {
tgt := KnownTargets["copilot"]
scoped := tgt.ForScope(false)
if scoped == nil {
t.Fatal("ForScope(false) returned nil")
}
if scoped.RootDir != ".github" {
t.Errorf("expected .github, got %s", scoped.RootDir)
}
}

func TestForScopeUserScopeCopilot(t *testing.T) {
tgt := KnownTargets["copilot"]
scoped := tgt.ForScope(true)
if scoped == nil {
t.Fatal("ForScope(true) returned nil")
}
if scoped.RootDir != ".copilot" {
t.Errorf("expected .copilot, got %s", scoped.RootDir)
}
// prompts and instructions should be filtered out
if scoped.Supports("prompts") {
t.Error("prompts should be filtered at user scope")
}
if scoped.Supports("instructions") {
t.Error("instructions should be filtered at user scope")
}
if !scoped.Supports("skills") {
t.Error("skills should remain at user scope")
}
}

func TestForScopeNoUserSupport(t *testing.T) {
tgt := &TargetProfile{
Name:          "fake",
RootDir:       ".fake",
UserSupported: false,
Primitives:    map[string]PrimitiveMapping{},
}
if scoped := tgt.ForScope(true); scoped != nil {
t.Error("expected nil for unsupported user scope")
}
}

func TestApplyLegacySkillPaths(t *testing.T) {
profiles := []*TargetProfile{KnownTargets["copilot"], KnownTargets["claude"]}
result := ApplyLegacySkillPaths(profiles)
for _, p := range result {
if pm, ok := p.Primitives["skills"]; ok {
if pm.DeployRoot != "" {
t.Errorf("target %s: expected empty deploy_root after legacy, got %s", p.Name, pm.DeployRoot)
}
}
}
}

func TestGetIntegrationPrefixes(t *testing.T) {
prefixes := GetIntegrationPrefixes(nil)
found := false
for _, p := range prefixes {
if p == ".github/" {
found = true
break
}
}
if !found {
t.Error("expected .github/ in prefixes")
}
}

func TestActiveTargetsFallback(t *testing.T) {
// Non-existent project root -> should fallback to copilot
targets := ActiveTargets("/nonexistent/path", nil)
if len(targets) != 1 || targets[0].Name != "copilot" {
t.Errorf("expected fallback to copilot, got %v", targets)
}
}
