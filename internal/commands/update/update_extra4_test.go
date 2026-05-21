package update

import "testing"

func TestUpdateOptions_YesField_Extra4(t *testing.T) {
opts := UpdateOptions{Yes: true}
if !opts.Yes {
t.Error("expected Yes true")
}
}

func TestUpdateOptions_DryRunField_Extra4(t *testing.T) {
opts := UpdateOptions{DryRun: true}
if !opts.DryRun {
t.Error("expected DryRun true")
}
}

func TestUpdateOptions_PackageFilter_Extra4(t *testing.T) {
opts := UpdateOptions{Packages: []string{"org/pkg1", "org/pkg2"}}
if len(opts.Packages) != 2 {
t.Errorf("expected 2 packages, got %d", len(opts.Packages))
}
}

func TestUpdateResult_AppliedCount_Extra4(t *testing.T) {
r := UpdateResult{Applied: []PlanEntry{{ChangeType: "updated"}}}
if len(r.Applied) != 1 {
t.Errorf("expected 1 applied, got %d", len(r.Applied))
}
}

func TestUpdateResult_SkippedCount_Extra4(t *testing.T) {
r := UpdateResult{Skipped: []PlanEntry{{}, {}}}
if len(r.Skipped) != 2 {
t.Errorf("expected 2 skipped, got %d", len(r.Skipped))
}
}

func TestUpdateResult_DryRunField_Extra4(t *testing.T) {
r := UpdateResult{DryRun: true}
if !r.DryRun {
t.Error("expected DryRun true")
}
}

func TestPlanEntry_ChangeTypeUpdated_Extra4(t *testing.T) {
e := PlanEntry{ChangeType: "updated"}
if e.ChangeType != "updated" {
t.Errorf("unexpected ChangeType: %s", e.ChangeType)
}
}

func TestPlanEntry_ChangeTypeAdded_Extra4(t *testing.T) {
e := PlanEntry{ChangeType: "added"}
if e.ChangeType != "added" {
t.Errorf("unexpected ChangeType: %s", e.ChangeType)
}
}

func TestPlanEntry_ChangeTypeRemoved_Extra4(t *testing.T) {
e := PlanEntry{ChangeType: "removed"}
if e.ChangeType != "removed" {
t.Errorf("unexpected ChangeType: %s", e.ChangeType)
}
}

func TestPlanEntry_PackageField_Extra4(t *testing.T) {
e := PlanEntry{Package: "org/myrepo"}
if e.Package != "org/myrepo" {
t.Errorf("unexpected Package: %s", e.Package)
}
}

func TestPlanEntry_OldAndNewRef_Extra4b(t *testing.T) {
e := PlanEntry{OldRef: "v1.0", NewRef: "v2.0"}
if e.OldRef != "v1.0" || e.NewRef != "v2.0" {
t.Errorf("unexpected refs: old=%s new=%s", e.OldRef, e.NewRef)
}
}

func TestPlanEntry_OldAndNewSHA_Extra4(t *testing.T) {
e := PlanEntry{OldSHA: "abc1234", NewSHA: "def5678"}
if e.OldSHA != "abc1234" || e.NewSHA != "def5678" {
t.Errorf("unexpected SHAs: old=%s new=%s", e.OldSHA, e.NewSHA)
}
}
