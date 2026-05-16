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
