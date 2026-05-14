package targetdetection

import "testing"

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
