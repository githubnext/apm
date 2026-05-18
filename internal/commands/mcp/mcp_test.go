package mcp

import "testing"

func TestTruncate(t *testing.T) {
	tests := []struct {
		s    string
		n    int
		want string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "he..."},
		{"", 5, ""},
		{"abcdef", 6, "abcdef"},
		{"abcdefg", 6, "abc..."},
	}
	for _, tc := range tests {
		got := truncate(tc.s, tc.n)
		if got != tc.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tc.s, tc.n, got, tc.want)
		}
	}
}

func TestSearchOptions(t *testing.T) {
	opts := SearchOptions{
		Query:  "github",
		Format: "text",
		Limit:  10,
	}
	if opts.Query != "github" {
		t.Errorf("unexpected Query %q", opts.Query)
	}
	if opts.Limit != 10 {
		t.Errorf("unexpected Limit %d", opts.Limit)
	}
}

func TestInstallOptions(t *testing.T) {
	opts := InstallOptions{
		ServerRef:   "github/models",
		ProjectRoot: "/tmp/proj",
		UserScope:   true,
	}
	if opts.ServerRef != "github/models" {
		t.Errorf("unexpected ServerRef %q", opts.ServerRef)
	}
	if !opts.UserScope {
		t.Error("expected UserScope true")
	}
}

func TestMCPRegistryEnv(t *testing.T) {
	if MCPRegistryEnv != "MCP_REGISTRY_URL" {
		t.Errorf("MCPRegistryEnv = %q, want %q", MCPRegistryEnv, "MCP_REGISTRY_URL")
	}
}

func TestInfoOptions_Fields(t *testing.T) {
opts := InfoOptions{
ServerRef:   "github/copilot",
RegistryURL: "https://registry.example.com",
Format:      "json",
}
if opts.ServerRef != "github/copilot" {
t.Errorf("unexpected ServerRef %q", opts.ServerRef)
}
if opts.Format != "json" {
t.Errorf("unexpected Format %q", opts.Format)
}
}

func TestInstallOptions_ForceFlag(t *testing.T) {
opts := InstallOptions{
ServerRef: "github/models",
Force:     true,
}
if !opts.Force {
t.Error("expected Force true")
}
}

func TestInstallOptions_RuntimeField(t *testing.T) {
opts := InstallOptions{
ServerRef: "server-ref",
Runtime:   "node",
}
if opts.Runtime != "node" {
t.Errorf("unexpected Runtime %q", opts.Runtime)
}
}

func TestTruncate_ExactLength(t *testing.T) {
got := truncate("abc", 3)
if got != "abc" {
t.Errorf("truncate at exact length: want %q, got %q", "abc", got)
}
}

func TestTruncate_SmallN(t *testing.T) {
// n >= 3 is the minimum meaningful value (3 chars for "...")
got := truncate("hello", 3)
if got != "..." {
t.Errorf("truncate to 3: want ellipsis, got %q", got)
}
}

func TestSearchOptions_DefaultLimit(t *testing.T) {
opts := SearchOptions{Query: "test"}
if opts.Limit != 0 {
t.Errorf("default Limit should be 0, got %d", opts.Limit)
}
}
