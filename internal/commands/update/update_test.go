package update

import (
	"testing"
)

func TestUpdateOptionsFields(t *testing.T) {
	opts := UpdateOptions{
		ProjectRoot: "/proj",
		Yes:         true,
		DryRun:      false,
		Verbose:     true,
		Packages:    []string{"pkg-a", "pkg-b"},
	}
	if opts.ProjectRoot != "/proj" {
		t.Errorf("ProjectRoot = %q", opts.ProjectRoot)
	}
	if !opts.Yes {
		t.Error("Yes should be true")
	}
	if len(opts.Packages) != 2 {
		t.Errorf("Packages len = %d, want 2", len(opts.Packages))
	}
}

func TestUpdateResultFields(t *testing.T) {
	entries := []PlanEntry{
		{Package: "pkg", OldRef: "v1", NewRef: "v2", ChangeType: "updated"},
	}
	r := &UpdateResult{Applied: entries, DryRun: false}
	if len(r.Applied) != 1 {
		t.Errorf("Applied len = %d, want 1", len(r.Applied))
	}
	if r.DryRun {
		t.Error("DryRun should be false")
	}
}

func TestUpdateResult_DryRun(t *testing.T) {
	entries := []PlanEntry{
		{Package: "pkg", NewRef: "v2", ChangeType: "added"},
	}
	r := &UpdateResult{Skipped: entries, DryRun: true}
	if !r.DryRun {
		t.Error("DryRun should be true")
	}
	if len(r.Skipped) != 1 {
		t.Errorf("Skipped len = %d, want 1", len(r.Skipped))
	}
}

func TestPlanEntryFields(t *testing.T) {
	e := PlanEntry{
		Package:    "mypkg",
		OldRef:     "v1.0.0",
		NewRef:     "v2.0.0",
		OldSHA:     "deadbeef1234567",
		NewSHA:     "cafebabe1234567",
		ChangeType: "updated",
	}
	if e.Package != "mypkg" {
		t.Errorf("Package = %q", e.Package)
	}
	if e.ChangeType != "updated" {
		t.Errorf("ChangeType = %q", e.ChangeType)
	}
}

func TestShortSHA(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"abc1234def", "abc1234"},
		{"abc1234", "abc1234"},
		{"abc12", "abc12"},
		{"", ""},
	}
	for _, tc := range tests {
		got := shortSHA(tc.in)
		if got != tc.want {
			t.Errorf("shortSHA(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRenderPlanEntry(t *testing.T) {
	tests := []struct {
		e    PlanEntry
		want string
	}{
		{
			PlanEntry{Package: "mypkg", NewRef: "v1.0.0", ChangeType: "added"},
			"[+] mypkg  (new: v1.0.0)",
		},
		{
			PlanEntry{Package: "mypkg", OldRef: "v1.0.0", ChangeType: "removed"},
			"[-] mypkg  (was: v1.0.0)",
		},
		{
			PlanEntry{Package: "mypkg", OldRef: "v1.0.0", NewRef: "v2.0.0", ChangeType: "updated"},
			"[~] mypkg  v1.0.0  ->  v2.0.0",
		},
		{
			PlanEntry{Package: "mypkg", OldRef: "main", NewRef: "main", OldSHA: "abc1234def", NewSHA: "xyz5678abc", ChangeType: "updated"},
			"[~] mypkg  abc1234  ->  xyz5678",
		},
	}
	for _, tc := range tests {
		got := renderPlanEntry(tc.e)
		if got != tc.want {
			t.Errorf("renderPlanEntry(%+v) = %q, want %q", tc.e, got, tc.want)
		}
	}
}
