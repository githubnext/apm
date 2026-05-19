package targetscmd

import (
	"encoding/json"
	"testing"
)

func TestTargetRowJSONRoundtrip(t *testing.T) {
	r := TargetRow{
		Target:    "copilot",
		Status:    "active",
		Source:    "apm.yml",
		DeployDir: ".github/copilot-instructions.md",
		Needs:     "",
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var r2 TargetRow
	if err := json.Unmarshal(data, &r2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if r2.Target != r.Target || r2.Status != r.Status || r2.DeployDir != r.DeployDir {
		t.Errorf("roundtrip mismatch: got %+v", r2)
	}
}

func TestTargetRowJSONInactiveNeeds(t *testing.T) {
	r := TargetRow{
		Target:    "claude",
		Status:    "inactive",
		DeployDir: "CLAUDE.md",
		Needs:     "CLAUDE.md",
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(data)
	if len(s) == 0 {
		t.Error("empty JSON")
	}
	// source omitted when empty
	var m map[string]any
	_ = json.Unmarshal(data, &m)
	if _, ok := m["source"]; ok {
		t.Error("source should be omitted when empty")
	}
}

func TestTargetRowAllFieldsExtra(t *testing.T) {
	r := TargetRow{
		Target:    "cursor",
		Status:    "active",
		Source:    "detected",
		DeployDir: ".cursorrules",
		Needs:     "",
	}
	if r.Target != "cursor" {
		t.Errorf("Target mismatch")
	}
	if r.Status != "active" {
		t.Errorf("Status mismatch")
	}
	if r.Source != "detected" {
		t.Errorf("Source mismatch")
	}
	if r.DeployDir != ".cursorrules" {
		t.Errorf("DeployDir mismatch")
	}
	if r.Needs != "" {
		t.Errorf("Needs should be empty")
	}
}

func TestTargetRowManyTargets(t *testing.T) {
	targets := []string{"vscode", "cursor", "copilot", "claude", "agents", "codex", "gemini"}
	for _, name := range targets {
		r := TargetRow{Target: name, Status: "inactive", DeployDir: "."}
		if r.Target != name {
			t.Errorf("Target mismatch for %s", name)
		}
	}
}
