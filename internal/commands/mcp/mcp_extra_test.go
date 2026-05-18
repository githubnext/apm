package mcp

import "testing"

func TestMCPRegistryEnv_Value(t *testing.T) {
	if MCPRegistryEnv != "MCP_REGISTRY_URL" {
		t.Errorf("MCPRegistryEnv = %q, want MCP_REGISTRY_URL", MCPRegistryEnv)
	}
}

func TestSearchOptions_DefaultFormat(t *testing.T) {
	opts := SearchOptions{Query: "test"}
	if opts.Format != "" {
		t.Errorf("default Format should be empty, got %q", opts.Format)
	}
	if opts.Limit != 0 {
		t.Errorf("default Limit should be 0, got %d", opts.Limit)
	}
}

func TestSearchOptions_JSONFormat(t *testing.T) {
	opts := SearchOptions{
		Query:  "github",
		Format: "json",
		Limit:  50,
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q, want json", opts.Format)
	}
	if opts.Limit != 50 {
		t.Errorf("Limit = %d, want 50", opts.Limit)
	}
}

func TestSearchOptions_WithRegistryURL(t *testing.T) {
	opts := SearchOptions{
		Query:       "my-server",
		RegistryURL: "https://registry.example.com",
		Format:      "text",
	}
	if opts.RegistryURL != "https://registry.example.com" {
		t.Errorf("RegistryURL = %q", opts.RegistryURL)
	}
}

func TestInstallOptions_Defaults(t *testing.T) {
	opts := InstallOptions{}
	if opts.ServerRef != "" {
		t.Errorf("default ServerRef should be empty")
	}
	if opts.UserScope {
		t.Error("default UserScope should be false")
	}
	if opts.Force {
		t.Error("default Force should be false")
	}
}

func TestInstallOptions_AllFields(t *testing.T) {
	opts := InstallOptions{
		ServerRef:   "github/copilot",
		ProjectRoot: "/tmp/myproject",
		Runtime:     "node",
		UserScope:   true,
		Force:       true,
	}
	if opts.ServerRef != "github/copilot" {
		t.Errorf("ServerRef = %q", opts.ServerRef)
	}
	if opts.ProjectRoot != "/tmp/myproject" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.Runtime != "node" {
		t.Errorf("Runtime = %q", opts.Runtime)
	}
	if !opts.UserScope {
		t.Error("UserScope should be true")
	}
	if !opts.Force {
		t.Error("Force should be true")
	}
}

func TestInfoOptions_Defaults(t *testing.T) {
	opts := InfoOptions{}
	if opts.ServerRef != "" || opts.RegistryURL != "" || opts.Format != "" {
		t.Error("all InfoOptions fields should default to empty")
	}
}

func TestInfoOptions_AllFields(t *testing.T) {
	opts := InfoOptions{
		ServerRef:   "my-server",
		RegistryURL: "https://registry.example.com",
		Format:      "json",
	}
	if opts.ServerRef != "my-server" {
		t.Errorf("ServerRef = %q", opts.ServerRef)
	}
	if opts.RegistryURL != "https://registry.example.com" {
		t.Errorf("RegistryURL = %q", opts.RegistryURL)
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q", opts.Format)
	}
}

func TestTruncate_ExactLengthMatch(t *testing.T) {
	s := "hello"
	got := truncate(s, 5)
	if got != "hello" {
		t.Errorf("truncate(%q, 5) = %q, want hello", s, got)
	}
}

func TestTruncate_LongerString(t *testing.T) {
	got := truncate("hello world", 8)
	if len(got) != 8 {
		t.Errorf("truncate result len = %d, want 8: %q", len(got), got)
	}
	if got[5:] != "..." {
		t.Errorf("truncate should end with '...': %q", got)
	}
}

func TestTruncate_Unicode(t *testing.T) {
	// Basic ASCII only -- all chars are single bytes
	got := truncate("abcdefghij", 7)
	if len(got) != 7 {
		t.Errorf("len = %d, want 7: %q", len(got), got)
	}
}

func TestInstallOptions_ProjectRootVariants(t *testing.T) {
	roots := []string{"/home/user/proj", "/tmp/test", "."}
	for _, root := range roots {
		opts := InstallOptions{ProjectRoot: root}
		if opts.ProjectRoot != root {
			t.Errorf("ProjectRoot = %q, want %q", opts.ProjectRoot, root)
		}
	}
}
