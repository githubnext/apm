package targets

import (
	"testing"
)

func TestTargetProfile_Prefix(t *testing.T) {
	p := &TargetProfile{Name: "copilot", RootDir: ".github"}
	if p.Prefix() != ".github/" {
		t.Errorf("expected .github/, got %s", p.Prefix())
	}
}

func TestTargetProfile_Supports(t *testing.T) {
	p := &TargetProfile{
		Primitives: map[string]PrimitiveMapping{
			"instructions": {Subdir: "instructions", Extension: ".md"},
		},
	}
	if !p.Supports("instructions") {
		t.Error("expected instructions to be supported")
	}
	if p.Supports("hooks") {
		t.Error("hooks should not be supported")
	}
}

func TestTargetProfile_EffectiveRoot(t *testing.T) {
	p := &TargetProfile{RootDir: ".github", UserRootDir: ".copilot"}
	if p.EffectiveRoot(false) != ".github" {
		t.Errorf("project scope should return RootDir")
	}
	if p.EffectiveRoot(true) != ".copilot" {
		t.Errorf("user scope should return UserRootDir")
	}
	p2 := &TargetProfile{RootDir: ".claude"}
	if p2.EffectiveRoot(true) != ".claude" {
		t.Errorf("user scope with no UserRootDir should fall back to RootDir")
	}
}

func TestTargetProfile_SupportsAtUserScope(t *testing.T) {
	p := &TargetProfile{
		UserSupported:            "partial",
		UnsupportedUserPrimitives: []string{"prompts"},
		Primitives: map[string]PrimitiveMapping{
			"instructions": {},
			"prompts":      {},
		},
	}
	if !p.SupportsAtUserScope("instructions") {
		t.Error("instructions should be supported at user scope")
	}
	if p.SupportsAtUserScope("prompts") {
		t.Error("prompts should not be supported at user scope")
	}
}

func TestTargetProfile_SupportsAtUserScope_notSupported(t *testing.T) {
	p := &TargetProfile{UserSupported: false}
	if p.SupportsAtUserScope("instructions") {
		t.Error("target with UserSupported=false should not support any primitive at user scope")
	}
}

func TestTargetProfile_EffectivePackPrefixes_default(t *testing.T) {
	p := &TargetProfile{RootDir: ".github"}
	pp := p.EffectivePackPrefixes()
	if len(pp) != 1 || pp[0] != ".github/" {
		t.Errorf("expected ['.github/'], got %v", pp)
	}
}

func TestTargetProfile_EffectivePackPrefixes_override(t *testing.T) {
	p := &TargetProfile{
		RootDir:     ".codex",
		PackPrefixes: []string{".codex/", ".agents/"},
	}
	pp := p.EffectivePackPrefixes()
	if len(pp) != 2 {
		t.Errorf("expected 2 pack prefixes, got %v", pp)
	}
}

func TestTargetProfile_DeployPath(t *testing.T) {
	p := &TargetProfile{RootDir: ".github"}
	got := p.DeployPath("/repo", "instructions", "foo.md")
	if got == "" {
		t.Error("expected non-empty deploy path")
	}
}

func TestKnownTargets_copilot(t *testing.T) {
	p, ok := KnownTargets["copilot"]
	if !ok {
		t.Fatal("copilot target missing from KnownTargets")
	}
	if p.RootDir != ".github" {
		t.Errorf("copilot root should be .github, got %s", p.RootDir)
	}
}

func TestKnownTargets_claude(t *testing.T) {
	p, ok := KnownTargets["claude"]
	if !ok {
		t.Fatal("claude target missing")
	}
	if !p.Supports("instructions") {
		t.Error("claude should support instructions")
	}
}

func TestKnownTargets_cursor(t *testing.T) {
	p, ok := KnownTargets["cursor"]
	if !ok {
		t.Fatal("cursor target missing")
	}
	if p.CompileFamily != "agents" {
		t.Errorf("expected agents compile family, got %s", p.CompileFamily)
	}
}

func TestGetIntegrationPrefixes_noNils(t *testing.T) {
	pp := GetIntegrationPrefixes(nil)
	if len(pp) == 0 {
		t.Error("expected at least one integration prefix")
	}
}

func TestActiveTargets_fallback(t *testing.T) {
	targets := ActiveTargets("/nonexistent/path/xyz", nil)
	if len(targets) == 0 {
		t.Error("expected fallback target")
	}
}

func TestActiveTargets_explicit(t *testing.T) {
	targets := ActiveTargets("/repo", []string{"claude"})
	if len(targets) != 1 || targets[0].Name != "claude" {
		t.Errorf("expected [claude], got %v", targets)
	}
}

func TestActiveTargets_all(t *testing.T) {
	targets := ActiveTargets("/repo", []string{"all"})
	if len(targets) < 5 {
		t.Errorf("expected many targets for 'all', got %d", len(targets))
	}
}

func TestActiveTargets_vscode_alias(t *testing.T) {
	targets := ActiveTargets("/repo", []string{"vscode"})
	if len(targets) != 1 || targets[0].Name != "copilot" {
		t.Errorf("vscode alias should resolve to copilot, got %v", targets)
	}
}

func TestShouldUseLegacySkillPaths_default(t *testing.T) {
	result := ShouldUseLegacySkillPaths()
	_ = result // just verify it doesn't panic
}

func TestApplyLegacySkillPaths_noChange(t *testing.T) {
	p := &TargetProfile{
		Name:    "claude",
		RootDir: ".claude",
		Primitives: map[string]PrimitiveMapping{
			"instructions": {Subdir: "rules"},
		},
	}
	result := ApplyLegacySkillPaths([]*TargetProfile{p})
	if len(result) != 1 {
		t.Errorf("expected 1 profile, got %d", len(result))
	}
}

func TestApplyLegacySkillPaths_clearsDeployRoot(t *testing.T) {
	p := &TargetProfile{
		Name:    "copilot",
		RootDir: ".github",
		Primitives: map[string]PrimitiveMapping{
			"skills": {Subdir: "skills", DeployRoot: ".agents"},
		},
	}
	result := ApplyLegacySkillPaths([]*TargetProfile{p})
	if result[0].Primitives["skills"].DeployRoot != "" {
		t.Errorf("expected DeployRoot cleared, got %q", result[0].Primitives["skills"].DeployRoot)
	}
}
