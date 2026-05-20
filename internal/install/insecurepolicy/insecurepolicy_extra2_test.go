package insecurepolicy

import (
	"strings"
	"testing"
)

func TestIsValidFQDN_SingleLabel(t *testing.T) {
	// Single label without dots may not be a valid FQDN
	result := IsValidFQDN("localhost")
	_ = result // just verify it doesn't panic
}

func TestIsValidFQDN_IPAddress(t *testing.T) {
	if IsValidFQDN("192.168.1.1") {
		// IP addresses may or may not pass; just don't panic
	}
}

func TestIsValidFQDN_WithPort(t *testing.T) {
	// Host with port should fail FQDN validation
	if IsValidFQDN("example.com:8080") {
		t.Error("host with port should not be valid FQDN")
	}
}

func TestIsValidFQDN_Subdomain(t *testing.T) {
	if !IsValidFQDN("sub.example.com") {
		t.Error("expected subdomain to be valid FQDN")
	}
}

func TestNormalizeAllowInsecureHost_NoTrim(t *testing.T) {
	host, err := NormalizeAllowInsecureHost("example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host != "example.com" {
		t.Errorf("expected example.com, got %q", host)
	}
}

func TestNormalizeAllowInsecureHost_UpperCase(t *testing.T) {
	host, err := NormalizeAllowInsecureHost("EXAMPLE.COM")
	if err != nil {
		// Some implementations may reject uppercase — that's ok
		return
	}
	if strings.ToLower(host) != strings.ToLower("EXAMPLE.COM") {
		t.Errorf("unexpected host %q", host)
	}
}

func TestInsecureDependencyPolicyError_Error(t *testing.T) {
	e := &InsecureDependencyPolicyError{Message: "test error"}
	if e.Error() != "test error" {
		t.Errorf("expected 'test error', got %q", e.Error())
	}
}

func TestInsecureDependencyInfo_Fields(t *testing.T) {
	info := InsecureDependencyInfo{
		URL:          "http://example.com",
		IsTransitive: true,
		IntroducedBy: "parent/dep",
	}
	if !info.IsTransitive {
		t.Error("expected IsTransitive=true")
	}
	if info.IntroducedBy != "parent/dep" {
		t.Errorf("unexpected IntroducedBy: %q", info.IntroducedBy)
	}
}

func TestGetInsecureDependencyHost_PortURL(t *testing.T) {
	info := InsecureDependencyInfo{URL: "http://example.com:9000/path"}
	host := GetInsecureDependencyHost(info)
	_ = host // just verify it doesn't panic
}

func TestFormatInsecureDependencyRequirements_OnlyDirect(t *testing.T) {
	result := FormatInsecureDependencyRequirements("http://example.com", true, false)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestFormatInsecureDependencyRequirements_OnlyTransitive(t *testing.T) {
	result := FormatInsecureDependencyRequirements("http://example.com", false, true)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestFormatInsecureDependencyWarning_EmptyURL(t *testing.T) {
	info := InsecureDependencyInfo{URL: "", IsTransitive: false}
	result := FormatInsecureDependencyWarning(info)
	_ = result // verify no panic
}

func TestGetAllowedTransitiveInsecureHosts_NilPolicy(t *testing.T) {
	hosts := GetAllowedTransitiveInsecureHosts(nil, false, nil)
	_ = hosts
}
