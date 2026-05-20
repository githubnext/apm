package hostbackends

import (
	"testing"
)

func TestBackendForHost_GitHubE3(t *testing.T) {
	b := BackendForHost("github.com", nil)
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestBackendForHost_GitLabE3(t *testing.T) {
	b := BackendForHost("gitlab.com", nil)
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestBackendForHost_GenericHostE3(t *testing.T) {
	b := BackendForHost("myhost.example.com", nil)
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestBackendForHost_ADO(t *testing.T) {
	b := BackendForHost("dev.azure.com", nil)
	if b == nil {
		t.Fatal("expected non-nil backend for ADO")
	}
}

func TestBackendForHost_WithPort(t *testing.T) {
	port := 8080
	b := BackendForHost("custom.host.com", &port)
	if b == nil {
		t.Fatal("expected non-nil backend with port")
	}
}

func TestBackendForHost_EmptyHost(t *testing.T) {
	b := BackendForHost("", nil)
	if b == nil {
		t.Fatal("expected non-nil backend for empty host")
	}
}

func TestBackendForHost_IsGitHubFamily_True(t *testing.T) {
	b := BackendForHost("github.com", nil)
	if b.IsGitHubFamily() != true {
		t.Error("github.com should be github family")
	}
}

func TestBackendForHost_IsGitHubFamily_False(t *testing.T) {
	b := BackendForHost("gitlab.com", nil)
	if b.IsGitHubFamily() {
		t.Error("gitlab.com should not be github family")
	}
}

func TestBackendForHost_GetHostInfo_NonPanic(t *testing.T) {
	b := BackendForHost("github.com", nil)
	_ = b.GetHostInfo()
}

func TestBackendForHost_TwoDistinctHosts(t *testing.T) {
	b1 := BackendForHost("github.com", nil)
	b2 := BackendForHost("gitlab.com", nil)
	if b1 == b2 {
		t.Error("expected distinct backends for different hosts")
	}
}
