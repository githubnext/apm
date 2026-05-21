package installphase_test

import (
	"testing"

	installphase "github.com/githubnext/apm/internal/install/phases/installphase"
)

func TestParseTargetsField_NilWhenAbsent(t *testing.T) {
	got := installphase.ParseTargetsField(map[string]interface{}{})
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestParseTargetsField_SingleStringTarget(t *testing.T) {
	got := installphase.ParseTargetsField(map[string]interface{}{"targets": "claude"})
	if len(got) != 1 || got[0] != "claude" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestParseTargetsField_SliceTarget(t *testing.T) {
	got := installphase.ParseTargetsField(map[string]interface{}{"targets": []interface{}{"claude", "vscode"}})
	if len(got) != 2 {
		t.Errorf("expected 2 targets, got %v", got)
	}
}

func TestValidateTargets_Empty_v4(t *testing.T) {
	got := installphase.ValidateTargets([]string{})
	if len(got) != 0 {
		t.Errorf("expected no unknown targets for empty slice, got %v", got)
	}
}

func TestFormatProvenance_CLI_v4(t *testing.T) {
	got := installphase.FormatProvenance(installphase.TargetSourceCLI, "claude")
	if got == "" {
		t.Error("expected non-empty provenance string")
	}
}

func TestFormatProvenance_Detect_v4(t *testing.T) {
	got := installphase.FormatProvenance(installphase.TargetSourceDetect, "")
	if got == "" {
		t.Error("expected non-empty provenance for detect")
	}
}

func TestExpandAllTarget_NoAll_v4(t *testing.T) {
	got := installphase.ExpandAllTarget([]string{"claude"})
	if len(got) == 0 {
		t.Error("expected at least one target from ExpandAllTarget")
	}
}
