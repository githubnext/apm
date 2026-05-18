package plan

import (
	"strings"
	"testing"
)

func TestBuildUpdatePlan_AllUnchanged(t *testing.T) {
	old := map[string]*LockedDependency{
		"owner/repo": {Key: "owner/repo", RepoURL: "owner/repo", ResolvedRef: "main", ResolvedCommit: "abc123"},
	}
	resolved := []DependencyReference{
		{RepoURL: "owner/repo", ResolvedRefName: "main", ResolvedCommit: "abc123"},
	}
	plan := BuildUpdatePlan(old, resolved)
	if plan.HasChanges() {
		t.Error("expected no changes when commits match")
	}
	counts := plan.SummaryCounts()
	if counts[ActionUnchanged] != 1 {
		t.Errorf("expected 1 unchanged, got %d", counts[ActionUnchanged])
	}
}

func TestBuildUpdatePlan_AddAndRemove(t *testing.T) {
	old := map[string]*LockedDependency{
		"old/pkg": {Key: "old/pkg", RepoURL: "old/pkg"},
	}
	resolved := []DependencyReference{
		{RepoURL: "new/pkg", ResolvedRefName: "main", ResolvedCommit: "deadbeef"},
	}
	plan := BuildUpdatePlan(old, resolved)
	if !plan.HasChanges() {
		t.Error("expected changes")
	}
	counts := plan.SummaryCounts()
	if counts[ActionAdd] != 1 {
		t.Errorf("expected 1 add, got %d", counts[ActionAdd])
	}
	if counts[ActionRemove] != 1 {
		t.Errorf("expected 1 remove, got %d", counts[ActionRemove])
	}
}

func TestBuildUpdatePlan_Update(t *testing.T) {
	old := map[string]*LockedDependency{
		"owner/repo": {RepoURL: "owner/repo", ResolvedRef: "main", ResolvedCommit: "oldsha"},
	}
	resolved := []DependencyReference{
		{RepoURL: "owner/repo", ResolvedRefName: "main", ResolvedCommit: "newsha"},
	}
	plan := BuildUpdatePlan(old, resolved)
	if !plan.HasChanges() {
		t.Fatal("expected changes for commit update")
	}
	counts := plan.SummaryCounts()
	if counts[ActionUpdate] != 1 {
		t.Errorf("expected 1 update, got %d", counts[ActionUpdate])
	}
}

func TestBuildUpdatePlan_Empty(t *testing.T) {
	plan := BuildUpdatePlan(nil, nil)
	if plan.HasChanges() {
		t.Error("empty plan should have no changes")
	}
	if len(plan.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(plan.Entries))
	}
}

func TestRenderPlanText_NoChanges(t *testing.T) {
	plan := UpdatePlan{}
	out := RenderPlanText(plan, false)
	if out != "" {
		t.Errorf("expected empty output for no-change plan without verbose, got %q", out)
	}
}

func TestRenderPlanText_NoChangesVerbose(t *testing.T) {
	plan := UpdatePlan{}
	out := RenderPlanText(plan, true)
	if !strings.Contains(out, "[i]") {
		t.Errorf("expected [i] header in verbose output, got %q", out)
	}
}

func TestRenderPlanText_WithAdd(t *testing.T) {
	plan := UpdatePlan{Entries: []PlanEntry{
		{DepKey: "x/y", Action: ActionAdd, DisplayName: "x/y", NewResolvedRef: "main", NewResolvedCommit: "abc1234"},
	}}
	out := RenderPlanText(plan, false)
	if !strings.Contains(out, "[+]") {
		t.Errorf("expected [+] in add plan output, got %q", out)
	}
	if !strings.Contains(out, "x/y") {
		t.Errorf("expected dep name in output, got %q", out)
	}
}

func TestLockfileSatisfyManifest_AllPresent(t *testing.T) {
	locked := map[string]bool{"a/b": true, "c/d": true}
	deps := []DependencyReference{
		{RepoURL: "a/b"},
		{RepoURL: "c/d"},
	}
	ok, reasons := LockfileSatisfiesManifest(locked, deps)
	if !ok {
		t.Errorf("expected satisfied, got reasons: %v", reasons)
	}
}

func TestLockfileSatisfyManifest_Missing(t *testing.T) {
	locked := map[string]bool{"a/b": true}
	deps := []DependencyReference{
		{RepoURL: "a/b"},
		{RepoURL: "missing/pkg"},
	}
	ok, reasons := LockfileSatisfiesManifest(locked, deps)
	if ok {
		t.Error("expected unsatisfied manifest")
	}
	if len(reasons) != 1 {
		t.Errorf("expected 1 reason, got %d", len(reasons))
	}
}

func TestLockfileSatisfyManifest_LocalDepsSkipped(t *testing.T) {
	locked := map[string]bool{}
	deps := []DependencyReference{
		{RepoURL: "x/y", IsLocal: true, LocalPath: "/local/path"},
	}
	ok, reasons := LockfileSatisfiesManifest(locked, deps)
	if !ok {
		t.Errorf("local deps should be skipped, got reasons: %v", reasons)
	}
}

func TestChangedEntries_FilterWorks(t *testing.T) {
	plan := UpdatePlan{Entries: []PlanEntry{
		{Action: ActionAdd},
		{Action: ActionUnchanged},
		{Action: ActionRemove},
	}}
	changed := plan.ChangedEntries()
	if len(changed) != 2 {
		t.Errorf("expected 2 changed entries, got %d", len(changed))
	}
}
