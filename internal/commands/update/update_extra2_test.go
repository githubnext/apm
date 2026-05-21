package update

import (
	"strings"
	"testing"
)

func TestRenderPlanEntry_Added(t *testing.T) {
	e := PlanEntry{Package: "mypkg", NewRef: "v2.0", ChangeType: "added"}
	got := renderPlanEntry(e)
	if !strings.Contains(got, "[+]") {
		t.Errorf("added entry should contain [+], got %q", got)
	}
	if !strings.Contains(got, "mypkg") {
		t.Errorf("expected package name in output, got %q", got)
	}
}

func TestRenderPlanEntry_Removed(t *testing.T) {
	e := PlanEntry{Package: "oldpkg", OldRef: "v1.0", ChangeType: "removed"}
	got := renderPlanEntry(e)
	if !strings.Contains(got, "[-]") {
		t.Errorf("removed entry should contain [-], got %q", got)
	}
}

func TestRenderPlanEntry_UpdatedDiffRefs(t *testing.T) {
	e := PlanEntry{Package: "pkg", OldRef: "v1", NewRef: "v2", ChangeType: "updated"}
	got := renderPlanEntry(e)
	if !strings.Contains(got, "v1") || !strings.Contains(got, "v2") {
		t.Errorf("updated diff refs should show both, got %q", got)
	}
}

func TestRenderPlanEntry_UpdatedSameRefShowsSHA(t *testing.T) {
	e := PlanEntry{Package: "pkg", OldRef: "main", NewRef: "main", OldSHA: "abcdef1234567890", NewSHA: "1234567890abcdef", ChangeType: "updated"}
	got := renderPlanEntry(e)
	if !strings.Contains(got, "abcdef1") {
		t.Errorf("same ref update should show short SHA, got %q", got)
	}
}

func TestShortSHA_Short(t *testing.T) {
	got := shortSHA("abc")
	if got != "abc" {
		t.Errorf("short SHA should be unchanged, got %q", got)
	}
}

func TestShortSHA_Empty(t *testing.T) {
	got := shortSHA("")
	if got != "" {
		t.Errorf("empty SHA should return empty, got %q", got)
	}
}

func TestShortSHA_Exactly7(t *testing.T) {
	got := shortSHA("1234567")
	if got != "1234567" {
		t.Errorf("exactly 7 chars should be unchanged, got %q", got)
	}
}

func TestPlanEntry_ZeroValue(t *testing.T) {
	var e PlanEntry
	got := renderPlanEntry(e)
	if got == "" {
		t.Error("zero-value PlanEntry should produce non-empty output")
	}
}

func TestUpdateOptions_ZeroValue(t *testing.T) {
	var opts UpdateOptions
	if opts.Yes || opts.DryRun || opts.Verbose {
		t.Error("zero-value UpdateOptions should have false booleans")
	}
}

func TestUpdateResult_Fields(t *testing.T) {
	r := UpdateResult{
		Applied: []PlanEntry{{Package: "p", ChangeType: "added"}},
		DryRun:  true,
	}
	if len(r.Applied) != 1 {
		t.Errorf("expected 1 applied entry, got %d", len(r.Applied))
	}
	if !r.DryRun {
		t.Error("DryRun should be true")
	}
}
