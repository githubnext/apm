// Package mcpwriter persists MCP entries into apm.yml.
// Mirrors src/apm_cli/install/mcp/writer.py.
package mcpwriter

import (
	"fmt"
	"os"
)

// AddOutcome describes what add_mcp_to_apm_yml did with the entry.
type AddOutcome int

const (
	OutcomeAdded    AddOutcome = iota
	OutcomeReplaced AddOutcome = iota
	OutcomeSkipped  AddOutcome = iota
)

// DiffLine is one human-readable "key: old -> new" change line.
type DiffLine struct {
	Key      string
	OldValue interface{}
	NewValue interface{}
}

// DiffEntry computes the diff between two MCP entries for display.
// old and new are the raw map or string representations.
func DiffEntry(old, new interface{}) []DiffLine {
	oldMap := entryToMap(old)
	newMap := entryToMap(new)

	// Collect keys in order: old keys first, then new-only keys.
	seen := map[string]bool{}
	var keys []string
	for k := range oldMap {
		keys = append(keys, k)
		seen[k] = true
	}
	for k := range newMap {
		if !seen[k] {
			keys = append(keys, k)
		}
	}

	var diffs []DiffLine
	for _, k := range keys {
		ov := oldMap[k]
		nv := newMap[k]
		if fmt.Sprintf("%v", ov) != fmt.Sprintf("%v", nv) {
			diffs = append(diffs, DiffLine{Key: k, OldValue: ov, NewValue: nv})
		}
	}
	return diffs
}

func entryToMap(v interface{}) map[string]interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		return t
	case string:
		return map[string]interface{}{"name": t}
	default:
		return map[string]interface{}{}
	}
}

// ApmYMLData is the minimal representation of apm.yml for MCP writer operations.
type ApmYMLData struct {
	Dependencies    map[string]interface{}
	DevDependencies map[string]interface{}
}

// MCPListSection returns the mcp list from the appropriate section.
func MCPListSection(data *ApmYMLData, dev bool) []interface{} {
	var section map[string]interface{}
	if dev {
		section = data.DevDependencies
	} else {
		section = data.Dependencies
	}
	if section == nil {
		return nil
	}
	mcpRaw, ok := section["mcp"]
	if !ok {
		return nil
	}
	if mcpList, ok := mcpRaw.([]interface{}); ok {
		return mcpList
	}
	return nil
}

// FindExistingMCPEntry returns the index of an MCP entry with the given name,
// or -1 if not found.
func FindExistingMCPEntry(entries []interface{}, name string) int {
	for i, e := range entries {
		switch t := e.(type) {
		case string:
			if t == name {
				return i
			}
		case map[string]interface{}:
			if n, ok := t["name"].(string); ok && n == name {
				return i
			}
		}
	}
	return -1
}

// IsInteractiveTTY returns true when stdout is a TTY (interactive session).
func IsInteractiveTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
