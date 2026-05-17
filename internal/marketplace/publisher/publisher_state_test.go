package publisher

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// redactToken
// ---------------------------------------------------------------------------

func TestRedactToken_empty(t *testing.T) {
	got := redactToken("hello world", "")
	if got != "hello world" {
		t.Errorf("redactToken with empty token changed string: %q", got)
	}
}

func TestRedactToken_present(t *testing.T) {
	got := redactToken("error: bad token abc123 rejected", "abc123")
	if strings.Contains(got, "abc123") {
		t.Errorf("redactToken did not redact token, got: %q", got)
	}
	if !strings.Contains(got, "[REDACTED]") {
		t.Errorf("redactToken missing [REDACTED] marker, got: %q", got)
	}
}

func TestRedactToken_multiple_occurrences(t *testing.T) {
	got := redactToken("token abc123 and again abc123", "abc123")
	if strings.Contains(got, "abc123") {
		t.Errorf("redactToken left token in string: %q", got)
	}
}

func TestRedactToken_no_match(t *testing.T) {
	got := redactToken("no secret here", "abc123")
	if got != "no secret here" {
		t.Errorf("unexpected modification: %q", got)
	}
}

// ---------------------------------------------------------------------------
// DefaultOptions
// ---------------------------------------------------------------------------

func TestDefaultOptions_values(t *testing.T) {
	opts := DefaultOptions()
	if opts.DryRun {
		t.Error("DryRun should default to false")
	}
	if opts.Concurrency <= 0 {
		t.Errorf("Concurrency should be positive, got %d", opts.Concurrency)
	}
}

// ---------------------------------------------------------------------------
// PublishReport.OK
// ---------------------------------------------------------------------------

func TestPublishReport_OK_empty(t *testing.T) {
	r := &PublishReport{}
	if !r.OK() {
		t.Error("empty report should be OK")
	}
}

func TestPublishReport_OK_only_success(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Status: StatusSuccess},
			{Status: StatusSkipped},
		},
	}
	if !r.OK() {
		t.Error("all-success report should be OK")
	}
}

func TestPublishReport_OK_with_failure(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Status: StatusSuccess},
			{Status: StatusFailed},
		},
	}
	if r.OK() {
		t.Error("report with failure should not be OK")
	}
}

// ---------------------------------------------------------------------------
// LoadPublishState / SavePublishState round-trip
// ---------------------------------------------------------------------------

func TestLoadPublishState_missing(t *testing.T) {
	tmp := t.TempDir()
	state, err := LoadPublishState(filepath.Join(tmp, "nonexistent.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state == nil {
		t.Fatal("expected non-nil state for missing file")
	}
	if len(state.Consumers) != 0 {
		t.Errorf("expected empty consumers, got %v", state.Consumers)
	}
}

func TestLoadPublishState_invalid_json(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.json")
	if err := os.WriteFile(path, []byte("{invalid"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadPublishState(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSaveAndLoadPublishState(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.json")

	now := time.Now().UTC().Truncate(time.Second)
	state := &PublishState{
		LastPublished: now,
		Consumers:     map[string]string{"owner/repo-a": "apm/update-1.2.3"},
		Version:       "1.2.3",
	}

	if err := SavePublishState(path, state); err != nil {
		t.Fatalf("SavePublishState: %v", err)
	}

	loaded, err := LoadPublishState(path)
	if err != nil {
		t.Fatalf("LoadPublishState: %v", err)
	}

	if loaded.Version != "1.2.3" {
		t.Errorf("Version: got %q, want %q", loaded.Version, "1.2.3")
	}
	if branch, ok := loaded.Consumers["owner/repo-a"]; !ok || branch != "apm/update-1.2.3" {
		t.Errorf("Consumers: unexpected %v", loaded.Consumers)
	}
	if !loaded.LastPublished.Equal(now) {
		t.Errorf("LastPublished: got %v, want %v", loaded.LastPublished, now)
	}
}

func TestSavePublishState_creates_parent_dirs(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "subdir", "nested", "state.json")

	state := &PublishState{
		Consumers: map[string]string{},
		Version:   "2.0.0",
	}
	if err := SavePublishState(path, state); err != nil {
		t.Fatalf("SavePublishState with deep path: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("state file not created: %v", err)
	}
}

func TestLoadPublishState_nil_consumers_initialized(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.json")
	// Write state without consumers field
	if err := os.WriteFile(path, []byte(`{"version":"1.0.0"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	state, err := LoadPublishState(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Consumers == nil {
		t.Error("Consumers should be initialized to empty map, not nil")
	}
}

// ---------------------------------------------------------------------------
// BumpPatch edge cases
// ---------------------------------------------------------------------------

func TestBumpPatch_with_prerelease(t *testing.T) {
	got, err := BumpPatch("v1.2.3-beta.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// pre-release suffix is preserved but patch is bumped
	if !strings.HasPrefix(got, "v1.2.4") {
		t.Errorf("BumpPatch with prerelease: got %q, want prefix v1.2.4", got)
	}
}

func TestBumpPatch_zero(t *testing.T) {
	got, err := BumpPatch("0.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "0.0.1" {
		t.Errorf("BumpPatch(0.0.0) = %q, want 0.0.1", got)
	}
}

func TestBumpPatch_large_numbers(t *testing.T) {
	got, err := BumpPatch("100.200.999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "100.200.1000" {
		t.Errorf("BumpPatch(100.200.999) = %q, want 100.200.1000", got)
	}
}

// ---------------------------------------------------------------------------
// RenderTag edge cases
// ---------------------------------------------------------------------------

func TestRenderTag_multiple_placeholders(t *testing.T) {
	got := RenderTag("v{version}-{version}", "1.0.0")
	if got != "v1.0.0-1.0.0" {
		t.Errorf("RenderTag with double placeholder: got %q", got)
	}
}

func TestRenderTag_empty_pattern(t *testing.T) {
	got := RenderTag("", "1.0.0")
	if got != "" {
		t.Errorf("RenderTag empty pattern: got %q", got)
	}
}

// ---------------------------------------------------------------------------
// PublishStatus constants
// ---------------------------------------------------------------------------

func TestPublishStatus_values(t *testing.T) {
	if StatusSuccess == "" {
		t.Error("StatusSuccess should not be empty string")
	}
	if StatusFailed == "" {
		t.Error("StatusFailed should not be empty string")
	}
	if StatusSkipped == "" {
		t.Error("StatusSkipped should not be empty string")
	}
	if StatusSuccess == StatusFailed {
		t.Error("StatusSuccess and StatusFailed should differ")
	}
	if StatusSuccess == StatusSkipped {
		t.Error("StatusSuccess and StatusSkipped should differ")
	}
}

// ---------------------------------------------------------------------------
// RenderReport edge cases
// ---------------------------------------------------------------------------

func TestRenderReport_empty_results(t *testing.T) {
	r := &PublishReport{Results: []PublishResult{}}
	got := RenderReport(r)
	// Should return something, even if just a header
	_ = got
}

func TestRenderReport_mixed(t *testing.T) {
	r := &PublishReport{
		Results: []PublishResult{
			{Repo: "a/b", Status: StatusSuccess, Branch: "apm/v1.0.1"},
			{Repo: "c/d", Status: StatusFailed, Error: os.ErrPermission},
			{Repo: "e/f", Status: StatusSkipped, Reason: "no change"},
		},
	}
	got := RenderReport(r)
	if !strings.Contains(got, "a/b") {
		t.Errorf("report missing success repo: %q", got)
	}
	if !strings.Contains(got, "c/d") {
		t.Errorf("report missing failed repo: %q", got)
	}
	if !strings.Contains(got, "e/f") {
		t.Errorf("report missing skipped repo: %q", got)
	}
}
