package mcpcommand

import (
	"testing"
)

func TestParseEnvPair_EmptyVal(t *testing.T) {
	k, v, ok := ParseEnvPair("KEY=")
	if !ok || k != "KEY" || v != "" {
		t.Errorf("unexpected result: k=%q v=%q ok=%v", k, v, ok)
	}
}

func TestParseEnvPair_ValHasEquals(t *testing.T) {
	k, v, ok := ParseEnvPair("KEY=a=b=c")
	if !ok || k != "KEY" || v != "a=b=c" {
		t.Errorf("unexpected result: k=%q v=%q ok=%v", k, v, ok)
	}
}

func TestParseEnvPair_MissingEquals(t *testing.T) {
	_, _, ok := ParseEnvPair("NOEQ")
	if ok {
		t.Error("expected ok=false for pair without '='")
	}
}

func TestParseEnvPairs_SkipsInvalidPairs(t *testing.T) {
	m := ParseEnvPairs([]string{"A=1", "INVALID", "B=2"})
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected map: %v", m)
	}
	if len(m) != 2 {
		t.Errorf("expected 2 entries, got %d", len(m))
	}
}

func TestParseEnvPairs_NilInput(t *testing.T) {
	m := ParseEnvPairs(nil)
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestParseHeaderPair_ColonSpaceSep(t *testing.T) {
	k, v, ok := ParseHeaderPair("Content-Type: application/json")
	if !ok || k != "Content-Type" || v != "application/json" {
		t.Errorf("unexpected result: k=%q v=%q ok=%v", k, v, ok)
	}
}

func TestParseHeaderPair_EqualsSep(t *testing.T) {
	k, v, ok := ParseHeaderPair("X-Token=abc123")
	if !ok || k != "X-Token" || v != "abc123" {
		t.Errorf("unexpected result: k=%q v=%q ok=%v", k, v, ok)
	}
}

func TestParseHeaderPair_NeitherSep(t *testing.T) {
	_, _, ok := ParseHeaderPair("invalidddd")
	if ok {
		t.Error("expected ok=false for header without separator")
	}
}

func TestParseHeaderPairs_MixedSeps(t *testing.T) {
	m := ParseHeaderPairs([]string{"A: 1", "B=2"})
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected header map: %v", m)
	}
}

func TestMCPInstallRequest_ZeroFields(t *testing.T) {
	var r MCPInstallRequest
	if r.MCPName != "" || r.Transport != "" || r.Dev || r.Force {
		t.Errorf("unexpected zero-value MCPInstallRequest: %+v", r)
	}
}

func TestMCPInstallResult_ZeroFields(t *testing.T) {
	var r MCPInstallResult
	if r.Outcome != "" || r.EntryKey != "" || r.Integrated {
		t.Errorf("unexpected zero-value MCPInstallResult: %+v", r)
	}
}

func TestMCPInstallResult_OutcomeStrings(t *testing.T) {
	for _, outcome := range []string{"added", "replaced", "skipped"} {
		r := MCPInstallResult{Outcome: outcome, EntryKey: "key1", Integrated: true}
		if r.Outcome != outcome {
			t.Errorf("expected outcome %q, got %q", outcome, r.Outcome)
		}
	}
}
