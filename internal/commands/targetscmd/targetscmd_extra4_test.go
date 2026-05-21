package targetscmd

import (
"encoding/json"
"testing"
)

func TestTargetRow_SourceField_Extra4(t *testing.T) {
row := TargetRow{Source: "file.yaml"}
if row.Source != "file.yaml" {
t.Errorf("expected 'file.yaml', got %q", row.Source)
}
}

func TestTargetRow_TargetField_Extra4(t *testing.T) {
row := TargetRow{Target: "mymodule"}
if row.Target != "mymodule" {
t.Errorf("expected 'mymodule', got %q", row.Target)
}
}

func TestTargetRow_AllFieldsSet_Extra4(t *testing.T) {
row := TargetRow{
Target:    "n",
Source:    "s",
DeployDir: "d",
Status:    "active",
}
if row.Target != "n" || row.Source != "s" || row.DeployDir != "d" {
t.Error("field mismatch")
}
if row.Status != "active" {
t.Error("expected Status='active'")
}
}

func TestTargetRow_JSONMarshal_Extra4(t *testing.T) {
row := TargetRow{Target: "mod", Status: "active"}
data, err := json.Marshal(row)
if err != nil {
t.Fatalf("marshal error: %v", err)
}
if len(data) == 0 {
t.Error("expected non-empty JSON")
}
}

func TestTargetRow_JSONUnmarshal_Extra4(t *testing.T) {
payload := `{"target":"x","status":"inactive"}`
var row TargetRow
if err := json.Unmarshal([]byte(payload), &row); err != nil {
t.Fatalf("unmarshal error: %v", err)
}
if row.Target != "x" {
t.Errorf("expected target 'x', got %q", row.Target)
}
}

func TestTargetRow_NeedsDefault_Extra4(t *testing.T) {
var row TargetRow
if row.Needs != "" {
t.Error("zero Needs should be empty string")
}
}

func TestTargetRow_StatusDefault_Extra4(t *testing.T) {
var row TargetRow
if row.Status != "" {
t.Error("zero Status should be empty")
}
}

func TestTargetRow_DeployDirDefault_Extra4(t *testing.T) {
var row TargetRow
if row.DeployDir != "" {
t.Error("zero DeployDir should be empty")
}
}

func TestTargetRow_SliceFields_Extra4(t *testing.T) {
rows := []TargetRow{{Target: "a"}, {Target: "b"}}
if len(rows) != 2 {
t.Error("expected 2 rows")
}
}

func TestTargetRow_NeedsField_Extra4(t *testing.T) {
row := TargetRow{Needs: "pkg1,pkg2"}
if row.Needs != "pkg1,pkg2" {
t.Errorf("expected 'pkg1,pkg2', got %q", row.Needs)
}
}
