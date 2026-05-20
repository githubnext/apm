package results

import "testing"

func TestInstallResult_PackageTypesMap(t *testing.T) {
	r := InstallResult{
		InstalledCount:    2,
		PromptsIntegrated: 1,
		AgentsIntegrated:  1,
		PackageTypes:      map[string]string{"dep1": "apm", "dep2": "mcp"},
	}
	if r.PackageTypes["dep1"] != "apm" {
		t.Errorf("expected PackageTypes[dep1]=apm, got %q", r.PackageTypes["dep1"])
	}
	if r.PackageTypes["dep2"] != "mcp" {
		t.Errorf("expected PackageTypes[dep2]=mcp, got %q", r.PackageTypes["dep2"])
	}
}

func TestInstallResult_NilPackageTypesIsNil(t *testing.T) {
	r := InstallResult{}
	if r.PackageTypes != nil {
		t.Error("zero-value PackageTypes should be nil")
	}
}

func TestPrimitiveCounts_ZeroValue(t *testing.T) {
	var c PrimitiveCounts
	if c.Prompts != 0 || c.Agents != 0 || c.Instructions != 0 ||
		c.Skills != 0 || c.Hooks != 0 || c.Commands != 0 {
		t.Error("all zero-value PrimitiveCounts fields should be 0")
	}
}

func TestPrimitiveCounts_SetAll(t *testing.T) {
	c := PrimitiveCounts{Prompts: 1, Agents: 2, Instructions: 3, Skills: 4, Hooks: 5, Commands: 6}
	if c.Prompts != 1 {
		t.Errorf("Prompts = %d", c.Prompts)
	}
	if c.Agents != 2 {
		t.Errorf("Agents = %d", c.Agents)
	}
	if c.Instructions != 3 {
		t.Errorf("Instructions = %d", c.Instructions)
	}
	if c.Skills != 4 {
		t.Errorf("Skills = %d", c.Skills)
	}
	if c.Hooks != 5 {
		t.Errorf("Hooks = %d", c.Hooks)
	}
	if c.Commands != 6 {
		t.Errorf("Commands = %d", c.Commands)
	}
}

func TestInstallResult_LargeCounters(t *testing.T) {
	r := InstallResult{InstalledCount: 9999, PromptsIntegrated: 500, AgentsIntegrated: 300}
	if r.InstalledCount != 9999 {
		t.Errorf("InstalledCount = %d", r.InstalledCount)
	}
}
