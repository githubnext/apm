package heal

import "testing"

func TestHealMessageLevel_Constants(t *testing.T) {
	if HealMessageInfo != 0 {
		t.Error("HealMessageInfo should be 0")
	}
	if HealMessageWarn != 1 {
		t.Error("HealMessageWarn should be 1")
	}
}

func TestHealMessage_ZeroValue(t *testing.T) {
	var m HealMessage
	if m.Level != HealMessageInfo || m.Text != "" || m.PackageKey != "" {
		t.Error("HealMessage zero value should have empty fields")
	}
}

func TestNewHealContext_Fields(t *testing.T) {
	ctx := NewHealContext("pkg-a", true, false, true, false)
	if ctx.PackageKey != "pkg-a" {
		t.Error("PackageKey mismatch")
	}
	if !ctx.LockfileMatch {
		t.Error("LockfileMatch should be true")
	}
	if ctx.LockfileMatchViaContentHashOnly {
		t.Error("LockfileMatchViaContentHashOnly should be false")
	}
	if !ctx.UpdateRefs {
		t.Error("UpdateRefs should be true")
	}
	if ctx.BypassKeys == nil {
		t.Error("BypassKeys should be initialized")
	}
	if ctx.FiredGroups == nil {
		t.Error("FiredGroups should be initialized")
	}
}

func TestHealContext_AddWarn_Multi(t *testing.T) {
	ctx := NewHealContext("pkg", false, false, false, false)
	ctx.AddWarn("warn1", "pkg")
	ctx.AddWarn("warn2", "pkg")
	if len(ctx.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(ctx.Messages))
	}
	if ctx.Messages[0].Level != HealMessageWarn {
		t.Error("first message level should be warn")
	}
	if ctx.Messages[1].Text != "warn2" {
		t.Error("second message text mismatch")
	}
}

func TestHealContext_AddInfo_Multi(t *testing.T) {
	ctx := NewHealContext("pkg", false, false, false, false)
	ctx.AddInfo("info-msg", "pkg")
	if len(ctx.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(ctx.Messages))
	}
	if ctx.Messages[0].Level != HealMessageInfo {
		t.Error("message level should be info")
	}
}

func TestRunHealChain_NoHealers(t *testing.T) {
	ctx := NewHealContext("pkg", false, false, false, false)
	lm, rc := RunHealChain(nil, &ctx)
	if lm || rc {
		t.Error("RunHealChain with no healers should return false, false")
	}
}

func TestHealContext_BypassKeys_Usage(t *testing.T) {
	ctx := NewHealContext("pkg", false, false, false, false)
	ctx.BypassKeys["key1"] = true
	if !ctx.BypassKeys["key1"] {
		t.Error("BypassKeys set/get mismatch")
	}
	if ctx.BypassKeys["key2"] {
		t.Error("key2 should not be in BypassKeys")
	}
}

func TestHealContext_FiredGroups_Usage(t *testing.T) {
	ctx := NewHealContext("pkg", false, false, false, false)
	ctx.FiredGroups["group-a"] = true
	if !ctx.FiredGroups["group-a"] {
		t.Error("FiredGroups set/get mismatch")
	}
}

func TestHealContext_Messages_MixedLevels(t *testing.T) {
	ctx := NewHealContext("pkg", false, false, false, false)
	ctx.AddWarn("w1", "pkg")
	ctx.AddInfo("i1", "pkg")
	ctx.AddWarn("w2", "pkg")
	warns := 0
	infos := 0
	for _, m := range ctx.Messages {
		switch m.Level {
		case HealMessageWarn:
			warns++
		case HealMessageInfo:
			infos++
		}
	}
	if warns != 2 || infos != 1 {
		t.Errorf("expected 2 warns 1 info, got warns=%d infos=%d", warns, infos)
	}
}
