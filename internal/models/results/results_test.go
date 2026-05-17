package results

import "testing"

func TestInstallResult(t *testing.T) {
	r := InstallResult{
		InstalledCount:    3,
		PromptsIntegrated: 2,
		AgentsIntegrated:  1,
		PackageTypes:      map[string]string{"foo": "skill", "bar": "agent"},
	}
	if r.InstalledCount != 3 {
		t.Errorf("InstalledCount = %d, want 3", r.InstalledCount)
	}
	if r.PackageTypes["foo"] != "skill" {
		t.Errorf("PackageTypes[foo] = %q, want skill", r.PackageTypes["foo"])
	}
}

func TestPrimitiveCounts(t *testing.T) {
	p := PrimitiveCounts{
		Prompts:      1,
		Agents:       2,
		Instructions: 3,
		Skills:       4,
		Hooks:        5,
		Commands:     6,
	}
	if p.Prompts != 1 || p.Commands != 6 {
		t.Error("PrimitiveCounts fields not set correctly")
	}
}

func TestInstallResult_Zero(t *testing.T) {
	var r InstallResult
	if r.InstalledCount != 0 || r.PromptsIntegrated != 0 || r.AgentsIntegrated != 0 {
		t.Error("zero-value InstallResult should have zero counts")
	}
}

func TestInstallResult_SinglePackage(t *testing.T) {
	r := InstallResult{
		InstalledCount:    1,
		PromptsIntegrated: 1,
		AgentsIntegrated:  0,
		PackageTypes:      map[string]string{"mypkg": "prompt"},
	}
	if len(r.PackageTypes) != 1 {
		t.Errorf("expected 1 package type, got %d", len(r.PackageTypes))
	}
	if r.PackageTypes["mypkg"] != "prompt" {
		t.Errorf("PackageTypes[mypkg] = %q, want prompt", r.PackageTypes["mypkg"])
	}
}

func TestInstallResult_ManyPackages(t *testing.T) {
	pkgs := map[string]string{
		"a": "skill",
		"b": "agent",
		"c": "prompt",
		"d": "skill",
		"e": "agent",
	}
	r := InstallResult{
		InstalledCount: 5,
		PackageTypes:   pkgs,
	}
	if r.InstalledCount != 5 {
		t.Errorf("InstalledCount = %d, want 5", r.InstalledCount)
	}
	for k, v := range pkgs {
		if r.PackageTypes[k] != v {
			t.Errorf("PackageTypes[%q] = %q, want %q", k, r.PackageTypes[k], v)
		}
	}
}

func TestPrimitiveCounts_Zero(t *testing.T) {
	var p PrimitiveCounts
	if p.Prompts != 0 || p.Agents != 0 || p.Instructions != 0 ||
		p.Skills != 0 || p.Hooks != 0 || p.Commands != 0 {
		t.Error("zero-value PrimitiveCounts should have all zeros")
	}
}

func TestPrimitiveCounts_AllFields(t *testing.T) {
	p := PrimitiveCounts{
		Prompts:      10,
		Agents:       20,
		Instructions: 30,
		Skills:       40,
		Hooks:        50,
		Commands:     60,
	}
	if p.Prompts != 10 {
		t.Errorf("Prompts = %d, want 10", p.Prompts)
	}
	if p.Agents != 20 {
		t.Errorf("Agents = %d, want 20", p.Agents)
	}
	if p.Instructions != 30 {
		t.Errorf("Instructions = %d, want 30", p.Instructions)
	}
	if p.Skills != 40 {
		t.Errorf("Skills = %d, want 40", p.Skills)
	}
	if p.Hooks != 50 {
		t.Errorf("Hooks = %d, want 50", p.Hooks)
	}
	if p.Commands != 60 {
		t.Errorf("Commands = %d, want 60", p.Commands)
	}
}

func TestInstallResult_NilPackageTypes(t *testing.T) {
	r := InstallResult{InstalledCount: 0}
	if r.PackageTypes != nil {
		// nil is valid; just ensure no panic on lookup
		_ = r.PackageTypes["key"]
	}
}
