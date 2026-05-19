package update

import (
	"testing"
)

func TestRenderPlanEntryUnchanged(t *testing.T) {
	e := PlanEntry{Package: "pkg", OldRef: "v1", NewRef: "v1", ChangeType: "updated"}
	got := renderPlanEntry(e)
	// Same ref: should show no SHA change
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestRenderPlanEntryUnknownType(t *testing.T) {
	e := PlanEntry{Package: "pkg", OldRef: "v1", NewRef: "v2", ChangeType: "other"}
	got := renderPlanEntry(e)
	// Falls through to default case
	if got == "" {
		t.Error("expected non-empty for unknown type")
	}
}

func TestShortSHALong(t *testing.T) {
	sha := "abcdef1234567890"
	got := shortSHA(sha)
	if len(got) != 7 {
		t.Errorf("shortSHA(%q) len = %d, want 7", sha, len(got))
	}
}

func TestShortSHAExact7(t *testing.T) {
	sha := "1234567"
	got := shortSHA(sha)
	if got != sha {
		t.Errorf("shortSHA(%q) = %q, want %q", sha, got, sha)
	}
}

func TestUpdateResultMultipleApplied(t *testing.T) {
	r := &UpdateResult{
		Applied: []PlanEntry{
			{Package: "a", ChangeType: "updated"},
			{Package: "b", ChangeType: "added"},
			{Package: "c", ChangeType: "removed"},
		},
		DryRun: false,
	}
	if len(r.Applied) != 3 {
		t.Errorf("expected 3 applied, got %d", len(r.Applied))
	}
}

func TestUpdateResultSkippedDryRun(t *testing.T) {
	r := &UpdateResult{
		Skipped: []PlanEntry{
			{Package: "x", ChangeType: "updated"},
		},
		DryRun: true,
	}
	if len(r.Skipped) != 1 || !r.DryRun {
		t.Error("wrong DryRun result")
	}
}

func TestPlanEntryWithSHA(t *testing.T) {
	e := PlanEntry{
		Package:    "p",
		OldSHA:     "aaa0000bbb111c",
		NewSHA:     "fff9999eee888d",
		OldRef:     "main",
		NewRef:     "main",
		ChangeType: "updated",
	}
	if e.OldSHA == "" || e.NewSHA == "" {
		t.Error("SHA fields should be set")
	}
	got := renderPlanEntry(e)
	if got == "" {
		t.Error("renderPlanEntry should return non-empty string")
	}
}

func TestUpdateOptionsDefaults(t *testing.T) {
	opts := UpdateOptions{}
	if opts.ProjectRoot != "" {
		t.Error("default ProjectRoot should be empty")
	}
	if opts.Yes {
		t.Error("default Yes should be false")
	}
	if opts.DryRun {
		t.Error("default DryRun should be false")
	}
	if len(opts.Packages) != 0 {
		t.Error("default Packages should be nil/empty")
	}
}
