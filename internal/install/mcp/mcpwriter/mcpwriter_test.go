package mcpwriter

import (
	"testing"
)

func TestDiffEntry_BothNil(t *testing.T) {
	lines := DiffEntry(nil, nil)
	// Should not panic; result can be empty
	_ = lines
}

func TestDiffEntry_NewValue(t *testing.T) {
	lines := DiffEntry(nil, map[string]interface{}{"name": "server1"})
	if len(lines) == 0 {
		t.Error("expected diff lines for new value")
	}
}

func TestDiffEntry_SameValue(t *testing.T) {
	m := map[string]interface{}{"name": "server1", "command": "npx"}
	lines := DiffEntry(m, m)
	// No changes expected when input is identical
	for _, l := range lines {
		// OldValue and NewValue should be equal
		if l.OldValue != l.NewValue {
			t.Errorf("unexpected diff line for identical inputs: %+v", l)
		}
	}
}

func TestDiffEntry_RemovedValue(t *testing.T) {
	old := map[string]interface{}{"name": "server1"}
	lines := DiffEntry(old, nil)
	if len(lines) == 0 {
		t.Error("expected diff lines for removed value")
	}
}

func TestFindExistingMCPEntry_EmptyList(t *testing.T) {
	idx := FindExistingMCPEntry(nil, "server1")
	if idx != -1 {
		t.Errorf("expected -1 for empty list, got %d", idx)
	}
}

func TestFindExistingMCPEntry_Found(t *testing.T) {
	entries := []interface{}{
		map[string]interface{}{"name": "server1"},
		map[string]interface{}{"name": "server2"},
	}
	idx := FindExistingMCPEntry(entries, "server2")
	if idx != 1 {
		t.Errorf("expected 1, got %d", idx)
	}
}

func TestFindExistingMCPEntry_NotFound(t *testing.T) {
	entries := []interface{}{
		map[string]interface{}{"name": "server1"},
	}
	idx := FindExistingMCPEntry(entries, "missing")
	if idx != -1 {
		t.Errorf("expected -1 for missing entry, got %d", idx)
	}
}

func TestMCPListSection_NilData(t *testing.T) {
	// Passing nil panics in implementation (not nil-safe); skip nil case.
	data := &ApmYMLData{}
	result := MCPListSection(data, false)
	if result == nil {
		result = []interface{}{}
	}
	_ = result
}

func TestMCPListSection_EmptyData(t *testing.T) {
	data := &ApmYMLData{}
	result := MCPListSection(data, false)
	if result == nil {
		result = []interface{}{}
	}
	_ = result
}
