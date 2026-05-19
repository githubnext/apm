package auth

import (
	"testing"
)

func TestDetectTokenType_FinedGrained(t *testing.T) {
	if got := DetectTokenType("github_pat_abc123"); got != "fine-grained" {
		t.Errorf("expected fine-grained, got %q", got)
	}
}

func TestDetectTokenType_Classic(t *testing.T) {
	if got := DetectTokenType("ghp_ABC123"); got != "classic" {
		t.Errorf("expected classic, got %q", got)
	}
}

func TestDetectTokenType_OAuthGhu(t *testing.T) {
	if got := DetectTokenType("ghu_abc"); got != "oauth" {
		t.Errorf("expected oauth, got %q", got)
	}
}

func TestDetectTokenType_OAuthGho(t *testing.T) {
	if got := DetectTokenType("gho_xyz"); got != "oauth" {
		t.Errorf("expected oauth, got %q", got)
	}
}

func TestDetectTokenType_GitHubApp_Ghs(t *testing.T) {
	if got := DetectTokenType("ghs_secret"); got != "github-app" {
		t.Errorf("expected github-app, got %q", got)
	}
}

func TestDetectTokenType_GitHubApp_Ghr(t *testing.T) {
	if got := DetectTokenType("ghr_token"); got != "github-app" {
		t.Errorf("expected github-app, got %q", got)
	}
}

func TestDetectTokenType_Unknown(t *testing.T) {
	for _, tok := range []string{"", "abc", "token123", "GITHUB_TOKEN"} {
		if got := DetectTokenType(tok); got != "unknown" {
			t.Errorf("DetectTokenType(%q) = %q, want unknown", tok, got)
		}
	}
}

func TestGitLabRESTHeaders_EmptyToken(t *testing.T) {
	h := GitLabRESTHeaders("", false)
	if len(h) != 0 {
		t.Errorf("expected empty headers for empty token, got %v", h)
	}
}

func TestGitLabRESTHeaders_PrivateToken(t *testing.T) {
	h := GitLabRESTHeaders("mytoken", false)
	if h["PRIVATE-TOKEN"] != "mytoken" {
		t.Errorf("expected PRIVATE-TOKEN header, got %v", h)
	}
}

func TestGitLabRESTHeaders_OAuthBearer(t *testing.T) {
	h := GitLabRESTHeaders("mytoken", true)
	if h["Authorization"] != "Bearer mytoken" {
		t.Errorf("expected Bearer auth, got %v", h)
	}
}

func TestClassifyHost_GHES(t *testing.T) {
	// A non-standard host that's not github.com, ghe.com, gitlab.com, or dev.azure.com
	// should fall through to ghes or generic
	info := ClassifyHost("myenterprise.example.com", nil)
	if info.Kind == "" {
		t.Error("ClassifyHost should return a non-empty Kind")
	}
}

func TestClassifyHost_GitLabSelf(t *testing.T) {
	info := ClassifyHost("gitlab.myco.com", nil)
	// Either gitlab or generic is acceptable
	if info.Kind == "github" || info.Kind == "ado" {
		t.Errorf("unexpected kind %q for self-hosted GitLab-like host", info.Kind)
	}
}

func TestNewAuthResolver_NotNil(t *testing.T) {
	r := NewAuthResolver(nil)
	if r == nil {
		t.Fatal("NewAuthResolver(nil) returned nil")
	}
}

func TestNewAuthResolver_NilTokenManager(t *testing.T) {
	// Should not panic with nil token manager
	r := NewAuthResolver(nil)
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestHostInfo_DisplayName_NoPort(t *testing.T) {
	h := HostInfo{Host: "example.com"}
	if h.DisplayName() != "example.com" {
		t.Errorf("unexpected: %s", h.DisplayName())
	}
}

func TestHostInfo_DisplayName_Port80Hidden(t *testing.T) {
	p := 80
	h := HostInfo{Host: "example.com", Port: &p}
	if h.DisplayName() != "example.com" {
		t.Errorf("port 80 should be hidden, got %s", h.DisplayName())
	}
}

func TestHostInfo_DisplayName_Port22Hidden(t *testing.T) {
	p := 22
	h := HostInfo{Host: "git.example.com", Port: &p}
	if h.DisplayName() != "git.example.com" {
		t.Errorf("port 22 should be hidden, got %s", h.DisplayName())
	}
}

func TestHostInfo_DisplayName_NonStandardPort(t *testing.T) {
	p := 4433
	h := HostInfo{Host: "ghe.example.com", Port: &p}
	want := "ghe.example.com:4433"
	if h.DisplayName() != want {
		t.Errorf("expected %s, got %s", want, h.DisplayName())
	}
}

func TestClassifyHost_GitHubCom_HasPublicRepos(t *testing.T) {
	info := ClassifyHost("github.com", nil)
	if !info.HasPublicRepos {
		t.Error("github.com should have public repos")
	}
}

func TestClassifyHost_ADO_APIBase(t *testing.T) {
	info := ClassifyHost("dev.azure.com", nil)
	if info.APIBase == "" {
		t.Error("expected non-empty APIBase for ADO")
	}
}
