package heal_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/phases/heal"
)

// countingHealer records how many times Applies and Execute are called.
type countingHealer struct {
	group        string
	appliesReturn bool
	appliesCalls  int
	executeCalls  int
}

func (h *countingHealer) ExclusiveGroup() string            { return h.group }
func (h *countingHealer) Applies(_ *heal.HealContext) bool  { h.appliesCalls++; return h.appliesReturn }
func (h *countingHealer) Execute(_ *heal.HealContext)       { h.executeCalls++ }

func TestHealContext_UpdateRefs(t *testing.T) {
	hctx := heal.NewHealContext("dep", false, false, true, false)
	if !hctx.UpdateRefs {
		t.Error("expected UpdateRefs=true")
	}
}

func TestHealContext_RefChanged(t *testing.T) {
	hctx := heal.NewHealContext("dep", false, false, false, true)
	if !hctx.RefChanged {
		t.Error("expected RefChanged=true")
	}
}

func TestHealContext_LockfileMatchViaContentHashOnly(t *testing.T) {
	hctx := heal.NewHealContext("dep", false, true, false, false)
	if !hctx.LockfileMatchViaContentHashOnly {
		t.Error("expected LockfileMatchViaContentHashOnly=true")
	}
}

func TestHealContext_AddMultipleMessages(t *testing.T) {
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	hctx.AddWarn("w1", "pkg")
	hctx.AddInfo("i1", "pkg")
	hctx.AddWarn("w2", "pkg")
	if len(hctx.Messages) != 3 {
		t.Errorf("expected 3 messages, got %d", len(hctx.Messages))
	}
	if hctx.Messages[0].Level != heal.HealMessageWarn {
		t.Error("first should be Warn")
	}
	if hctx.Messages[1].Level != heal.HealMessageInfo {
		t.Error("second should be Info")
	}
}

func TestHealContext_BypassKeys(t *testing.T) {
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	hctx.BypassKeys["bypass-a"] = true
	if !hctx.BypassKeys["bypass-a"] {
		t.Error("expected bypass key set")
	}
	if hctx.BypassKeys["bypass-b"] {
		t.Error("expected bypass-b unset")
	}
}

func TestRunHealChain_AppliesFalseShouldNotExecute(t *testing.T) {
	h := &countingHealer{appliesReturn: false}
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	heal.RunHealChain([]heal.Healer{h}, &hctx)
	if h.appliesCalls != 1 {
		t.Errorf("Applies should be called once, got %d", h.appliesCalls)
	}
	if h.executeCalls != 0 {
		t.Errorf("Execute should not be called when Applies=false, got %d", h.executeCalls)
	}
}

func TestRunHealChain_AppliesTrueShouldExecute(t *testing.T) {
	h := &countingHealer{appliesReturn: true}
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	heal.RunHealChain([]heal.Healer{h}, &hctx)
	if h.executeCalls != 1 {
		t.Errorf("Execute should be called once, got %d", h.executeCalls)
	}
}

func TestRunHealChain_MultipleHealersNoGroup(t *testing.T) {
	healers := []*countingHealer{
		{appliesReturn: true},
		{appliesReturn: false},
		{appliesReturn: true},
	}
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	chain := make([]heal.Healer, len(healers))
	for i, h := range healers {
		chain[i] = h
	}
	heal.RunHealChain(chain, &hctx)
	if healers[0].executeCalls != 1 {
		t.Errorf("healer[0] expected 1 execute, got %d", healers[0].executeCalls)
	}
	if healers[1].executeCalls != 0 {
		t.Errorf("healer[1] should not execute (Applies=false), got %d", healers[1].executeCalls)
	}
	if healers[2].executeCalls != 1 {
		t.Errorf("healer[2] expected 1 execute, got %d", healers[2].executeCalls)
	}
}

func TestRunHealChain_DifferentGroups(t *testing.T) {
	h1 := &countingHealer{group: "g1", appliesReturn: true}
	h2 := &countingHealer{group: "g2", appliesReturn: true}
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	heal.RunHealChain([]heal.Healer{h1, h2}, &hctx)
	if h1.executeCalls != 1 {
		t.Errorf("h1 (g1) should execute once, got %d", h1.executeCalls)
	}
	if h2.executeCalls != 1 {
		t.Errorf("h2 (g2) should execute once, got %d", h2.executeCalls)
	}
}

func TestRunHealChain_FiredGroupsPropagated(t *testing.T) {
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	// Check FiredGroups starts empty
	if len(hctx.FiredGroups) != 0 {
		t.Errorf("expected empty FiredGroups, got %v", hctx.FiredGroups)
	}
	h := &countingHealer{group: "exclusive", appliesReturn: true}
	heal.RunHealChain([]heal.Healer{h}, &hctx)
	if !hctx.FiredGroups["exclusive"] {
		t.Error("FiredGroups should contain 'exclusive' after healer ran")
	}
}

func TestHealMessage_PackageKey(t *testing.T) {
	hctx := heal.NewHealContext("mypkg", false, false, false, false)
	hctx.AddWarn("warning text", "mypkg")
	if hctx.Messages[0].PackageKey != "mypkg" {
		t.Errorf("expected PackageKey='mypkg', got %q", hctx.Messages[0].PackageKey)
	}
}

func TestHealMessage_Text(t *testing.T) {
	hctx := heal.NewHealContext("pkg", false, false, false, false)
	hctx.AddInfo("info text here", "pkg")
	if hctx.Messages[0].Text != "info text here" {
		t.Errorf("expected 'info text here', got %q", hctx.Messages[0].Text)
	}
}
