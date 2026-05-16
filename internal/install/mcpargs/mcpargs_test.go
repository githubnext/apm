package mcpargs_test

import (
"testing"

"github.com/githubnext/apm/internal/install/mcpargs"
)

func TestParseKVPairs_valid(t *testing.T) {
pairs := []string{"KEY=value", "FOO=bar=baz", "EMPTY="}
got, err := mcpargs.ParseKVPairs(pairs, "--test")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if got["KEY"] != "value" {
t.Errorf("KEY: got %q, want %q", got["KEY"], "value")
}
if got["FOO"] != "bar=baz" {
t.Errorf("FOO: got %q, want %q", got["FOO"], "bar=baz")
}
if got["EMPTY"] != "" {
t.Errorf("EMPTY: got %q, want %q", got["EMPTY"], "")
}
}

func TestParseKVPairs_missingEquals(t *testing.T) {
_, err := mcpargs.ParseKVPairs([]string{"NOEQUALS"}, "--test")
if err == nil {
t.Fatal("expected error for missing '='")
}
}

func TestParseKVPairs_emptyKey(t *testing.T) {
_, err := mcpargs.ParseKVPairs([]string{"=value"}, "--test")
if err == nil {
t.Fatal("expected error for empty key")
}
}

func TestParseEnvPairs(t *testing.T) {
pairs := []string{"HOME=/root", "PATH=/usr/bin"}
got, err := mcpargs.ParseEnvPairs(pairs)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if got["HOME"] != "/root" {
t.Errorf("HOME: got %q, want /root", got["HOME"])
}
}

func TestParseHeaderPairs(t *testing.T) {
pairs := []string{"Authorization=Bearer token123"}
got, err := mcpargs.ParseHeaderPairs(pairs)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if got["Authorization"] != "Bearer token123" {
t.Errorf("Authorization: got %q, want %q", got["Authorization"], "Bearer token123")
}
}

func TestParseKVPairs_empty(t *testing.T) {
got, err := mcpargs.ParseKVPairs(nil, "--test")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 0 {
t.Errorf("expected empty map, got %v", got)
}
}
