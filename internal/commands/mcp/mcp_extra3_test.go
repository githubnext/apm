package mcp

import "testing"

func TestInstallOptions_ZeroValue(t *testing.T) {
	var opts InstallOptions
	if opts.ServerRef != "" {
		t.Errorf("ServerRef should be empty, got %q", opts.ServerRef)
	}
	if opts.Force {
		t.Error("Force should default to false")
	}
	if opts.UserScope {
		t.Error("UserScope should default to false")
	}
}

func TestInfoOptions_ZeroValue(t *testing.T) {
	var opts InfoOptions
	if opts.ServerRef != "" {
		t.Errorf("ServerRef should be empty, got %q", opts.ServerRef)
	}
	if opts.Format != "" {
		t.Errorf("Format should be empty, got %q", opts.Format)
	}
}

func TestSearchOptions_ZeroValue(t *testing.T) {
	var opts SearchOptions
	if opts.Query != "" {
		t.Errorf("Query should be empty, got %q", opts.Query)
	}
	if opts.Limit != 0 {
		t.Errorf("Limit should be 0, got %d", opts.Limit)
	}
}

func TestMCPRegistryEnv_NotEmpty(t *testing.T) {
	if MCPRegistryEnv == "" {
		t.Error("MCPRegistryEnv constant should not be empty")
	}
}

func TestMCPRegistryEnv_IsString(t *testing.T) {
	var _ string = MCPRegistryEnv
}

func TestSearchOptions_AssignFields(t *testing.T) {
	opts := SearchOptions{
		Query:       "test-query",
		RegistryURL: "https://registry.example.com",
		Format:      "json",
		Limit:       50,
	}
	if opts.Query != "test-query" {
		t.Errorf("Query = %q, want %q", opts.Query, "test-query")
	}
	if opts.Limit != 50 {
		t.Errorf("Limit = %d, want 50", opts.Limit)
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q, want json", opts.Format)
	}
}

func TestInstallOptions_AssignFields(t *testing.T) {
	opts := InstallOptions{
		ServerRef:   "my-server@1.0",
		ProjectRoot: "/tmp/project",
		Runtime:     "node",
		UserScope:   true,
		Force:       true,
	}
	if !opts.UserScope {
		t.Error("UserScope should be true")
	}
	if !opts.Force {
		t.Error("Force should be true")
	}
	if opts.Runtime != "node" {
		t.Errorf("Runtime = %q, want node", opts.Runtime)
	}
}

func TestTruncate_LongString_Extra3(t *testing.T) {
	s := "abcdefghij"
	got := truncate(s, 5)
	if len(got) > 5 {
		t.Errorf("truncate should shorten string, got %q (len=%d)", got, len(got))
	}
}

func TestTruncate_ExactLength_Extra3(t *testing.T) {
	s := "hello!"
	got := truncate(s, len(s))
	if got != s {
		t.Errorf("truncate at exact length: got %q, want %q", got, s)
	}
}
