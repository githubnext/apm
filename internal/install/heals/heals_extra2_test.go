package heals

import (
	"testing"
)

func TestHealMessageLevel_InfoIsZero(t *testing.T) {
	if HealMessageInfo != 0 {
		t.Error("expected HealMessageInfo to be 0")
	}
}

func TestHealMessageLevel_WarnGreaterThanInfo(t *testing.T) {
	if HealMessageWarn <= HealMessageInfo {
		t.Error("expected Warn > Info")
	}
}

func TestNewHealContext_ZeroPackageKey(t *testing.T) {
	hctx := NewHealContext("", false, false, false)
	if hctx.PackageKey != "" {
		t.Errorf("expected empty PackageKey, got %q", hctx.PackageKey)
	}
}

func TestNewHealContext_RefTypes(t *testing.T) {
	hctx := NewHealContext("pkg", false, false, false)
	_ = hctx.ResolvedRefType // just verify field exists
}

func TestHealContext_BypassKeys_Initialized(t *testing.T) {
	hctx := NewHealContext("pkg", false, false, false)
	if hctx.BypassKeys == nil {
		t.Error("expected non-nil BypassKeys")
	}
}

func TestHealContext_FiredGroups_Initialized(t *testing.T) {
	hctx := NewHealContext("pkg", false, false, false)
	if hctx.FiredGroups == nil {
		t.Error("expected non-nil FiredGroups")
	}
}

func TestHealContext_Emit_Count(t *testing.T) {
	hctx := NewHealContext("pkg", false, false, false)
	hctx.Emit(HealMessageInfo, "msg1")
	hctx.Emit(HealMessageWarn, "msg2")
	if len(hctx.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(hctx.Messages))
	}
}

func TestHealContext_Emit_Level(t *testing.T) {
	hctx := NewHealContext("pkg", false, false, false)
	hctx.Emit(HealMessageWarn, "warning message")
	if hctx.Messages[0].Level != HealMessageWarn {
		t.Errorf("expected Warn level, got %v", hctx.Messages[0].Level)
	}
}

func TestRunHealChain_NoOp(t *testing.T) {
	hctx := NewHealContext("pkg", false, false, false)
	RunHealChain(&hctx, []Heal{})
	if len(hctx.Messages) != 0 {
		t.Error("expected no messages after empty chain")
	}
}

func TestDefaultHealChain_NotEmpty(t *testing.T) {
	chain := DefaultHealChain()
	if len(chain) == 0 {
		t.Error("expected non-empty default heal chain")
	}
}

func TestBranchRefDriftHeal_Interface(t *testing.T) {
	var h Heal = BranchRefDriftHeal{}
	if h.Name() == "" {
		t.Error("expected non-empty name")
	}
	if h.Order() < 0 {
		t.Error("expected non-negative order")
	}
}

func TestBuggyLockfileRecoveryHeal_Interface(t *testing.T) {
	var h Heal = BuggyLockfileRecoveryHeal{}
	if h.Name() == "" {
		t.Error("expected non-empty name")
	}
	if h.ExclusiveGroup() == "" {
		t.Error("expected non-empty exclusive group")
	}
}

func TestHealMessage_Fields(t *testing.T) {
	msg := HealMessage{Level: HealMessageInfo, Text: "test"}
	if msg.Text != "test" {
		t.Errorf("expected 'test', got %q", msg.Text)
	}
}
