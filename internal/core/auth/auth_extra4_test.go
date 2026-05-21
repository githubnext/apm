package auth_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/core/auth"
)

func TestHostInfo_ZeroValue_Extra4(t *testing.T) {
	var h auth.HostInfo
	if h.Host != "" || h.Kind != "" {
		t.Fatal("expected zero HostInfo")
	}
}

func TestHostInfo_DisplayName_NilPort_Extra4(t *testing.T) {
	h := auth.HostInfo{Host: "github.com"}
	if h.DisplayName() != "github.com" {
		t.Fatalf("expected github.com, got %s", h.DisplayName())
	}
}

func TestHostInfo_DisplayName_NonStandardPort_Extra4(t *testing.T) {
	port := 8443
	h := auth.HostInfo{Host: "ghe.corp", Port: &port}
	dn := h.DisplayName()
	if !strings.Contains(dn, "8443") {
		t.Fatalf("expected port in display name, got %s", dn)
	}
}

func TestAuthContext_ZeroValue_Extra4(t *testing.T) {
	var c auth.AuthContext
	if c.Token != nil {
		t.Fatal("expected nil Token")
	}
}

func TestBearerFallbackOutcome_ZeroValue_Extra4(t *testing.T) {
	var o auth.BearerFallbackOutcome
	if o.BearerAttempted {
		t.Fatal("expected BearerAttempted=false")
	}
}

func TestClassifyHost_GitHub_Extra4(t *testing.T) {
	h := auth.ClassifyHost("github.com", nil)
	if h.Kind != "github" {
		t.Fatalf("expected github, got %s", h.Kind)
	}
	if h.Host != "github.com" {
		t.Fatalf("expected github.com, got %s", h.Host)
	}
}

func TestClassifyHost_ADO_Extra4(t *testing.T) {
	h := auth.ClassifyHost("dev.azure.com", nil)
	if h.Kind != "ado" {
		t.Fatalf("expected ado, got %s", h.Kind)
	}
}

func TestDetectTokenType_GHU_Extra4(t *testing.T) {
	tok := "ghu_" + strings.Repeat("a", 36)
	typ := auth.DetectTokenType(tok)
	if typ == "" {
		t.Fatal("expected non-empty token type")
	}
}

func TestDetectTokenType_Classic_Extra4(t *testing.T) {
	tok := "ghp_" + strings.Repeat("b", 36)
	typ := auth.DetectTokenType(tok)
	if typ == "" {
		t.Fatal("expected non-empty token type")
	}
}

func TestNewAuthResolver_NotNil_Extra4(t *testing.T) {
	r := auth.NewAuthResolver(nil)
	if r == nil {
		t.Fatal("expected non-nil AuthResolver")
	}
}
