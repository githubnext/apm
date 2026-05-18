package insecurepolicy_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/insecurepolicy"
)

func TestIsValidFQDN(t *testing.T) {
	tests := []struct {
		host  string
		valid bool
	}{
		{"mirror.example.com", true},
		{"example.com", true},
		{"sub.domain.example.org", true},
		{"localhost", false},
		{"192.168.1.1", false},
		{"", false},
		{"invalid", false},
		{"-bad.com", false},
	}
	for _, tt := range tests {
		got := insecurepolicy.IsValidFQDN(tt.host)
		if got != tt.valid {
			t.Errorf("IsValidFQDN(%q) = %v, want %v", tt.host, got, tt.valid)
		}
	}
}

func TestNormalizeAllowInsecureHost_Valid(t *testing.T) {
	got, err := insecurepolicy.NormalizeAllowInsecureHost("MIRROR.Example.Com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "mirror.example.com" {
		t.Fatalf("expected lowercase, got %q", got)
	}
}

func TestNormalizeAllowInsecureHost_Invalid(t *testing.T) {
	_, err := insecurepolicy.NormalizeAllowInsecureHost("localhost")
	if err == nil {
		t.Fatal("expected error for localhost")
	}
}

func TestGetInsecureDependencyHost(t *testing.T) {
	info := insecurepolicy.InsecureDependencyInfo{URL: "http://mirror.example.com/repo.git"}
	got := insecurepolicy.GetInsecureDependencyHost(info)
	if got != "mirror.example.com" {
		t.Fatalf("expected 'mirror.example.com', got %q", got)
	}
}

func TestGetInsecureDependencyHost_Empty(t *testing.T) {
	info := insecurepolicy.InsecureDependencyInfo{URL: "not-a-url"}
	got := insecurepolicy.GetInsecureDependencyHost(info)
	if got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestFormatInsecureDependencyRequirements_BothMissing(t *testing.T) {
	msg := insecurepolicy.FormatInsecureDependencyRequirements("http://example.com/r", true, true)
	if !strings.Contains(msg, "allow_insecure") {
		t.Error("expected allow_insecure step")
	}
	if !strings.Contains(msg, "--allow-insecure") {
		t.Error("expected --allow-insecure step")
	}
}

func TestFormatInsecureDependencyRequirements_OnlyCLI(t *testing.T) {
	msg := insecurepolicy.FormatInsecureDependencyRequirements("http://example.com/r", false, true)
	if strings.Contains(msg, "allow_insecure: true") {
		t.Error("should not include allow_insecure step")
	}
	if !strings.Contains(msg, "--allow-insecure") {
		t.Error("expected --allow-insecure step")
	}
}

func TestFormatInsecureDependencyWarning_Direct(t *testing.T) {
	info := insecurepolicy.InsecureDependencyInfo{URL: "http://mirror.example.com/r"}
	msg := insecurepolicy.FormatInsecureDependencyWarning(info)
	if !strings.Contains(msg, "http://mirror.example.com/r") {
		t.Errorf("expected URL in warning, got %q", msg)
	}
}

func TestFormatInsecureDependencyWarning_Transitive(t *testing.T) {
	info := insecurepolicy.InsecureDependencyInfo{
		URL:          "http://mirror.example.com/r",
		IsTransitive: true,
		IntroducedBy: "owner/root-dep",
	}
	msg := insecurepolicy.FormatInsecureDependencyWarning(info)
	if !strings.Contains(msg, "transitive") {
		t.Errorf("expected 'transitive', got %q", msg)
	}
	if !strings.Contains(msg, "owner/root-dep") {
		t.Errorf("expected introduced-by, got %q", msg)
	}
}

func TestGuardTransitiveInsecureDependencies_Allowed(t *testing.T) {
	infos := []insecurepolicy.InsecureDependencyInfo{
		{URL: "http://mirror.example.com/r", IsTransitive: false},
		{URL: "http://mirror.example.com/dep", IsTransitive: true},
	}
	// allowInsecure=true opens all direct-dep hosts transitively
	err := insecurepolicy.GuardTransitiveInsecureDependencies(infos, true, nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestGuardTransitiveInsecureDependencies_Blocked(t *testing.T) {
	infos := []insecurepolicy.InsecureDependencyInfo{
		{URL: "http://blocked.example.com/dep", IsTransitive: true},
	}
	err := insecurepolicy.GuardTransitiveInsecureDependencies(infos, false, nil)
	if err == nil {
		t.Fatal("expected policy error")
	}
	if !strings.Contains(err.Error(), "blocked.example.com") {
		t.Errorf("error should mention host, got %v", err)
	}
}

func TestGuardTransitiveInsecureDependencies_AllowedByHost(t *testing.T) {
	infos := []insecurepolicy.InsecureDependencyInfo{
		{URL: "http://mirror.example.com/dep", IsTransitive: true},
	}
	err := insecurepolicy.GuardTransitiveInsecureDependencies(infos, false, []string{"mirror.example.com"})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
