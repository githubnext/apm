package auth

import (
	"strings"
	"testing"
)

func TestDetectTokenType_EmptyToken(t *testing.T) {
	got := DetectTokenType("")
	if got == "" {
		t.Error("expected non-empty type for empty token")
	}
}

func TestDetectTokenType_ShortString(t *testing.T) {
	got := DetectTokenType("abc")
	if got == "" {
		t.Error("expected non-empty type")
	}
}

func TestGitLabRESTHeaders_BothModes(t *testing.T) {
	h1 := GitLabRESTHeaders("mytoken", false)
	h2 := GitLabRESTHeaders("mytoken", true)
	if len(h1) == 0 || len(h2) == 0 {
		t.Error("expected non-empty headers for both modes")
	}
	// Header key may be uppercase; check case-insensitively
	found := false
	for k, v := range h1 {
		if strings.EqualFold(k, "Private-Token") && v == "mytoken" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected Private-Token header, got %v", h1)
	}
}

func TestClassifyHost_Generic(t *testing.T) {
	hi := ClassifyHost("example.com", nil)
	if hi.Host == "" {
		t.Error("expected non-empty host")
	}
}

func TestClassifyHost_GitHub(t *testing.T) {
	hi := ClassifyHost("github.com", nil)
	if hi.Kind != "github" {
		t.Errorf("expected kind=github, got %q", hi.Kind)
	}
}

func TestClassifyHost_ADO(t *testing.T) {
	hi := ClassifyHost("dev.azure.com", nil)
	if hi.Kind != "ado" {
		t.Errorf("expected kind=ado, got %q", hi.Kind)
	}
}

func TestClassifyHost_PortPreserved(t *testing.T) {
	port := 8443
	hi := ClassifyHost("mygithub.corp.com", &port)
	if hi.Port == nil || *hi.Port != 8443 {
		t.Errorf("expected port 8443, got %v", hi.Port)
	}
}

func TestHostInfo_ZeroValue(t *testing.T) {
	var hi HostInfo
	if hi.Host != "" || hi.Kind != "" || hi.HasPublicRepos || hi.APIBase != "" {
		t.Error("expected zero value")
	}
}

func TestAuthContext_ZeroValue(t *testing.T) {
	var ac AuthContext
	if ac.Token != nil || ac.Source != "" || ac.TokenType != "" {
		t.Error("expected zero value")
	}
}

func TestAuthContext_TokenSet(t *testing.T) {
	tok := "secret"
	ac := AuthContext{
		Token:     &tok,
		Source:    "GITHUB_TOKEN",
		TokenType: "classic",
	}
	if *ac.Token != "secret" {
		t.Errorf("unexpected token")
	}
	if ac.Source != "GITHUB_TOKEN" {
		t.Errorf("unexpected source: %q", ac.Source)
	}
}

func TestHostInfo_DisplayName_GitLabSelf(t *testing.T) {
	hi := HostInfo{Host: "gitlab.corp.com", Kind: "gitlab"}
	dn := hi.DisplayName()
	if dn == "" {
		t.Error("expected non-empty display name")
	}
	if !strings.Contains(dn, "gitlab.corp.com") {
		t.Errorf("expected host in display name, got %q", dn)
	}
}

func TestHostInfo_DisplayName_GHECloud(t *testing.T) {
	hi := HostInfo{Host: "myorg.ghe.com", Kind: "ghe_cloud"}
	dn := hi.DisplayName()
	if !strings.Contains(dn, "myorg.ghe.com") {
		t.Errorf("expected host in display name, got %q", dn)
	}
}
