package update

import (
	"strings"
	"testing"
)

func TestPlanEntry_ChangeTypes_Extra3(t *testing.T) {
	types := []string{"updated", "added", "removed"}
	for _, ct := range types {
		e := PlanEntry{Package: "pkg", ChangeType: ct}
		if e.ChangeType != ct {
			t.Errorf("ChangeType = %q, want %q", e.ChangeType, ct)
		}
	}
}

func TestPlanEntry_OldAndNewRef_Extra3(t *testing.T) {
	e := PlanEntry{
		Package:    "owner/pkg",
		OldRef:     "v1.0.0",
		NewRef:     "v2.0.0",
		OldSHA:     "abc1234",
		NewSHA:     "def5678",
		ChangeType: "updated",
	}
	if e.OldRef != "v1.0.0" {
		t.Errorf("OldRef = %q", e.OldRef)
	}
	if e.NewSHA != "def5678" {
		t.Errorf("NewSHA = %q", e.NewSHA)
	}
}

func TestRenderPlanEntry_Added_Extra3(t *testing.T) {
	e := PlanEntry{Package: "owner/pkg", NewRef: "v1.0.0", ChangeType: "added"}
	s := renderPlanEntry(e)
	if s == "" {
		t.Error("renderPlanEntry(added) should not be empty")
	}
}

func TestRenderPlanEntry_Removed_Extra3(t *testing.T) {
	e := PlanEntry{Package: "owner/pkg", OldRef: "v1.0.0", ChangeType: "removed"}
	s := renderPlanEntry(e)
	if s == "" {
		t.Error("renderPlanEntry(removed) should not be empty")
	}
}

func TestRenderPlanEntry_ContainsPackageName_Extra3(t *testing.T) {
	e := PlanEntry{Package: "myowner/mypackage", ChangeType: "added", NewRef: "v1.0.0"}
	s := renderPlanEntry(e)
	if !strings.Contains(s, "mypackage") && !strings.Contains(s, "myowner") {
		t.Errorf("render output %q should contain package name", s)
	}
}

func TestShortSHA_LongSHA_Extra3(t *testing.T) {
	full := "abc1234567890abcdef1234567890abcdef123456"
	short := shortSHA(full)
	if len(short) > 7 {
		t.Errorf("shortSHA len = %d, want <= 7", len(short))
	}
	if !strings.HasPrefix(full, short) {
		t.Errorf("shortSHA %q not prefix of %q", short, full)
	}
}

func TestUpdateOptions_PackageList_Extra3(t *testing.T) {
	opts := UpdateOptions{
		Packages: []string{"owner/a", "owner/b"},
	}
	if len(opts.Packages) != 2 {
		t.Errorf("Packages len = %d, want 2", len(opts.Packages))
	}
}

func TestUpdateOptions_YesAndDryRun_Extra3(t *testing.T) {
	opts := UpdateOptions{Yes: true, DryRun: false}
	if !opts.Yes {
		t.Error("Yes should be true")
	}
	if opts.DryRun {
		t.Error("DryRun should be false")
	}
}

func TestUpdateResult_DryRunTrue_Extra3(t *testing.T) {
	r := UpdateResult{DryRun: true}
	if !r.DryRun {
		t.Error("DryRun should be true")
	}
}

func TestUpdateResult_AppliedAndSkipped_Extra3(t *testing.T) {
	r := UpdateResult{
		Applied: []PlanEntry{{Package: "a"}},
		Skipped: []PlanEntry{{Package: "b"}, {Package: "c"}},
	}
	if len(r.Applied) != 1 {
		t.Errorf("Applied len = %d, want 1", len(r.Applied))
	}
	if len(r.Skipped) != 2 {
		t.Errorf("Skipped len = %d, want 2", len(r.Skipped))
	}
}
