package installphase_test

import (
	"strings"
	"testing"

	"github.com/githubnext/apm/internal/install/phases/installphase"
)

func TestValidateTargets_AllKnownSet(t *testing.T) {
	unknown := installphase.ValidateTargets([]string{"claude", "vscode", "cursor"})
	if len(unknown) != 0 {
		t.Errorf("expected no unknown targets, got %v", unknown)
	}
}

func TestValidateTargets_SomeUnknown(t *testing.T) {
	unknown := installphase.ValidateTargets([]string{"claude", "notreal"})
	if len(unknown) != 1 || unknown[0] != "notreal" {
		t.Errorf("expected ['notreal'], got %v", unknown)
	}
}

func TestValidateTargets_CaseInsensitive(t *testing.T) {
	unknown := installphase.ValidateTargets([]string{"Claude", "VSCode"})
	if len(unknown) != 0 {
		t.Errorf("target validation should be case-insensitive, got unknowns: %v", unknown)
	}
}

func TestValidateTargets_EmptyList(t *testing.T) {
	unknown := installphase.ValidateTargets([]string{})
	if len(unknown) != 0 {
		t.Errorf("expected no unknown for empty list, got %v", unknown)
	}
}

func TestExpandAllTarget_All(t *testing.T) {
	result := installphase.ExpandAllTarget([]string{"all"})
	if len(result) < 2 {
		t.Errorf("expected multiple targets from 'all', got %v", result)
	}
	for _, t2 := range result {
		if t2 == "all" {
			t.Error("'all' should not appear in expanded list")
		}
	}
}

func TestExpandAllTarget_NoAllTargets(t *testing.T) {
	result := installphase.ExpandAllTarget([]string{"claude", "vscode"})
	if len(result) != 2 {
		t.Errorf("expected unchanged list, got %v", result)
	}
}

func TestFormatProvenance_CLI(t *testing.T) {
	got := installphase.FormatProvenance(installphase.TargetSourceCLI, "claude")
	if !strings.Contains(got, "claude") || !strings.Contains(got, "--target") {
		t.Errorf("expected CLI provenance format, got %q", got)
	}
}

func TestFormatProvenance_YAML(t *testing.T) {
	got := installphase.FormatProvenance(installphase.TargetSourceYAML, "vscode")
	if !strings.Contains(got, "apm.yml") {
		t.Errorf("expected YAML provenance format, got %q", got)
	}
}

func TestFormatProvenance_Env(t *testing.T) {
	got := installphase.FormatProvenance(installphase.TargetSourceEnv, "cursor")
	if !strings.Contains(got, "APM_TARGET") {
		t.Errorf("expected Env provenance format, got %q", got)
	}
}

func TestFormatProvenance_Detect(t *testing.T) {
	got := installphase.FormatProvenance(installphase.TargetSourceDetect, "codex")
	if !strings.Contains(got, "auto-detected") {
		t.Errorf("expected Detect provenance format, got %q", got)
	}
}

func TestParseTargetsField_NilMap(t *testing.T) {
	got := installphase.ParseTargetsField(map[string]interface{}{})
	if got != nil {
		t.Errorf("expected nil for empty map, got %v", got)
	}
}

func TestTarget_Fields(t *testing.T) {
	tgt := installphase.Target{Name: "claude", ConfigDir: "/home/user/.claude"}
	if tgt.Name != "claude" {
		t.Errorf("expected name 'claude', got %q", tgt.Name)
	}
}
