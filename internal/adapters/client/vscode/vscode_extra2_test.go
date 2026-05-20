package vscode

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTranslateEnvValueForVSCode_AngleVarUnderscore(t *testing.T) {
	got := translateEnvValueForVSCode("MY_TOKEN")
	if got != "MY_TOKEN" {
		t.Errorf("plain string unchanged: %q", got)
	}
}

func TestTranslateEnvValueForVSCode_InputSubstitution(t *testing.T) {
	got := translateEnvValueForVSCode("${input:API_KEY}")
	if got != "${input:API_KEY}" {
		t.Errorf("input ref should be unchanged: %q", got)
	}
}

func TestTranslateEnvValueForVSCode_LegacyAngleBracketVar(t *testing.T) {
	got := translateEnvValueForVSCode("<MY_SECRET>")
	if !strings.Contains(got, "MY_SECRET") {
		t.Errorf("angle bracket var should expand to contain var name: %q", got)
	}
}

func TestTranslateEnvValueForVSCode_EmptyString(t *testing.T) {
	got := translateEnvValueForVSCode("")
	if got != "" {
		t.Errorf("empty string should remain empty: %q", got)
	}
}

func TestExtractPackageArgs_WithPackageAndArgs(t *testing.T) {
	pkg := map[string]interface{}{
		"package": "my-pkg@1.0.0",
		"args":    []interface{}{"--flag", "value"},
	}
	args := extractPackageArgs(pkg)
	_ = args // result may be empty or contain args; just ensure no panic
}

func TestExtractPackageArgs_EmptyPackage(t *testing.T) {
	pkg := map[string]interface{}{}
	args := extractPackageArgs(pkg)
	_ = args // just ensure no panic for empty input
}

func TestToStringSlice_StringValue(t *testing.T) {
	got := toStringSlice("single-string")
	if len(got) != 0 {
		t.Errorf("single string should not be coerced to slice: %v", got)
	}
}

func TestToStringSlice_MixedTypes(t *testing.T) {
	input := []interface{}{"a", 42, "b"}
	got := toStringSlice(input)
	count := 0
	for _, s := range got {
		if s != "" {
			count++
		}
	}
	if count == 0 {
		t.Error("expected some string elements")
	}
}

func TestStrField_IntValue(t *testing.T) {
	m := map[string]interface{}{"key": 42}
	got := strField(m, "key")
	if got != "" {
		t.Errorf("int value should return empty string: %q", got)
	}
}

func TestStrField_NilMap(t *testing.T) {
	got := strField(nil, "key")
	if got != "" {
		t.Errorf("nil map should return empty string: %q", got)
	}
}

func TestToInterfaceSlice_Empty(t *testing.T) {
	got := toInterfaceSlice(nil)
	if got == nil {
		t.Error("should return non-nil for nil input")
	}
}

func TestToInterfaceSlice_SingleElement(t *testing.T) {
	got := toInterfaceSlice([]string{"only"})
	if len(got) != 1 {
		t.Errorf("expected 1 element, got %d", len(got))
	}
	if got[0] != "only" {
		t.Errorf("element: %v", got[0])
	}
}

func TestToSliceOfMaps_ValidInput(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{"key": "val"},
	}
	got := toSliceOfMaps(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 map, got %d", len(got))
	}
	if got[0]["key"] != "val" {
		t.Errorf("key: %v", got[0]["key"])
	}
}

func TestUpdateConfig_WritesServersSection(t *testing.T) {
	dir := t.TempDir()
	a := New(dir, false)
	err := a.UpdateConfig(map[string]interface{}{
		"servers": map[string]interface{}{
			"my-srv": map[string]interface{}{"command": "npx"},
		},
	})
	if err != nil {
		t.Fatalf("UpdateConfig error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, ".vscode", "mcp.json"))
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	servers, ok := cfg["servers"]
	if !ok {
		t.Error("servers key missing from written config")
	}
	sm, ok := servers.(map[string]interface{})
	if !ok || sm["my-srv"] == nil {
		t.Error("my-srv missing from servers")
	}
}

func TestGetConfigPath_HasVSCodeInPath(t *testing.T) {
	a := New("/my/project", false)
	p := a.GetConfigPath()
	if !strings.Contains(p, ".vscode") {
		t.Errorf("expected .vscode in path: %q", p)
	}
	if !strings.HasSuffix(p, "mcp.json") {
		t.Errorf("expected mcp.json suffix: %q", p)
	}
}

func TestGetCurrentConfig_InvalidJSON_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	vsDir := filepath.Join(dir, ".vscode")
	if err := os.MkdirAll(vsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vsDir, "mcp.json"), []byte("{bad json}"), 0o644); err != nil {
		t.Fatal(err)
	}
	a := New(dir, false)
	cfg := a.GetCurrentConfig()
	if len(cfg) != 0 {
		t.Errorf("expected empty map for invalid JSON, got %v", cfg)
	}
}

func TestFilterOut_MultipleMatches(t *testing.T) {
	got := filterOut([]string{"a", "b", "a", "c"}, "a")
	for _, s := range got {
		if s == "a" {
			t.Error("filterOut should remove all occurrences of target")
		}
	}
	if len(got) != 2 {
		t.Errorf("expected 2 elements, got %d", len(got))
	}
}

func TestFilterOut_NoMatch(t *testing.T) {
	got := filterOut([]string{"x", "y"}, "z")
	if len(got) != 2 {
		t.Errorf("no-match filter should return original length: %d", len(got))
	}
}
