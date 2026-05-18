package heals_test

import (
	"testing"

	"github.com/githubnext/apm/internal/install/heals"
)

func TestNewHealContext(t *testing.T) {
	hctx := heals.NewHealContext("owner/repo@pkg", true, false, false)
	if hctx.PackageKey != "owner/repo@pkg" {
		t.Errorf("PackageKey: got %q", hctx.PackageKey)
	}
	if !hctx.LockfileMatch {
		t.Error("LockfileMatch should be true")
	}
	if hctx.LockfileMatchViaContentHashOnly {
		t.Error("LockfileMatchViaContentHashOnly should be false")
	}
	if hctx.BypassKeys == nil {
		t.Error("BypassKeys should be initialized")
	}
	if hctx.FiredGroups == nil {
		t.Error("FiredGroups should be initialized")
	}
	if len(hctx.Messages) != 0 {
		t.Error("Messages should be empty")
	}
}

func TestHealContextEmit(t *testing.T) {
	hctx := heals.NewHealContext("pkg", false, false, false)
	hctx.Emit(heals.HealMessageInfo, "info message")
	hctx.Emit(heals.HealMessageWarn, "warn message")
	if len(hctx.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(hctx.Messages))
	}
	if hctx.Messages[0].Level != heals.HealMessageInfo {
		t.Errorf("first message level: got %d", hctx.Messages[0].Level)
	}
	if hctx.Messages[0].Text != "info message" {
		t.Errorf("first message text: got %q", hctx.Messages[0].Text)
	}
	if hctx.Messages[0].PackageKey != "pkg" {
		t.Errorf("message PackageKey: got %q", hctx.Messages[0].PackageKey)
	}
	if hctx.Messages[1].Level != heals.HealMessageWarn {
		t.Errorf("second message level: got %d", hctx.Messages[1].Level)
	}
}

func TestHealContextAddBypassKey(t *testing.T) {
	hctx := heals.NewHealContext("pkg", false, false, false)
	hctx.AddBypassKey("dep/a")
	hctx.AddBypassKey("dep/b")
	if !hctx.BypassKeys["dep/a"] {
		t.Error("dep/a should be a bypass key")
	}
	if !hctx.BypassKeys["dep/b"] {
		t.Error("dep/b should be a bypass key")
	}
	if hctx.BypassKeys["dep/c"] {
		t.Error("dep/c should not be a bypass key")
	}
}

// mockHeal implements heals.Heal for testing.
type mockHeal struct {
	name      string
	order     int
	group     string
	applies   bool
	executed  bool
}

func (m *mockHeal) Name() string           { return m.name }
func (m *mockHeal) Order() int             { return m.order }
func (m *mockHeal) ExclusiveGroup() string { return m.group }
func (m *mockHeal) Applies(hctx *heals.HealContext) bool { return m.applies }
func (m *mockHeal) Execute(hctx *heals.HealContext)      { m.executed = true }

func TestRunHealChain_ExclusiveGroup(t *testing.T) {
	h1 := &mockHeal{name: "h1", order: 1, group: "grp", applies: true}
	h2 := &mockHeal{name: "h2", order: 2, group: "grp", applies: true}
	hctx := heals.NewHealContext("pkg", false, false, false)
	heals.RunHealChain(&hctx, []heals.Heal{h1, h2})
	if !h1.executed {
		t.Error("h1 should have executed")
	}
	if h2.executed {
		t.Error("h2 should NOT have executed (exclusive group already fired)")
	}
}

func TestRunHealChain_NotApplies(t *testing.T) {
	h1 := &mockHeal{name: "h1", order: 1, group: "", applies: false}
	hctx := heals.NewHealContext("pkg", false, false, false)
	heals.RunHealChain(&hctx, []heals.Heal{h1})
	if h1.executed {
		t.Error("h1 should NOT execute when Applies returns false")
	}
}

func TestRunHealChain_MultipleGroups(t *testing.T) {
	h1 := &mockHeal{name: "h1", order: 1, group: "grp1", applies: true}
	h2 := &mockHeal{name: "h2", order: 2, group: "grp2", applies: true}
	hctx := heals.NewHealContext("pkg", false, false, false)
	heals.RunHealChain(&hctx, []heals.Heal{h1, h2})
	if !h1.executed {
		t.Error("h1 should have executed")
	}
	if !h2.executed {
		t.Error("h2 should have executed (different group)")
	}
}
