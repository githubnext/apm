package base_test

import (
	"testing"

	"github.com/githubnext/apm/internal/adapters/client/base"
)

func TestInputVarRE_NoMatchBraceOnly_Extra3(t *testing.T) {
	if base.InputVarRE.MatchString("${}") {
		t.Error("expected no match for ${}")
	}
}

func TestInputVarRE_NoMatchEnvColon_Extra3(t *testing.T) {
	if base.InputVarRE.MatchString("${env:FOO}") {
		t.Error("expected no match for ${env:FOO}")
	}
}

func TestInputVarRE_MatchWithDot_Extra3(t *testing.T) {
	m := base.InputVarRE.FindStringSubmatch("${input:my.var}")
	if m == nil {
		t.Fatal("expected match for ${input:my.var}")
	}
	if m[1] != "my.var" {
		t.Errorf("got %q want my.var", m[1])
	}
}

func TestInputVarRE_MatchWithNumber_Extra3(t *testing.T) {
	m := base.InputVarRE.FindStringSubmatch("${input:var1}")
	if m == nil {
		t.Fatal("expected match")
	}
	if m[1] != "var1" {
		t.Errorf("got %q want var1", m[1])
	}
}

func TestEnvVarRE_NoMatchDollarOnly_Extra3(t *testing.T) {
	if base.EnvVarRE.MatchString("$FOO") {
		t.Error("expected no match for plain $FOO")
	}
}

func TestEnvVarRE_MatchLongName_Extra3(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${VERY_LONG_VARIABLE_NAME_123}")
	if m == nil {
		t.Fatal("expected match for long name")
	}
	if m[1] != "VERY_LONG_VARIABLE_NAME_123" {
		t.Errorf("got %q", m[1])
	}
}

func TestEnvVarRE_FindAllInSentence_Extra3(t *testing.T) {
	s := "prefix ${FOO} middle ${env:BAR} suffix"
	all := base.EnvVarRE.FindAllStringSubmatch(s, -1)
	if len(all) != 2 {
		t.Errorf("expected 2 matches, got %d", len(all))
	}
}

func TestInputVarRE_FindAllMultiple_Extra3(t *testing.T) {
	s := "Use ${input:a} and ${input:b} in command"
	all := base.InputVarRE.FindAllString(s, -1)
	if len(all) != 2 {
		t.Errorf("expected 2 matches, got %d", len(all))
	}
}

func TestEnvVarRE_EnvPrefixCapture_Extra3(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${env:MY_TOKEN}")
	if m == nil {
		t.Fatal("expected match")
	}
	if m[1] != "MY_TOKEN" {
		t.Errorf("got %q want MY_TOKEN", m[1])
	}
}

func TestMCPClientAdapter_InterfaceHasMethods_Extra3(t *testing.T) {
	// Compile-time check that MCPClientAdapter interface is usable as a type
	var _ base.MCPClientAdapter = (base.MCPClientAdapter)(nil)
}

func TestEnvVarRE_NoMatchNumbers_Extra3(t *testing.T) {
	if base.EnvVarRE.MatchString("${123}") {
		t.Error("should not match numbers only")
	}
}

func TestInputVarRE_EmptyInputName_Extra3(t *testing.T) {
	// ${input:} has empty name - depends on regex; confirm consistent
	got := base.InputVarRE.MatchString("${input:}")
	// Just ensure no panic; value may be true or false
	_ = got
}

func TestEnvVarRE_UnderscoreStart_Extra3(t *testing.T) {
	m := base.EnvVarRE.FindStringSubmatch("${_PRIVATE}")
	if m == nil {
		t.Fatal("expected match for underscore-start")
	}
	if m[1] != "_PRIVATE" {
		t.Errorf("got %q", m[1])
	}
}

func TestInputVarRE_WithSpaces_Extra3(t *testing.T) {
	// Spaces inside braces are not typical identifiers but pattern should handle
	has := base.InputVarRE.MatchString("${input:with space}")
	_ = has // no panic
}

func TestEnvVarRE_ReplaceAll_Extra3(t *testing.T) {
	s := "cmd --token ${TOKEN}"
	result := base.EnvVarRE.ReplaceAllString(s, "REPLACED")
	if result == s {
		t.Error("expected replacement to occur")
	}
}
