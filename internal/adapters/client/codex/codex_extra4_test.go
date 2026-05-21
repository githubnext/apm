package codex

import (
	"strings"
	"testing"
)

func TestMCPServersKey_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.MCPServersKey() != "mcp_servers" {
		t.Errorf("expected mcp_servers, got %s", a.MCPServersKey())
	}
}

func TestSupportsUserScope_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if !a.SupportsUserScope() {
		t.Error("expected SupportsUserScope true")
	}
}

func TestGetConfigPath_Project_Extra4(t *testing.T) {
	a := New("/myproject", false)
	p := a.GetConfigPath()
	if !strings.Contains(p, ".codex") {
		t.Errorf("expected .codex in path, got %s", p)
	}
	if !strings.HasSuffix(p, "config.toml") {
		t.Errorf("expected config.toml suffix, got %s", p)
	}
}

func TestGetConfigPath_User_Extra4(t *testing.T) {
	a := New("", true)
	p := a.GetConfigPath()
	if !strings.HasSuffix(p, "config.toml") {
		t.Errorf("expected config.toml suffix, got %s", p)
	}
	if !strings.Contains(p, ".codex") {
		t.Errorf("expected .codex in path, got %s", p)
	}
}

func TestGetConfigPath_NotEmpty_Extra4(t *testing.T) {
	for _, root := range []string{"/a", "/b/c", ""} {
		a := New(root, false)
		p := a.GetConfigPath()
		if p == "" {
			t.Errorf("expected non-empty path for root=%q", root)
		}
	}
}

func TestTargetName_IsCodex_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.TargetName() != "codex" {
		t.Errorf("expected codex, got %s", a.TargetName())
	}
}

func TestNew_RuntimeEnvSubstitutionFalse_Extra4(t *testing.T) {
	a := New("/tmp", false)
	if a.Adapter.SupportsRuntimeEnvSubstitution {
		t.Error("expected SupportsRuntimeEnvSubstitution=false for codex")
	}
}

func TestGetCurrentConfig_MissingFile_Extra4(t *testing.T) {
	a := New("/nonexistent/path/xyz", false)
	cfg := a.GetCurrentConfig()
	if cfg == nil {
		t.Error("expected non-nil map for missing file")
	}
}

func TestNormalizeProjectArg_Dollar_Extra4(t *testing.T) {
	cases := []struct{ in, out string }{
		{"$PROJECT", "."},
		{"${PROJECT}", "."},
		{"other", "other"},
		{"", ""},
		{"$PROJECT_X", "$PROJECT_X"},
	}
	for _, tc := range cases {
		got := normalizeProjectArg(tc.in)
		if got != tc.out {
			t.Errorf("normalizeProjectArg(%q) = %q, want %q", tc.in, got, tc.out)
		}
	}
}

func TestEnsureDockerEnvFlags_AddsMissing_Extra4(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	out := ensureDockerEnvFlags([]string{"run", "image"}, env)
	found := false
	for i, v := range out {
		if v == "-e" && i+1 < len(out) && out[i+1] == "FOO=bar" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected -e FOO=bar in %v", out)
	}
}

func TestEnsureDockerEnvFlags_Nodup_Extra4(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	args := []string{"-e", "FOO=existing"}
	out := ensureDockerEnvFlags(args, env)
	count := 0
	for _, v := range out {
		if v == "FOO=existing" || v == "FOO=bar" {
			count++
		}
	}
	if count > 1 {
		t.Errorf("duplicate -e flags: %v", out)
	}
}

func TestServerKeyFor_NameTakesPrecedence_Extra4(t *testing.T) {
	k := serverKeyFor("owner/pkg", "myname")
	if k != "myname" {
		t.Errorf("expected myname, got %s", k)
	}
}

func TestServerKeyFor_SlashURL_Extra4(t *testing.T) {
	k := serverKeyFor("owner/pkg", "")
	if k != "pkg" {
		t.Errorf("expected pkg, got %s", k)
	}
}

func TestServerKeyFor_NoSlash_Extra4(t *testing.T) {
	k := serverKeyFor("simplepkg", "")
	if k != "simplepkg" {
		t.Errorf("expected simplepkg, got %s", k)
	}
}

func TestFilterOut_RemovesTarget_Extra4(t *testing.T) {
	out := filterOut([]string{"a", "-y", "b", "-y"}, "-y")
	for _, v := range out {
		if v == "-y" {
			t.Error("filterOut should have removed -y")
		}
	}
	if len(out) != 2 {
		t.Errorf("expected 2 items, got %d", len(out))
	}
}

func TestCond_Preferred_Extra4(t *testing.T) {
	if cond("npx", "fallback") != "npx" {
		t.Error("expected npx")
	}
}

func TestCond_Fallback_Extra4(t *testing.T) {
	if cond("", "fallback") != "fallback" {
		t.Error("expected fallback")
	}
}

func TestInferRegistryName_NPM_Extra4(t *testing.T) {
	pkg := map[string]interface{}{"registry": "npm"}
	if inferRegistryName(pkg) != "npm" {
		t.Error("expected npm")
	}
}

func TestInferRegistryName_Docker_Extra4(t *testing.T) {
	pkg := map[string]interface{}{"registry": "docker-hub"}
	if inferRegistryName(pkg) != "docker" {
		t.Error("expected docker")
	}
}

func TestInferRegistryName_Default_Extra4(t *testing.T) {
	pkg := map[string]interface{}{}
	if inferRegistryName(pkg) != "npm" {
		t.Error("expected default npm")
	}
}

func TestToTOMLString_Escaping_Extra4(t *testing.T) {
	s := toTOMLString(`foo"bar\baz`)
	if !strings.Contains(s, `\"`) {
		t.Errorf("expected escaped quote in %s", s)
	}
	if !strings.Contains(s, `\\`) {
		t.Errorf("expected escaped backslash in %s", s)
	}
}

func TestParseSimpleTOML_Empty_Extra4(t *testing.T) {
	m, err := parseSimpleTOML([]byte{})
	if err != nil {
		t.Error(err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestParseSimpleTOML_JSON_Extra4(t *testing.T) {
	m, err := parseSimpleTOML([]byte(`{"key":"val"}`))
	if err != nil {
		t.Error(err)
	}
	if m["key"] != "val" {
		t.Errorf("expected val, got %v", m["key"])
	}
}
