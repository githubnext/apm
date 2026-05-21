package hostbackends_test

import (
	"testing"

	"github.com/githubnext/apm/internal/deps/hostbackends"
)

func TestBackendForHost_GitHubKind(t *testing.T) {
	b := hostbackends.BackendForHost("github.com", nil)
	if b.Kind() != "github" {
		t.Errorf("expected github, got %s", b.Kind())
	}
}

func TestBackendForHost_GitHubIsGitHubFamily(t *testing.T) {
	b := hostbackends.BackendForHost("github.com", nil)
	if !b.IsGitHubFamily() {
		t.Error("expected IsGitHubFamily=true for github.com")
	}
}

func TestBackendForHost_GitHubNotGeneric(t *testing.T) {
	b := hostbackends.BackendForHost("github.com", nil)
	if b.IsGeneric() {
		t.Error("expected IsGeneric=false for github.com")
	}
}

func TestBackendForHost_ADOKind(t *testing.T) {
	b := hostbackends.BackendForHost("dev.azure.com", nil)
	if b.Kind() != "ado" {
		t.Errorf("expected ado, got %s", b.Kind())
	}
}

func TestBackendForHost_GitLabKind(t *testing.T) {
	b := hostbackends.BackendForHost("gitlab.com", nil)
	if b.Kind() != "gitlab" {
		t.Errorf("expected gitlab, got %s", b.Kind())
	}
}

func TestBackendForHost_GHEKind(t *testing.T) {
	b := hostbackends.BackendForHost("ghe.example.com", nil)
	if b == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestBackendForHost_GenericKind(t *testing.T) {
	b := hostbackends.BackendForHost("bitbucket.org", nil)
	if b.Kind() == "" {
		t.Error("expected non-empty kind")
	}
}

func TestBackendForHost_GetHostInfoReturnsValue(t *testing.T) {
	b := hostbackends.BackendForHost("github.com", nil)
	hi := b.GetHostInfo()
	_ = hi // just verify it doesn't panic
}

func TestBackendForHost_ADONotGitHubFamily(t *testing.T) {
	b := hostbackends.BackendForHost("dev.azure.com", nil)
	if b.IsGitHubFamily() {
		t.Error("expected IsGitHubFamily=false for ADO")
	}
}
