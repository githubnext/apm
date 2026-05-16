package installphase_test

import (
	"sort"
	"testing"

	"github.com/githubnext/apm/internal/install/phases/installphase"
)

func TestParseTargetsField_StringSingle(t *testing.T) {
	data := map[string]interface{}{"targets": "claude"}
	got := installphase.ParseTargetsField(data)
	if len(got) != 1 || got[0] != "claude" {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestParseTargetsField_StringCSV(t *testing.T) {
	data := map[string]interface{}{"targets": "claude, vscode, cursor"}
	got := installphase.ParseTargetsField(data)
	if len(got) != 3 {
		t.Fatalf("expected 3 targets, got %v", got)
	}
}

func TestParseTargetsField_Slice(t *testing.T) {
	data := map[string]interface{}{"targets": []interface{}{"vscode", "claude"}}
	got := installphase.ParseTargetsField(data)
	if len(got) != 2 {
		t.Fatalf("expected 2 targets, got %v", got)
	}
}

func TestParseTargetsField_FallbackTargetKey(t *testing.T) {
	data := map[string]interface{}{"target": "cursor"}
	got := installphase.ParseTargetsField(data)
	if len(got) != 1 || got[0] != "cursor" {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestParseTargetsField_Missing(t *testing.T) {
	data := map[string]interface{}{}
	got := installphase.ParseTargetsField(data)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestValidateTargets_AllKnown(t *testing.T) {
	unknown := installphase.ValidateTargets([]string{"claude", "vscode", "cursor"})
	if len(unknown) != 0 {
		t.Fatalf("expected no unknown targets, got %v", unknown)
	}
}

func TestValidateTargets_Unknown(t *testing.T) {
	unknown := installphase.ValidateTargets([]string{"claude", "unknowntool"})
	if len(unknown) != 1 || unknown[0] != "unknowntool" {
		t.Fatalf("expected [unknowntool], got %v", unknown)
	}
}

func TestExpandAllTarget_NoAll(t *testing.T) {
	in := []string{"claude", "vscode"}
	got := installphase.ExpandAllTarget(in)
	if len(got) != 2 {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestExpandAllTarget_WithAll(t *testing.T) {
	got := installphase.ExpandAllTarget([]string{"all"})
	// Should expand to all known targets except "all" itself
	if len(got) == 0 {
		t.Fatal("expected non-empty expansion of 'all'")
	}
	sort.Strings(got)
	for _, t2 := range got {
		if t2 == "all" {
			t.Fatal("'all' should not appear in expanded list")
		}
	}
}

func TestFormatProvenance(t *testing.T) {
	cases := []struct {
		src  installphase.TargetSource
		val  string
		want string
	}{
		{installphase.TargetSourceCLI, "claude", "from --target flag: claude"},
		{installphase.TargetSourceYAML, "vscode", "from apm.yml targets field: vscode"},
		{installphase.TargetSourceEnv, "cursor", "from APM_TARGET environment variable: cursor"},
		{installphase.TargetSourceDetect, "codex", "auto-detected: codex"},
	}
	for _, c := range cases {
		got := installphase.FormatProvenance(c.src, c.val)
		if got != c.want {
			t.Errorf("FormatProvenance(%v,%q) = %q; want %q", c.src, c.val, got, c.want)
		}
	}
}
