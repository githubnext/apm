package publisher

import (
	"testing"
	"time"
)

func TestBumpPatch_WithPrerelease(t *testing.T) {
	got, err := BumpPatch("1.2.3-beta")
	if err != nil {
		t.Fatalf("BumpPatch error: %v", err)
	}
	if got != "1.2.4-beta" {
		t.Errorf("BumpPatch(1.2.3-beta) = %q, want 1.2.4-beta", got)
	}
}

func TestBumpPatch_LargeNumbers(t *testing.T) {
	got, err := BumpPatch("100.200.999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "100.200.1000" {
		t.Errorf("got %q, want 100.200.1000", got)
	}
}

func TestBumpPatch_ZeroPatch(t *testing.T) {
	got, err := BumpPatch("2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "2.0.1" {
		t.Errorf("got %q, want 2.0.1", got)
	}
}

func TestBumpPatch_EmptyString(t *testing.T) {
	_, err := BumpPatch("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestBumpPatch_OnlyMajorMinor(t *testing.T) {
	_, err := BumpPatch("1.2")
	if err == nil {
		t.Error("expected error for non-semver 1.2")
	}
}

func TestRenderTag_NoPlaceholder(t *testing.T) {
	got := RenderTag("release", "1.0.0")
	if got != "release" {
		t.Errorf("RenderTag = %q, want release", got)
	}
}

func TestRenderTag_MultiplePlaceholders(t *testing.T) {
	got := RenderTag("{version}-{version}", "2.3.4")
	if got != "2.3.4-2.3.4" {
		t.Errorf("RenderTag multiple = %q, want 2.3.4-2.3.4", got)
	}
}

func TestRenderTag_EmptyVersion(t *testing.T) {
	got := RenderTag("v{version}", "")
	if got != "v" {
		t.Errorf("RenderTag empty version = %q, want v", got)
	}
}

func TestPublishReport_OKAllSuccess(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Repo: "a", Status: StatusSuccess},
			{Repo: "b", Status: StatusSkipped},
		},
	}
	if !r.OK() {
		t.Error("report with all success/skipped should be OK")
	}
}

func TestPublishReport_OKWithFailure(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Repo: "a", Status: StatusSuccess},
			{Repo: "b", Status: StatusFailed},
		},
	}
	if r.OK() {
		t.Error("report with failure should not be OK")
	}
}

func TestPublishReport_OKEmptyResults(t *testing.T) {
	r := &PublishReport{}
	if !r.OK() {
		t.Error("empty report should be OK")
	}
}

func TestPublishReport_Duration(t *testing.T) {
	r := &PublishReport{
		StartedAt: time.Now(),
		Duration:  5 * time.Second,
	}
	if r.Duration != 5*time.Second {
		t.Errorf("Duration = %v, want 5s", r.Duration)
	}
}

func TestPublishStatus_Constants(t *testing.T) {
	if StatusSuccess != "success" {
		t.Errorf("StatusSuccess = %q, want success", StatusSuccess)
	}
	if StatusSkipped != "skipped" {
		t.Errorf("StatusSkipped = %q, want skipped", StatusSkipped)
	}
	if StatusFailed != "failed" {
		t.Errorf("StatusFailed = %q, want failed", StatusFailed)
	}
}

func TestRenderReport_Nil(t *testing.T) {
	got := RenderReport(nil)
	if got != "" {
		t.Errorf("RenderReport(nil) = %q, want empty", got)
	}
}

func TestRenderReport_EmptyResults(t *testing.T) {
	got := RenderReport(&PublishReport{})
	if got == "" {
		// Just verify it doesn't panic; empty is acceptable
	}
	_ = got
}

func TestConsumerUpdate_Fields(t *testing.T) {
	u := ConsumerUpdate{
		Repo:        "owner/repo",
		BranchName:  "update/mypkg-1.0.1",
		CommitMsg:   "chore: update mypkg to 1.0.1",
		PackageName: "mypkg",
		NewVersion:  "1.0.1",
		OldVersion:  "1.0.0",
	}
	if u.Repo != "owner/repo" {
		t.Errorf("Repo = %q", u.Repo)
	}
	if u.NewVersion != "1.0.1" {
		t.Errorf("NewVersion = %q", u.NewVersion)
	}
}

func TestPublishResult_Fields(t *testing.T) {
	r := PublishResult{
		Repo:    "owner/repo",
		Status:  StatusSuccess,
		Branch:  "update/pkg-1.0.1",
		Skipped: false,
	}
	if r.Status != StatusSuccess {
		t.Errorf("Status = %v", r.Status)
	}
	if r.Skipped {
		t.Error("Skipped should be false")
	}
}
