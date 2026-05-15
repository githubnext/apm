package auth

import (
	"testing"
)

func TestClassifyHost_github(t *testing.T) {
	info := ClassifyHost("github.com", nil)
	if info.Kind != "github" {
		t.Errorf("expected github, got %s", info.Kind)
	}
	if !info.HasPublicRepos {
		t.Error("github.com should have public repos")
	}
	if info.APIBase != "https://api.github.com" {
		t.Errorf("unexpected APIBase: %s", info.APIBase)
	}
}

func TestClassifyHost_ghe_cloud(t *testing.T) {
	info := ClassifyHost("myorg.ghe.com", nil)
	if info.Kind != "ghe_cloud" {
		t.Errorf("expected ghe_cloud, got %s", info.Kind)
	}
}

func TestClassifyHost_gitlab(t *testing.T) {
	info := ClassifyHost("gitlab.com", nil)
	if info.Kind != "gitlab" {
		t.Errorf("expected gitlab, got %s", info.Kind)
	}
	if info.APIBase != "https://gitlab.com/api/v4" {
		t.Errorf("unexpected APIBase: %s", info.APIBase)
	}
}

func TestClassifyHost_ado(t *testing.T) {
	info := ClassifyHost("dev.azure.com", nil)
	if info.Kind != "ado" {
		t.Errorf("expected ado, got %s", info.Kind)
	}
}

func TestClassifyHost_generic(t *testing.T) {
	info := ClassifyHost("bitbucket.example.com", nil)
	if info.Kind != "generic" {
		t.Errorf("expected generic, got %s", info.Kind)
	}
}

func TestClassifyHost_case_insensitive(t *testing.T) {
	info := ClassifyHost("GitHub.COM", nil)
	if info.Kind != "github" {
		t.Errorf("expected github for uppercase, got %s", info.Kind)
	}
}

func TestHostInfo_DisplayName_no_port(t *testing.T) {
	h := HostInfo{Host: "github.com"}
	if h.DisplayName() != "github.com" {
		t.Errorf("unexpected display name: %s", h.DisplayName())
	}
}

func TestHostInfo_DisplayName_with_nonstandard_port(t *testing.T) {
	p := 8080
	h := HostInfo{Host: "myghe.com", Port: &p}
	got := h.DisplayName()
	if got != "myghe.com:8080" {
		t.Errorf("expected myghe.com:8080, got %s", got)
	}
}

func TestHostInfo_DisplayName_standard_port_443(t *testing.T) {
	p := 443
	h := HostInfo{Host: "github.com", Port: &p}
	if h.DisplayName() != "github.com" {
		t.Errorf("port 443 should be hidden, got %s", h.DisplayName())
	}
}

func TestDetectTokenType_fine_grained(t *testing.T) {
	tt := DetectTokenType("github_pat_ABCDEF")
	if tt != "fine-grained" {
		t.Errorf("expected fine-grained, got %s", tt)
	}
}

func TestDetectTokenType_classic(t *testing.T) {
	tt := DetectTokenType("ghp_ABCDEF")
	if tt != "classic" {
		t.Errorf("expected classic, got %s", tt)
	}
}

func TestDetectTokenType_unknown(t *testing.T) {
	tt := DetectTokenType("someothertoken")
	if tt != "unknown" {
		t.Errorf("expected unknown, got %s", tt)
	}
}

func TestNewAuthResolver_not_nil(t *testing.T) {
	r := NewAuthResolver(nil)
	if r == nil {
		t.Fatal("NewAuthResolver returned nil")
	}
}

func TestGitLabRESTHeaders_with_token(t *testing.T) {
	headers := GitLabRESTHeaders("mytoken", false)
	if headers["PRIVATE-TOKEN"] != "mytoken" {
		t.Errorf("expected PRIVATE-TOKEN header, got %v", headers)
	}
}

func TestGitLabRESTHeaders_oauth_bearer(t *testing.T) {
	headers := GitLabRESTHeaders("mytoken", true)
	if headers["Authorization"] != "Bearer mytoken" {
		t.Errorf("expected Bearer Authorization header, got %v", headers)
	}
}
