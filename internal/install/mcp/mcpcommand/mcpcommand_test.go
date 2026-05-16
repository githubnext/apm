package mcpcommand

import (
	"testing"
)

func TestParseEnvPair_Valid(t *testing.T) {
	k, v, ok := ParseEnvPair("FOO=bar")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if k != "FOO" {
		t.Errorf("expected key 'FOO', got %q", k)
	}
	if v != "bar" {
		t.Errorf("expected value 'bar', got %q", v)
	}
}

func TestParseEnvPair_ValueWithEquals(t *testing.T) {
	k, v, ok := ParseEnvPair("URL=http://host?a=b&c=d")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if k != "URL" {
		t.Errorf("expected key 'URL', got %q", k)
	}
	if v != "http://host?a=b&c=d" {
		t.Errorf("unexpected value: %q", v)
	}
}

func TestParseEnvPair_NoEquals(t *testing.T) {
	_, _, ok := ParseEnvPair("NOEQUALS")
	if ok {
		t.Error("expected ok=false for pair without '='")
	}
}

func TestParseEnvPair_Empty(t *testing.T) {
	_, _, ok := ParseEnvPair("")
	if ok {
		t.Error("expected ok=false for empty pair")
	}
}

func TestParseEnvPairs_Multiple(t *testing.T) {
	pairs := []string{"A=1", "B=2", "C=three"}
	got := ParseEnvPairs(pairs)
	if got["A"] != "1" || got["B"] != "2" || got["C"] != "three" {
		t.Errorf("unexpected pairs map: %v", got)
	}
}

func TestParseEnvPairs_Empty(t *testing.T) {
	got := ParseEnvPairs(nil)
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestParseHeaderPair_Valid(t *testing.T) {
	k, v, ok := ParseHeaderPair("Authorization=Bearer token123")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if k != "Authorization" {
		t.Errorf("expected key 'Authorization', got %q", k)
	}
	if v != "Bearer token123" {
		t.Errorf("expected 'Bearer token123', got %q", v)
	}
}

func TestParseHeaderPair_NoEquals(t *testing.T) {
	_, _, ok := ParseHeaderPair("NoHeader")
	if ok {
		t.Error("expected ok=false")
	}
}

func TestParseHeaderPairs_Multiple(t *testing.T) {
	pairs := []string{"X-Token=abc", "Accept=application/json"}
	got := ParseHeaderPairs(pairs)
	if got["X-Token"] != "abc" {
		t.Errorf("unexpected map: %v", got)
	}
}

func TestTransportDefault_SSE(t *testing.T) {
	// TransportDefault returns "http" for URL-only input (no special sse detection)
	result := TransportDefault("http://localhost/sse", nil, "")
	if result != "http" {
		t.Errorf("expected 'http' for URL-only input, got %q", result)
	}
}

func TestTransportDefault_ExplicitTransport(t *testing.T) {
	result := TransportDefault("http://host", nil, "stdio")
	if result != "stdio" {
		t.Errorf("expected 'stdio' when explicit, got %q", result)
	}
}

func TestTransportDefault_CommandArgv(t *testing.T) {
	result := TransportDefault("", []string{"npx", "server"}, "")
	if result == "" {
		t.Error("expected a non-empty default transport for argv")
	}
}
