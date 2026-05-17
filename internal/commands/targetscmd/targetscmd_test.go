package targetscmd

import "testing"

func TestTargetRowStruct(t *testing.T) {
	r := TargetRow{
		Target:    "vscode",
		Status:    "active",
		Source:    "apm.yml",
		DeployDir: "/home/user/.vscode",
	}
	if r.Target != "vscode" {
		t.Errorf("unexpected Target %q", r.Target)
	}
	if r.Status != "active" {
		t.Errorf("unexpected Status %q", r.Status)
	}
}
