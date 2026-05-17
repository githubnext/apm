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

func TestParseKVPairs_multipleEquals(t *testing.T) {
pairs := []string{"URL=https://example.com/path?a=1&b=2"}
got, err := mcpargs.ParseKVPairs(pairs, "--test")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if got["URL"] != "https://example.com/path?a=1&b=2" {
t.Errorf("URL: got %q", got["URL"])
}
}

func TestParseKVPairs_duplicateKey(t *testing.T) {
pairs := []string{"KEY=first", "KEY=second"}
got, err := mcpargs.ParseKVPairs(pairs, "--test")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if got["KEY"] != "second" {
t.Errorf("expected last value wins, got %q", got["KEY"])
}
}

func TestParseEnvPairs_emptyInput(t *testing.T) {
got, err := mcpargs.ParseEnvPairs(nil)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 0 {
t.Errorf("expected empty, got %v", got)
}
}

func TestParseEnvPairs_multipleVars(t *testing.T) {
pairs := []string{"HOME=/root", "PATH=/usr/bin:/usr/local/bin", "EMPTY="}
got, err := mcpargs.ParseEnvPairs(pairs)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if got["HOME"] != "/root" {
t.Errorf("HOME: got %q", got["HOME"])
}
if got["PATH"] != "/usr/bin:/usr/local/bin" {
t.Errorf("PATH: got %q", got["PATH"])
}
if got["EMPTY"] != "" {
t.Errorf("EMPTY: got %q", got["EMPTY"])
}
}

func TestParseHeaderPairs_multipleHeaders(t *testing.T) {
pairs := []string{"Authorization=Bearer tok", "X-Custom=value=with=equals"}
got, err := mcpargs.ParseHeaderPairs(pairs)
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if got["Authorization"] != "Bearer tok" {
t.Errorf("Authorization: got %q", got["Authorization"])
}
if got["X-Custom"] != "value=with=equals" {
t.Errorf("X-Custom: got %q", got["X-Custom"])
}
}

func TestParseKVPairs_flagNameInError(t *testing.T) {
_, err := mcpargs.ParseKVPairs([]string{"noequals"}, "--env")
if err == nil {
t.Fatal("expected error")
}
// Just ensure the error is non-empty
if err.Error() == "" {
t.Error("expected non-empty error message")
}
}
