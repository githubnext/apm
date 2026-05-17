package targetdetection

import (
	"testing"
)

func TestDetectTarget_explicit(t *testing.T) {
	target, reason := DetectTarget("/tmp", "copilot", "")
	if target != "vscode" {
		t.Errorf("expected vscode got %s", target)
	}
	if reason != "explicit --target flag" {
		t.Errorf("unexpected reason: %s", reason)
	}
}

func TestNormalizeTarget(t *testing.T) {
	cases := map[string]string{
		"copilot": "vscode",
		"agents":  "vscode",
		"vscode":  "vscode",
		"claude":  "claude",
		"cursor":  "cursor",
	}
	for in, want := range cases {
		got := NormalizeTarget(in)
		if got != want {
			t.Errorf("NormalizeTarget(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFormatProvenance(t *testing.T) {
	r := ResolvedTargets{Targets: []string{"claude", "copilot"}, Source: "--target flag"}
	got := FormatProvenance(r)
	want := "Targets: claude, copilot  (source: --target flag)"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestNormalizeTarget_CanonicalTargets(t *testing.T) {
	canonical := []string{"claude", "cursor", "codex", "gemini", "opencode", "windsurf", "all", "minimal"}
	for _, t2 := range canonical {
		got := NormalizeTarget(t2)
		if got != t2 {
			t.Errorf("NormalizeTarget(%q) = %q, want %q (canonical should pass through)", t2, got, t2)
		}
	}
}

func TestFormatProvenance_SingleTarget(t *testing.T) {
	r := ResolvedTargets{Targets: []string{"claude"}, Source: "apm.yml"}
	got := FormatProvenance(r)
	want := "Targets: claude  (source: apm.yml)"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestResolveTargets_Flag(t *testing.T) {
	dir := t.TempDir()
	r, err := ResolveTargets(dir, []string{"claude"}, nil)
	if err != nil {
		t.Fatalf("ResolveTargets: %v", err)
	}
	if len(r.Targets) != 1 || r.Targets[0] != "claude" {
		t.Errorf("unexpected targets: %v", r.Targets)
	}
	if r.Source != "--target flag" {
		t.Errorf("unexpected source: %q", r.Source)
	}
}
