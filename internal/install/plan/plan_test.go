package plan

import (
	"testing"
)

func TestPlanEntryHasChanges(t *testing.T) {
	unchanged := PlanEntry{Action: ActionUnchanged}
	if unchanged.HasChanges() {
		t.Error("ActionUnchanged should not have changes")
	}

	for _, action := range []string{ActionAdd, ActionUpdate, ActionRemove} {
		e := PlanEntry{Action: action}
		if !e.HasChanges() {
			t.Errorf("action=%s should have changes", action)
		}
	}
}

func TestShortCommit(t *testing.T) {
	e := PlanEntry{
		Action:            ActionAdd,
		OldResolvedCommit: "",
		NewResolvedCommit: "abcdef1234567890",
	}
	if got := e.ShortOldCommit(); got != "-" {
		t.Errorf("empty old commit: want '-', got %q", got)
	}
	if got := e.ShortNewCommit(); got != "abcdef1" {
		t.Errorf("short new commit: want %q, got %q", "abcdef1", got)
	}
}

func TestShortCommit_Short(t *testing.T) {
	e := PlanEntry{NewResolvedCommit: "abc"}
	if got := e.ShortNewCommit(); got != "abc" {
		t.Errorf("want 'abc', got %q", got)
	}
}

func TestDepRefKey_Local(t *testing.T) {
	d := DependencyReference{IsLocal: true, LocalPath: "/path/to/local"}
	if got := depRefKey(d); got != "/path/to/local" {
		t.Errorf("local key: want %q, got %q", "/path/to/local", got)
	}
}

func TestDepRefKey_Virtual(t *testing.T) {
	d := DependencyReference{IsVirtual: true, RepoURL: "https://github.com/org/repo", VirtualPath: "sub"}
	want := "https://github.com/org/repo/sub"
	if got := depRefKey(d); got != want {
		t.Errorf("virtual key: want %q, got %q", want, got)
	}
}

func TestDepRefKey_Regular(t *testing.T) {
	d := DependencyReference{RepoURL: "https://github.com/org/repo"}
	if got := depRefKey(d); got != "https://github.com/org/repo" {
		t.Errorf("regular key: got %q", got)
	}
}

func TestLockfileSatisfiesManifest_AllPresent(t *testing.T) {
	locked := map[string]bool{"https://github.com/org/repo": true}
	deps := []DependencyReference{{RepoURL: "https://github.com/org/repo"}}
	ok, reasons := LockfileSatisfiesManifest(locked, deps)
	if !ok || len(reasons) != 0 {
		t.Errorf("expected satisfied, got reasons=%v", reasons)
	}
}

func TestLockfileSatisfiesManifest_Missing(t *testing.T) {
	locked := map[string]bool{}
	deps := []DependencyReference{{RepoURL: "https://github.com/org/repo"}}
	ok, reasons := LockfileSatisfiesManifest(locked, deps)
	if ok || len(reasons) == 0 {
		t.Error("expected unsatisfied with a reason")
	}
}

func TestLockfileSatisfiesManifest_LocalSkipped(t *testing.T) {
	locked := map[string]bool{}
	deps := []DependencyReference{{IsLocal: true, LocalPath: "/local"}}
	ok, _ := LockfileSatisfiesManifest(locked, deps)
	if !ok {
		t.Error("local deps should be skipped in manifest check")
	}
}
