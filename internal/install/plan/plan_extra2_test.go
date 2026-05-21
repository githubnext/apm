package plan

import (
	"strings"
	"testing"
)

func TestActionConstants(t *testing.T) {
	if ActionUpdate != "update" {
		t.Errorf("unexpected ActionUpdate %q", ActionUpdate)
	}
	if ActionAdd != "add" {
		t.Errorf("unexpected ActionAdd %q", ActionAdd)
	}
	if ActionRemove != "remove" {
		t.Errorf("unexpected ActionRemove %q", ActionRemove)
	}
	if ActionUnchanged != "unchanged" {
		t.Errorf("unexpected ActionUnchanged %q", ActionUnchanged)
	}
}

func TestPlanEntry_HasChanges_Update(t *testing.T) {
	e := PlanEntry{Action: ActionUpdate}
	if !e.HasChanges() {
		t.Error("expected HasChanges true for update")
	}
}

func TestPlanEntry_HasChanges_Add(t *testing.T) {
	e := PlanEntry{Action: ActionAdd}
	if !e.HasChanges() {
		t.Error("expected HasChanges true for add")
	}
}

func TestPlanEntry_HasChanges_Remove(t *testing.T) {
	e := PlanEntry{Action: ActionRemove}
	if !e.HasChanges() {
		t.Error("expected HasChanges true for remove")
	}
}

func TestPlanEntry_HasChanges_Unchanged(t *testing.T) {
	e := PlanEntry{Action: ActionUnchanged}
	if e.HasChanges() {
		t.Error("expected HasChanges false for unchanged")
	}
}

func TestPlanEntry_ShortOldCommit_Empty(t *testing.T) {
	e := PlanEntry{OldResolvedCommit: ""}
	if e.ShortOldCommit() != "-" {
		t.Errorf("expected '-' for empty commit, got %q", e.ShortOldCommit())
	}
}

func TestPlanEntry_ShortOldCommit_Long(t *testing.T) {
	e := PlanEntry{OldResolvedCommit: "abcdef1234567890"}
	short := e.ShortOldCommit()
	if len(short) != 7 {
		t.Errorf("expected 7 chars, got %d: %q", len(short), short)
	}
	if short != "abcdef1" {
		t.Errorf("expected abcdef1, got %q", short)
	}
}

func TestPlanEntry_ShortNewCommit_Short(t *testing.T) {
	e := PlanEntry{NewResolvedCommit: "abc"}
	if e.ShortNewCommit() != "abc" {
		t.Errorf("expected abc, got %q", e.ShortNewCommit())
	}
}

func TestUpdatePlan_HasChanges_WhenAdd(t *testing.T) {
	p := UpdatePlan{Entries: []PlanEntry{{Action: ActionAdd}}}
	if !p.HasChanges() {
		t.Error("expected HasChanges true")
	}
}

func TestUpdatePlan_HasChanges_AllUnchanged(t *testing.T) {
	p := UpdatePlan{Entries: []PlanEntry{{Action: ActionUnchanged}}}
	if p.HasChanges() {
		t.Error("expected HasChanges false")
	}
}

func TestUpdatePlan_ChangedEntries_FiltersUnchanged(t *testing.T) {
	p := UpdatePlan{Entries: []PlanEntry{
		{Action: ActionUnchanged, DepKey: "a"},
		{Action: ActionAdd, DepKey: "b"},
	}}
	changed := p.ChangedEntries()
	if len(changed) != 1 {
		t.Errorf("expected 1 changed entry, got %d", len(changed))
	}
	if changed[0].DepKey != "b" {
		t.Errorf("expected depkey b, got %q", changed[0].DepKey)
	}
}

func TestUpdatePlan_SummaryCounts_Add(t *testing.T) {
	p := UpdatePlan{Entries: []PlanEntry{
		{Action: ActionAdd},
		{Action: ActionAdd},
		{Action: ActionRemove},
	}}
	counts := p.SummaryCounts()
	if counts[ActionAdd] != 2 {
		t.Errorf("expected add=2, got %d", counts[ActionAdd])
	}
	if counts[ActionRemove] != 1 {
		t.Errorf("expected remove=1, got %d", counts[ActionRemove])
	}
}

func TestRenderPlanText_ContainsAdd(t *testing.T) {
	p := UpdatePlan{Entries: []PlanEntry{
		{Action: ActionAdd, DepKey: "org/pkg", DisplayName: "org/pkg"},
	}}
	text := RenderPlanText(p, false)
	if !strings.Contains(text, "[+]") {
		t.Errorf("expected [+] in plan text, got: %q", text)
	}
}

func TestLockedDependency_ZeroValue(t *testing.T) {
	var d LockedDependency
	if d.Key != "" || d.RepoURL != "" {
		t.Error("expected zero value fields to be empty")
	}
}

func TestDependencyReference_IsLocal(t *testing.T) {
	d := DependencyReference{LocalPath: "/local/path", IsLocal: true}
	if !d.IsLocal {
		t.Error("expected IsLocal true")
	}
	if d.LocalPath != "/local/path" {
		t.Errorf("unexpected LocalPath %q", d.LocalPath)
	}
}
