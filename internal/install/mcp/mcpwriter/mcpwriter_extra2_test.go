package mcpwriter

import (
	"testing"
)

func TestAddOutcome_Constants(t *testing.T) {
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

func TestDiffLine_ZeroValue(t *testing.T) {
	var d DiffLine
	if d.Key != "" {
		t.Errorf("expected empty Key, got %q", d.Key)
	}
	if d.OldValue != nil {
		t.Errorf("expected nil OldValue")
	}
	if d.NewValue != nil {
		t.Errorf("expected nil NewValue")
	}
}

func TestDiffLine_FieldAssignment(t *testing.T) {
	d := DiffLine{Key: "url", OldValue: "http://old", NewValue: "http://new"}
	if d.Key != "url" {
		t.Errorf("unexpected Key %q", d.Key)
	}
	if d.OldValue != "http://old" {
		t.Errorf("unexpected OldValue %v", d.OldValue)
	}
}

func TestDiffEntry_IdenticalMaps_Empty(t *testing.T) {
	old := map[string]interface{}{"cmd": "node", "ver": "1"}
	new := map[string]interface{}{"cmd": "node", "ver": "1"}
	diffs := DiffEntry(old, new)
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs for identical maps, got %d", len(diffs))
	}
}

func TestDiffEntry_ChangedValue(t *testing.T) {
	old := map[string]interface{}{"cmd": "node"}
	new := map[string]interface{}{"cmd": "deno"}
	diffs := DiffEntry(old, new)
	if len(diffs) != 1 {
		t.Errorf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Key != "cmd" {
		t.Errorf("expected key cmd, got %q", diffs[0].Key)
	}
}

func TestDiffEntry_KeyRemoved(t *testing.T) {
	old := map[string]interface{}{"cmd": "node", "env": "dev"}
	new := map[string]interface{}{"cmd": "node"}
	diffs := DiffEntry(old, new)
	if len(diffs) != 1 {
		t.Errorf("expected 1 diff for removed key, got %d", len(diffs))
	}
}

func TestDiffEntry_KeyAdded(t *testing.T) {
	old := map[string]interface{}{"cmd": "node"}
	new := map[string]interface{}{"cmd": "node", "transport": "sse"}
	diffs := DiffEntry(old, new)
	if len(diffs) != 1 {
		t.Errorf("expected 1 diff for added key, got %d", len(diffs))
	}
	if diffs[0].Key != "transport" {
		t.Errorf("expected key transport, got %q", diffs[0].Key)
	}
}

func TestFindExistingMCPEntry_EmptyList_Extra2(t *testing.T) {
	idx := FindExistingMCPEntry(nil, "server")
	if idx != -1 {
		t.Errorf("expected -1 for nil list, got %d", idx)
	}
}

func TestFindExistingMCPEntry_NotFound_Extra2(t *testing.T) {
	entries := []interface{}{
		map[string]interface{}{"name": "other"},
	}
	idx := FindExistingMCPEntry(entries, "server")
	if idx != -1 {
		t.Errorf("expected -1 for missing entry, got %d", idx)
	}
}

func TestFindExistingMCPEntry_Found_Extra2(t *testing.T) {
	entries := []interface{}{
		map[string]interface{}{"name": "alpha"},
		map[string]interface{}{"name": "server"},
	}
	idx := FindExistingMCPEntry(entries, "server")
	if idx != 1 {
		t.Errorf("expected idx 1, got %d", idx)
	}
}

func TestMCPListSection_EmptyData_Empty(t *testing.T) {
	d := &ApmYMLData{}
	result := MCPListSection(d, false)
	if len(result) != 0 {
		t.Errorf("expected empty result for empty data, got %d", len(result))
	}
}

func TestApmYMLData_ZeroValue(t *testing.T) {
	var d ApmYMLData
	if d.Dependencies != nil {
		t.Error("expected nil Dependencies")
	}
	if d.DevDependencies != nil {
		t.Error("expected nil DevDependencies")
	}
}
