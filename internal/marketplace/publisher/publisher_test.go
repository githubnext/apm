package publisher

import (
	"fmt"
	"strings"
	"testing"
)

func TestBumpPatch_basic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.0.0", "1.0.1"},
		{"v2.3.4", "v2.3.5"},
		{"0.0.0", "0.0.1"},
		{"10.20.30", "10.20.31"},
	}
	for _, tc := range tests {
		got, err := BumpPatch(tc.input)
		if err != nil {
			t.Errorf("BumpPatch(%q) error: %v", tc.input, err)
			continue
		}
		if got != tc.expected {
			t.Errorf("BumpPatch(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestBumpPatch_invalid(t *testing.T) {
	_, err := BumpPatch("not-semver")
	if err == nil {
		t.Error("expected error for invalid semver")
	}
}

func TestRenderTag_substitution(t *testing.T) {
	got := RenderTag("v{version}", "1.2.3")
	if got != "v1.2.3" {
		t.Errorf("RenderTag = %q, want %q", got, "v1.2.3")
	}
}

func TestRenderTag_no_placeholder(t *testing.T) {
	got := RenderTag("release", "1.0.0")
	if got != "release" {
		t.Errorf("RenderTag without placeholder = %q", got)
	}
}

func TestRenderReport_nil(t *testing.T) {
	got := RenderReport(nil)
	if got != "" {
		t.Errorf("RenderReport(nil) should be empty, got %q", got)
	}
}

func TestRenderReport_success(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Repo: "owner/repo", Branch: "apm/update-1.0.1", Status: StatusSuccess},
		},
	}
	got := RenderReport(r)
	if !strings.Contains(got, "owner/repo") {
		t.Errorf("report should contain repo name, got %q", got)
	}
	if !strings.Contains(got, "[+]") {
		t.Errorf("success should have [+] prefix, got %q", got)
	}
}

func TestRenderReport_failed(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Repo: "owner/repo2", Status: StatusFailed, Error: fmt.Errorf("push failed")},
		},
	}
	got := RenderReport(r)
	if !strings.Contains(got, "[x]") {
		t.Errorf("failure should have [x] prefix, got %q", got)
	}
}

func TestRenderReport_skipped(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Repo: "owner/repo3", Status: StatusSkipped, Reason: "already up-to-date"},
		},
	}
	got := RenderReport(r)
	if !strings.Contains(got, "[i]") {
		t.Errorf("skipped should have [i] prefix, got %q", got)
	}
	if !strings.Contains(got, "already up-to-date") {
		t.Errorf("reason should be in report, got %q", got)
	}
}

func TestPublishReport_OK(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Status: StatusSuccess},
			{Status: StatusSkipped},
		},
	}
	if !r.OK() {
		t.Error("expected OK() = true when no failures")
	}

	r.Results = append(r.Results, PublishResult{Status: StatusFailed})
	if r.OK() {
		t.Error("expected OK() = false when there is a failure")
	}
}
