package tagpattern_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace/tagpattern"
)

func TestRenderTag_EmptyPattern(t *testing.T) {
	result := tagpattern.RenderTag("", "myname", "1.0.0")
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestRenderTag_NoPlaceholdersPreserved(t *testing.T) {
	pattern := "v-release"
	result := tagpattern.RenderTag(pattern, "x", "1.0")
	if result != pattern {
		t.Errorf("expected %q, got %q", pattern, result)
	}
}

func TestRenderTag_VersionFirst(t *testing.T) {
	result := tagpattern.RenderTag("{version}-{name}", "pkg", "2.0.1")
	if result != "2.0.1-pkg" {
		t.Errorf("unexpected: %q", result)
	}
}

func TestBuildTagRegex_VersionOnly_ExtractWorks(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("v{version}")
	if err != nil {
		t.Fatalf("build regex: %v", err)
	}
	v, ok := tagpattern.ExtractVersion(re, "v1.2.3")
	if !ok || v != "1.2.3" {
		t.Errorf("expected '1.2.3', got %q ok=%v", v, ok)
	}
}

func TestBuildTagRegex_PrefixSuffix(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("release-{version}-stable")
	if err != nil {
		t.Fatalf("build regex: %v", err)
	}
	v, ok := tagpattern.ExtractVersion(re, "release-3.0.0-stable")
	if !ok || v != "3.0.0" {
		t.Errorf("expected '3.0.0', got %q ok=%v", v, ok)
	}
}

func TestExtractVersion_WrongPrefix(t *testing.T) {
	re, _ := tagpattern.BuildTagRegex("v{version}")
	_, ok := tagpattern.ExtractVersion(re, "x1.0.0")
	if ok {
		t.Error("should not match wrong prefix")
	}
}

func TestBuildTagRegex_ReturnsNonNil(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("{version}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if re == nil {
		t.Error("expected non-nil regex")
	}
}

func TestRenderTag_BothInReverseOrder(t *testing.T) {
	result := tagpattern.RenderTag("{name}@{version}", "lib", "0.1")
	if result != "lib@0.1" {
		t.Errorf("unexpected: %q", result)
	}
}

func TestExtractVersion_CorrectVersion(t *testing.T) {
	re, _ := tagpattern.BuildTagRegex("pkg-{version}")
	v, ok := tagpattern.ExtractVersion(re, "pkg-4.5.6")
	if !ok || v != "4.5.6" {
		t.Errorf("expected '4.5.6', got %q ok=%v", v, ok)
	}
}

func TestBuildTagRegex_NoVersionPlaceholder(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := tagpattern.ExtractVersion(re, "latest")
	// No version named group -- should not panic
	_ = ok
}

func TestRenderTag_RepeatedPlaceholders(t *testing.T) {
	result := tagpattern.RenderTag("{version}-{version}", "n", "9.9")
	if result != "9.9-9.9" {
		t.Errorf("unexpected: %q", result)
	}
}
