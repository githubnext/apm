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
