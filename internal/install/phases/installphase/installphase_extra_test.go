package installphase_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/phases/installphase"
)

func TestParseTargetsField_EmptyString(t *testing.T) {
	data := map[string]interface{}{"targets": ""}
	got := installphase.ParseTargetsField(data)
	// Empty string may result in a slice with one empty entry or nil -- just no panic
	_ = got
}

func TestParseTargetsField_MultipleSpaces(t *testing.T) {
	data := map[string]interface{}{"targets": "claude,  vscode,  cursor"}
	got := installphase.ParseTargetsField(data)
	if len(got) < 2 {
		t.Fatalf("expected at least 2 targets, got %v", got)
	}
}

func TestParseTargetsField_SliceOfInterfaces(t *testing.T) {
	data := map[string]interface{}{"targets": []interface{}{"claude", "cursor", "vscode"}}
	got := installphase.ParseTargetsField(data)
	if len(got) != 3 {
		t.Fatalf("expected 3 targets, got %v", got)
	}
}

func TestValidateTargets_Empty(t *testing.T) {
	unknown := installphase.ValidateTargets(nil)
	if len(unknown) != 0 {
		t.Fatalf("expected no unknowns for nil input, got %v", unknown)
	}
}

func TestValidateTargets_AllUnknown(t *testing.T) {
	unknown := installphase.ValidateTargets([]string{"bogus1", "bogus2"})
	if len(unknown) != 2 {
		t.Fatalf("expected 2 unknowns, got %v", unknown)
	}
}

func TestExpandAllTarget_Empty(t *testing.T) {
	got := installphase.ExpandAllTarget(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty for nil, got %v", got)
	}
}

func TestExpandAllTarget_NoAllPassthrough(t *testing.T) {
	in := []string{"cursor"}
	got := installphase.ExpandAllTarget(in)
	if len(got) != 1 || got[0] != "cursor" {
		t.Fatalf("expected [cursor], got %v", got)
	}
}

func TestFormatProvenance_CLISource(t *testing.T) {
	got := installphase.FormatProvenance(installphase.TargetSourceCLI, "vscode")
	if !strings.Contains(got, "vscode") {
		t.Errorf("FormatProvenance CLI: expected vscode in %q", got)
	}
}

func TestFormatProvenance_AllSources(t *testing.T) {
	sources := []installphase.TargetSource{
		installphase.TargetSourceCLI,
		installphase.TargetSourceYAML,
		installphase.TargetSourceEnv,
		installphase.TargetSourceDetect,
	}
	for _, s := range sources {
		got := installphase.FormatProvenance(s, "testval")
		if got == "" {
			t.Errorf("FormatProvenance(%v) returned empty string", s)
		}
		if !strings.Contains(got, "testval") {
			t.Errorf("FormatProvenance(%v) should include value 'testval', got %q", s, got)
		}
	}
}

func TestDetectTargetsFromEnv_NoEnv(t *testing.T) {
	// Without APM_TARGET set, should return empty or nil
	got := installphase.DetectTargetsFromEnv()
	_ = got // Just check no panic
}

func TestExpandAllTarget_AlreadyExpanded(t *testing.T) {
	// If the list contains "all", the expanded list should not contain "all"
	got := installphase.ExpandAllTarget([]string{"all", "cursor"})
	for _, t2 := range got {
		if t2 == "all" {
			t.Error("'all' should not appear in expanded list")
		}
	}
}
