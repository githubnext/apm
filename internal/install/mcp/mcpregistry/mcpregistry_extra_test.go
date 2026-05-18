package mcpregistry_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/mcp/mcpregistry"
)

func TestValidateRegistryURL_DecimalIPLoopback(t *testing.T) {
	// 2130706433 == 127.0.0.1 (decimal int form)
	_, warn, err := mcpregistry.ValidateRegistryURL("http://2130706433:8080/reg")
	if err != nil {
		t.Fatalf("unexpected error for decimal loopback: %v", err)
	}
	if warn == "" {
		t.Error("expected local host warning for decimal 127.0.0.1")
	}
}

func TestValidateRegistryURL_CloudMetadataIP(t *testing.T) {
	_, warn, err := mcpregistry.ValidateRegistryURL("http://169.254.169.254/latest/meta-data")
	if err != nil {
		t.Fatalf("unexpected error for cloud metadata: %v", err)
	}
	if warn == "" {
		t.Error("expected warning for cloud metadata IP")
	}
}

func TestValidateRegistryURL_RFC1918Private(t *testing.T) {
	// 10.0.0.1 is RFC1918 private
	_, warn, err := mcpregistry.ValidateRegistryURL("http://10.0.0.1/registry")
	if err != nil {
		t.Fatalf("unexpected error for RFC1918: %v", err)
	}
	if warn == "" {
		t.Error("expected local host warning for 10.0.0.1")
	}
}

func TestValidateRegistryURL_RFC1918_192(t *testing.T) {
	_, warn, err := mcpregistry.ValidateRegistryURL("http://192.168.1.50/registry")
	if err != nil {
		t.Fatalf("unexpected error for 192.168.x: %v", err)
	}
	if warn == "" {
		t.Error("expected local host warning for 192.168.1.50")
	}
}

func TestValidateRegistryURL_PublicIPNoWarning(t *testing.T) {
	_, warn, err := mcpregistry.ValidateRegistryURL("https://8.8.8.8/registry")
	if err != nil {
		t.Fatalf("unexpected error for public IP: %v", err)
	}
	if warn != "" {
		t.Errorf("unexpected warning for public IP: %q", warn)
	}
}

func TestValidateRegistryURL_NormalizedOutput(t *testing.T) {
	norm, _, err := mcpregistry.ValidateRegistryURL("https://registry.example.com/mcp?v=1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(norm, "registry.example.com") {
		t.Errorf("normalized URL missing host: %q", norm)
	}
}

func TestValidateRegistryURL_ExactlyMaxLength(t *testing.T) {
	// Exactly 2048 chars should pass
	url := "https://example.com/" + strings.Repeat("a", 2028)
	if len(url) != 2048 {
		t.Fatalf("test setup: len=%d, want 2048", len(url))
	}
	_, _, err := mcpregistry.ValidateRegistryURL(url)
	if err != nil {
		t.Errorf("URL of exactly 2048 chars should be valid, got: %v", err)
	}
}

func TestValidateRegistryURL_OnePastMaxLength(t *testing.T) {
	url := "https://example.com/" + strings.Repeat("a", 2029)
	if len(url) != 2049 {
		t.Fatalf("test setup: len=%d, want 2049", len(url))
	}
	_, _, err := mcpregistry.ValidateRegistryURL(url)
	if err == nil {
		t.Error("URL of 2049 chars should fail")
	}
}

func TestRedactURLCredentials_InvalidURL(t *testing.T) {
	bad := "://not-a-url"
	got := mcpregistry.RedactURLCredentials(bad)
	if got != bad {
		t.Errorf("invalid URL should be returned as-is, got %q", got)
	}
}

func TestRedactURLCredentials_UsernameOnly(t *testing.T) {
	u := "https://user@example.com/reg"
	got := mcpregistry.RedactURLCredentials(u)
	if strings.Contains(got, "user") {
		t.Errorf("username should be stripped: %q", got)
	}
}

func TestRegistryEnvOverride_HTTPSNoAllowHTTP(t *testing.T) {
	env, allow := mcpregistry.RegistryEnvOverride("https://registry.example.com")
	if allow {
		t.Error("https should not allow HTTP")
	}
	if _, ok := env["MCP_REGISTRY_ALLOW_HTTP"]; ok {
		t.Error("MCP_REGISTRY_ALLOW_HTTP should not be set for HTTPS")
	}
}

func TestResolveRegistryURL_BothEmpty(t *testing.T) {
	got := mcpregistry.ResolveRegistryURL("", "")
	if got != "" {
		t.Errorf("both empty should return empty, got %q", got)
	}
}

func TestResolveRegistryURL_OnlyFlag(t *testing.T) {
	got := mcpregistry.ResolveRegistryURL("https://flag.example.com", "")
	if got != "https://flag.example.com" {
		t.Errorf("got %q", got)
	}
}

func TestAllowedSchemes(t *testing.T) {
	if !mcpregistry.AllowedSchemes["https"] {
		t.Error("https should be allowed")
	}
	if !mcpregistry.AllowedSchemes["http"] {
		t.Error("http should be allowed")
	}
	if mcpregistry.AllowedSchemes["ftp"] {
		t.Error("ftp should not be allowed")
	}
}

func TestValidationError_Message(t *testing.T) {
	_, _, err := mcpregistry.ValidateRegistryURL("ftp://example.com")
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}
	if !strings.Contains(msg, "ftp") {
		t.Errorf("error should mention ftp scheme: %q", msg)
	}
}
