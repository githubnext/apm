package mcpcommand

import (
	"testing"
)

func TestParseEnvPair_valid(t *testing.T) {
	k, v, ok := ParseEnvPair("FOO=bar")
	if !ok || k != "FOO" || v != "bar" {
		t.Errorf("expected FOO=bar, got %s=%s ok=%v", k, v, ok)
	}
}

func TestParseEnvPair_emptyValue(t *testing.T) {
	k, v, ok := ParseEnvPair("FOO=")
	if !ok || k != "FOO" || v != "" {
		t.Errorf("empty value: expected ok, got k=%s v=%q ok=%v", k, v, ok)
	}
}

func TestParseEnvPair_noEquals(t *testing.T) {
	_, _, ok := ParseEnvPair("NOEQUALSSIGN")
	if ok {
		t.Error("expected false for pair without =")
	}
}

func TestParseEnvPair_valueWithEquals(t *testing.T) {
	k, v, ok := ParseEnvPair("URL=http://host?a=1&b=2")
	if !ok || k != "URL" || v != "http://host?a=1&b=2" {
		t.Errorf("expected URL=http://host?a=1&b=2, got %s=%s ok=%v", k, v, ok)
	}
}

func TestParseEnvPairs_multiple(t *testing.T) {
	result := ParseEnvPairs([]string{"A=1", "B=2", "C=three"})
	if result["A"] != "1" || result["B"] != "2" || result["C"] != "three" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestParseEnvPairs_skipsInvalid(t *testing.T) {
	result := ParseEnvPairs([]string{"VALID=ok", "badformat", "X=y"})
	if len(result) != 2 {
		t.Errorf("expected 2 valid pairs, got %d", len(result))
	}
}

func TestParseEnvPairs_empty(t *testing.T) {
	result := ParseEnvPairs(nil)
	if len(result) != 0 {
		t.Errorf("expected empty map for nil input")
	}
}

func TestParseHeaderPair_colonSpace(t *testing.T) {
	k, v, ok := ParseHeaderPair("Authorization: Bearer token123")
	if !ok || k != "Authorization" || v != "Bearer token123" {
		t.Errorf("expected Authorization: Bearer token123, got %s: %s ok=%v", k, v, ok)
	}
}

func TestParseHeaderPair_equals(t *testing.T) {
	k, v, ok := ParseHeaderPair("X-Custom=value")
	if !ok || k != "X-Custom" || v != "value" {
		t.Errorf("expected X-Custom=value, got %s=%s ok=%v", k, v, ok)
	}
}

func TestParseHeaderPair_invalid(t *testing.T) {
	_, _, ok := ParseHeaderPair("nodelimiter")
	if ok {
		t.Error("expected false for header without delimiter")
	}
}

func TestParseHeaderPairs_multiple(t *testing.T) {
	result := ParseHeaderPairs([]string{"Content-Type: application/json", "Accept: text/plain"})
	if result["Content-Type"] != "application/json" {
		t.Errorf("unexpected Content-Type: %v", result["Content-Type"])
	}
	if result["Accept"] != "text/plain" {
		t.Errorf("unexpected Accept: %v", result["Accept"])
	}
}

func TestTransportDefault_stdio(t *testing.T) {
	got := TransportDefault("", []string{"node", "server.js"}, "")
	if got != "stdio" {
		t.Errorf("expected stdio, got %s", got)
	}
}

func TestTransportDefault_http(t *testing.T) {
	got := TransportDefault("http://localhost:3000/mcp", nil, "")
	if got != "http" {
		t.Errorf("expected http, got %s", got)
	}
}

func TestTransportDefault_explicit(t *testing.T) {
	got := TransportDefault("http://x", []string{"cmd"}, "sse")
	if got != "sse" {
		t.Errorf("expected explicit sse, got %s", got)
	}
}

func TestTransportDefault_empty(t *testing.T) {
	got := TransportDefault("", nil, "")
	if got != "" {
		t.Errorf("expected empty transport, got %s", got)
	}
}

func TestMCPInstallRequest_fields(t *testing.T) {
	req := MCPInstallRequest{
		MCPName:   "my-server",
		Transport: "stdio",
		Verbose:   true,
	}
	if req.MCPName != "my-server" {
		t.Errorf("unexpected MCPName: %s", req.MCPName)
	}
	if !req.Verbose {
		t.Error("expected Verbose=true")
	}
}

func TestMCPInstallResult_fields(t *testing.T) {
	result := MCPInstallResult{
		Outcome:    "added",
		EntryKey:   "my-server",
		Integrated: true,
	}
	if result.Outcome != "added" {
		t.Errorf("unexpected Outcome: %s", result.Outcome)
	}
	if !result.Integrated {
		t.Error("expected Integrated=true")
	}
}
