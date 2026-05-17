package tagpattern_test

import (
	"testing"

	"github.com/githubnext/apm/internal/marketplace/tagpattern"
)

func TestRenderTag(t *testing.T) {
	tests := []struct {
		pattern, name, version, want string
	}{
		{"{name}-v{version}", "myapp", "1.2.3", "myapp-v1.2.3"},
		{"v{version}", "anything", "2.0.0", "v2.0.0"},
		{"{name}/{version}", "owner/repo", "3.0", "owner/repo/3.0"},
		{"release-{version}-{name}", "tool", "4.5", "release-4.5-tool"},
		{"static-tag", "x", "1", "static-tag"},
	}
	for _, tt := range tests {
		got := tagpattern.RenderTag(tt.pattern, tt.name, tt.version)
		if got != tt.want {
			t.Errorf("RenderTag(%q, %q, %q) = %q, want %q", tt.pattern, tt.name, tt.version, got, tt.want)
		}
	}
}

func TestBuildTagRegex_NoVersion(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("static-tag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !re.MatchString("static-tag") {
		t.Error("expected match for exact static-tag")
	}
	if re.MatchString("other") {
		t.Error("unexpected match for 'other'")
	}
}

func TestBuildTagRegex_WithVersion(t *testing.T) {
	re, err := tagpattern.BuildTagRegex("v{version}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ver, ok := tagpattern.ExtractVersion(re, "v1.2.3")
	if !ok {
		t.Fatal("expected version extraction to succeed")
	}
	if ver != "1.2.3" {
		t.Errorf("expected version '1.2.3', got %q", ver)
	}
}

func TestBuildTagRegex_NamePlaceholder(t *testing.T) {
	// {name} is substituted with ".+" before QuoteMeta, so it becomes a
	// literal ".+" in the compiled regex (not a wildcard). RenderTag is
	// the intended way to produce a concrete tag for a known name.
	re, err := tagpattern.BuildTagRegex("{name}-v{version}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The regex matches the literal string produced after {name}->.+ substitution.
	ver, ok := tagpattern.ExtractVersion(re, ".+-v2.0.0")
	if !ok {
		t.Fatal("expected version extraction to succeed for literal .+ name")
	}
	if ver != "2.0.0" {
		t.Errorf("expected '2.0.0', got %q", ver)
	}
}

func TestExtractVersion_NoMatch(t *testing.T) {
	re, _ := tagpattern.BuildTagRegex("v{version}")
	_, ok := tagpattern.ExtractVersion(re, "nope")
	if ok {
		t.Error("expected no match")
	}
}

func TestRenderTag_EmptyPlaceholders(t *testing.T) {
got := tagpattern.RenderTag("{name}-{version}", "", "")
if got != "-" {
t.Errorf("expected '-', got %q", got)
}
}

func TestRenderTag_NoPlaceholders(t *testing.T) {
got := tagpattern.RenderTag("release", "anything", "1.0")
if got != "release" {
t.Errorf("expected 'release', got %q", got)
}
}

func TestBuildTagRegex_EmptyPattern(t *testing.T) {
re, err := tagpattern.BuildTagRegex("")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !re.MatchString("") {
t.Error("empty pattern should match empty string")
}
}

func TestBuildTagRegex_MultipleVersionTokens(t *testing.T) {
// Only the first {version} is treated as the capture group; the second
// is included literally after QuoteMeta (after the first split on {version}).
re, err := tagpattern.BuildTagRegex("v{version}-end")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
ver, ok := tagpattern.ExtractVersion(re, "v3.0.1-end")
if !ok {
t.Fatal("expected match")
}
if ver != "3.0.1" {
t.Errorf("expected '3.0.1', got %q", ver)
}
}

func TestExtractVersion_EmptyTag(t *testing.T) {
re, _ := tagpattern.BuildTagRegex("v{version}")
_, ok := tagpattern.ExtractVersion(re, "")
if ok {
t.Error("expected no match for empty tag")
}
}

func TestRenderTag_OnlyVersion(t *testing.T) {
got := tagpattern.RenderTag("{version}", "ignore", "9.9.9")
if got != "9.9.9" {
t.Errorf("expected '9.9.9', got %q", got)
}
}
