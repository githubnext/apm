package mcpwriter

import (
	"testing"
)

func TestDiffEntry_ModifiedValue(t *testing.T) {
	old := map[string]interface{}{"name": "srv", "command": "old-cmd"}
	new := map[string]interface{}{"name": "srv", "command": "new-cmd"}
	lines := DiffEntry(old, new)
	if len(lines) == 0 {
		t.Fatal("expected diff lines for modified entry")
	}
	found := false
	for _, l := range lines {
		if l.Key == "command" {
			found = true
			if l.OldValue != "old-cmd" {
				t.Errorf("OldValue: got %v want old-cmd", l.OldValue)
			}
			if l.NewValue != "new-cmd" {
				t.Errorf("NewValue: got %v want new-cmd", l.NewValue)
			}
		}
	}
	if !found {
		t.Error("expected a diff line for 'command'")
	}
}

func TestFindExistingMCPEntry_MultipleEntries(t *testing.T) {
	entries := []interface{}{
		map[string]interface{}{"name": "alpha"},
		map[string]interface{}{"name": "beta"},
		map[string]interface{}{"name": "gamma"},
	}
	if idx := FindExistingMCPEntry(entries, "alpha"); idx != 0 {
		t.Errorf("expected 0, got %d", idx)
	}
	if idx := FindExistingMCPEntry(entries, "gamma"); idx != 2 {
		t.Errorf("expected 2, got %d", idx)
	}
}

func TestMCPListSection_DevTrue(t *testing.T) {
	data := &ApmYMLData{
		DevDependencies: map[string]interface{}{
			"mcp": []interface{}{map[string]interface{}{"name": "dev-srv"}},
		},
	}
	result := MCPListSection(data, true)
	if len(result) == 0 {
		t.Fatal("expected non-empty mcp list for dev=true")
	}
}

func TestMCPListSection_ProdDeps(t *testing.T) {
	data := &ApmYMLData{
		Dependencies: map[string]interface{}{
			"mcp": []interface{}{
				map[string]interface{}{"name": "prod-srv"},
			},
		},
	}
	result := MCPListSection(data, false)
	if len(result) == 0 {
		t.Fatal("expected non-empty mcp list for dev=false with prod deps")
	}
}

func TestOutcomeConstants_Distinct(t *testing.T) {
	if OutcomeAdded == OutcomeReplaced {
		t.Error("OutcomeAdded and OutcomeReplaced must be distinct")
	}
	if OutcomeAdded == OutcomeSkipped {
		t.Error("OutcomeAdded and OutcomeSkipped must be distinct")
	}
	if OutcomeReplaced == OutcomeSkipped {
		t.Error("OutcomeReplaced and OutcomeSkipped must be distinct")
	}
}

func TestDiffEntry_NoChange(t *testing.T) {
entry := map[string]interface{}{"name": "srv", "command": "cmd"}
lines := DiffEntry(entry, entry)
if len(lines) != 0 {
t.Errorf("expected no diff lines for identical entries, got %d", len(lines))
}
}

func TestDiffEntry_NewKeyAdded(t *testing.T) {
old := map[string]interface{}{"name": "srv"}
new := map[string]interface{}{"name": "srv", "args": []string{"--flag"}}
lines := DiffEntry(old, new)
if len(lines) == 0 {
t.Fatal("expected diff lines when new key is added")
}
found := false
for _, l := range lines {
if l.Key == "args" {
found = true
}
}
if !found {
t.Error("expected diff line for 'args'")
}
}

func TestDiffEntry_StringEntry(t *testing.T) {
// string entries are treated as {name: value}
lines := DiffEntry("old-name", "new-name")
if len(lines) == 0 {
t.Fatal("expected diff for string entries")
}
}

func TestFindExistingMCPEntry_SingleMissing(t *testing.T) {
	entries := []interface{}{
		map[string]interface{}{"name": "alpha"},
	}
	if idx := FindExistingMCPEntry(entries, "missing"); idx != -1 {
		t.Errorf("expected -1 for missing entry, got %d", idx)
	}
}

func TestFindExistingMCPEntry_NilList(t *testing.T) {
	if idx := FindExistingMCPEntry(nil, "any"); idx != -1 {
		t.Errorf("expected -1 for nil list, got %d", idx)
	}
}

func TestFindExistingMCPEntry_StringEntry(t *testing.T) {
entries := []interface{}{"alpha", "beta"}
if idx := FindExistingMCPEntry(entries, "beta"); idx != 1 {
t.Errorf("expected 1 for string entry 'beta', got %d", idx)
}
}

func TestMCPListSection_NilSection(t *testing.T) {
data := &ApmYMLData{}
result := MCPListSection(data, false)
if result != nil {
t.Error("expected nil for missing deps section")
}
}

func TestMCPListSection_NoMCPKey(t *testing.T) {
data := &ApmYMLData{
Dependencies: map[string]interface{}{"other": "value"},
}
result := MCPListSection(data, false)
if result != nil {
t.Error("expected nil when no mcp key in deps")
}
}
