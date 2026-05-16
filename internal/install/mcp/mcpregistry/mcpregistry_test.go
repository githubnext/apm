package mcpregistry_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/mcp/mcpregistry"
)

func TestRedactURLCredentials_NoUser(t *testing.T) {
	u := "https://example.com/registry"
	if got := mcpregistry.RedactURLCredentials(u); got != u {
		t.Errorf("no-op expected, got %q", got)
	}
}

func TestRedactURLCredentials_WithUser(t *testing.T) {
	u := "https://user:pass@example.com/registry"
	got := mcpregistry.RedactURLCredentials(u)
	if strings.Contains(got, "pass") {
		t.Errorf("password not redacted: %q", got)
	}
	if !strings.Contains(got, "example.com") {
		t.Errorf("host missing: %q", got)
	}
}

func TestValidateRegistryURL_Valid(t *testing.T) {
	norm, warn, err := mcpregistry.ValidateRegistryURL("https://registry.example.com/mcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if warn != "" {
		t.Errorf("unexpected warning: %q", warn)
	}
	if !strings.HasPrefix(norm, "https://") {
		t.Errorf("unexpected normalized url: %q", norm)
	}
}

func TestValidateRegistryURL_HTTPAllowed(t *testing.T) {
	_, _, err := mcpregistry.ValidateRegistryURL("http://registry.example.com/mcp")
	if err != nil {
		t.Errorf("http should be allowed, got error: %v", err)
	}
}

func TestValidateRegistryURL_InvalidScheme(t *testing.T) {
	_, _, err := mcpregistry.ValidateRegistryURL("ftp://example.com/mcp")
	if err == nil {
		t.Error("expected error for ftp scheme")
	}
}

func TestValidateRegistryURL_NoHost(t *testing.T) {
	_, _, err := mcpregistry.ValidateRegistryURL("https:///path")
	if err == nil {
		t.Error("expected error for missing host")
	}
}

func TestValidateRegistryURL_TooLong(t *testing.T) {
	long := "https://example.com/" + strings.Repeat("a", 2048)
	_, _, err := mcpregistry.ValidateRegistryURL(long)
	if err == nil {
		t.Error("expected error for too-long URL")
	}
}

func TestValidateRegistryURL_LocalhostWarning(t *testing.T) {
	_, warn, err := mcpregistry.ValidateRegistryURL("http://localhost:8080/registry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if warn == "" {
		t.Error("expected local host warning")
	}
}

func TestResolveRegistryURL_FlagTakesPrecedence(t *testing.T) {
	got := mcpregistry.ResolveRegistryURL("https://flag.example.com", "https://env.example.com")
	if got != "https://flag.example.com" {
		t.Errorf("flag should take precedence, got %q", got)
	}
}

func TestResolveRegistryURL_FallbackToEnv(t *testing.T) {
	got := mcpregistry.ResolveRegistryURL("", "https://env.example.com")
	if got != "https://env.example.com" {
		t.Errorf("should fall back to env, got %q", got)
	}
}

func TestRegistryEnvOverride_Empty(t *testing.T) {
	env, allow := mcpregistry.RegistryEnvOverride("")
	if env != nil || allow {
		t.Errorf("empty URL should return nil,false; got %v,%v", env, allow)
	}
}

func TestRegistryEnvOverride_HTTPS(t *testing.T) {
	env, allow := mcpregistry.RegistryEnvOverride("https://registry.example.com")
	if env["MCP_REGISTRY_URL"] != "https://registry.example.com" {
		t.Errorf("unexpected env: %v", env)
	}
	if allow {
		t.Error("https should not set allowHTTP")
	}
}

func TestRegistryEnvOverride_HTTP(t *testing.T) {
	env, allow := mcpregistry.RegistryEnvOverride("http://localhost:9090")
	if !allow {
		t.Error("http URL should set allowHTTP=true")
	}
	if env["MCP_REGISTRY_ALLOW_HTTP"] != "1" {
		t.Errorf("MCP_REGISTRY_ALLOW_HTTP not set: %v", env)
	}
}
