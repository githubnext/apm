package core_test

import (
	"os"
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/core"
)

// ---------------------------------------------------------------------------
// Parity: HostInfo.DisplayName
// ---------------------------------------------------------------------------

func TestParityHostInfoDisplayName(t *testing.T) {
	cases := []struct {
		host string
		port int
		want string
	}{
		{"github.com", 0, "github.com"},
		{"github.com", 443, "github.com"},
		{"github.com", 80, "github.com"},
		{"bitbucket.example.com", 7999, "bitbucket.example.com:7999"},
		{"bitbucket.example.com", 7990, "bitbucket.example.com:7990"},
	}
	for _, c := range cases {
		h := core.HostInfo{Host: c.host, Port: c.port}
		if got := h.DisplayName(); got != c.want {
			t.Errorf("DisplayName(%q, %d) = %q, want %q", c.host, c.port, got, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Parity: ClassifyHost
// ---------------------------------------------------------------------------

func TestParityClassifyHostGitHub(t *testing.T) {
	info := core.ClassifyHost("github.com", 0)
	if info.Kind != "github" {
		t.Errorf("github.com kind = %q, want github", info.Kind)
	}
	if !info.HasPublicRepos {
		t.Error("github.com should have public repos")
	}
	if info.APIBase != "https://api.github.com" {
		t.Errorf("github.com APIBase = %q", info.APIBase)
	}
}

func TestParityClassifyHostGHECloud(t *testing.T) {
	info := core.ClassifyHost("myenterprise.ghe.com", 0)
	if info.Kind != "ghe_cloud" {
		t.Errorf("*.ghe.com kind = %q, want ghe_cloud", info.Kind)
	}
	if info.HasPublicRepos {
		t.Error("ghe_cloud should NOT have public repos")
	}
}

func TestParityClassifyHostADO(t *testing.T) {
	info := core.ClassifyHost("dev.azure.com", 0)
	if info.Kind != "ado" {
		t.Errorf("dev.azure.com kind = %q, want ado", info.Kind)
	}
}

func TestParityClassifyHostVisualStudio(t *testing.T) {
	info := core.ClassifyHost("myorg.visualstudio.com", 0)
	if info.Kind != "ado" {
		t.Errorf("*.visualstudio.com kind = %q, want ado", info.Kind)
	}
}

func TestParityClassifyHostGitLab(t *testing.T) {
	info := core.ClassifyHost("gitlab.com", 0)
	if info.Kind != "gitlab" {
		t.Errorf("gitlab.com kind = %q, want gitlab", info.Kind)
	}
	if info.APIBase != "https://gitlab.com/api/v4" {
		t.Errorf("gitlab.com APIBase = %q", info.APIBase)
	}
}

func TestParityClassifyHostGHES(t *testing.T) {
	t.Setenv("GITHUB_HOST", "ghes.example.com")
	info := core.ClassifyHost("ghes.example.com", 0)
	if info.Kind != "ghes" {
		t.Errorf("GITHUB_HOST=ghes.example.com kind = %q, want ghes", info.Kind)
	}
}

func TestParityClassifyHostGeneric(t *testing.T) {
	os.Unsetenv("GITHUB_HOST")
	info := core.ClassifyHost("bitbucket.example.com", 0)
	if info.Kind != "generic" {
		t.Errorf("generic host kind = %q, want generic", info.Kind)
	}
}

func TestParityClassifyHostPort(t *testing.T) {
	info := core.ClassifyHost("bitbucket.example.com", 7999)
	if info.Port != 7999 {
		t.Errorf("port = %d, want 7999", info.Port)
	}
}

// ---------------------------------------------------------------------------
// Parity: DetectTokenType
// ---------------------------------------------------------------------------

func TestParityDetectTokenType(t *testing.T) {
	cases := []struct {
		token string
		want  string
	}{
		{"github_pat_abc123", "fine-grained"},
		{"ghp_abc123", "classic"},
		{"ghu_abc123", "oauth"},
		{"gho_abc123", "oauth"},
		{"ghs_abc123", "github-app"},
		{"ghr_abc123", "github-app"},
		{"sometoken", "unknown"},
		{"", "unknown"},
	}
	for _, c := range cases {
		got := core.DetectTokenType(c.token)
		if got != c.want {
			t.Errorf("DetectTokenType(%q) = %q, want %q", c.token, got, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Parity: GitLabRESTHeaders
// ---------------------------------------------------------------------------

func TestParityGitLabRESTHeaders(t *testing.T) {
	headers := core.GitLabRESTHeaders("mytoken", false)
	if headers["PRIVATE-TOKEN"] != "mytoken" {
		t.Errorf("PAT header = %q", headers["PRIVATE-TOKEN"])
	}
	bearer := core.GitLabRESTHeaders("mytoken", true)
	if bearer["Authorization"] != "Bearer mytoken" {
		t.Errorf("bearer header = %q", bearer["Authorization"])
	}
	empty := core.GitLabRESTHeaders("", false)
	if len(empty) != 0 {
		t.Error("empty token should return empty map")
	}
}

// ---------------------------------------------------------------------------
// Parity: AuthResolver.Resolve -- token resolution from env
// ---------------------------------------------------------------------------

func TestParityAuthResolverResolveGitHub(t *testing.T) {
	t.Setenv("GITHUB_APM_PAT", "ghp_testtoken")
	defer os.Unsetenv("GITHUB_APM_PAT")

	r := core.NewAuthResolver()
	ctx := r.Resolve("github.com", "", 0)
	if ctx.Token != "ghp_testtoken" {
		t.Errorf("token = %q, want ghp_testtoken", ctx.Token)
	}
	if ctx.Source != "GITHUB_APM_PAT" {
		t.Errorf("source = %q, want GITHUB_APM_PAT", ctx.Source)
	}
	if ctx.HostInfo.Kind != "github" {
		t.Errorf("kind = %q, want github", ctx.HostInfo.Kind)
	}
}

func TestParityAuthResolverResolveADO(t *testing.T) {
	t.Setenv("ADO_APM_PAT", "adotoken")
	defer os.Unsetenv("ADO_APM_PAT")

	r := core.NewAuthResolver()
	ctx := r.Resolve("dev.azure.com", "", 0)
	if ctx.Token != "adotoken" {
		t.Errorf("token = %q, want adotoken", ctx.Token)
	}
	if ctx.Source != "ADO_APM_PAT" {
		t.Errorf("source = %q, want ADO_APM_PAT", ctx.Source)
	}
}

func TestParityAuthResolverResolveNoToken(t *testing.T) {
	os.Unsetenv("GITHUB_APM_PAT")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("GH_TOKEN")

	r := core.NewAuthResolver()
	ctx := r.Resolve("github.com", "", 0)
	// In CI without gh CLI or git credentials, token should be "" with source "none" or "gh-auth-token"
	if ctx.Token == "" && ctx.Source != "none" && ctx.Source != "git-credential-fill" && ctx.Source != "gh-auth-token" {
		t.Errorf("unexpected source %q when no token env set", ctx.Source)
	}
}

func TestParityAuthResolverResolvePerOrgToken(t *testing.T) {
	t.Setenv("GITHUB_APM_PAT_MYORG", "ghp_orgtoken")
	defer os.Unsetenv("GITHUB_APM_PAT_MYORG")

	r := core.NewAuthResolver()
	ctx := r.Resolve("github.com", "myorg", 0)
	if ctx.Token != "ghp_orgtoken" {
		t.Errorf("per-org token = %q, want ghp_orgtoken", ctx.Token)
	}
	if ctx.Source != "GITHUB_APM_PAT_MYORG" {
		t.Errorf("source = %q, want GITHUB_APM_PAT_MYORG", ctx.Source)
	}
}

func TestParityAuthResolverCacheHit(t *testing.T) {
	t.Setenv("GITHUB_APM_PAT", "ghp_cached")
	defer os.Unsetenv("GITHUB_APM_PAT")

	r := core.NewAuthResolver()
	ctx1 := r.Resolve("github.com", "", 0)
	ctx2 := r.Resolve("github.com", "", 0)
	if ctx1 != ctx2 {
		t.Error("second resolve should return cached pointer")
	}
}

func TestParityAuthResolverOrgToEnvSuffix(t *testing.T) {
	// Verify per-org env var naming: hyphens -> underscores, upper-case.
	t.Setenv("GITHUB_APM_PAT_MY_COOL_ORG", "ghp_orgtoken2")
	defer os.Unsetenv("GITHUB_APM_PAT_MY_COOL_ORG")

	r := core.NewAuthResolver()
	ctx := r.Resolve("github.com", "my-cool-org", 0)
	if ctx.Token != "ghp_orgtoken2" {
		t.Errorf("hyphen org token = %q, want ghp_orgtoken2", ctx.Token)
	}
}

// ---------------------------------------------------------------------------
// Parity: token_manager utilities
// ---------------------------------------------------------------------------

func TestParityTokenManagerValidateTokensPass(t *testing.T) {
	env := map[string]string{"GITHUB_APM_PAT": "ghp_test"}
	mgr := core.NewGitHubTokenManager()
	ok, _ := mgr.ValidateTokens(env)
	if !ok {
		t.Error("ValidateTokens with GITHUB_APM_PAT should pass")
	}
}

func TestParityTokenManagerValidateTokensFail(t *testing.T) {
	env := map[string]string{}
	mgr := core.NewGitHubTokenManager()
	ok, msg := mgr.ValidateTokens(env)
	if ok {
		t.Error("ValidateTokens with no tokens should fail")
	}
	if !strings.Contains(msg, "No tokens found") {
		t.Errorf("message = %q", msg)
	}
}

func TestParityTokenManagerGetTokenForPurpose(t *testing.T) {
	env := map[string]string{"GITHUB_APM_PAT": "ghp_test"}
	mgr := core.NewGitHubTokenManager()
	token, ok := mgr.GetTokenForPurpose("modules", env)
	if !ok || token != "ghp_test" {
		t.Errorf("modules token = %q, ok = %v", token, ok)
	}
}

func TestParityTokenManagerGetTokenUnknownPurpose(t *testing.T) {
	env := map[string]string{"GITHUB_APM_PAT": "ghp_test"}
	mgr := core.NewGitHubTokenManager()
	_, ok := mgr.GetTokenForPurpose("nonexistent", env)
	if ok {
		t.Error("unknown purpose should return false")
	}
}
