package mcpregistry_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/mcp/mcpregistry"
)

func TestRedactURLCredentials_NoCredentials(t *testing.T) {
	u := "https://registry.example.com/api"
	got := mcpregistry.RedactURLCredentials(u)
	if got != u {
		t.Errorf("URL without credentials should be unchanged: got %q", got)
	}
}

func TestRedactURLCredentials_WithPassword(t *testing.T) {
	u := "https://user:secret@registry.example.com/api"
	got := mcpregistry.RedactURLCredentials(u)
	if strings.Contains(got, "secret") {
		t.Errorf("expected password redacted, got %q", got)
	}
	if !strings.Contains(got, "registry.example.com") {
		t.Errorf("expected host preserved, got %q", got)
	}
}

func TestValidateRegistryURL_ValidHTTPS(t *testing.T) {
	url, warn, err := mcpregistry.ValidateRegistryURL("https://registry.example.com/v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url == "" {
		t.Error("expected non-empty normalized URL")
	}
	_ = warn
}

func TestValidateRegistryURL_HTTPScheme_NoError(t *testing.T) {
	_, _, err := mcpregistry.ValidateRegistryURL("http://internal.corp.com/registry")
	// http is allowed (though may warn); not an error
	if err != nil {
		t.Fatalf("http should be allowed: %v", err)
	}
}

func TestValidateRegistryURL_EmptyURL_Error(t *testing.T) {
	_, _, err := mcpregistry.ValidateRegistryURL("")
	if err == nil {
		t.Error("expected error for empty URL")
	}
}

func TestValidateRegistryURL_Localhost_Warning(t *testing.T) {
	_, warn, err := mcpregistry.ValidateRegistryURL("http://localhost:8080/reg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if warn == "" {
		t.Error("expected warning for localhost URL")
	}
}

func TestValidateRegistryURL_IPv6Loopback_Warning(t *testing.T) {
	_, warn, err := mcpregistry.ValidateRegistryURL("http://[::1]:8080/reg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if warn == "" {
		t.Error("expected warning for IPv6 loopback")
	}
}

func TestRegistryEnvOverride_EmptyURL(t *testing.T) {
	env, ok := mcpregistry.RegistryEnvOverride("")
	if ok {
		t.Error("empty URL should not produce override")
	}
	if len(env) != 0 {
		t.Errorf("expected empty env, got %v", env)
	}
}

func TestRegistryEnvOverride_ValidURL(t *testing.T) {
	// RegistryEnvOverride only returns ok=true for http (allowHTTP), not https.
	env, _ := mcpregistry.RegistryEnvOverride("https://mcp.example.com/v1")
	if len(env) == 0 {
		t.Error("expected non-empty env map for valid URL")
	}
	// MCP_REGISTRY_URL should always be set for non-empty URL.
	if env["MCP_REGISTRY_URL"] != "https://mcp.example.com/v1" {
		t.Errorf("MCP_REGISTRY_URL = %q", env["MCP_REGISTRY_URL"])
	}
}

func TestResolveRegistryURL_EnvFallback(t *testing.T) {
	got := mcpregistry.ResolveRegistryURL("", "https://env.example.com")
	if got != "https://env.example.com" {
		t.Errorf("expected env fallback, got %q", got)
	}
}

func TestResolveRegistryURL_FlagTakesPrecedence_extra2(t *testing.T) {
	got := mcpregistry.ResolveRegistryURL("https://flag.example.com", "https://env.example.com")
	if got != "https://flag.example.com" {
		t.Errorf("expected flag to take precedence, got %q", got)
	}
}
