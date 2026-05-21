package template_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/template"
)

// extraMockDiag is a diagnostics counter for extra tests.
type extraMockDiag struct {
	errors []string
	counts map[string]int
}

func (m *extraMockDiag) CountForPackage(depKey, kind string) int { return m.counts[depKey+":"+kind] }
func (m *extraMockDiag) AddError(msg, pkg string)               { m.errors = append(m.errors, msg) }

// extraMockLogger captures verbose warnings.
type extraMockLogger struct {
	verbose  bool
	warnings []string
}

func (m *extraMockLogger) Verbose() bool                   { return m.verbose }
func (m *extraMockLogger) PackageInlineWarning(msg string) { m.warnings = append(m.warnings, msg) }

func makeExtraConfig(secGate template.SecurityGateFunc, integrateFn template.IntegrateFunc) *template.Config {
	diag := &extraMockDiag{counts: map[string]int{}}
	return &template.Config{
		SecurityGate:         secGate,
		Integrate:            integrateFn,
		Diagnostics:          diag,
		Logger:               &extraMockLogger{},
		ProjectRoot:          "/project",
		HasTargets:           true,
		IntegrateErrorPrefix: "err prefix",
		PackageDeployedFiles: map[string][]string{},
	}
}

func TestRunIntegrationTemplate_NilPackageInfo(t *testing.T) {
	cfg := makeExtraConfig(nil, nil)
	m := &template.Materialization{
		InstallPath: "/install",
		DepKey:      "pkg/a",
		PackageInfo: nil,
	}
	result := template.RunIntegrationTemplate(m, cfg)
	if result == nil {
		t.Fatal("expected non-nil deltas even when PackageInfo is nil")
	}
	if files := cfg.PackageDeployedFiles["pkg/a"]; len(files) != 0 {
		t.Errorf("expected empty deployed files, got %v", files)
	}
}

func TestRunIntegrationTemplate_HasTargetsFalse(t *testing.T) {
	cfg := makeExtraConfig(nil, nil)
	cfg.HasTargets = false
	m := &template.Materialization{
		InstallPath: "/install",
		DepKey:      "pkg/b",
		PackageInfo: &template.PackageInfo{Name: "pkg/b", Path: "/install"},
	}
	result := template.RunIntegrationTemplate(m, cfg)
	if result == nil {
		t.Fatal("expected non-nil deltas")
	}
	if files := cfg.PackageDeployedFiles["pkg/b"]; len(files) != 0 {
		t.Errorf("expected empty deployed files when HasTargets=false, got %v", files)
	}
}

func TestRunIntegrationTemplate_IntegrateReturnsAllDeltas(t *testing.T) {
	cfg := makeExtraConfig(nil, func(info *template.PackageInfo, root string) (*template.IntegrationResult, error) {
		return &template.IntegrationResult{
			Prompts:       1,
			Agents:        2,
			Skills:        3,
			SubSkills:     4,
			Instructions:  5,
			Commands:      6,
			Hooks:         7,
			LinksResolved: 8,
			DeployedFiles: []string{"a.txt", "b.txt"},
		}, nil
	})
	m := &template.Materialization{
		InstallPath: "/install",
		DepKey:      "pkg/c",
		PackageInfo: &template.PackageInfo{Name: "pkg/c", Path: "/install"},
	}
	deltas := template.RunIntegrationTemplate(m, cfg)
	if deltas["prompts"] != 1 {
		t.Errorf("prompts = %d, want 1", deltas["prompts"])
	}
	if deltas["agents"] != 2 {
		t.Errorf("agents = %d, want 2", deltas["agents"])
	}
	if deltas["skills"] != 3 {
		t.Errorf("skills = %d, want 3", deltas["skills"])
	}
	if deltas["links_resolved"] != 8 {
		t.Errorf("links_resolved = %d, want 8", deltas["links_resolved"])
	}
	if len(cfg.PackageDeployedFiles["pkg/c"]) != 2 {
		t.Errorf("deployed files = %v, want 2 entries", cfg.PackageDeployedFiles["pkg/c"])
	}
}

func TestRunIntegrationTemplate_PreexistingDeltas(t *testing.T) {
	cfg := makeExtraConfig(nil, nil)
	m := &template.Materialization{
		InstallPath: "/install",
		DepKey:      "pkg/d",
		PackageInfo: nil,
		Deltas:      template.Deltas{"pre": 42},
	}
	result := template.RunIntegrationTemplate(m, cfg)
	if result["pre"] != 42 {
		t.Errorf("pre-existing delta should be preserved, got %v", result)
	}
}

func TestRunIntegrationTemplate_VerboseCollisionWarning(t *testing.T) {
	diag := &extraMockDiag{counts: map[string]int{"pkg/e:collision": 2}}
	log := &extraMockLogger{verbose: true}
	cfg := &template.Config{
		SecurityGate:         nil,
		Integrate:            nil,
		Diagnostics:          diag,
		Logger:               log,
		ProjectRoot:          "/proj",
		HasTargets:           true,
		IntegrateErrorPrefix: "err",
		PackageDeployedFiles: map[string][]string{},
	}
	m := &template.Materialization{
		InstallPath: "/install",
		DepKey:      "pkg/e",
		PackageInfo: &template.PackageInfo{Name: "pkg/e"},
	}
	template.RunIntegrationTemplate(m, cfg)
	if len(log.warnings) == 0 {
		t.Error("expected verbose warning for collision count > 0")
	}
}

func TestMaterialization_Fields(t *testing.T) {
	m := template.Materialization{
		InstallPath: "/install",
		DepKey:      "owner/repo",
		PackageInfo: &template.PackageInfo{Name: "repo", Path: "/install"},
		Deltas:      template.Deltas{"prompts": 3},
	}
	if m.InstallPath != "/install" {
		t.Errorf("InstallPath = %q", m.InstallPath)
	}
	if m.DepKey != "owner/repo" {
		t.Errorf("DepKey = %q", m.DepKey)
	}
	if m.PackageInfo.Name != "repo" {
		t.Errorf("PackageInfo.Name = %q", m.PackageInfo.Name)
	}
	if m.Deltas["prompts"] != 3 {
		t.Errorf("Deltas[prompts] = %d", m.Deltas["prompts"])
	}
}

func TestIntegrationResult_Fields(t *testing.T) {
	r := template.IntegrationResult{
		Prompts:       10,
		Agents:        20,
		Skills:        30,
		SubSkills:     40,
		Instructions:  50,
		Commands:      60,
		Hooks:         70,
		LinksResolved: 80,
		DeployedFiles: []string{"f1.txt"},
	}
	if r.Prompts != 10 || r.Hooks != 70 || r.LinksResolved != 80 {
		t.Errorf("IntegrationResult fields mismatch: %+v", r)
	}
}
