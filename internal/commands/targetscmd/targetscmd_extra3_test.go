package targetscmd

import (
	"encoding/json"
	"testing"
)

func TestTargetRow_ExtraFields_Extra3(t *testing.T) {
	row := TargetRow{
		Target:    "copilot",
		Status:    "active",
		Source:    "apm.yml",
		DeployDir: ".github",
		Needs:     "",
	}
	if row.Target != "copilot" {
		t.Errorf("Target = %q, want copilot", row.Target)
	}
	if row.DeployDir != ".github" {
		t.Errorf("DeployDir = %q, want .github", row.DeployDir)
	}
}

func TestTargetRow_JSONNeedsOmitted_Extra3(t *testing.T) {
	row := TargetRow{Target: "vscode", Status: "active"}
	b, err := json.Marshal(row)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := m["needs"]; ok {
		t.Error("needs should be omitted when empty")
	}
	if _, ok := m["source"]; ok {
		t.Error("source should be omitted when empty")
	}
}

func TestTargetRow_JSONNeedsPresent_Extra3(t *testing.T) {
	row := TargetRow{Target: "cursor", Status: "active", Needs: "node"}
	b, err := json.Marshal(row)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["needs"] != "node" {
		t.Errorf("needs = %v, want node", m["needs"])
	}
}

func TestTargetRow_DeployDirJSON_Extra3(t *testing.T) {
	row := TargetRow{DeployDir: ".vscode"}
	b, _ := json.Marshal(row)
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	if m["deploy_dir"] != ".vscode" {
		t.Errorf("deploy_dir = %v, want .vscode", m["deploy_dir"])
	}
}

func TestTargetRow_Slice_Extra3(t *testing.T) {
	rows := []TargetRow{
		{Target: "a", Status: "active"},
		{Target: "b", Status: "inactive"},
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].Target != "a" {
		t.Errorf("first row target = %q", rows[0].Target)
	}
}

func TestTargetRow_PointerEquality_Extra3(t *testing.T) {
	r1 := TargetRow{Target: "x"}
	r2 := r1
	r2.Target = "y"
	if r1.Target != "x" {
		t.Error("struct copy should not affect original")
	}
}
