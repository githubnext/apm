package auth

import (
	"testing"
)

func TestClassifyHost_GitHub2(t *testing.T) {
	h := ClassifyHost("github.com", nil)
	if h.Kind != "github" {
		t.Fatalf("expected 'github', got %q", h.Kind)
	}
}

func TestClassifyHost_ADO2(t *testing.T) {
	h := ClassifyHost("dev.azure.com", nil)
	if h.Kind != "ado" {
		t.Fatalf("expected 'ado', got %q", h.Kind)
	}
}

func TestClassifyHost_GitLab(t *testing.T) {
	h := ClassifyHost("gitlab.com", nil)
	if h.Kind != "gitlab" && h.Kind != "generic" {
		t.Fatalf("unexpected kind for gitlab.com: %q", h.Kind)
	}
}

func TestClassifyHost_Generic2(t *testing.T) {
	h := ClassifyHost("mycompany.internal", nil)
	if h.Kind == "" {
		t.Fatal("kind should not be empty for generic host")
	}
}

func TestClassifyHost_HostPreserved(t *testing.T) {
	h := ClassifyHost("github.com", nil)
	if h.Host != "github.com" {
		t.Fatalf("expected host 'github.com', got %q", h.Host)
	}
}

func TestDetectTokenType_GHSPrefix(t *testing.T) {
	tt := DetectTokenType("ghs_ABCDEFGHIJ1234567890")
	if tt == "" {
		t.Fatal("token type should not be empty")
	}
}

func TestDetectTokenType_GHPPrefix(t *testing.T) {
	tt := DetectTokenType("ghp_ABCDEFGHIJ1234567890")
	if tt == "" {
		t.Fatal("token type should not be empty")
	}
}

func TestDetectTokenType_GHFPrefix(t *testing.T) {
	tt := DetectTokenType("github_pat_ABCDEFGHIJ")
	if tt == "" {
		t.Fatal("token type should not be empty for fine-grained pat")
	}
}

func TestDetectTokenType_EmptyString2(t *testing.T) {
	tt := DetectTokenType("")
	if tt == "" {
		t.Fatal("should return a type string even for empty token")
	}
}

func TestHostInfo_DisplayName_StandardPort(t *testing.T) {
	port := 443
	h := HostInfo{Host: "github.com", Port: &port}
	name := h.DisplayName()
	if name != "github.com" {
		t.Fatalf("standard port 443 should be hidden, got %q", name)
	}
}

func TestHostInfo_DisplayName_CustomPort(t *testing.T) {
	port := 8443
	h := HostInfo{Host: "ghe.company.com", Port: &port}
	name := h.DisplayName()
	if name == "ghe.company.com" {
		t.Fatalf("non-standard port should appear in display name, got %q", name)
	}
}

func TestHostInfo_ZeroValue2(t *testing.T) {
	var h HostInfo
	if h.Host != "" || h.Kind != "" {
		t.Fatal("zero value HostInfo should have empty fields")
	}
}

func TestAuthContext_ZeroValue2(t *testing.T) {
	var a AuthContext
	if a.Token != nil {
		t.Fatal("zero AuthContext should have nil Token")
	}
}

func TestGitLabRESTHeaders_TokenKey(t *testing.T) {
	headers := GitLabRESTHeaders("mytoken", false)
	if len(headers) == 0 {
		t.Fatal("expected at least one header")
	}
}

func TestGitLabRESTHeaders_OAuthMode(t *testing.T) {
	headers := GitLabRESTHeaders("oauthtoken", true)
	if len(headers) == 0 {
		t.Fatal("expected at least one header in oauth mode")
	}
}
