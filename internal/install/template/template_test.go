package template_test

import (
	"errors"
	"testing"

	"github.com/githubnext/apm/internal/install/template"
)

// mockDiag implements DiagnosticsCounter.
type mockDiag struct {
	errors []string
	counts map[string]int
}

func (m *mockDiag) CountForPackage(depKey, kind string) int { return m.counts[depKey+":"+kind] }
func (m *mockDiag) AddError(msg, pkg string)                { m.errors = append(m.errors, msg) }

// mockLogger implements Logger.
type mockLogger struct {
	verbose  bool
	warnings []string
}

func (m *mockLogger) Verbose() bool                   { return m.verbose }
func (m *mockLogger) PackageInlineWarning(msg string) { m.warnings = append(m.warnings, msg) }

func makeConfig(secGate template.SecurityGateFunc, integrateFn template.IntegrateFunc) *template.Config {
	diag := &mockDiag{counts: map[string]int{}}
	return &template.Config{
		SecurityGate:        secGate,
		Integrate:           integrateFn,
		Diagnostics:         diag,
		Logger:              &mockLogger{},
		ProjectRoot:         "/project",
		HasTargets:          true,
		Force:               false,
		IntegrateErrorPrefix: "integrate error",
		PackageDeployedFiles: map[string][]string{},
	}
}

func TestRunIntegrationTemplate_NilMaterialization(t *testing.T) {
	cfg := makeConfig(nil, nil)
	result := template.RunIntegrationTemplate(nil, cfg)
	if result != nil {
		t.Errorf("expected nil result for nil materialization, got %v", result)
	}
}

func TestRunIntegrationTemplate_NoTargets(t *testing.T) {
	cfg := makeConfig(nil, nil)
	cfg.HasTargets = false
	m := &template.Materialization{
		InstallPath: "/pkg",
		DepKey:      "owner/repo",
		PackageInfo: &template.PackageInfo{Name: "repo", Path: "/pkg"},
		Deltas:      template.Deltas{},
	}
	result := template.RunIntegrationTemplate(m, cfg)
	if result == nil {
		t.Error("expected non-nil deltas")
	}
	if len(cfg.PackageDeployedFiles["owner/repo"]) != 0 {
		t.Error("expected empty deployed files when no targets")
	}
}

func TestRunIntegrationTemplate_SecurityGateBlocks(t *testing.T) {
	blockGate := func(installPath, pkgName string, force bool) bool { return false }
	integrated := false
	integrateFn := func(info *template.PackageInfo, projectRoot string) (*template.IntegrationResult, error) {
		integrated = true
		return &template.IntegrationResult{Prompts: 1}, nil
	}
	cfg := makeConfig(blockGate, integrateFn)
	m := &template.Materialization{
		InstallPath: "/pkg",
		DepKey:      "owner/repo",
		PackageInfo: &template.PackageInfo{Name: "repo", Path: "/pkg"},
	}
	template.RunIntegrationTemplate(m, cfg)
	if integrated {
		t.Error("expected integrate to NOT be called when security gate blocks")
	}
}

func TestRunIntegrationTemplate_IntegrateError(t *testing.T) {
	allowGate := func(installPath, pkgName string, force bool) bool { return true }
	failIntegrate := func(info *template.PackageInfo, projectRoot string) (*template.IntegrationResult, error) {
		return nil, errors.New("failed to integrate")
	}
	diag := &mockDiag{counts: map[string]int{}}
	cfg := &template.Config{
		SecurityGate:        allowGate,
		Integrate:           failIntegrate,
		Diagnostics:         diag,
		Logger:              &mockLogger{},
		ProjectRoot:         "/project",
		HasTargets:          true,
		IntegrateErrorPrefix: "prefix",
		PackageDeployedFiles: map[string][]string{},
	}
	m := &template.Materialization{
		InstallPath: "/pkg",
		DepKey:      "owner/repo",
		PackageInfo: &template.PackageInfo{Name: "repo", Path: "/pkg"},
	}
	template.RunIntegrationTemplate(m, cfg)
	if len(diag.errors) == 0 {
		t.Error("expected an error to be recorded in diagnostics")
	}
}

func TestRunIntegrationTemplate_SuccessfulIntegration(t *testing.T) {
	allowGate := func(installPath, pkgName string, force bool) bool { return true }
	integrateFn := func(info *template.PackageInfo, projectRoot string) (*template.IntegrationResult, error) {
		return &template.IntegrationResult{
			Prompts:       3,
			Skills:        2,
			DeployedFiles: []string{"f1.txt", "f2.txt"},
		}, nil
	}
	cfg := makeConfig(allowGate, integrateFn)
	m := &template.Materialization{
		InstallPath: "/pkg",
		DepKey:      "owner/repo",
		PackageInfo: &template.PackageInfo{Name: "repo", Path: "/pkg"},
	}
	deltas := template.RunIntegrationTemplate(m, cfg)
	if deltas["prompts"] != 3 {
		t.Errorf("prompts = %d, want 3", deltas["prompts"])
	}
	if deltas["skills"] != 2 {
		t.Errorf("skills = %d, want 2", deltas["skills"])
	}
	deployed := cfg.PackageDeployedFiles["owner/repo"]
	if len(deployed) != 2 {
		t.Errorf("expected 2 deployed files, got %d", len(deployed))
	}
}
