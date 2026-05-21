package publisher

import (
	"testing"
	"time"
)

func TestPublishStatusConstants_Extra2(t *testing.T) {
	if StatusSuccess == StatusFailed {
		t.Error("StatusSuccess and StatusFailed should differ")
	}
	if StatusSkipped == StatusSuccess {
		t.Error("StatusSkipped and StatusSuccess should differ")
	}
	if StatusSkipped == StatusFailed {
		t.Error("StatusSkipped and StatusFailed should differ")
	}
}

func TestDefaultOptions_Concurrency_Extra2(t *testing.T) {
	opts := DefaultOptions()
	if opts.Concurrency <= 0 {
		t.Errorf("Concurrency = %d, want > 0", opts.Concurrency)
	}
}

func TestPublishReport_OK_AllSuccess_Extra2(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Status: StatusSuccess, Repo: "repo1"},
			{Status: StatusSuccess, Repo: "repo2"},
		},
	}
	if !r.OK() {
		t.Error("all-success report should be OK")
	}
}

func TestPublishReport_OK_WithFailed_Extra2(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Status: StatusSuccess, Repo: "repo1"},
			{Status: StatusFailed, Repo: "repo2"},
		},
	}
	if r.OK() {
		t.Error("report with failures should not be OK")
	}
}

func TestPublishReport_Empty_Extra2(t *testing.T) {
	r := &PublishReport{}
	if !r.OK() {
		t.Error("empty report should be OK (no failures)")
	}
}

func TestConsumerUpdate_Fields_Extra2(t *testing.T) {
	cu := ConsumerUpdate{
		Repo:       "org/consumer-repo",
		BranchName: "autoloop/updates",
		OldVersion: "1.0.0",
		NewVersion: "1.1.0",
	}
	if cu.Repo != "org/consumer-repo" {
		t.Errorf("Repo = %q", cu.Repo)
	}
	if cu.OldVersion != "1.0.0" {
		t.Errorf("OldVersion = %q", cu.OldVersion)
	}
}

func TestPublishResult_Fields_Extra2(t *testing.T) {
	pr := PublishResult{
		Status:  StatusSkipped,
		Repo:    "org/repo",
		Skipped: true,
		Reason:  "already up to date",
	}
	if pr.Status != StatusSkipped {
		t.Errorf("Status = %q", pr.Status)
	}
	if pr.Reason != "already up to date" {
		t.Errorf("Reason = %q", pr.Reason)
	}
}

func TestMarketplaceYML_Fields(t *testing.T) {
	m := MarketplaceYML{
		Name:      "my-marketplace",
		Version:   "2.0.0",
		Consumers: []string{"org/c1", "org/c2"},
	}
	if m.Name != "my-marketplace" {
		t.Errorf("Name = %q", m.Name)
	}
	if len(m.Consumers) != 2 {
		t.Errorf("Consumers len = %d", len(m.Consumers))
	}
}

func TestPublishState_Fields(t *testing.T) {
	ps := PublishState{
		Version:   "3.0.0",
		Consumers: map[string]string{"org/repo": "autoloop/v3"},
	}
	if ps.Version != "3.0.0" {
		t.Errorf("Version = %q", ps.Version)
	}
	if ps.Consumers["org/repo"] != "autoloop/v3" {
		t.Errorf("Consumers = %v", ps.Consumers)
	}
}

func TestPublishState_LastPublished_Zero(t *testing.T) {
	ps := PublishState{}
	if !ps.LastPublished.IsZero() {
		t.Error("zero PublishState should have zero LastPublished")
	}
}

func TestPublishReport_DurationExtra2(t *testing.T) {
	r := &PublishReport{Duration: 42 * time.Second}
	if r.Duration != 42*time.Second {
		t.Errorf("Duration = %v", r.Duration)
	}
}

func TestPublishOptions_Fields(t *testing.T) {
	o := PublishOptions{Concurrency: 8, DryRun: true}
	if o.Concurrency != 8 {
		t.Errorf("Concurrency = %d", o.Concurrency)
	}
	if !o.DryRun {
		t.Error("DryRun should be true")
	}
}
