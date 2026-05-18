package heals_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/heals"
)

func TestNewHealContext_UpdateRefs(t *testing.T) {
	hctx := heals.NewHealContext("owner/repo@main", false, false, true)
	if !hctx.UpdateRefs {
		t.Error("UpdateRefs should be true")
	}
}

func TestNewHealContext_LockfileMatchFalse(t *testing.T) {
	hctx := heals.NewHealContext("pkg", false, false, false)
	if hctx.LockfileMatch {
		t.Error("LockfileMatch should be false")
	}
}

func TestNewHealContext_ContentHashOnly(t *testing.T) {
	hctx := heals.NewHealContext("pkg", true, true, false)
	if !hctx.LockfileMatchViaContentHashOnly {
		t.Error("LockfileMatchViaContentHashOnly should be true")
	}
}

func TestHealContext_AddBypassKey_Multiple(t *testing.T) {
	hctx := heals.NewHealContext("pkg", true, false, false)
	hctx.AddBypassKey("key1")
	hctx.AddBypassKey("key2")
	hctx.AddBypassKey("key3")
	if len(hctx.BypassKeys) != 3 {
		t.Errorf("expected 3 bypass keys, got %d", len(hctx.BypassKeys))
	}
	if !hctx.BypassKeys["key2"] {
		t.Error("BypassKeys should contain key2")
	}
}

func TestHealContext_AddBypassKey_Idempotent(t *testing.T) {
	hctx := heals.NewHealContext("pkg", true, false, false)
	hctx.AddBypassKey("k")
	hctx.AddBypassKey("k")
	if len(hctx.BypassKeys) != 1 {
		t.Errorf("duplicate add should not increase count, got %d", len(hctx.BypassKeys))
	}
}

func TestHealContext_Emit_Info(t *testing.T) {
	hctx := heals.NewHealContext("pkg", true, false, false)
	hctx.Emit(heals.HealMessageInfo, "info message")
	if len(hctx.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(hctx.Messages))
	}
	if hctx.Messages[0].Level != heals.HealMessageInfo {
		t.Errorf("Level = %v, want Info", hctx.Messages[0].Level)
	}
	if hctx.Messages[0].Text != "info message" {
		t.Errorf("Text = %q, want info message", hctx.Messages[0].Text)
	}
}

func TestHealContext_Emit_Warn(t *testing.T) {
	hctx := heals.NewHealContext("pkg", true, false, false)
	hctx.Emit(heals.HealMessageWarn, "warn message")
	if hctx.Messages[0].Level != heals.HealMessageWarn {
		t.Errorf("Level = %v, want Warn", hctx.Messages[0].Level)
	}
}

func TestHealContext_Emit_PackageKey(t *testing.T) {
	hctx := heals.NewHealContext("owner/repo@1.0", true, false, false)
	hctx.Emit(heals.HealMessageInfo, "msg")
	if hctx.Messages[0].PackageKey != "owner/repo@1.0" {
		t.Errorf("PackageKey = %q, want owner/repo@1.0", hctx.Messages[0].PackageKey)
	}
}

func TestHealContext_Emit_Multiple(t *testing.T) {
	hctx := heals.NewHealContext("pkg", true, false, false)
	hctx.Emit(heals.HealMessageInfo, "msg1")
	hctx.Emit(heals.HealMessageWarn, "msg2")
	hctx.Emit(heals.HealMessageInfo, "msg3")
	if len(hctx.Messages) != 3 {
		t.Errorf("expected 3 messages, got %d", len(hctx.Messages))
	}
	if hctx.Messages[1].Text != "msg2" {
		t.Errorf("Messages[1].Text = %q", hctx.Messages[1].Text)
	}
}

func TestBranchRefDriftHeal_Metadata(t *testing.T) {
	h := heals.BranchRefDriftHeal{}
	if h.Name() != "branch_ref_drift" {
		t.Errorf("Name() = %q", h.Name())
	}
	if h.Order() != 10 {
		t.Errorf("Order() = %d, want 10", h.Order())
	}
	if h.ExclusiveGroup() != "branch_drift" {
		t.Errorf("ExclusiveGroup() = %q, want branch_drift", h.ExclusiveGroup())
	}
}

func TestBuggyLockfileRecoveryHeal_Metadata(t *testing.T) {
	h := heals.BuggyLockfileRecoveryHeal{}
	if h.Name() != "buggy_lockfile_recovery" {
		t.Errorf("Name() = %q", h.Name())
	}
	if h.Order() != 20 {
		t.Errorf("Order() = %d, want 20", h.Order())
	}
	if h.ExclusiveGroup() != "branch_drift" {
		t.Errorf("ExclusiveGroup() = %q, want branch_drift", h.ExclusiveGroup())
	}
}

func TestRunHealChain_EmptyChain(t *testing.T) {
	hctx := heals.NewHealContext("pkg", true, false, false)
	heals.RunHealChain(&hctx, nil)
	if len(hctx.Messages) != 0 {
		t.Errorf("empty chain should produce no messages, got %d", len(hctx.Messages))
	}
}

func TestHealContext_FiredGroupsInitialized(t *testing.T) {
	hctx := heals.NewHealContext("pkg", true, false, false)
	if hctx.FiredGroups == nil {
		t.Error("FiredGroups should be initialized")
	}
	hctx.FiredGroups["branch_drift"] = true
	if !hctx.FiredGroups["branch_drift"] {
		t.Error("FiredGroups should store values")
	}
}

func TestHealMessageLevel_Values(t *testing.T) {
	if heals.HealMessageInfo == heals.HealMessageWarn {
		t.Error("Info and Warn levels should be distinct")
	}
}
