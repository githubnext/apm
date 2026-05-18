package schema

import (
	"testing"
)

func TestApmPolicy_FetchFailure(t *testing.T) {
	p := ApmPolicy{FetchFailure: "block"}
	if p.FetchFailure != "block" {
		t.Errorf("FetchFailure = %q, want block", p.FetchFailure)
	}
}

func TestApmPolicy_Remote(t *testing.T) {
	p := ApmPolicy{Remote: "https://policy.example.com/policy.yml"}
	if p.Remote != "https://policy.example.com/policy.yml" {
		t.Errorf("Remote = %q", p.Remote)
	}
}

func TestApmPolicy_Version(t *testing.T) {
	p := ApmPolicy{Version: "1"}
	if p.Version != "1" {
		t.Errorf("Version = %q, want 1", p.Version)
	}
}

func TestApmPolicy_Cache(t *testing.T) {
	p := ApmPolicy{Cache: PolicyCache{TTL: 7200}}
	if p.Cache.TTL != 7200 {
		t.Errorf("Cache.TTL = %d, want 7200", p.Cache.TTL)
	}
}

func TestApmPolicy_FullConstruct(t *testing.T) {
	p := ApmPolicy{
		Version:     "2",
		Remote:      "https://policy.corp/apm-policy.yml",
		Cache:       PolicyCache{TTL: 3600},
		Enforcement: "block",
		FetchFailure: "warn",
		Deps: DependencyPolicy{
			Allow:             []string{"*"},
			Deny:              []string{"bad-pkg"},
			RequireResolution: "policy-wins",
			MaxDepth:          10,
		},
		MCP: McpPolicy{
			Allow:       []string{"github/*"},
			SelfDefined: "deny",
		},
	}
	if p.Version != "2" {
		t.Errorf("Version = %q", p.Version)
	}
	if p.Deps.MaxDepth != 10 {
		t.Errorf("Deps.MaxDepth = %d", p.Deps.MaxDepth)
	}
	if p.MCP.SelfDefined != "deny" {
		t.Errorf("MCP.SelfDefined = %q", p.MCP.SelfDefined)
	}
}

func TestDependencyPolicy_AllowDenyRequire(t *testing.T) {
	p := DependencyPolicy{
		Allow:   []string{"safe-pkg", "another-pkg"},
		Deny:    []string{"unsafe-pkg"},
		Require: []string{"required-pkg"},
	}
	if len(p.Allow) != 2 {
		t.Errorf("Allow len = %d, want 2", len(p.Allow))
	}
	if p.Allow[0] != "safe-pkg" {
		t.Errorf("Allow[0] = %q", p.Allow[0])
	}
	if len(p.Deny) != 1 || p.Deny[0] != "unsafe-pkg" {
		t.Errorf("Deny = %v", p.Deny)
	}
	if len(p.Require) != 1 || p.Require[0] != "required-pkg" {
		t.Errorf("Require = %v", p.Require)
	}
}

func TestDependencyPolicy_RequireResolutionVariants(t *testing.T) {
	for _, resolution := range []string{"project-wins", "policy-wins", "block"} {
		p := DependencyPolicy{RequireResolution: resolution}
		if p.RequireResolution != resolution {
			t.Errorf("RequireResolution = %q, want %q", p.RequireResolution, resolution)
		}
	}
}

func TestMcpPolicy_SelfDefinedVariants(t *testing.T) {
	for _, v := range []string{"deny", "warn", "allow"} {
		p := McpPolicy{SelfDefined: v}
		if p.SelfDefined != v {
			t.Errorf("SelfDefined = %q, want %q", p.SelfDefined, v)
		}
	}
}

func TestCompilationPolicy_TargetAllow(t *testing.T) {
	p := CompilationPolicy{
		Targets: CompilationTargetPolicy{Allow: []string{"vscode", "claude", "all"}},
	}
	if len(p.Targets.Allow) != 3 {
		t.Errorf("want 3 targets, got %d", len(p.Targets.Allow))
	}
}

func TestCompilationPolicy_EnforceVariants(t *testing.T) {
	for _, enforce := range []string{"distributed", "single-file"} {
		p := CompilationPolicy{Strategy: CompilationStrategyPolicy{Enforce: enforce}}
		if p.Strategy.Enforce != enforce {
			t.Errorf("Enforce = %q, want %q", p.Strategy.Enforce, enforce)
		}
	}
}

func TestApmPolicy_EnforcementOff(t *testing.T) {
	p := ApmPolicy{Enforcement: "off"}
	if p.Enforcement != "off" {
		t.Errorf("Enforcement = %q, want off", p.Enforcement)
	}
}

func TestDefaultDependencyPolicy_AllowNil(t *testing.T) {
	p := DefaultDependencyPolicy()
	if p.Allow != nil {
		t.Errorf("expected nil Allow slice, got %v", p.Allow)
	}
	if p.Deny != nil {
		t.Errorf("expected nil Deny slice, got %v", p.Deny)
	}
}

func TestMcpTransportPolicy_MultipleAllows(t *testing.T) {
	tp := McpTransportPolicy{Allow: []string{"stdio", "sse", "http", "streamable-http"}}
	if len(tp.Allow) != 4 {
		t.Errorf("want 4 transports, got %d", len(tp.Allow))
	}
}
