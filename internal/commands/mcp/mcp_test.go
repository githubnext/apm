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
