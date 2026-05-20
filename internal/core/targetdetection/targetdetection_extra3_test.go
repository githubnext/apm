package targetdetection

import (
	"testing"
)

func TestNormalizeTarget_Vscode(t *testing.T) {
	if NormalizeTarget("vscode") != "vscode" {
		t.Errorf("expected vscode, got %q", NormalizeTarget("vscode"))
	}
}

func TestNormalizeTarget_Agents(t *testing.T) {
	if NormalizeTarget("agents") != "vscode" {
		t.Errorf("expected vscode for 'agents', got %q", NormalizeTarget("agents"))
	}
}

func TestNormalizeTarget_Claude(t *testing.T) {
	if NormalizeTarget("claude") != "claude" {
		t.Errorf("expected claude, got %q", NormalizeTarget("claude"))
	}
}

func TestNormalizeTarget_Gemini(t *testing.T) {
	if NormalizeTarget("gemini") != "gemini" {
		t.Errorf("expected gemini, got %q", NormalizeTarget("gemini"))
	}
}

func TestValidTargets_ContainsClaude(t *testing.T) {
	if !ValidTargets["claude"] {
		t.Error("expected claude in ValidTargets")
	}
}

func TestValidTargets_ContainsAll(t *testing.T) {
	if !ValidTargets["all"] {
		t.Error("expected 'all' in ValidTargets")
	}
}

func TestCanonicalTargetsOrdered_NotEmpty(t *testing.T) {
	if len(CanonicalTargetsOrdered) == 0 {
		t.Error("expected non-empty CanonicalTargetsOrdered")
	}
}

func TestCanonicalDeployDirs_ContainsClaude(t *testing.T) {
	if CanonicalDeployDirs["claude"] == "" {
		t.Error("expected deploy dir for claude")
	}
}

func TestFormatProvenance_NoTargets(t *testing.T) {
	resolved := ResolvedTargets{}
	result := FormatProvenance(resolved)
	_ = result
}
