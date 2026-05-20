package template_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/template"
)

func TestDeltas_ZeroValue(t *testing.T) {
	var d template.Deltas
	if d != nil {
		t.Error("Deltas zero value should be nil")
	}
}

func TestDeltas_SetAndGet(t *testing.T) {
	d := template.Deltas{"prompts": 3, "agents": 1}
	if d["prompts"] != 3 {
		t.Error("prompts delta mismatch")
	}
	if d["agents"] != 1 {
		t.Error("agents delta mismatch")
	}
}

func TestPackageInfo_ZeroValue(t *testing.T) {
	var p template.PackageInfo
	if p.Name != "" || p.Path != "" {
		t.Error("PackageInfo zero value should have empty fields")
	}
}

func TestPackageInfo_Fields(t *testing.T) {
	p := template.PackageInfo{Name: "my-pkg", Path: "/opt/pkg"}
	if p.Name != "my-pkg" || p.Path != "/opt/pkg" {
		t.Error("PackageInfo field mismatch")
	}
}

func TestMaterialization_ZeroValue(t *testing.T) {
	var m template.Materialization
	if m.InstallPath != "" || m.DepKey != "" || m.PackageInfo != nil {
		t.Error("Materialization zero value should have empty/nil fields")
	}
}

func TestIntegrationResult_ZeroValue(t *testing.T) {
	var r template.IntegrationResult
	if r.Prompts != 0 || r.Agents != 0 || r.Skills != 0 || len(r.DeployedFiles) != 0 {
		t.Error("IntegrationResult zero value should have zero fields")
	}
}

func TestIntegrationResult_AllFields(t *testing.T) {
	r := template.IntegrationResult{
		Prompts:       2,
		Agents:        1,
		Skills:        3,
		SubSkills:     4,
		Instructions:  5,
		Commands:      6,
		Hooks:         7,
		LinksResolved: 8,
		DeployedFiles: []string{"a.txt", "b.txt"},
	}
	if r.Prompts != 2 || r.Commands != 6 || len(r.DeployedFiles) != 2 {
		t.Error("IntegrationResult field mismatch")
	}
}

func TestRunIntegrationTemplate_NilMaterializationExtra2(t *testing.T) {
	result := template.RunIntegrationTemplate(nil, &template.Config{})
	if result != nil {
		t.Error("RunIntegrationTemplate with nil materialization should return nil")
	}
}

func TestConfig_ZeroValue(t *testing.T) {
	var cfg template.Config
	if cfg.ProjectRoot != "" || cfg.Force || cfg.IsLocal {
		t.Error("Config zero value should have empty/false fields")
	}
}
