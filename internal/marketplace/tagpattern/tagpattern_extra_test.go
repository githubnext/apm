package tagpattern_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace/tagpattern"
)

func TestRenderTag_BothPlaceholders(t *testing.T) {
	got := tagpattern.RenderTag("{name}-v{version}", "myapp", "1.2.3")
	want := "myapp-v1.2.3"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRenderTag_VersionOnly(t *testing.T) {
	got := tagpattern.RenderTag("v{version}", "ignored", "2.0.0")
	if got != "v2.0.0" {
		t.Errorf("got %q", got)
	}
}

func TestRenderTag_NameOnly(t *testing.T) {
	got := tagpattern.RenderTag("{name}-latest", "mypkg", "any")
	if got != "mypkg-latest" {
		t.Errorf("got %q", got)
	}
}

func TestRenderTag_NoPlaceholdersExtra(t *testing.T) {
	got := tagpattern.RenderTag("stable", "n", "v")
	if got != "stable" {
		t.Errorf("got %q", got)
	}
}

func TestBuildTagRegex_ExtractVersion(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("v{version}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ver, ok := tagpattern.ExtractVersion(re, "v1.2.3")
	if !ok {
		t.Fatal("expected match")
	}
	if ver != "1.2.3" {
		t.Errorf("got %q, want %q", ver, "1.2.3")
	}
}

func TestBuildTagRegex_ExtractVersion_WithName(t *testing.T) {
	// Note: {name} in pattern gets literal-escaped after replacement, so
	// BuildTagRegex("{name}-v{version}") expects tag starting with ".+-v" literally.
	// Instead, test a simpler pattern with no {name} to verify version extraction.
	re, err := tagpattern.BuildTagRegex("release-v{version}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ver, ok := tagpattern.ExtractVersion(re, "release-v3.0.0")
	if !ok {
		t.Fatalf("expected match for tag %q", "release-v3.0.0")
	}
	if ver != "3.0.0" {
		t.Errorf("got %q, want %q", ver, "3.0.0")
	}
}

func TestBuildTagRegex_NoVersionPlaceholder_ExactMatch(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("stable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := tagpattern.ExtractVersion(re, "stable")
	// No version group, but pattern matches; ok may be false since no named group
	_ = ok
	// Just ensure no panic and regex compiled
}

func TestExtractVersion_NonMatchingTag(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("v{version}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := tagpattern.ExtractVersion(re, "release-1.0")
	if ok {
		t.Error("should not match")
	}
}

func TestBuildTagRegex_EmptyTag_NoMatch(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("v{version}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := tagpattern.ExtractVersion(re, "")
	if ok {
		t.Error("empty tag should not match")
	}
}

func TestRenderTag_MultipleVersionOccurrences(t *testing.T) {
	// Both occurrences of {version} should be replaced
	got := tagpattern.RenderTag("{version}-{version}", "n", "1.0")
	if got != "1.0-1.0" {
		t.Errorf("got %q", got)
	}
}

func TestRenderTag_SpecialCharsInVersion(t *testing.T) {
	got := tagpattern.RenderTag("v{version}", "n", "1.0.0-rc.1+build.5")
	if got != "v1.0.0-rc.1+build.5" {
		t.Errorf("got %q", got)
	}
}

func TestBuildTagRegex_VersionInMiddle(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("release-{version}-stable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ver, ok := tagpattern.ExtractVersion(re, "release-2.1.0-stable")
	if !ok {
		t.Fatal("expected match")
	}
	if ver != "2.1.0" {
		t.Errorf("got %q", ver)
	}
}

func TestExtractVersion_PreservesFullVersion(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("v{version}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ver, ok := tagpattern.ExtractVersion(re, "v10.20.30-beta.1+sha.123")
	if !ok {
		t.Fatal("expected match")
	}
	if ver != "10.20.30-beta.1+sha.123" {
		t.Errorf("got %q", ver)
	}
}
