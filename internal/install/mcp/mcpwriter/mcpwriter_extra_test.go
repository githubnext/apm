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
