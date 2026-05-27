package bundle_test

import (
	"testing"

	"github.com/githubnext/apm/internal/bundle"
)

// TestParityPackResultDefaults verifies PackResult zero value mirrors Python defaults.
func TestParityPackResultDefaults(t *testing.T) {
	r := bundle.PackResult{}
	if r.BundlePath != "" {
		t.Error("expected empty BundlePath")
	}
	if len(r.Files) != 0 {
		t.Error("expected empty Files")
	}
	if r.LockfileEnriched {
		t.Error("expected LockfileEnriched=false")
	}
	if r.MappedCount != 0 {
		t.Error("expected MappedCount=0")
	}
}

// TestParityUnpackResultDefaults verifies UnpackResult zero value mirrors Python defaults.
func TestParityUnpackResultDefaults(t *testing.T) {
	r := bundle.UnpackResult{}
	if r.Verified {
		t.Error("expected Verified=false")
	}
	if r.SkippedCount != 0 {
		t.Error("expected SkippedCount=0")
	}
	if r.SecurityWarnings != 0 {
		t.Error("expected SecurityWarnings=0")
	}
	if r.SecurityCritical != 0 {
		t.Error("expected SecurityCritical=0")
	}
}

// TestParityExtractPackTargetsNil returns empty for nil lockfile.
func TestParityExtractPackTargetsNil(t *testing.T) {
	targets := bundle.ExtractPackTargets(nil)
	if len(targets) != 0 {
		t.Errorf("expected empty targets, got %v", targets)
	}
}

// TestParityExtractPackTargetsNoPack returns empty for lockfile without pack.
func TestParityExtractPackTargetsNoPack(t *testing.T) {
	lf := map[string]interface{}{"deps": []string{}}
	targets := bundle.ExtractPackTargets(lf)
	if len(targets) != 0 {
		t.Errorf("expected empty targets, got %v", targets)
	}
}

// TestParityExtractPackTargetsStringTarget returns single-element slice.
func TestParityExtractPackTargetsStringTarget(t *testing.T) {
	lf := map[string]interface{}{
		"pack": map[string]interface{}{"target": "copilot"},
	}
	targets := bundle.ExtractPackTargets(lf)
	if len(targets) != 1 || targets[0] != "copilot" {
		t.Errorf("expected [copilot], got %v", targets)
	}
}

// TestParityExtractPackTargetsListTarget returns multi-element slice.
func TestParityExtractPackTargetsListTarget(t *testing.T) {
	lf := map[string]interface{}{
		"pack": map[string]interface{}{
			"target": []interface{}{"copilot", "claude"},
		},
	}
	targets := bundle.ExtractPackTargets(lf)
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %v", targets)
	}
	if targets[0] != "copilot" || targets[1] != "claude" {
		t.Errorf("unexpected targets: %v", targets)
	}
}

// TestParityCheckTargetMismatchEmpty returns empty when bundleTargets empty.
func TestParityCheckTargetMismatchEmpty(t *testing.T) {
	msg := bundle.CheckTargetMismatch([]string{}, []string{"copilot"})
	if msg != "" {
		t.Errorf("expected empty warning, got %q", msg)
	}
}

// TestParityCheckTargetMismatchAll returns empty for "all" bundle target.
func TestParityCheckTargetMismatchAll(t *testing.T) {
	msg := bundle.CheckTargetMismatch([]string{"all"}, []string{"copilot"})
	if msg != "" {
		t.Errorf("expected empty warning for 'all', got %q", msg)
	}
}

// TestParityCheckTargetMismatchCovered returns empty when covered.
func TestParityCheckTargetMismatchCovered(t *testing.T) {
	msg := bundle.CheckTargetMismatch([]string{"copilot"}, []string{"copilot", "claude"})
	if msg != "" {
		t.Errorf("expected empty when covered, got %q", msg)
	}
}

// TestParityCheckTargetMismatchMissing returns warning string when missing.
func TestParityCheckTargetMismatchMissing(t *testing.T) {
	msg := bundle.CheckTargetMismatch([]string{"copilot", "claude"}, []string{"copilot"})
	if msg == "" {
		t.Error("expected non-empty warning for missing target")
	}
	// Should mention 'claude' as missing
	if len(msg) < 10 {
		t.Errorf("warning too short: %q", msg)
	}
}

// TestParityIsSafeRelPathValid passes for normal relative paths.
func TestParityIsSafeRelPathValid(t *testing.T) {
	cases := []string{
		"agents/foo.md",
		"skills/bar.json",
		"apm.lock.yaml",
	}
	for _, c := range cases {
		if !bundle.IsSafeRelPath(c) {
			t.Errorf("expected safe: %q", c)
		}
	}
}

// TestParityIsSafeRelPathInvalid rejects unsafe paths.
func TestParityIsSafeRelPathInvalid(t *testing.T) {
	cases := []string{
		"/absolute/path",
		"../traversal",
		"foo/../../etc",
		"C:\\windows\\path",
	}
	for _, c := range cases {
		if bundle.IsSafeRelPath(c) {
			t.Errorf("expected unsafe: %q", c)
		}
	}
}

// TestParityLocalBundleInfoFields verifies LocalBundleInfo struct fields.
func TestParityLocalBundleInfoFields(t *testing.T) {
	info := bundle.LocalBundleInfo{
		SourceDir:   "/tmp/bundle",
		PackageID:   "my-plugin",
		IsArchive:   true,
		PackTargets: []string{"copilot"},
		PluginJSON:  map[string]interface{}{"id": "my-plugin"},
		Lockfile:    nil,
		TempDir:     "/tmp/extract-123",
	}
	if info.PackageID != "my-plugin" {
		t.Errorf("unexpected PackageID: %s", info.PackageID)
	}
	if !info.IsArchive {
		t.Error("expected IsArchive=true")
	}
	if len(info.PackTargets) != 1 || info.PackTargets[0] != "copilot" {
		t.Errorf("unexpected PackTargets: %v", info.PackTargets)
	}
}

// TestParityLocalBundleInfoDefaultEmpty verifies zero LocalBundleInfo is safe.
func TestParityLocalBundleInfoDefaultEmpty(t *testing.T) {
	info := bundle.LocalBundleInfo{}
	if info.IsArchive {
		t.Error("expected IsArchive=false")
	}
	if len(info.PackTargets) != 0 {
		t.Error("expected empty PackTargets")
	}
}

// TestParityCheckTargetMismatchNoInstallTargets warns when install targets empty.
func TestParityCheckTargetMismatchNoInstallTargets(t *testing.T) {
	msg := bundle.CheckTargetMismatch([]string{"copilot"}, []string{})
	if msg == "" {
		t.Error("expected warning when install targets empty")
	}
}

// TestParityExtractPackTargetsEmptyString returns empty for empty-string target.
func TestParityExtractPackTargetsEmptyString(t *testing.T) {
	lf := map[string]interface{}{
		"pack": map[string]interface{}{"target": ""},
	}
	targets := bundle.ExtractPackTargets(lf)
	if len(targets) != 0 {
		t.Errorf("expected empty targets for empty string, got %v", targets)
	}
}
