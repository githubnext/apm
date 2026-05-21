package schema

import "testing"

func TestApmPolicy_ZeroValue(t *testing.T) {
	p := ApmPolicy{}
	if p.Version != "" || p.Remote != "" || p.Enforcement != "" {
		t.Error("zero value should have empty strings")
	}
}

func TestDefaultDependencyPolicy_MaxDepth(t *testing.T) {
	dp := DefaultDependencyPolicy()
	if dp.MaxDepth != 50 {
		t.Errorf("expected MaxDepth=50, got %d", dp.MaxDepth)
	}
}

func TestDefaultDependencyPolicy_RequireResolution(t *testing.T) {
	dp := DefaultDependencyPolicy()
	if dp.RequireResolution != "project-wins" {
		t.Errorf("expected project-wins, got %q", dp.RequireResolution)
	}
}

func TestDefaultDependencyPolicy_NilSlices(t *testing.T) {
	dp := DefaultDependencyPolicy()
	if len(dp.Allow) != 0 || len(dp.Deny) != 0 || len(dp.Require) != 0 {
		t.Error("default slices should be empty")
	}
}

func TestPolicyCache_ZeroValue(t *testing.T) {
	pc := PolicyCache{}
	if pc.TTL != 0 {
		t.Errorf("zero value TTL should be 0, got %d", pc.TTL)
	}
}

func TestPolicyCache_SetTTL(t *testing.T) {
	pc := PolicyCache{TTL: 3600}
	if pc.TTL != 3600 {
		t.Errorf("expected 3600, got %d", pc.TTL)
	}
}

func TestMcpPolicy_TrustTransitive(t *testing.T) {
	p := McpPolicy{TrustTransitive: true}
	if !p.TrustTransitive {
		t.Error("expected TrustTransitive=true")
	}
}

func TestCompilationTargetPolicy_AllowList(t *testing.T) {
	p := CompilationTargetPolicy{Allow: []string{"vscode", "claude"}, Enforce: "block"}
	if len(p.Allow) != 2 {
		t.Errorf("expected 2 allow entries, got %d", len(p.Allow))
	}
	if p.Enforce != "block" {
		t.Errorf("expected block, got %q", p.Enforce)
	}
}

func TestCompilationStrategyPolicy_EnforceVariants(t *testing.T) {
	for _, v := range []string{"distributed", "single-file", ""} {
		p := CompilationStrategyPolicy{Enforce: v}
		if p.Enforce != v {
			t.Errorf("expected %q, got %q", v, p.Enforce)
		}
	}
}

func TestApmPolicy_EnforcementBlock(t *testing.T) {
	p := ApmPolicy{Enforcement: "block"}
	if p.Enforcement != "block" {
		t.Errorf("expected block, got %q", p.Enforcement)
	}
}
