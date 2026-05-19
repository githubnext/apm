package insecurepolicy_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/insecurepolicy"
)

func TestIsValidFQDN_ValidCases(t *testing.T) {
	valid := []string{
		"mirror.example.com",
		"a.b.c.d.example.org",
		"my-registry.internal.corp",
		"x.io",
	}
	for _, h := range valid {
		if !insecurepolicy.IsValidFQDN(h) {
			t.Errorf("IsValidFQDN(%q) should be true", h)
		}
	}
}

func TestIsValidFQDN_InvalidCases(t *testing.T) {
	invalid := []string{
		"",
		"localhost",
		"192.168.1.1",
		"single",
		"-starts-dash.com",
		"has space.com",
		"has/slash.com",
	}
	for _, h := range invalid {
		if insecurepolicy.IsValidFQDN(h) {
			t.Errorf("IsValidFQDN(%q) should be false", h)
		}
	}
}

func TestNormalizeAllowInsecureHost_ValidCase(t *testing.T) {
	norm, err := insecurepolicy.NormalizeAllowInsecureHost("Mirror.Example.COM")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if norm != "mirror.example.com" {
		t.Errorf("expected lowercase, got %q", norm)
	}
}

func TestNormalizeAllowInsecureHost_InvalidCase(t *testing.T) {
	_, err := insecurepolicy.NormalizeAllowInsecureHost("localhost")
	if err == nil {
		t.Error("expected error for 'localhost'")
	}
}

func TestNormalizeAllowInsecureHost_Whitespace(t *testing.T) {
	norm, err := insecurepolicy.NormalizeAllowInsecureHost("  mirror.example.com  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if norm != "mirror.example.com" {
		t.Errorf("expected trimmed+lowercased, got %q", norm)
	}
}

func TestGetInsecureDependencyHost_ValidURL(t *testing.T) {
	info := insecurepolicy.InsecureDependencyInfo{URL: "http://mirror.example.com/pkg.tar.gz"}
	got := insecurepolicy.GetInsecureDependencyHost(info)
	if got != "mirror.example.com" {
		t.Errorf("expected mirror.example.com, got %q", got)
	}
}

func TestGetInsecureDependencyHost_EmptyURL(t *testing.T) {
	info := insecurepolicy.InsecureDependencyInfo{URL: ""}
	got := insecurepolicy.GetInsecureDependencyHost(info)
	if got != "" {
		t.Errorf("expected empty for empty URL, got %q", got)
	}
}

func TestFormatInsecureDependencyRequirements_BothFlags(t *testing.T) {
	out := insecurepolicy.FormatInsecureDependencyRequirements(
		"http://mirror.example.com/pkg",
		true, true,
	)
	if !strings.Contains(out, "allow_insecure") {
		t.Errorf("expected allow_insecure in output, got %q", out)
	}
	if !strings.Contains(out, "--allow-insecure") {
		t.Errorf("expected --allow-insecure in output, got %q", out)
	}
}

func TestFormatInsecureDependencyRequirements_NeitherFlag(t *testing.T) {
	out := insecurepolicy.FormatInsecureDependencyRequirements(
		"http://mirror.example.com/pkg",
		false, false,
	)
	if strings.Contains(out, "allow_insecure") {
		t.Errorf("unexpected allow_insecure when both false, got %q", out)
	}
}

func TestFormatInsecureDependencyWarning_DirectURL(t *testing.T) {
	info := insecurepolicy.InsecureDependencyInfo{
		URL:          "http://mirror.example.com/pkg",
		IsTransitive: false,
	}
	out := insecurepolicy.FormatInsecureDependencyWarning(info)
	if !strings.Contains(out, "mirror.example.com") {
		t.Errorf("expected URL in warning, got %q", out)
	}
}

func TestFormatInsecureDependencyWarning_TransitiveDep(t *testing.T) {
	info := insecurepolicy.InsecureDependencyInfo{
		URL:          "http://mirror.example.com/pkg",
		IsTransitive: true,
		IntroducedBy: "root-pkg",
	}
	out := insecurepolicy.FormatInsecureDependencyWarning(info)
	if !strings.Contains(out, "root-pkg") {
		t.Errorf("expected introducedBy in warning, got %q", out)
	}
	if !strings.Contains(out, "transitive") {
		t.Errorf("expected 'transitive' in warning, got %q", out)
	}
}
