package targetscmd

import (
	"encoding/json"
	"testing"
)

func TestTargetRow_EmptyFields(t *testing.T) {
	row := TargetRow{}
	if row.Target != "" || row.Status != "" || row.Source != "" || row.DeployDir != "" || row.Needs != "" {
		t.Error("zero value TargetRow should have all empty fields")
	}
}

func TestTargetRow_ActiveStatus(t *testing.T) {
	row := TargetRow{
		Target:    "vscode",
		Status:    "active",
		Source:    ".vscode/settings.json",
		DeployDir: ".vscode/",
	}
	if row.Status != "active" {
		t.Errorf("expected active, got %q", row.Status)
	}
	if row.Needs != "" {
		t.Error("active row should have empty Needs")
	}
}

func TestTargetRow_InactiveStatus(t *testing.T) {
	row := TargetRow{
		Target:    "claude",
		Status:    "inactive",
		Needs:     "CLAUDE.md",
		DeployDir: ".claude/",
	}
	if row.Source != "" {
		t.Error("inactive row should have empty Source")
	}
	if row.Needs != "CLAUDE.md" {
		t.Errorf("expected CLAUDE.md, got %q", row.Needs)
	}
}

func TestTargetRow_JSONOmitempty(t *testing.T) {
	row := TargetRow{
		Target:    "copilot",
		Status:    "inactive",
		DeployDir: ".github/",
	}
	b, err := json.Marshal(row)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["source"]; ok {
		t.Error("source should be omitted when empty")
	}
	if _, ok := m["needs"]; ok {
		t.Error("needs should be omitted when empty")
	}
}

func TestTargetRow_JSONFields(t *testing.T) {
	row := TargetRow{
		Target:    "cursor",
		Status:    "active",
		Source:    ".cursor/settings.json",
		DeployDir: ".cursor/",
		Needs:     "",
	}
	b, err := json.Marshal(row)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m["target"] != "cursor" {
		t.Errorf("unexpected target %v", m["target"])
	}
	if m["deploy_dir"] != ".cursor/" {
		t.Errorf("unexpected deploy_dir %v", m["deploy_dir"])
	}
}

func TestTargetRow_MultipleRows(t *testing.T) {
	rows := []TargetRow{
		{Target: "vscode", Status: "active", DeployDir: ".vscode/"},
		{Target: "claude", Status: "inactive", Needs: "CLAUDE.md", DeployDir: ".claude/"},
		{Target: "cursor", Status: "active", DeployDir: ".cursor/"},
	}
	b, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	var result []TargetRow
	if err := json.Unmarshal(b, &result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 rows, got %d", len(result))
	}
	if result[1].Needs != "CLAUDE.md" {
		t.Errorf("expected CLAUDE.md, got %q", result[1].Needs)
	}
}
