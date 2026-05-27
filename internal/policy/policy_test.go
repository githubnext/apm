package policy_test

import (
	"testing"

	"github.com/githubnext/apm/internal/policy"
)

// TestParityCheckResult validates CheckResult structure mirrors Python models.CheckResult.
func TestParityCheckResult(t *testing.T) {
	cr := policy.CheckResult{
		Name:    "lockfile-exists",
		Passed:  true,
		Message: "Lockfile found",
		Details: []string{},
	}
	if cr.Name != "lockfile-exists" {
		t.Errorf("expected name lockfile-exists, got %s", cr.Name)
	}
	if !cr.Passed {
		t.Errorf("expected passed=true")
	}
}

// TestParityCIAuditResultPassed verifies all-pass aggregate.
func TestParityCIAuditResultPassed(t *testing.T) {
	r := &policy.CIAuditResult{
		Checks: []policy.CheckResult{
			{Name: "lockfile-exists", Passed: true, Message: "ok"},
			{Name: "ref-consistency", Passed: true, Message: "ok"},
		},
	}
	if !r.Passed() {
		t.Errorf("expected all checks to pass")
	}
	if len(r.FailedChecks()) != 0 {
		t.Errorf("expected 0 failed checks")
	}
}

// TestParityCIAuditResultFailed verifies failed checks aggregation.
func TestParityCIAuditResultFailed(t *testing.T) {
	r := &policy.CIAuditResult{
		Checks: []policy.CheckResult{
			{Name: "lockfile-exists", Passed: true, Message: "ok"},
			{Name: "ref-consistency", Passed: false, Message: "mismatch", Details: []string{"sha mismatch on pkg-a"}},
		},
	}
	if r.Passed() {
		t.Errorf("expected not all passed")
	}
	failed := r.FailedChecks()
	if len(failed) != 1 {
		t.Errorf("expected 1 failed check, got %d", len(failed))
	}
	if failed[0].Name != "ref-consistency" {
		t.Errorf("expected ref-consistency to fail")
	}
}

// TestParityCIAuditResultToJSON verifies JSON serialization mirrors Python.
func TestParityCIAuditResultToJSON(t *testing.T) {
	r := &policy.CIAuditResult{
		Checks: []policy.CheckResult{
			{Name: "lockfile-exists", Passed: true, Message: "ok", Details: []string{}},
			{Name: "ref-consistency", Passed: false, Message: "mismatch", Details: []string{"detail1"}},
		},
	}
	j := r.ToJSON()
	if j["passed"] != false {
		t.Errorf("expected passed=false")
	}
	summary, ok := j["summary"].(map[string]interface{})
	if !ok {
		t.Fatalf("summary not a map")
	}
	if summary["total"] != 2 {
		t.Errorf("expected total=2, got %v", summary["total"])
	}
	if summary["passed"] != 1 {
		t.Errorf("expected passed=1")
	}
	if summary["failed"] != 1 {
		t.Errorf("expected failed=1")
	}
}

// TestParityCheckArtifactMap validates the artifact mapping mirrors Python.
func TestParityCheckArtifactMap(t *testing.T) {
	cases := map[string]string{
		"lockfile-exists":        "apm.lock.yaml",
		"dependency-allowlist":   "apm.yml",
		"mcp-transport":          "apm.yml",
		"ref-consistency":        "apm.lock.yaml",
		"compilation-target":     "apm.yml",
		"required-manifest-fields": "apm.yml",
	}
	for check, want := range cases {
		got, ok := policy.CheckArtifactMap[check]
		if !ok {
			t.Errorf("check %s missing from artifact map", check)
			continue
		}
		if got != want {
			t.Errorf("check %s: want %s, got %s", check, want, got)
		}
	}
}

// TestParityDefaultDependencyPolicy mirrors Python DependencyPolicy defaults.
func TestParityDefaultDependencyPolicy(t *testing.T) {
	dp := policy.DefaultDependencyPolicy()
	if dp.MaxDepth != 50 {
		t.Errorf("expected MaxDepth=50, got %d", dp.MaxDepth)
	}
	if dp.RequireResolution != policy.RequireResolutionProjectWins {
		t.Errorf("expected project-wins resolution")
	}
}

// TestParityDependencyPolicyEffectiveDeny verifies nil -> empty semantics.
func TestParityDependencyPolicyEffectiveDeny(t *testing.T) {
	dp := policy.DependencyPolicy{Deny: nil}
	if len(dp.EffectiveDeny()) != 0 {
		t.Errorf("nil deny should return empty slice")
	}
	dp2 := policy.DependencyPolicy{Deny: []string{"bad-pkg"}}
	if len(dp2.EffectiveDeny()) != 1 {
		t.Errorf("expected 1 deny entry")
	}
}

// TestParityDependencyPolicyEffectiveRequire mirrors Python require semantics.
func TestParityDependencyPolicyEffectiveRequire(t *testing.T) {
	dp := policy.DependencyPolicy{Require: nil}
	if len(dp.EffectiveRequire()) != 0 {
		t.Errorf("nil require should return empty slice")
	}
	dp2 := policy.DependencyPolicy{Require: []string{"required-pkg"}}
	if len(dp2.EffectiveRequire()) != 1 {
		t.Errorf("expected 1 require entry")
	}
}

// TestParityDefaultPolicyCache verifies TTL default of 3600.
func TestParityDefaultPolicyCache(t *testing.T) {
	c := policy.DefaultPolicyCache()
	if c.TTL != 3600 {
		t.Errorf("expected TTL=3600, got %d", c.TTL)
	}
}

// TestParityDefaultOutcomeRouting verifies block-by-default routing.
func TestParityDefaultOutcomeRouting(t *testing.T) {
	o := policy.DefaultOutcomeRouting()
	if o.Default != policy.OutcomeBlock {
		t.Errorf("expected default outcome=block")
	}
}

// TestParityOutcomeRoutingActionFor verifies per-check routing overrides.
func TestParityOutcomeRoutingActionFor(t *testing.T) {
	o := policy.OutcomeRouting{
		Default: policy.OutcomeBlock,
		Checks: map[string]policy.OutcomeAction{
			"lockfile-exists": policy.OutcomeWarn,
		},
	}
	if o.ActionFor("lockfile-exists") != policy.OutcomeWarn {
		t.Errorf("expected warn for lockfile-exists")
	}
	if o.ActionFor("ref-consistency") != policy.OutcomeBlock {
		t.Errorf("expected block for unknown check")
	}
}

// TestParityDefaultPolicyDocument verifies all defaults are set.
func TestParityDefaultPolicyDocument(t *testing.T) {
	doc := policy.DefaultPolicyDocument()
	if doc.Cache.TTL != 3600 {
		t.Errorf("expected cache TTL=3600")
	}
	if doc.Dependencies.MaxDepth != 50 {
		t.Errorf("expected deps MaxDepth=50")
	}
	if doc.Outcomes.Default != policy.OutcomeBlock {
		t.Errorf("expected outcomes default=block")
	}
}

// TestParityMatcherMatchesAny verifies glob matching behavior.
func TestParityMatcherMatchesAny(t *testing.T) {
	cases := []struct {
		candidate string
		patterns  []string
		want      bool
	}{
		{"pkg-a", []string{"pkg-*"}, true},
		{"pkg-a", []string{"other-*"}, false},
		{"pkg-a", []string{"pkg-b", "pkg-a"}, true},
		{"pkg-a", []string{}, false},
		{"anything", []string{"*"}, true},
	}
	for _, tc := range cases {
		got := policy.MatchesAny(tc.candidate, tc.patterns)
		if got != tc.want {
			t.Errorf("MatchesAny(%q, %v) = %v, want %v", tc.candidate, tc.patterns, got, tc.want)
		}
	}
}

// TestParityMatcherAllowList verifies nil = allow-all semantics.
func TestParityMatcherAllowList(t *testing.T) {
	if !policy.MatchesAllowList("anything", nil) {
		t.Errorf("nil allowlist should permit everything")
	}
	if policy.MatchesAllowList("pkg-a", []string{}) {
		t.Errorf("empty allowlist should permit nothing")
	}
	if !policy.MatchesAllowList("pkg-a", []string{"pkg-*"}) {
		t.Errorf("pkg-a should match pkg-*")
	}
}

// TestParityMatcherDenyList verifies nil = deny-nothing semantics.
func TestParityMatcherDenyList(t *testing.T) {
	if policy.MatchesDenyList("pkg-a", nil) {
		t.Errorf("nil denylist should deny nothing")
	}
	if policy.MatchesDenyList("pkg-a", []string{}) {
		t.Errorf("empty denylist should deny nothing")
	}
	if !policy.MatchesDenyList("bad-pkg", []string{"bad-*"}) {
		t.Errorf("bad-pkg should match bad-*")
	}
}
