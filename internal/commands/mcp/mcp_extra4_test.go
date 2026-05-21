package mcp

import "testing"

func TestSearchOptions_AllFields_Extra4(t *testing.T) {
	opts := SearchOptions{
		Query:       "git tools",
		RegistryURL: "https://registry.example.com",
		Format:      "json",
		Limit:       50,
	}
	if opts.Query != "git tools" {
		t.Errorf("Query = %q", opts.Query)
	}
	if opts.RegistryURL != "https://registry.example.com" {
		t.Errorf("RegistryURL = %q", opts.RegistryURL)
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q", opts.Format)
	}
	if opts.Limit != 50 {
		t.Errorf("Limit = %d", opts.Limit)
	}
}

func TestSearchOptions_ZeroValue_Extra4(t *testing.T) {
	var opts SearchOptions
	if opts.Query != "" {
		t.Errorf("zero Query = %q", opts.Query)
	}
	if opts.Limit != 0 {
		t.Errorf("zero Limit = %d", opts.Limit)
	}
}

func TestInstallOptions_AllFields_Extra4(t *testing.T) {
	opts := InstallOptions{
		ServerRef:   "org/server@v1",
		ProjectRoot: "/project",
		Runtime:     "node",
		UserScope:   true,
		Force:       false,
	}
	if opts.ServerRef != "org/server@v1" {
		t.Errorf("ServerRef = %q", opts.ServerRef)
	}
	if opts.ProjectRoot != "/project" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if opts.Runtime != "node" {
		t.Errorf("Runtime = %q", opts.Runtime)
	}
	if !opts.UserScope {
		t.Error("UserScope should be true")
	}
	if opts.Force {
		t.Error("Force should be false")
	}
}

func TestInstallOptions_ZeroValue_Extra4(t *testing.T) {
	var opts InstallOptions
	if opts.ServerRef != "" {
		t.Errorf("zero ServerRef = %q", opts.ServerRef)
	}
	if opts.Force {
		t.Error("zero Force should be false")
	}
}

func TestInfoOptions_Fields_Extra4(t *testing.T) {
	opts := InfoOptions{
		ServerRef:   "tools/server",
		RegistryURL: "https://reg.example.com",
		Format:      "text",
	}
	if opts.ServerRef != "tools/server" {
		t.Errorf("ServerRef = %q", opts.ServerRef)
	}
	if opts.Format != "text" {
		t.Errorf("Format = %q", opts.Format)
	}
}

func TestInfoOptions_ZeroValue_Extra4(t *testing.T) {
	var opts InfoOptions
	if opts.ServerRef != "" {
		t.Errorf("zero ServerRef = %q", opts.ServerRef)
	}
	if opts.Format != "" {
		t.Errorf("zero Format = %q", opts.Format)
	}
}

func TestMCPRegistryEnv_Extra4(t *testing.T) {
	if MCPRegistryEnv == "" {
		t.Error("MCPRegistryEnv should be non-empty")
	}
}
