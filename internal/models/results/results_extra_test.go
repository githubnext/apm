package results

import "testing"

func TestInstallResult_AllFieldsZero(t *testing.T) {
	var r InstallResult
	if r.InstalledCount != 0 {
		t.Errorf("InstalledCount = %d", r.InstalledCount)
	}
	if r.PromptsIntegrated != 0 {
		t.Errorf("PromptsIntegrated = %d", r.PromptsIntegrated)
	}
	if r.AgentsIntegrated != 0 {
		t.Errorf("AgentsIntegrated = %d", r.AgentsIntegrated)
	}
	if r.PackageTypes != nil {
		t.Errorf("PackageTypes should be nil")
	}
}

func TestInstallResult_OverwritePackageType(t *testing.T) {
	r := InstallResult{
		PackageTypes: map[string]string{"pkg": "skill"},
	}
	r.PackageTypes["pkg"] = "agent"
	if r.PackageTypes["pkg"] != "agent" {
		t.Errorf("overwrite failed, got %q", r.PackageTypes["pkg"])
	}
}

func TestInstallResult_MissingKey(t *testing.T) {
	r := InstallResult{
		PackageTypes: map[string]string{"a": "skill"},
	}
	v := r.PackageTypes["nonexistent"]
	if v != "" {
		t.Errorf("missing key should return empty string, got %q", v)
	}
}

func TestInstallResult_NegativeCounts(t *testing.T) {
	r := InstallResult{InstalledCount: -1, PromptsIntegrated: -5}
	if r.InstalledCount != -1 {
		t.Errorf("InstalledCount = %d", r.InstalledCount)
	}
	if r.PromptsIntegrated != -5 {
		t.Errorf("PromptsIntegrated = %d", r.PromptsIntegrated)
	}
}

func TestPrimitiveCounts_Partial(t *testing.T) {
	p := PrimitiveCounts{Prompts: 3, Skills: 7}
	if p.Prompts != 3 {
		t.Errorf("Prompts = %d", p.Prompts)
	}
	if p.Skills != 7 {
		t.Errorf("Skills = %d", p.Skills)
	}
	if p.Agents != 0 {
		t.Errorf("Agents should be 0, got %d", p.Agents)
	}
}

func TestPrimitiveCounts_LargeValues(t *testing.T) {
	p := PrimitiveCounts{
		Prompts:      1000,
		Agents:       2000,
		Instructions: 3000,
		Skills:       4000,
		Hooks:        5000,
		Commands:     6000,
	}
	total := p.Prompts + p.Agents + p.Instructions + p.Skills + p.Hooks + p.Commands
	if total != 21000 {
		t.Errorf("total = %d, want 21000", total)
	}
}

func TestInstallResult_EmptyPackageTypes(t *testing.T) {
	r := InstallResult{
		InstalledCount: 0,
		PackageTypes:   map[string]string{},
	}
	if len(r.PackageTypes) != 0 {
		t.Errorf("expected empty map, got %d entries", len(r.PackageTypes))
	}
}

func TestInstallResult_PromptsAndAgents(t *testing.T) {
	r := InstallResult{
		InstalledCount:    10,
		PromptsIntegrated: 7,
		AgentsIntegrated:  3,
	}
	if r.PromptsIntegrated+r.AgentsIntegrated != 10 {
		t.Errorf("prompts+agents = %d, want 10", r.PromptsIntegrated+r.AgentsIntegrated)
	}
}

func TestPrimitiveCounts_OnlyHooks(t *testing.T) {
	p := PrimitiveCounts{Hooks: 42}
	if p.Hooks != 42 {
		t.Errorf("Hooks = %d", p.Hooks)
	}
	if p.Prompts != 0 || p.Agents != 0 || p.Instructions != 0 || p.Skills != 0 || p.Commands != 0 {
		t.Error("other fields should be zero")
	}
}

func TestInstallResult_PackageTypeVariants(t *testing.T) {
	types := []string{"skill", "agent", "prompt", "instruction", "hook", "command"}
	pkgs := make(map[string]string)
	for i, typ := range types {
		key := "pkg" + string(rune('a'+i))
		pkgs[key] = typ
	}
	r := InstallResult{InstalledCount: len(types), PackageTypes: pkgs}
	if r.InstalledCount != 6 {
		t.Errorf("InstalledCount = %d", r.InstalledCount)
	}
	if r.PackageTypes["pkga"] != "skill" {
		t.Errorf("pkga = %q", r.PackageTypes["pkga"])
	}
}
