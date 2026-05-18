package heal_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/phases/heal"
)

func TestNewHealContext_Defaults(t *testing.T) {
	hctx := heal.NewHealContext("pkg-a", true, false, false, false)
	if hctx.PackageKey != "pkg-a" {
		t.Errorf("expected pkg-a, got %s", hctx.PackageKey)
	}
	if !hctx.LockfileMatch {
		t.Error("expected LockfileMatch true")
	}
	if hctx.BypassKeys == nil {
		t.Error("BypassKeys should be initialized")
	}
	if hctx.FiredGroups == nil {
		t.Error("FiredGroups should be initialized")
	}
	if len(hctx.Messages) != 0 {
		t.Error("Messages should start empty")
	}
}

func TestHealContext_AddWarn(t *testing.T) {
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	hctx.AddWarn("some warning", "pkg")
	if len(hctx.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(hctx.Messages))
	}
	if hctx.Messages[0].Level != heal.HealMessageWarn {
		t.Error("expected warn level")
	}
	if hctx.Messages[0].Text != "some warning" {
		t.Errorf("unexpected text: %s", hctx.Messages[0].Text)
	}
}

func TestHealContext_AddInfo(t *testing.T) {
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	hctx.AddInfo("info message", "pkg")
	if len(hctx.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(hctx.Messages))
	}
	if hctx.Messages[0].Level != heal.HealMessageInfo {
		t.Error("expected info level")
	}
}

// noopHealer always applies and does nothing.
type noopHealer struct {
	group  string
	called int
}

func (h *noopHealer) ExclusiveGroup() string           { return h.group }
func (h *noopHealer) Applies(_ *heal.HealContext) bool { h.called++; return true }
func (h *noopHealer) Execute(_ *heal.HealContext)      {}

// conditionalHealer only applies when LockfileMatch is false.
type conditionalHealer struct {
	group  string
	called int
}

func (h *conditionalHealer) ExclusiveGroup() string { return h.group }
func (h *conditionalHealer) Applies(hctx *heal.HealContext) bool {
	return !hctx.LockfileMatch
}
func (h *conditionalHealer) Execute(hctx *heal.HealContext) {
	h.called++
	hctx.LockfileMatch = true
}

func TestRunHealChain_AppliesAll(t *testing.T) {
	h1 := &noopHealer{}
	h2 := &noopHealer{}
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	heal.RunHealChain([]heal.Healer{h1, h2}, &hctx)
	if h1.called != 1 {
		t.Errorf("expected h1 called 1, got %d", h1.called)
	}
	if h2.called != 1 {
		t.Errorf("expected h2 called 1, got %d", h2.called)
	}
}

func TestRunHealChain_ExclusiveGroup(t *testing.T) {
	h1 := &noopHealer{group: "grp"}
	h2 := &noopHealer{group: "grp"}
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	heal.RunHealChain([]heal.Healer{h1, h2}, &hctx)
	// h1 fires (group not yet marked), h2 should not Applies again for same group
	if h1.called != 1 {
		t.Errorf("expected h1 called 1, got %d", h1.called)
	}
	// h2 is skipped because group was fired — so Applies is never called
	if h2.called != 0 {
		t.Errorf("expected h2 called 0 (exclusive group), got %d", h2.called)
	}
}

func TestRunHealChain_ConditionalNotApplied(t *testing.T) {
	h := &conditionalHealer{}
	hctx := heal.NewHealContext("pkg", true, false, false, false)
	// LockfileMatch=true so Applies returns false
	heal.RunHealChain([]heal.Healer{h}, &hctx)
	if h.called != 0 {
		t.Errorf("expected healer not executed, got %d", h.called)
	}
}

func TestRunHealChain_ConditionalApplied(t *testing.T) {
	h := &conditionalHealer{}
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	lm, rc := heal.RunHealChain([]heal.Healer{h}, &hctx)
	if h.called != 1 {
		t.Errorf("expected healer executed once, got %d", h.called)
	}
	if !lm {
		t.Error("expected lockfileMatch=true after heal")
	}
	_ = rc
}

func TestRunHealChain_EmptyChain(t *testing.T) {
	hctx := heal.NewHealContext("pkg", true, false, false, true)
	lm, rc := heal.RunHealChain(nil, &hctx)
	if !lm {
		t.Error("expected original lockfileMatch returned")
	}
	if !rc {
		t.Error("expected original refChanged returned")
	}
}

func TestHealMessageLevelConstants(t *testing.T) {
	if heal.HealMessageInfo == heal.HealMessageWarn {
		t.Error("Info and Warn levels should differ")
	}
}
